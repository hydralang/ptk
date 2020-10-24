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
)

// ArgPos specifies a particular character location within an argument
// list.  It is a component of the ArgLocation type.
type ArgPos struct {
	I int // The index of the argument in its list
	C int // The index of the character within the argument
}

// ArgLocation is an implementation of Location that identifies the
// location of an element within a list of strings, typically command
// line arguments.  It represents a full range, but tab stops and
// newlines are not treated specially.
type ArgLocation struct {
	B ArgPos // Beginning of the range
	E ArgPos // End of the range
}

// String constructs a string representation of the location.
func (l ArgLocation) String() string {
	buf := &bytes.Buffer{}

	// Add the basic prefix to the location buf
	fmt.Fprintf(buf, "%d:%d", l.B.I, l.B.C)

	// Is it split across arguments?  Wider than one character?
	if l.B.I != l.E.I {
		fmt.Fprintf(buf, "-%d:%d", l.E.I, l.E.C)
	} else if l.E.C-l.B.C > 1 {
		fmt.Fprintf(buf, "-%d", l.E.C)
	}

	return buf.String()
}

// Thru creates a new Location that ranges from the beginning of this
// location to the beginning of another Location.
func (l ArgLocation) Thru(other Location) (Location, error) {
	// Verify that other's compatible
	o, ok := other.(ArgLocation)
	if !ok {
		return nil, ErrSplitLocation
	}

	// Create and return a new location
	return ArgLocation{
		B: l.B,
		E: o.B,
	}, nil
}

// ThruEnd is similar to Thru, except that it creates a new Location
// that ranges from the beginning of this location to the ending of
// another location.
func (l ArgLocation) ThruEnd(other Location) (Location, error) {
	// Verify that other's compatible
	o, ok := other.(ArgLocation)
	if !ok {
		return nil, ErrSplitLocation
	}

	// Create and return a new location
	return ArgLocation{
		B: l.B,
		E: o.E,
	}, nil
}

// Incr increments the location by one character.  It is passed the
// character (a rune) and the tabstop size (for handling tabs).  It
// should return a new Location.
func (l ArgLocation) Incr(c rune, tabstop int) Location {
	// Begin by advancing the beginning
	l.B = l.E

	// Advance the ending column
	if c != EOF {
		l.E.C++
	}

	return l
}

// NewArgumentScanner constructs and returns a Scanner implementation
// that returns characters drawn from a provided list of argument
// strings.  This is intended for use with arguments taken from the
// command line, but could be useful in other contexts as well.  The
// strings are logically joined by spaces; to use a different joiner,
// pass that as an option.
func NewArgumentScanner(args []string, options ...ArgOption) Scanner {
	// Process the options
	opts := &argOptions{
		joiner: " ",
		opts:   []FileOption{LineEndings(NoLineStyle)},
	}
	for _, opt := range options {
		opt.argApply(opts)
	}

	// Construct the joiner scanner
	loc := ArgLocation{
		B: ArgPos{I: 0, C: 1},
		E: ArgPos{I: 0, C: 1},
	}
	joiner := NewMemoizingScanner(NewFileScanner(bytes.NewBufferString(opts.joiner), loc, opts.opts...))

	// Construct a list of scanners
	streams := []Scanner{}
	for i, arg := range args {
		if i != 0 {
			streams = append(streams, joiner)
		}

		loc := ArgLocation{
			B: ArgPos{I: i + 1, C: 1},
			E: ArgPos{I: i + 1, C: 1},
		}
		streams = append(streams, NewFileScanner(bytes.NewBufferString(arg), loc, opts.opts...))
	}

	return NewChainingScanner(streams)
}
