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

	"github.com/hydralang/ptk/common"
)

func TestAppState(t *testing.T) {
	s := &MockState{}
	s.On("PushAppState", "state")

	opt := AppState("state")
	opt(s)

	s.AssertExpectations(t)
}

func TestMockLexerImplementsLexer(t *testing.T) {
	assert.Implements(t, (*Lexer)(nil), &MockLexer{})
}

func TestMockLexerClassifierNil(t *testing.T) {
	obj := &MockLexer{}
	obj.On("Classifier").Return(nil)

	result := obj.Classifier()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockLexerClassifierNotNil(t *testing.T) {
	cls := &MockClassifier{}
	obj := &MockLexer{}
	obj.On("Classifier").Return(cls)

	result := obj.Classifier()

	assert.Equal(t, cls, result)
	obj.AssertExpectations(t)
}

func TestMockLexerLexNil(t *testing.T) {
	cs := &common.MockCharStream{}
	obj := &MockLexer{}
	obj.On("Lex", cs, mock.Anything).Return(nil)

	result := obj.Lex(cs)

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockLexerLexNotNil(t *testing.T) {
	stream := &common.MockTokenStream{}
	cs := &common.MockCharStream{}
	obj := &MockLexer{}
	obj.On("Lex", cs, mock.Anything).Return(stream)

	result := obj.Lex(cs)

	assert.Same(t, stream, result)
	obj.AssertExpectations(t)
}

func TestLexerImplementsLexer(t *testing.T) {
	assert.Implements(t, (*Lexer)(nil), &lexer{})
}

func TestNew(t *testing.T) {
	cls := &MockClassifier{}

	result := New(cls)

	assert.Equal(t, &lexer{
		cls: cls,
	}, result)
}

func TestLexerClassifier(t *testing.T) {
	cls := &MockClassifier{}
	obj := &lexer{
		cls: cls,
	}

	result := obj.Classifier()

	assert.Same(t, cls, result)
}

func TestLexerLex(t *testing.T) {
	obj := &lexer{}
	cs := &common.MockCharStream{}
	options := []Option{
		func(s State) {},
		func(s State) {},
	}
	state := &MockState{}
	defer patcher.SetVar(&newState, func(l Lexer, src common.CharStream, options []Option) State {
		assert.Same(t, obj, l)
		assert.Same(t, cs, src)
		assert.Len(t, options, 2)
		return state
	}).Install().Restore()

	result := obj.Lex(cs, options...)

	assert.Same(t, state, result)
}
