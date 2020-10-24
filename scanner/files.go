// Copyright (c) 2020 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package scanner

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)

// FilePos specifies a particular character location within a file.
// It is a component of the FileLocation type.
type FilePos struct {
	L int // The line number of the position (1-indexed)
	C int // The column number of the position (1-indexed)
}

// FileLocation is an implementation of Location that identifies the
// location of an element within a file.  It represents a full range,
// and has some additional utilities to simplify handling advancement
// within the file, including of tab stops.
type FileLocation struct {
	File string  // Name of the file
	B    FilePos // The beginning of the range
	E    FilePos // The end of the range
}

// String constructs a string representation of the location.
func (l FileLocation) String() string {
	buf := &bytes.Buffer{}

	// Add the basic prefix to the location buf
	fmt.Fprintf(buf, "%s:%d:%d", l.File, l.B.L, l.B.C)

	// Is it split across lines or wider than one column?
	if l.B.L != l.E.L {
		fmt.Fprintf(buf, "-%d:%d", l.E.L, l.E.C)
	} else if l.E.C-l.B.C > 1 {
		fmt.Fprintf(buf, "-%d", l.E.C)
	}

	return buf.String()
}

// Thru creates a new Location that ranges from the beginning of this
// location to the beginning of another Location.
func (l FileLocation) Thru(other Location) (Location, error) {
	// Location can't range across files
	o, ok := other.(FileLocation)
	if !ok || l.File != o.File {
		return nil, ErrSplitLocation
	}

	// Create and return a new location
	return FileLocation{
		File: l.File,
		B:    l.B,
		E:    o.B,
	}, nil
}

// ThruEnd is similar to Thru, except that it creates a new Location
// that ranges from the beginning of this location to the ending of
// another location.
func (l FileLocation) ThruEnd(other Location) (Location, error) {
	// Location can't range across files
	o, ok := other.(FileLocation)
	if !ok || l.File != o.File {
		return nil, ErrSplitLocation
	}

	// Create and return a new location
	return FileLocation{
		File: l.File,
		B:    l.B,
		E:    o.E,
	}, nil
}

// advance is a helper for Incr that advances the file location by the
// designated offset, returning a new file location.
func (l FileLocation) advance(offset FilePos) Location {
	// Begin by advancing the beginning
	l.B = l.E

	// Advance the ending by the offset
	if offset.L > 0 {
		l.E.L += offset.L
		l.E.C = 1
	}
	l.E.C += offset.C

	// Return the new location
	return l
}

// Incr increments the location by one character.  It is passed the
// character (a rune) and the tabstop size (for handling tabs).  It
// should return a new Location.
func (l FileLocation) Incr(c rune, tabstop int) Location {
	switch c {
	case EOF: // End of file
		return l.advance(FilePos{})

	case '\n': // Newline
		return l.advance(FilePos{L: 1})

	case '\t': // Tab
		return l.advance(FilePos{C: 1 + tabstop - l.E.C%tabstop})

	case '\f': // Skip formfeeds at the beginning of lines
		if l.B.C <= 1 {
			return l
		}
		fallthrough
	default: // Everything else advances by one column
		return l.advance(FilePos{C: 1})
	}
}

// DefaultTabStop is the default tab stop for the scanner.
const DefaultTabStop = 8

// scanBuf is the size of the read buffer to utilize.
const scanBuf = 4096

// FileScanner is an implementation of Scanner that wraps io.Reader
// objects.  FileScanner also contains logic to handle diverse line
// ending styles through the use of LineStyle objects; line endings
// will be converted into single newlines if this support is enabled.
type FileScanner struct {
	src   io.Reader            // The character source
	buf   [scanBuf + 1]byte    // The read buffer
	pos   int                  // The current index into the read buffer
	end   int                  // The end of the buffer
	err   error                // Deferred error
	ts    int                  // The tabstop in use
	ls    LineStyle            // Current line ending style
	saved rune                 // One-character pushback for line endings
	loc   Location             // Location of head of read buffer
	enc   EncodingErrorHandler // Handler for encoding errors
}

// NewFileScanner constructs a new instance of the FileScanner.
func NewFileScanner(r io.Reader, loc Location, options ...FileOption) *FileScanner {
	// Construct the scanner
	s := &FileScanner{
		src:   r,
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    UnknownLineStyle,
		saved: sentinel,
		loc:   loc,
	}

	// Apply the options
	for _, opt := range options {
		opt.fileApply(s)
	}

	return s
}

// next is the inner implementation of the scan algorithm.  It returns
// the next character read from the source as a rune.
func (s *FileScanner) next() rune {
	// Convert next byte into a rune; optimized for the common
	// case
	ch, width := rune(s.buf[s.pos]), 1

	// Is it part of a multi-byte sequence or the end of the
	// buffer?
	if ch >= utf8.RuneSelf {
		// Need to get more?
		for s.pos+utf8.UTFMax > s.end && !utf8.FullRune(s.buf[s.pos:s.end]) {
			// A nil src indicates EOF has been
			// encountered
			if s.src == nil {
				if s.err == nil {
					return EOF
				}
				return errRune
			}

			// Shift unread portion to beginning
			copy(s.buf[0:], s.buf[s.pos:s.end])

			// Read more
			bufLen := s.end - s.pos
			readLen, err := s.src.Read(s.buf[bufLen:scanBuf])
			s.pos = 0
			s.end = bufLen + readLen
			s.buf[s.end] = utf8.RuneSelf // mark end of buffer

			// Did we get an error?
			if err != nil {
				// Mark closed
				s.src = nil

				// Save the error
				if err != io.EOF {
					s.err = err
				}
			}
		}

		// OK, try a rune conversion again
		ch = rune(s.buf[s.pos])
		if ch >= utf8.RuneSelf {
			// Not ASCII subset of UTF8
			ch, width = utf8.DecodeRune(s.buf[s.pos:s.end])

			// Was it a decoding error?
			if ch == utf8.RuneError && width == 1 {
				err := LocationError(s.loc.Incr(ch, s.ts), ErrBadEncoding)

				// If we have a handler, call it
				if s.enc != nil {
					err = s.enc.Handle(err)
				}

				// If handler didn't handle it, return
				// the error
				if err != nil {
					s.err = err
					return errRune
				}
			}
		}
	}

	// Update buffer position
	s.pos += width

	return ch
}

// Next returns the next character from the stream as a Char, which
// will include the character's location.  If an error was
// encountered, that will also be returned.
func (s *FileScanner) Next() (Char, error) {
	// Select the next character to process
	var ch rune
	if s.saved != sentinel {
		ch, s.saved = s.saved, sentinel
	} else {
		// Get the next character
		ch = s.next()
	}

	// If it's CR or NL, do line ending handling
	if ch == '\r' || ch == '\n' {
		seq := []rune{ch}
		for {
			dis, ls := s.ls.Handle(seq)
			switch dis {
			case LineDisNewline:
				ch = '\n'
			case LineDisNewlineSave:
				ch = '\n'
				s.saved = seq[1]
			case LineDisSpace:
				ch = ' '
				if len(seq) > 1 {
					s.saved = seq[1]
				}
			case LineDisMore:
				seq = append(seq, s.next())
				continue
			}

			s.ls = ls
			break
		}
	}

	// Are we returning an error?
	var err error
	if ch == errRune {
		ch = EOF
		err = s.err
	}

	// Increment the location by the character
	s.loc = s.loc.Incr(ch, s.ts)

	return Char{
		Rune: ch,
		Loc:  s.loc,
	}, err
}
