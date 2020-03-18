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

package charstreams

import "github.com/hydralang/ptk/common"

// memoizingCharStream is a character stream that wraps another
// character stream and records the characters that it returns.  Once
// the wrapped character stream returns an EOF, the
// memoizingCharStream will loop back around and simply return the
// same characters over and over.
type memoizingCharStream struct {
	chars  []common.Char     // The characters to return
	idx    int               // Current index into chars
	src    common.CharStream // The wrapped character stream
	replay bool              // Flag indicating we're now replaying
}

// NewMemoizingCharStream constructs and returns a common.CharStream
// implementation that wraps another and simply returns all the
// characters it returns.  After the EOF character is returned--which
// is passed through unchanged--the character stream is replayed over
// and over.  If the source character stream returns an error, that
// error is reported immediately.
func NewMemoizingCharStream(src common.CharStream) common.CharStream {
	return &memoizingCharStream{
		chars: []common.Char{},
		src:   src,
	}
}

// Next returns the next character from the stream as a Char, which
// will include the character's location.  If an error was
// encountered, that will also be returned.
func (mcs *memoizingCharStream) Next() (common.Char, error) {
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
	if ch.Rune == common.EOF {
		mcs.replay = true
	}

	return ch, err
}
