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

package common

import (
	"bytes"
	"fmt"

	"github.com/stretchr/testify/mock"
)

// Location is an interface for location data.  Each token and node
// should have attached location data that reports its location.  This
// aids in finding the location of errors.
type Location interface {
	// String constructs a string representation of the location.
	String() string

	// Thru creates a new Location that ranges from the beginning
	// of this location to the beginning of another Location.
	Thru(other Location) (Location, error)

	// ThruEnd is similar to Thru, except that it creates a new
	// Location that ranges from the beginning of this location to
	// the ending of another location.
	ThruEnd(other Location) (Location, error)
}

// MockLocation is a mock for Location
type MockLocation struct {
	mock.Mock
}

// String constructs a string representation of the location.
func (m *MockLocation) String() string {
	args := m.MethodCalled("String")

	return args.String(0)
}

// Thru creates a new Location that ranges from the beginning of this
// location to the beginning of another Location.
func (m *MockLocation) Thru(other Location) (Location, error) {
	args := m.MethodCalled("Thru", other)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Location), args.Error(1)
	}

	return nil, args.Error(1)
}

// ThruEnd is similar to Thru, except that it creates a new Location
// that ranges from the beginning of this location to the ending of
// another location.
func (m *MockLocation) ThruEnd(other Location) (Location, error) {
	args := m.MethodCalled("ThruEnd", other)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Location), args.Error(1)
	}

	return nil, args.Error(1)
}

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

// Advance advances a FileLocation in place.  The current range end
// becomes the range beginning, and the range end is the sum of the
// new range beginning and the provided offset.  The offset is a
// FilePos; if the line number is incremented, the column is reset to
// 1 before offsetting the column number.
func (l *FileLocation) Advance(offset FilePos) {
	// Begin by advancing the beginning
	l.B = l.E

	// Now advance the ending by the offset
	if offset.L > 0 {
		l.E.L += offset.L
		l.E.C = 1
	}
	l.E.C += offset.C
}

// AdvanceTab advances a location in place, as if by a tab character.
// The argument specifies the size of a tab stop.
func (l *FileLocation) AdvanceTab(tabstop int) {
	l.Advance(FilePos{C: 1 + tabstop - l.E.C%tabstop})
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
