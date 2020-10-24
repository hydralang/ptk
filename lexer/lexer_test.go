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
	"testing"

	"github.com/klmitch/patcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/scanner"
)

func TestAppState(t *testing.T) {
	s := &mockState{}
	s.On("PushAppState", "state")

	opt := AppState("state")
	opt(s)

	s.AssertExpectations(t)
}

type mockLexer struct {
	mock.Mock
}

func (m *mockLexer) Classifier() Classifier {
	args := m.MethodCalled("Classifier")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

func (m *mockLexer) Lex(cs scanner.Scanner, options ...Option) TokenStream {
	args := m.MethodCalled("Lex", cs, options)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(TokenStream)
	}

	return nil
}

func TestLexerImplementsLexer(t *testing.T) {
	assert.Implements(t, (*Lexer)(nil), &lexer{})
}

func TestNew(t *testing.T) {
	cls := &mockClassifier{}

	result := New(cls)

	assert.Equal(t, &lexer{
		cls: cls,
	}, result)
}

func TestLexerClassifier(t *testing.T) {
	cls := &mockClassifier{}
	obj := &lexer{
		cls: cls,
	}

	result := obj.Classifier()

	assert.Same(t, cls, result)
}

func TestLexerLex(t *testing.T) {
	obj := &lexer{}
	cs := &mockScanner{}
	options := []Option{
		func(s State) {},
		func(s State) {},
	}
	state := &mockState{}
	defer patcher.SetVar(&newState, func(l Lexer, src scanner.Scanner, options []Option) State {
		assert.Same(t, obj, l)
		assert.Same(t, cs, src)
		assert.Len(t, options, 2)
		return state
	}).Install().Restore()

	result := obj.Lex(cs, options...)

	assert.Same(t, state, result)
}
