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

// MemoizingScanner is a scanner that wraps another scanner and
// records the characters that it returns.  Once the wrapped scanner
// returns an EOF, the MemoizingScanner will loop back around and
// simply return the same characters over and over.
type MemoizingScanner struct {
	chars  []Char  // The characters to return
	idx    int     // Current index into chars
	src    Scanner // The wrapped scanner
	replay bool    // Flag indicating we're now replaying
}

// NewMemoizingScanner constructs and returns a Scanner implementation
// that wraps another and simply returns all the characters it
// returns.  After the EOF character is returned--which is passed
// through unchanged--the scanner is replayed over and over.  If the
// source scanner returns an error, that error is reported
// immediately.
func NewMemoizingScanner(src Scanner) *MemoizingScanner {
	return &MemoizingScanner{
		chars: []Char{},
		src:   src,
	}
}

// Next returns the next character from the stream as a Char, which
// will include the character's location.  If an error was
// encountered, that will also be returned.
func (mcs *MemoizingScanner) Next() (Char, error) {
	// Are we replaying?
	if mcs.replay {
		defer func() {
			// Increment index with wrap-around
			mcs.idx++
			if mcs.idx >= len(mcs.chars) {
				mcs.idx = 0
			}
		}()

		return mcs.chars[mcs.idx], nil
	}

	// Get char from source
	ch, err := mcs.src.Next()

	// Save it
	mcs.chars = append(mcs.chars, ch)

	// End of stream?
	if ch.Rune == EOF {
		mcs.replay = true
	}

	return ch, err
}
