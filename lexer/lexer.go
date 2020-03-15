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
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/common"
)

// Lexer represents the actual lexer.
type Lexer interface {
	// Classifier returns the default classifier that will be used
	// to initialize the state.
	Classifier() Classifier

	// Lex returns an object that satisfies the common.TokenStream
	// interface and which reads the specified io.Reader and
	// converts it to tokens.  Tokens represent the "words" of the
	// language being parsed.
	Lex(cs common.CharStream, options ...Option) common.TokenStream
}

// MockLexer is a mock implementation of the Lexer interface.
type MockLexer struct {
	mock.Mock
}

// Classifier returns the default classifier that will be used to
// initialize the state.
func (m *MockLexer) Classifier() Classifier {
	args := m.MethodCalled("Classifier")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

// Lex returns an object that satisfies the common.TokenStream
// interface and which reads the specified io.Reader and converts it
// to tokens.  Tokens represent the "words" of the language being
// parsed.
func (m *MockLexer) Lex(cs common.CharStream, options ...Option) common.TokenStream {
	args := m.MethodCalled("Lex", cs, options)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(common.TokenStream)
	}

	return nil
}
