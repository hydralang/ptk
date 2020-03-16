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

import "github.com/hydralang/ptk/common"

// chainingCharStream is a character stream that chains together
// several character streams.  When one character stream returns EOF,
// the chainingCharacterStream skips it and proceeds to the next one.
// Note that if any call to Next returns an error, that error is
// returned immediately.
type chainingCharStream struct {
	streams []common.CharStream // The streams to chain over
	idx     int                 // Index of the stream currently being used
	last    common.Char         // Last character returned
}

// NewChainingCharStream constructs and returns a common.CharStream
// implementation that returns characters from each of the provided
// character streams in turn.  When one character stream returns EOF,
// that EOF is skipped and the chaining character stream begins
// drawing characters from the next one.  If any character stream
// returns an error, that error is returned immediately.
func NewChainingCharStream(streams []common.CharStream) common.CharStream {
	return &chainingCharStream{
		streams: streams,
		last: common.Char{
			Rune: common.EOF,
		},
	}
}

// Next returns the next character from the stream as a Char, which
// will include the character's location.  If an error was
// encountered, that will also be returned.
func (ccs *chainingCharStream) Next() (common.Char, error) {
	var next common.Char
	for {
		// Return EOF if we've exhausted character streams
		if ccs.idx >= len(ccs.streams) {
			return ccs.last, nil
		}

		// Get the next character from the current stream
		var err error
		next, err = ccs.streams[ccs.idx].Next()
		if err != nil {
			return next, err
		}

		// Was it an EOF?
		if next.Rune == common.EOF {
			ccs.last = next
			ccs.idx++
			continue
		}

		break
	}

	return next, nil
}
