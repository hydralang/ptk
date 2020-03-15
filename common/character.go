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
	"unicode"

	"github.com/stretchr/testify/mock"
)

// EOF is used to signal an end-of-file in the character stream.
const EOF rune = unicode.MaxRune + 1

// Char represents a character retrieved from the source input stream.
// It bundles together a rune and a location.
type Char struct {
	Rune rune     // The rune read from the source
	Loc  Location // The location of the rune within the stream
}

// CharStream presents a stream of characters.  The basic character
// stream does not provide backtracking or character push-back.  The
// Scanner is one implementation of CharStream.
type CharStream interface {
	// Next returns the next character from the stream as a Char,
	// which will include the character's location.  If an error
	// was encountered, that will also be returned.
	Next() (Char, error)
}

// MockCharStream is a mock implementation of the CharStream
// interface.
type MockCharStream struct {
	mock.Mock
}

// Next returns the next character from the stream as a Char, which
// will include the character's location.  If an error was
// encountered, that will also be returned.
func (m *MockCharStream) Next() (Char, error) {
	args := m.MethodCalled("Next")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Char), args.Error(1)
	}

	return Char{}, args.Error(1)
}

// TrackAll is a special value for the max argument to
// BackTracker.SetMax that indicates the desire to track all
// characters.
const TrackAll = -1

// BackTracker is an interface for a backtracker, a CharStream that
// also provides the ability to back up to an earlier character in the
// stream.
type BackTracker interface {
	CharStream

	// SetMax allows updating the maximum number of characters to
	// allow backtracking over.  Setting a TrackAll value will
	// allow all newly returned characters to be backtracked over.
	// If the new value for max is less than the previous value,
	// characters at the front of the backtracking queue will be
	// discarded to bring the size down to max.
	SetMax(max int)

	// Accept accepts characters from the backtracking queue,
	// leaving only the specified number of characters on the
	// queue.
	Accept(leave int)

	// Len returns the number of characters saved so far on the
	// backtracking queue.
	Len() int

	// Pos returns the position of the most recently returned
	// character within the saved character list.
	Pos() int

	// BackTrack resets to the beginning of the backtracking
	// queue.
	BackTrack()
}

// MockBackTracker is a mock implementation of the BackTracker
// interface.
type MockBackTracker struct {
	MockCharStream
}

// SetMax allows updating the maximum number of characters to allow
// backtracking over.  Setting a TrackAll value will allow all newly
// returned characters to be backtracked over.  If the new value for
// max is less than the previous value, characters at the front of the
// backtracking queue will be discarded to bring the size down to max.
func (m *MockBackTracker) SetMax(max int) {
	m.MethodCalled("SetMax", max)
}

// Accept accepts characters from the backtracking queue, leaving only
// the specified number of characters on the queue.
func (m *MockBackTracker) Accept(leave int) {
	m.MethodCalled("Accept", leave)
}

// Len returns the number of characters saved so far on the
// backtracking queue.
func (m *MockBackTracker) Len() int {
	args := m.MethodCalled("Len")

	return args.Int(0)
}

// Pos returns the position of the most recently returned character
// within the saved character list.
func (m *MockBackTracker) Pos() int {
	args := m.MethodCalled("Pos")

	return args.Int(0)
}

// BackTrack resets to the beginning of the backtracking queue.
func (m *MockBackTracker) BackTrack() {
	m.MethodCalled("BackTrack")
}
