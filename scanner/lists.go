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

// ListScanner is a scanner that returns characters from a simple
// list.  It is intended for testing in cases where MockScanner is not
// a good fit.
type ListScanner struct {
	chars []Char // Characters to return
	pos   int    // Position within the character list
	err   error  // Error to return on next EOF
}

// NewListScanner constructs and returns a Scanner implementation that
// returns characters from a list of characters.  The last character
// should be a EOF; this character will be returned with the error
// passed in.  This scanner is intended for testing.
func NewListScanner(chars []Char, err error) *ListScanner {
	return &ListScanner{
		chars: chars,
		err:   err,
	}
}

// Next returns the next character from the stream as a Char, which
// will include the character's location.  If an error was
// encountered, that will also be returned.
func (lcs *ListScanner) Next() (Char, error) {
	// Last character will be returned endlessly
	idx := lcs.pos
	if idx >= len(lcs.chars) {
		idx = len(lcs.chars) - 1
	} else {
		lcs.pos++
	}

	// Select an error to return on first EOF
	var err error
	if lcs.chars[idx].Rune == EOF {
		err = lcs.err
		lcs.err = nil
	}

	return lcs.chars[idx], err
}
