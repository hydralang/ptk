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

package nodes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/lexer"
	"github.com/hydralang/ptk/parser"
)

type mockParser struct {
	mock.Mock
}

func (m *mockParser) Expression(rbp int) (parser.Node, error) {
	args := m.MethodCalled("Expression", rbp)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(parser.Node), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockParser) Statement() (parser.Node, error) {
	args := m.MethodCalled("Statement")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(parser.Node), args.Error(1)
	}

	return nil, args.Error(1)
}

type mockPushBackLexer struct {
	mock.Mock
}

func (m *mockPushBackLexer) Next() *lexer.Token {
	args := m.MethodCalled("Next")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(*lexer.Token)
	}

	return nil
}

func (m *mockPushBackLexer) PushBack(tok *lexer.Token) {
	m.MethodCalled("PushBack", tok)
}

type mockState struct {
	mock.Mock
}

func (m *mockState) Table() parser.Table {
	args := m.MethodCalled("Table")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(parser.Table)
	}

	return nil
}

func TestLiteral(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	tok := &lexer.Token{}

	result, err := literal(p, s, lex, 42, tok)

	assert.NoError(t, err)
	assert.Equal(t, &parser.TokenNode{Token: *tok}, result)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestPrefixBase(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	exp := &mockNode{}
	p.On("Expression", 42).Return(exp, nil)
	node := &mockNode{}
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, o *lexer.Token, e parser.Node) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, op, o)
		assert.Same(t, exp, e)
		factoryCalled = true
		return node, nil
	}

	first := Prefix(factory, 42)
	result, err := first(p, s, lex, 17, op)

	assert.NoError(t, err)
	assert.Same(t, node, result)
	assert.True(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestPrefixExpressionFails(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	exp := &mockNode{}
	p.On("Expression", 42).Return(nil, assert.AnError)
	node := &mockNode{}
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, o *lexer.Token, e parser.Node) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, op, o)
		assert.Same(t, exp, e)
		factoryCalled = true
		return node, nil
	}

	first := Prefix(factory, 42)
	result, err := first(p, s, lex, 17, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.False(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestPrefixFactoryFails(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	exp := &mockNode{}
	p.On("Expression", 42).Return(exp, nil)
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, o *lexer.Token, e parser.Node) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, op, o)
		assert.Same(t, exp, e)
		factoryCalled = true
		return nil, assert.AnError
	}

	first := Prefix(factory, 42)
	result, err := first(p, s, lex, 17, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.True(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestInfixBase(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	left := &mockNode{}
	right := &mockNode{}
	p.On("Expression", 17).Return(right, nil)
	node := &mockNode{}
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, l, r parser.Node, o *lexer.Token) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		assert.Same(t, op, o)
		factoryCalled = true
		return node, nil
	}

	next := Infix(factory)
	result, err := next(p, s, lex, 17, left, op)

	assert.NoError(t, err)
	assert.Same(t, node, result)
	assert.True(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestInfixExpressionFails(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	left := &mockNode{}
	right := &mockNode{}
	p.On("Expression", 17).Return(nil, assert.AnError)
	node := &mockNode{}
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, l, r parser.Node, o *lexer.Token) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		assert.Same(t, op, o)
		factoryCalled = true
		return node, nil
	}

	next := Infix(factory)
	result, err := next(p, s, lex, 17, left, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.False(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestInfixFactoryFails(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	left := &mockNode{}
	right := &mockNode{}
	p.On("Expression", 17).Return(right, nil)
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, l, r parser.Node, o *lexer.Token) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, op, o)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		factoryCalled = true
		return nil, assert.AnError
	}

	next := Infix(factory)
	result, err := next(p, s, lex, 17, left, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.True(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestInfixRBase(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	left := &mockNode{}
	right := &mockNode{}
	p.On("Expression", 16).Return(right, nil)
	node := &mockNode{}
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, l, r parser.Node, o *lexer.Token) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		assert.Same(t, op, o)
		factoryCalled = true
		return node, nil
	}

	next := InfixR(factory)
	result, err := next(p, s, lex, 17, left, op)

	assert.NoError(t, err)
	assert.Same(t, node, result)
	assert.True(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestInfixRExpressionFails(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	left := &mockNode{}
	right := &mockNode{}
	p.On("Expression", 16).Return(nil, assert.AnError)
	node := &mockNode{}
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, l, r parser.Node, o *lexer.Token) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		assert.Same(t, op, o)
		factoryCalled = true
		return node, nil
	}

	next := InfixR(factory)
	result, err := next(p, s, lex, 17, left, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.False(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}

func TestInfixRFactoryFails(t *testing.T) {
	p := &mockParser{}
	s := &mockState{}
	lex := &mockPushBackLexer{}
	op := &lexer.Token{}
	left := &mockNode{}
	right := &mockNode{}
	p.On("Expression", 16).Return(right, nil)
	factoryCalled := false
	factory := func(fp parser.IParser, fs parser.State, fLex parser.IPushBackLexer, l, r parser.Node, o *lexer.Token) (parser.Node, error) {
		assert.Same(t, p, fp)
		assert.Same(t, s, fs)
		assert.Same(t, lex, fLex)
		assert.Same(t, op, o)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		factoryCalled = true
		return nil, assert.AnError
	}

	next := InfixR(factory)
	result, err := next(p, s, lex, 17, left, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.True(t, factoryCalled)
	p.AssertExpectations(t)
	s.AssertExpectations(t)
	lex.AssertExpectations(t)
}
