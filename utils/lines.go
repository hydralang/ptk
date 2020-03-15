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

import "github.com/stretchr/testify/mock"

// LineDis is a type for the return value of LineStyle.Handle.
type LineDis int

// Possible return value codes for LineStyle.Handle.
const (
	LineDisNewline     LineDis = iota // Newline sequence recognized
	LineDisNewlineSave                // Newline followed by newline
	LineDisSpace                      // Substitute a space
	LineDisMore                       // Need another character
)

// LineStyle represents a line ending style.  This controls how lines
// are handled when read from the source input stream, and is used for
// maintaining the location that is then attached to runes.
type LineStyle interface {
	// Handle checks to see if a line ending sequence has been
	// encountered.  It returns a LineDis value, which indicates
	// the disposition of the character; and a LineStyle object to
	// use next time around.
	Handle(chs []rune) (LineDis, LineStyle)
}

// MockLineStyle is a mock implementation of the LineStyle interface.
type MockLineStyle struct {
	mock.Mock
}

// Handle checks to see if a line ending sequence has been
// encountered.  It returns a LineDis value, which indicates the
// disposition of the character; and a LineStyle object to use next
// time around.
func (m *MockLineStyle) Handle(chs []rune) (LineDis, LineStyle) {
	args := m.MethodCalled("Handle", chs)

	if tmp := args.Get(1); tmp != nil {
		return args.Get(0).(LineDis), tmp.(LineStyle)
	}

	return args.Get(0).(LineDis), nil
}

// unixLineStyle is a style for handling the case of bare newline line
// endings, also known as UNIX line endings.
type unixLineStyle struct{}

// Handle checks to see if a line ending sequence has been
// encountered.  It returns a LineDis value, which indicates the
// disposition of the character; and a LineStyle object to use next
// time around.
func (ls *unixLineStyle) Handle(chs []rune) (LineDis, LineStyle) {
	if chs[0] == '\n' {
		return LineDisNewline, ls
	}

	return LineDisSpace, ls
}

// UNIXLineStyle is a style for handling the use of bare newline line
// endings, also known as UNIX line endings.
var UNIXLineStyle = &unixLineStyle{}

// macLineStyle is a style for handling the case of bare carriage
// return line endings, also known as Mac classic line endings.
type macLineStyle struct{}

// Handle checks to see if a line ending sequence has been
// encountered.  It returns a LineDis value, which indicates the
// disposition of the character; and a LineStyle object to use next
// time around.
func (ls *macLineStyle) Handle(chs []rune) (LineDis, LineStyle) {
	if chs[0] == '\r' {
		return LineDisNewline, ls
	}

	return LineDisSpace, ls
}

// MacLineStyle is a style for handling the use of bare carriage
// return line endings, also known as Mac classic line endings.
var MacLineStyle = &macLineStyle{}

// dosLineStyle is a style for handling the use of carriage return and
// newline line endings, also known as DOS line endings.
type dosLineStyle struct{}

// Handle checks to see if a line ending sequence has been
// encountered.  It returns a LineDis value, which indicates the
// disposition of the character; and a LineStyle object to use next
// time around.
func (ls *dosLineStyle) Handle(chs []rune) (LineDis, LineStyle) {
	if len(chs) == 1 {
		// Bare newline, call it good
		if chs[0] == '\n' {
			return LineDisNewline, ls
		}

		// Need another
		return LineDisMore, ls
	}

	// Line ending sequence
	if chs[0] == '\r' && chs[1] == '\n' {
		return LineDisNewline, ls
	}

	// Bare \r in text, convert to a space
	return LineDisSpace, ls
}

// DOSLineStyle is a style for handling the use of carriage return and
// newline line endings, also known as DOS line endings.
var DOSLineStyle = &dosLineStyle{}

// unknownLineStyle is a style for handling files where the line
// ending style is not yet known.  It guesses the line ending style
// from the first line ending encountered, then switches to the
// appropriate line ending style.
type unknownLineStyle struct{}

// Handle checks to see if a line ending sequence has been
// encountered.  It returns a LineDis value, which indicates the
// disposition of the character; and a LineStyle object to use next
// time around.
func (ls *unknownLineStyle) Handle(chs []rune) (LineDis, LineStyle) {
	// If we encountered a newline, we're done
	if chs[0] == '\n' {
		return LineDisNewline, UNIXLineStyle
	}

	// Has to be a carriage return; need to know what's next
	if len(chs) <= 1 {
		return LineDisMore, ls
	}

	// If it's a newline, then switch to DOS style
	if chs[1] == '\n' {
		return LineDisNewline, DOSLineStyle
	}

	// Must be Mac style
	return LineDisNewlineSave, MacLineStyle
}

// UnknownLineStyle is a style for handling files where the line
// ending style is not yet known.  It guesses the line ending style
// from the first line ending encountered, then switches to the
// appropriate line ending style.
var UnknownLineStyle = &unknownLineStyle{}
