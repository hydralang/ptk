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

package tokenstreams

import "github.com/hydralang/ptk/common"

// ChanTokenStreamSize is the size of the input channel.
const ChanTokenStreamSize = 20

// ChanTokenStream is a trivial implementation of TokenStream that
// uses a channel to retrieve tokens.  It implements an extra Push
// method, that allows pushing tokens onto the stream, as well as a
// Done method to signal the token stream that all tokens have been
// pushed.
type ChanTokenStream struct {
	Chan chan *common.Token // The input channel
}

// NewChanTokenStream returns a ChanTokenStream
func NewChanTokenStream() *ChanTokenStream {
	return &ChanTokenStream{
		Chan: make(chan *common.Token, ChanTokenStreamSize),
	}
}

// Next returns the next token.  At the end of the token stream, a nil
// should be returned.
func (q *ChanTokenStream) Next() *common.Token {
	return <-q.Chan
}

// Push pushes a token onto the stream.  It returns true if the push
// was successful; it will return false if Done has been called.
func (q *ChanTokenStream) Push(tok *common.Token) (ok bool) {
	// Panic means we sent to a closed channel
	defer func() {
		if panicData := recover(); panicData != nil {
			ok = false
		}
	}()

	// Send the token to the channel
	ok = true
	q.Chan <- tok

	return
}

// Done indicates to the stream that there will be no more tokens
// pushed onto the queue.
func (q *ChanTokenStream) Done() {
	close(q.Chan)
}
