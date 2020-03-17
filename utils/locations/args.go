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

package locations

import (
	"bytes"
	"fmt"

	"github.com/hydralang/ptk/common"
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
func (l ArgLocation) Thru(other common.Location) (common.Location, error) {
	// Verify that other's compatible
	o, ok := other.(ArgLocation)
	if !ok {
		return nil, common.ErrSplitLocation
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
func (l ArgLocation) ThruEnd(other common.Location) (common.Location, error) {
	// Verify that other's compatible
	o, ok := other.(ArgLocation)
	if !ok {
		return nil, common.ErrSplitLocation
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
func (l ArgLocation) Incr(c rune, tabstop int) common.Location {
	// Begin by advancing the beginning
	l.B = l.E

	// Advance the ending column
	if c != common.EOF {
		l.E.C++
	}

	return l
}
