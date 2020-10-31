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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hydralang/ptk/lexer"
	"github.com/hydralang/ptk/scanner"
)

func TestNewBase(t *testing.T) {
	l := &mockLexer{}
	state := &mockState{}

	result := New(l, state)

	require.NotNil(t, result.Lexer)
	assert.Same(t, l, result.Lexer.(*PushBackLexer).Lexer)
	assert.Same(t, state, result.State)
}

func TestNewWithPushBackLexer(t *testing.T) {
	l := &mockPushBackLexer{}
	state := &mockState{}

	result := New(l, state)

	assert.Same(t, l, result.Lexer)
	assert.Same(t, state, result.State)
}

type binaryNode struct {
	L  Node
	R  Node
	Op *lexer.Token
}

func (bn *binaryNode) Location() scanner.Location {
	return nil
}

func (bn *binaryNode) Children() []Node {
	return []Node{bn.L, bn.R}
}

func (bn *binaryNode) String() string {
	return bn.Op.String()
}

func TestParserExpressionBase(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	var first ExprFirst = func(p *Parser, pow int, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		assert.Equal(t, 0, pow)
		return &TokenNode{Token: tok}, nil
	}
	var next ExprNext = func(p *Parser, pow int, ln Node, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		assert.Equal(t, 10, pow)
		r, _ := p.Expression(pow)
		return &binaryNode{
			L:  ln,
			R:  r,
			Op: tok,
		}, nil
	}
	tab := Table{
		"n": Entry{
			Power: 0,
			First: first,
		},
		"+": Entry{
			Power: 10,
			Next:  next,
		},
	}
	state.On("Table").Return(tab)

	result, err := obj.Expression(0)

	assert.NoError(t, err)
	assert.Equal(t, &binaryNode{
		Op: &lexer.Token{Type: "+"},
		L: &binaryNode{
			Op: &lexer.Token{Type: "+"},
			L:  &TokenNode{Token: &lexer.Token{Type: "n", Value: 1}},
			R:  &TokenNode{Token: &lexer.Token{Type: "n", Value: 2}},
		},
		R: &TokenNode{Token: &lexer.Token{Type: "n", Value: 3}},
	}, result)
	state.AssertExpectations(t)
}

func TestParserExpressionPrecedence(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "*"},
		{Type: "n", Value: 3},
		{Type: "+"},
		{Type: "n", Value: 4},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	var first ExprFirst = func(p *Parser, pow int, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		assert.Equal(t, 0, pow)
		return &TokenNode{Token: tok}, nil
	}
	var next ExprNext = func(p *Parser, pow int, ln Node, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		switch tok.Type {
		case "+":
			assert.Equal(t, 10, pow)
		case "*":
			assert.Equal(t, 20, pow)
		}
		r, _ := p.Expression(pow)
		return &binaryNode{
			L:  ln,
			R:  r,
			Op: tok,
		}, nil
	}
	tab := Table{
		"n": Entry{
			Power: 0,
			First: first,
		},
		"+": Entry{
			Power: 10,
			Next:  next,
		},
		"*": Entry{
			Power: 20,
			Next:  next,
		},
	}
	state.On("Table").Return(tab)

	result, err := obj.Expression(0)

	assert.NoError(t, err)
	assert.Equal(t, &binaryNode{
		Op: &lexer.Token{Type: "+"},
		L: &binaryNode{
			Op: &lexer.Token{Type: "+"},
			L:  &TokenNode{Token: &lexer.Token{Type: "n", Value: 1}},
			R: &binaryNode{
				Op: &lexer.Token{Type: "*"},
				L:  &TokenNode{Token: &lexer.Token{Type: "n", Value: 2}},
				R:  &TokenNode{Token: &lexer.Token{Type: "n", Value: 3}},
			},
		},
		R: &TokenNode{Token: &lexer.Token{Type: "n", Value: 4}},
	}, result)
	state.AssertExpectations(t)
}

func TestParserExpressionNoTokens(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}

	result, err := obj.Expression(0)

	assert.True(t, errors.Is(err, ErrExpectedToken))
	assert.Nil(t, result)
	state.AssertExpectations(t)
}

func TestParserExpressionFirstEntryMissing(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	tab := Table{}
	state.On("Table").Return(tab)

	result, err := obj.Expression(0)

	assert.True(t, errors.Is(err, ErrUnknownTokenType))
	assert.Nil(t, result)
	state.AssertExpectations(t)
}

