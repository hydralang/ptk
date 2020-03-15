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

package lexer

import (
	"container/list"

	"github.com/hydralang/ptk/common"
)

// TrackAll is a special value for the max argument to NewBackTracker
// and BackTracker.SetMax that indicates the desire to track all
// characters.
const TrackAll = -1

// btElem is a struct type containing the returned character and error
// from the source character stream.
type btElem struct {
	ch  common.Char // The character returned
	err error       // The error returned
}

// BackTracker is an implementation of common.CharStream that includes
// backtracking capability.  A BackTracker wraps another character
// stream (including another instance of BackTracker), but provides
// additional methods for controlling backtracking.
type BackTracker struct {
	src   common.CharStream // The source character stream
	max   int               // Maximum length to backtrack by
	saved *list.List        // Saved characters
	next  *list.Element     // Next character to return
	pos   int               // Position within the saved characters
	last  btElem            // Last return from source
}

// NewBackTracker wraps another character stream (which may also be a
// BackTracker, if desired) in a BackTracker.  The max parameter
// indicates the maximum number of characters to track; use 0 to track
// no characters, and TrackAll to track all characters.
func NewBackTracker(src common.CharStream, max int) *BackTracker {
	return &BackTracker{
		src:   src,
		max:   max,
		saved: &list.List{},
		last: btElem{
			ch: common.Char{Rune: common.EOF},
		},
	}
}

// Next returns the next character from the stream as a Char, which
// will include the character's location.  If an error was
// encountered, that will also be returned.
func (bt *BackTracker) Next() (ch common.Char, err error) {
	// Check if we're revisiting old friends
	if bt.next != nil {
		ch = bt.next.Value.(btElem).ch
		err = bt.next.Value.(btElem).err
		bt.next = bt.next.Next()
		bt.pos++
		return
	}

	// Need to get a new one from the source
	if bt.src != nil {
		ch, err = bt.src.Next()

		// Save if we need to
		if bt.max != 0 {
			bt.saved.PushBack(btElem{
				ch:  ch,
				err: err,
			})

			// Do any required trimming
			if bt.max > TrackAll && bt.saved.Len() > bt.max {
				bt.saved.Remove(bt.saved.Front())
			} else {
				bt.pos++
			}

			// See if the source is exhausted
			if ch.Rune == common.EOF {
				bt.src = nil
				bt.last = btElem{
					ch: ch,
				}
			}
		}

		return
	}

	// No data to return
	return bt.last.ch, nil
}

// SetMax allows updating the maximum number of characters to allow
// backtracking over.  Setting a TrackAll value will allow all newly
// returned characters to be backtracked over.  If the new value for
// max is less than the previous value, characters at the front of the
// backtracking queue will be discarded to bring the size down to max.
func (bt *BackTracker) SetMax(max int) {
	bt.max = max

	// Do any required trimming
	if bt.max == 0 {
		bt.saved = &list.List{}
		bt.pos = 0
	} else {
		for bt.max > TrackAll && bt.saved.Len() > bt.max {
			bt.saved.Remove(bt.saved.Front())
			bt.pos--
		}
	}
}

// Accept accepts characters from the backtracking queue, leaving only
// the specified number of characters on the queue.
func (bt *BackTracker) Accept(leave int) {
	if bt.max == 0 {
		// Nothing saved
		return
	}

	// Handle the quick case first
	if bt.next == nil && leave == 0 {
		bt.saved = &list.List{}
		bt.pos = 0
		return
	}

	// Determine the stop element
	stop := bt.next
	for leave > 0 {
		// Back up one
		if stop == nil {
			stop = bt.saved.Back()
		} else {
			stop = stop.Prev()
		}

		if stop == nil {
			// Nothing to remove
			return
		}

		leave--
	}

	// Discard elements from the front of the queue up to but
	// excluding stop
	for bt.saved.Front() != stop {
		bt.saved.Remove(bt.saved.Front())
		bt.pos--
	}
}

// Len returns the number of characters saved so far on the
// backtracking queue.
func (bt *BackTracker) Len() int {
	return bt.saved.Len()
}

// Pos returns the position of the most recently returned character
// within the saved character list.
func (bt *BackTracker) Pos() int {
	return bt.pos - 1
}

// BackTrack resets to the beginning of the backtracking queue.
func (bt *BackTracker) BackTrack() {
	bt.next = bt.saved.Front()
	bt.pos = 0
}
