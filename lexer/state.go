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

	"github.com/hydralang/ptk/internal"
	"github.com/hydralang/ptk/scanner"
)

// State represents the state of the lexer.
type State interface {
	TokenStream

	// Lexer returns the lexer object.
	Lexer() Lexer

	// CharStream retrieves the character stream stored in the
	// state.  Most applications should prefer to utilize the
	// BackTracker instance passed to the Classify, Recognize, or
	// Error methods, but this can be used to obtain a reference
	// to the underlying scanner.Scanner, if that is desired.
	CharStream() scanner.Scanner

	// AppState returns the current application state.
	AppState() interface{}

	// PushAppState allows pushing an alternative application
	// state onto the application state stack.  Use this, paired
	// with PopAppState, when your grammar needs different state
	// in particular sections of the file.
	PushAppState(state interface{})

	// PopAppState allows popping an application state pushed with
	// PushAppState off the application state stack.  Use this,
	// paired with PushAppState, when your grammar needs different
	// state in particular sections of the file.
	PopAppState() interface{}

	// SetAppState allows changing the current application state
	// on the fly.  Its action is similar to a PopAppState
	// followed by a PushAppState, so the number of entries in the
	// application state stack remains the same.
	SetAppState(state interface{}) interface{}

	// Classifier returns the classifier currently in use.
	Classifier() Classifier

	// PushClassifier allows pushing an alternative classifier
	// onto the classifier stack.  Use this, paired with
	// PopClassifier, when your lexer must accommodate multiple
	// classes of tokens in particular sections of the file, e.g.,
	// like PHP embedded in a web page template.
	PushClassifier(cls Classifier)

	// PopClassifier allows popping a classifier pushed with
	// PushClassifier off the classifier stack.  Use this, paired
	// with PushClassifier, when your lexer must accommodate
	// multiple classes of takens in particular sections of the
	// file, e.g., like PHP embedded in a web page template.
	PopClassifier() Classifier

	// SetClassifier allows changing the current classifier on the
	// fly.  Its action is similar to a PopClassifier followed by
	// a PushClassifier, so the number of entries in the
	// classifier stack remains the same.
	SetClassifier(cls Classifier) Classifier

	// Push pushes a token onto the list of tokens to be returned
	// by the lexer.  Recognizers should call this method with the
	// token or tokens that they recognize from the input.
	Push(tok *Token) bool
}

// MockState is a mock implementation of the State interface.
type MockState struct {
	MockTokenStream
}

// Lexer returns the lexer object.
func (m *MockState) Lexer() Lexer {
	args := m.MethodCalled("Lexer")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Lexer)
	}

	return nil
}

// CharStream retrieves the character stream stored in the state.
// Most applications should prefer to utilize the BackTracker instance
// passed to the Classify, Recognize, or Error methods, but this can
// be used to obtain a reference to the underlying scanner.Scanner,
// if that is desired.
func (m *MockState) CharStream() scanner.Scanner {
	args := m.MethodCalled("CharStream")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(scanner.Scanner)
	}

	return nil
}

// AppState returns the current application state.
func (m *MockState) AppState() interface{} {
	args := m.MethodCalled("AppState")

	return args.Get(0)
}

// PushAppState allows pushing an alternative application state onto
// the application state stack.  Use this, paired with PopAppState,
// when your grammar needs different state in particular sections of
// the file.
func (m *MockState) PushAppState(state interface{}) {
	m.MethodCalled("PushAppState", state)
}

// PopAppState allows popping an application state pushed with
// PushAppState off the application state stack.  Use this, paired
// with PushAppState, when your grammar needs different state in
// particular sections of the file.
func (m *MockState) PopAppState() interface{} {
	args := m.MethodCalled("PopAppState")

	return args.Get(0)
}

// SetAppState allows changing the current application state on the
// fly.  Its action is similar to a PopAppState followed by a
// PushAppState, so the number of entries in the application state
// stack remains the same.
func (m *MockState) SetAppState(state interface{}) interface{} {
	args := m.MethodCalled("SetAppState", state)

	return args.Get(0)
}

