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

	"github.com/hydralang/ptk/common"
	"github.com/hydralang/ptk/lexer"
	"github.com/hydralang/ptk/parser"
)

func TestLiteral(t *testing.T) {
	s := &parser.MockState{}
	tok := &lexer.Token{}

	result, err := literal(s, 42, tok)

	assert.NoError(t, err)
	assert.Equal(t, &common.TokenNode{Token: *tok}, result)
}

func TestPrefixBase(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	exp := &common.MockNode{}
	s.On("Expression", 42).Return(exp, nil)
	node := &common.MockNode{}
	factoryCalled := false
	factory := func(fs parser.State, o *lexer.Token, e common.Node) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, op, o)
		assert.Same(t, exp, e)
		factoryCalled = true
		return node, nil
	}

	first := Prefix(factory, 42)
	result, err := first(s, 17, op)

	assert.NoError(t, err)
	assert.Same(t, node, result)
	assert.True(t, factoryCalled)
}

func TestPrefixExpressionFails(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	exp := &common.MockNode{}
	s.On("Expression", 42).Return(nil, assert.AnError)
	node := &common.MockNode{}
	factoryCalled := false
	factory := func(fs parser.State, o *lexer.Token, e common.Node) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, op, o)
		assert.Same(t, exp, e)
		factoryCalled = true
		return node, nil
	}

	first := Prefix(factory, 42)
	result, err := first(s, 17, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.False(t, factoryCalled)
}

func TestPrefixFactoryFails(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	exp := &common.MockNode{}
	s.On("Expression", 42).Return(exp, nil)
	factoryCalled := false
	factory := func(fs parser.State, o *lexer.Token, e common.Node) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, op, o)
		assert.Same(t, exp, e)
		factoryCalled = true
		return nil, assert.AnError
	}

	first := Prefix(factory, 42)
	result, err := first(s, 17, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.True(t, factoryCalled)
}

func TestInfixBase(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	left := &common.MockNode{}
	right := &common.MockNode{}
	s.On("Expression", 17).Return(right, nil)
	node := &common.MockNode{}
	factoryCalled := false
	factory := func(fs parser.State, l, r common.Node, o *lexer.Token) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		assert.Same(t, op, o)
		factoryCalled = true
		return node, nil
	}

	next := Infix(factory)
	result, err := next(s, 17, left, op)

	assert.NoError(t, err)
	assert.Same(t, node, result)
	assert.True(t, factoryCalled)
}

func TestInfixExpressionFails(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	left := &common.MockNode{}
	right := &common.MockNode{}
	s.On("Expression", 17).Return(nil, assert.AnError)
	node := &common.MockNode{}
	factoryCalled := false
	factory := func(fs parser.State, l, r common.Node, o *lexer.Token) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		assert.Same(t, op, o)
		factoryCalled = true
		return node, nil
	}

	next := Infix(factory)
	result, err := next(s, 17, left, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.False(t, factoryCalled)
}

func TestInfixFactoryFails(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	left := &common.MockNode{}
	right := &common.MockNode{}
	s.On("Expression", 17).Return(right, nil)
	factoryCalled := false
	factory := func(fs parser.State, l, r common.Node, o *lexer.Token) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, op, o)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		factoryCalled = true
		return nil, assert.AnError
	}

	next := Infix(factory)
	result, err := next(s, 17, left, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.True(t, factoryCalled)
}

func TestInfixRBase(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	left := &common.MockNode{}
	right := &common.MockNode{}
	s.On("Expression", 16).Return(right, nil)
	node := &common.MockNode{}
	factoryCalled := false
	factory := func(fs parser.State, l, r common.Node, o *lexer.Token) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		assert.Same(t, op, o)
		factoryCalled = true
		return node, nil
	}

	next := InfixR(factory)
	result, err := next(s, 17, left, op)

	assert.NoError(t, err)
	assert.Same(t, node, result)
	assert.True(t, factoryCalled)
}

func TestInfixRExpressionFails(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	left := &common.MockNode{}
	right := &common.MockNode{}
	s.On("Expression", 16).Return(nil, assert.AnError)
	node := &common.MockNode{}
	factoryCalled := false
	factory := func(fs parser.State, l, r common.Node, o *lexer.Token) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		assert.Same(t, op, o)
		factoryCalled = true
		return node, nil
	}

	next := InfixR(factory)
	result, err := next(s, 17, left, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.False(t, factoryCalled)
}

func TestInfixRFactoryFails(t *testing.T) {
	s := &parser.MockState{}
	op := &lexer.Token{}
	left := &common.MockNode{}
	right := &common.MockNode{}
	s.On("Expression", 16).Return(right, nil)
	factoryCalled := false
	factory := func(fs parser.State, l, r common.Node, o *lexer.Token) (common.Node, error) {
		assert.Same(t, s, fs)
		assert.Same(t, op, o)
		assert.Same(t, left, l)
		assert.Same(t, right, r)
		factoryCalled = true
		return nil, assert.AnError
	}

	next := InfixR(factory)
	result, err := next(s, 17, left, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.True(t, factoryCalled)
}
