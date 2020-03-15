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

import "github.com/hydralang/ptk/common"

// State represents the state of the lexer.
type State interface {
	common.TokenStream

	// Lexer returns the lexer object.
	Lexer() Lexer

	// CharStream retrieves the character stream stored in the
	// state.  Most applications should prefer to utilize the
	// BackTracker instance passed to the Classify, Recognize, or
	// Error methods, but this can be used to obtain a reference
	// to the underlying common.CharStream, if that is desired.
	CharStream() common.CharStream

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
	Push(tok *common.Token) bool
}

// MockState is a mock implementation of the State interface.
type MockState struct {
	common.MockTokenStream
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
// be used to obtain a reference to the underlying common.CharStream,
// if that is desired.
func (m *MockState) CharStream() common.CharStream {
	args := m.MethodCalled("CharStream")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(common.CharStream)
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
func (m *MockState) Push(tok *common.Token) bool {
	args := m.MethodCalled("Push", tok)

	return args.Bool(0)
}
