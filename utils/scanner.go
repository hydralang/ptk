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

package utils

import (
	"io"
	"unicode/utf8"

	"github.com/klmitch/kent"

	"github.com/hydralang/ptk/common"
)

// DefaultTabStop is the default tab stop for the scanner.
const DefaultTabStop = 8

// scanBuf is the size of the read buffer to utilize.
const scanBuf = 4096

// Special runes used for internal processing
const (
	errRune  rune = common.EOF + iota // Signals an error
	sentinel                          // Signals saved is unset
)

// Scanner is an implementation of common.CharStream that wraps
// io.Reader objects.  Scanner also contains logic to handle diverse
// line ending styles through the use of LineStyle objects; line
// endings will be converted into single newlines if this support is
// enabled.
type scanner struct {
	src   io.Reader         // The character source
	buf   [scanBuf + 1]byte // The read buffer
	pos   int               // The current index into the read buffer
	end   int               // The end of the buffer
	err   error             // Deferred error
	ts    int               // The tabstop in use
	ls    LineStyle         // Current line ending style
	saved rune              // One-character pushback for line endings
	loc   common.Location   // Location of head of read buffer
	rep   kent.Reporter     // A kent.Reporter for reporting encoding errors
}

// NewScanner constructs a new instance of the Scanner.  It returns a
// common.CharStream object.
func NewScanner(r io.Reader, loc common.Location, options ...ScannerOption) common.CharStream {
	// Construct the scanner
	s := &scanner{
		src:   r,
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    UnknownLineStyle,
		saved: sentinel,
		loc:   loc,
		rep:   kent.Root(),
	}

	// Apply the options
	for _, opt := range options {
		opt(s)
	}

	return s
}

// next is the inner implementation of the scan algorithm.  It returns
// the next character read from the source as a rune.
func (s *scanner) next() rune {
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
					return common.EOF
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
			if s.rep != nil && ch == utf8.RuneError && width == 1 {
				// Report encoding error
				s.rep.Report(common.LocationError(s.loc.Incr(ch, s.ts), ErrBadEncoding))
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
func (s *scanner) Next() (common.Char, error) {
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
		ch = common.EOF
		err = s.err
	}

	// Increment the location by the character
	s.loc = s.loc.Incr(ch, s.ts)

	return common.Char{
		Rune: ch,
		Loc:  s.loc,
	}, err
}