func TestParserExpressionFirstFails(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	var first ExprFirst = func(p *Parser, pow int, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		assert.Equal(t, 0, pow)
		return nil, assert.AnError
	}
	var next ExprNext = func(p *Parser, pow int, ln Node, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		assert.Equal(t, 10, pow)
		r, _ := p.Expression(pow)
		return &binaryNode{
			L:  ln,
			R:  r,
			Op: tok,
		}, nil
	}
	tab := Table{
		"n": Entry{
			Power: 0,
			First: first,
		},
		"+": Entry{
			Power: 10,
			Next:  next,
		},
	}
	state.On("Table").Return(tab)

	result, err := obj.Expression(0)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	state.AssertExpectations(t)
}

func TestParserExpressionNextEntryMissing(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	var first ExprFirst = func(p *Parser, pow int, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		assert.Equal(t, 0, pow)
		return &TokenNode{Token: tok}, nil
	}
	tab := Table{
		"n": Entry{
			Power: 0,
			First: first,
		},
	}
	state.On("Table").Return(tab)

	result, err := obj.Expression(0)

	assert.True(t, errors.Is(err, ErrUnknownTokenType))
	assert.Nil(t, result)
	state.AssertExpectations(t)
}

func TestParserExpressionNextFails(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	var first ExprFirst = func(p *Parser, pow int, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		assert.Equal(t, 0, pow)
		return &TokenNode{Token: tok}, nil
	}
	var next ExprNext = func(p *Parser, pow int, ln Node, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		assert.Equal(t, 10, pow)
		return nil, assert.AnError
	}
	tab := Table{
		"n": Entry{
			Power: 0,
			First: first,
		},
		"+": Entry{
			Power: 10,
			Next:  next,
		},
	}
	state.On("Table").Return(tab)

	result, err := obj.Expression(0)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	state.AssertExpectations(t)
}

type stmtNode struct {
	toks []*lexer.Token
}

func (sn *stmtNode) Location() scanner.Location {
	return nil
}

func (sn *stmtNode) Children() []Node {
	children := []Node{}
	for _, tok := range sn.toks[1:] {
		children = append(children, &TokenNode{Token: tok})
	}
	return children
}

func (sn *stmtNode) String() string {
	return sn.toks[0].String()
}

func TestStateStatementBase(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "stmt"},
		{Type: "kw", Value: "kw1"},
		{Type: "kw", Value: "kw2"},
		{Type: "kw", Value: "kw3"},
		{Type: "end"},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	var stmt Statement = func(p *Parser, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		node := &stmtNode{
			toks: []*lexer.Token{tok},
		}
		for nextTok := p.Lexer.Next(); nextTok != nil; nextTok = p.Lexer.Next() {
			node.toks = append(node.toks, nextTok)
			if tok.Type == "end" {
				break
			}
		}
		return node, nil
	}
	tab := Table{
		"stmt": Entry{
			Stmt: stmt,
		},
	}
	state.On("Table").Return(tab)

	result, err := obj.Statement()

	assert.NoError(t, err)
	assert.Equal(t, &stmtNode{
		toks: []*lexer.Token{
			{Type: "stmt"},
			{Type: "kw", Value: "kw1"},
			{Type: "kw", Value: "kw2"},
			{Type: "kw", Value: "kw3"},
			{Type: "end"},
		},
	}, result)
	state.AssertExpectations(t)
}

func TestStateStatementNoTokens(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}

	result, err := obj.Statement()

	assert.NoError(t, err)
	assert.Nil(t, result)
	state.AssertExpectations(t)
}

func TestStateStatementNoEntry(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "stmt"},
		{Type: "kw", Value: "kw1"},
		{Type: "kw", Value: "kw2"},
		{Type: "kw", Value: "kw3"},
		{Type: "end"},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	tab := Table{}
	state.On("Table").Return(tab)

	result, err := obj.Statement()

	assert.True(t, errors.Is(err, ErrUnknownTokenType))
	assert.Nil(t, result)
	state.AssertExpectations(t)
}

func TestStateStatementStmtFailed(t *testing.T) {
	l := NewPushBackLexer(lexer.NewListLexer([]*lexer.Token{
		{Type: "stmt"},
		{Type: "kw", Value: "kw1"},
		{Type: "kw", Value: "kw2"},
		{Type: "kw", Value: "kw3"},
		{Type: "end"},
	}))
	state := &mockState{}
	obj := &Parser{
		Lexer: l,
		State: state,
	}
	var stmt Statement = func(p *Parser, tok *lexer.Token) (Node, error) {
		assert.Same(t, obj, p)
		return nil, assert.AnError
	}
	tab := Table{
		"stmt": Entry{
			Stmt: stmt,
		},
	}
	state.On("Table").Return(tab)

	result, err := obj.Statement()

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	state.AssertExpectations(t)
}
