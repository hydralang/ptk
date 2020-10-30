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

package parser

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hydralang/ptk/lexer"
)

type mockLexer struct {
	mock.Mock
}

func (m *mockLexer) Next() *lexer.Token {
	args := m.MethodCalled("Next")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(*lexer.Token)
	}

	return nil
}

type mockPushBackLexer struct {
	mockLexer
}

func (m *mockPushBackLexer) PushBack(tok *lexer.Token) {
	m.MethodCalled("PushBack", tok)
}

func TestPushBackLexerImplementIPushBackLexer(t *testing.T) {
	assert.Implements(t, (*IPushBackLexer)(nil), &PushBackLexer{})
}

func TestNewPushBackLexer(t *testing.T) {
	l := &mockLexer{}

	result := NewPushBackLexer(l)

	assert.Equal(t, &PushBackLexer{
		Lexer: l,
		toks:  &list.List{},
	}, result)
}

func TestPushBackLexerNextBase(t *testing.T) {
	tok := &lexer.Token{}
	l := &mockLexer{}
	l.On("Next").Return(tok)
	obj := &PushBackLexer{
		Lexer: l,
		toks:  &list.List{},
	}

	result := obj.Next()

	assert.Same(t, tok, result)
	assert.Same(t, l, obj.Lexer)
	assert.Equal(t, 0, obj.toks.Len())
	l.AssertExpectations(t)
}

func TestPushBackLexerNextExhausted(t *testing.T) {
	l := &mockLexer{}
	l.On("Next").Return(nil)
	obj := &PushBackLexer{
		Lexer: l,
		toks:  &list.List{},
	}

	result := obj.Next()

	assert.Nil(t, result)
	assert.Nil(t, obj.Lexer)
	assert.Equal(t, 0, obj.toks.Len())
	l.AssertExpectations(t)
}

func TestPushBackLexerNextPushed(t *testing.T) {
	tok := &lexer.Token{}
	l := &mockLexer{}
	obj := &PushBackLexer{
		Lexer: l,
		toks:  &list.List{},
	}
	obj.toks.PushFront(tok)

	result := obj.Next()

	assert.Same(t, tok, result)
	assert.Same(t, l, obj.Lexer)
	assert.Equal(t, 0, obj.toks.Len())
	l.AssertExpectations(t)
}

func TestPushBackLexerPushBack(t *testing.T) {
	tok := &lexer.Token{}
	obj := &PushBackLexer{
		toks: &list.List{},
	}

	obj.PushBack(tok)

	require.Equal(t, 1, obj.toks.Len())
	assert.Same(t, tok, obj.toks.Front().Value)
}