// Classifier returns the classifier currently in use.
func (m *MockState) Classifier() Classifier {
	args := m.MethodCalled("Classifier")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

// PushClassifier allows pushing an alternative classifier onto the
// classifier stack.  Use this, paired with PopClassifier, when your
// lexer must accommodate multiple classes of tokens in particular
// sections of the file, e.g., like PHP embedded in a web page
// template.
func (m *MockState) PushClassifier(cls Classifier) {
	m.MethodCalled("PushClassifier", cls)
}

// PopClassifier allows popping a classifier pushed with
// PushClassifier off the classifier stack.  Use this, paired with
// PushClassifier, when your lexer must accommodate multiple classes
// of takens in particular sections of the file, e.g., like PHP
// embedded in a web page template.
func (m *MockState) PopClassifier() Classifier {
	args := m.MethodCalled("PopClassifier")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

// SetClassifier allows changing the current classifier on the fly.
// Its action is similar to a PopClassifier followed by a
// PushClassifier, so the number of entries in the classifier stack
// remains the same.
func (m *MockState) SetClassifier(cls Classifier) Classifier {
	args := m.MethodCalled("SetClassifier", cls)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

// Push pushes a token onto the list of tokens to be returned by the
// lexer.  Recognizers should call this method with the token or
// tokens that they recognize from the input.
func (m *MockState) Push(tok *Token) bool {
	args := m.MethodCalled("Push", tok)

	return args.Bool(0)
}

// state is an implementation of State.
type state struct {
	lexer    Lexer           // The lexer being used
	src      scanner.Scanner // The source CharStream
	bt       BackTracker     // Backtracker wrapping the source
	appState internal.Stack  // Stack for application state
	cls      internal.Stack  // Stack for classifier
	toks     *list.List      // List of tokens to produce
}

// NewState constructs and returns a new state, with the specified
// classifier and character stream.
func NewState(lexer Lexer, src scanner.Scanner, options []Option) State {
	obj := &state{
		lexer:    lexer,
		src:      src,
		bt:       NewBackTracker(src, TrackAll),
		appState: internal.NewStack(),
		cls:      internal.NewStack(),
		toks:     &list.List{},
	}

	// Push the initial classifier
	obj.cls.Push(lexer.Classifier())

	// Apply options
	for _, opt := range options {
		opt(obj)
	}

	return obj
}

// next is the actual implementation of the lexer.  This is the
// routine that calls the Classify, Recognize, and Error methods
// provided by the user.
func (s *state) next() {
	// Reset the backtracker
	s.bt.SetMax(TrackAll)

	// Classify the contents
	for _, rec := range s.Classifier().Classify(s, s.bt) {
		s.bt.BackTrack()
		if rec.Recognize(s, s.bt) {
			s.bt.Accept(0)
			return
		}
	}

	// None of the recognizers recognized the contents
	s.bt.BackTrack()
	s.Classifier().Error(s, s.bt)
	s.bt.Accept(0)
}

// Next returns the next token.  At the end of the token stream, a nil
// should be returned.
func (s *state) Next() *Token {
	// Loop until we have a token or all characters have been
	// processed
	for s.toks.Len() <= 0 {
		if !s.bt.More() {
			return nil
		}

		s.next()
	}

	// Return a token off the token queue
	defer func() {
		s.toks.Remove(s.toks.Front())
	}()
	return s.toks.Front().Value.(*Token)
}

// Lexer returns the lexer object.
func (s *state) Lexer() Lexer {
	return s.lexer
}

// CharStream retrieves the character stream stored in the state.
// Most applications should prefer to utilize the BackTracker instance
// passed to the Classify, Recognize, or Error methods, but this can
// be used to obtain a reference to the underlying scanner.Scanner,
// if that is desired.
func (s *state) CharStream() scanner.Scanner {
	return s.src
}

// AppState returns the current application state.
func (s *state) AppState() interface{} {
	return s.appState.Get()
}

// PushAppState allows pushing an alternative application state onto
// the application state stack.  Use this, paired with PopAppState,
// when your grammar needs different state in particular sections of
// the file.
func (s *state) PushAppState(state interface{}) {
	s.appState.Push(state)
}

// PopAppState allows popping an application state pushed with
// PushAppState off the application state stack.  Use this, paired
// with PushAppState, when your grammar needs different state in
// particular sections of the file.
func (s *state) PopAppState() interface{} {
	return s.appState.Pop()
}

// SetAppState allows changing the current application state on the
// fly.  Its action is similar to a PopAppState followed by a
// PushAppState, so the number of entries in the application state
// stack remains the same.
func (s *state) SetAppState(state interface{}) interface{} {
	return s.appState.Set(state)
}

// Classifier returns the classifier currently in use.
func (s *state) Classifier() Classifier {
	if tmp := s.cls.Get(); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

// PushClassifier allows pushing an alternative classifier onto the
// classifier stack.  Use this, paired with PopClassifier, when your
// lexer must accommodate multiple classes of tokens in particular
// sections of the file, e.g., like PHP embedded in a web page
// template.
func (s *state) PushClassifier(cls Classifier) {
	s.cls.Push(cls)
}

// PopClassifier allows popping a classifier pushed with
// PushClassifier off the classifier stack.  Use this, paired with
// PushClassifier, when your lexer must accommodate multiple classes
// of takens in particular sections of the file, e.g., like PHP
// embedded in a web page template.
func (s *state) PopClassifier() Classifier {
	if tmp := s.cls.Pop(); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

// SetClassifier allows changing the current classifier on the fly.
// Its action is similar to a PopClassifier followed by a
// PushClassifier, so the number of entries in the classifier stack
// remains the same.
func (s *state) SetClassifier(cls Classifier) Classifier {
	if tmp := s.cls.Set(cls); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

// Push pushes a token onto the list of tokens to be returned by the
// lexer.  Recognizers should call this method with the token or
// tokens that they recognize from the input.
func (s *state) Push(tok *Token) bool {
	s.toks.PushBack(tok)
	return true
}
