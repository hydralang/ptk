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
	"testing"

	"github.com/klmitch/patcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/lexer"
)

func TestAppState(t *testing.T) {
	s := &MockState{}
	s.On("PushAppState", "state")

	opt := AppState("state")
	opt(s)

	s.AssertExpectations(t)
}

func TestMockParserImplementsParser(t *testing.T) {
	assert.Implements(t, (*Parser)(nil), &MockParser{})
}

func TestMockParserTableNil(t *testing.T) {
	obj := &MockParser{}
	obj.On("Table").Return(nil)

	result := obj.Table()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockParserTableNotNil(t *testing.T) {
	tab := Table{
		"ent": Entry{},
	}
	obj := &MockParser{}
	obj.On("Table").Return(tab)

	result := obj.Table()

	assert.Equal(t, tab, result)
	obj.AssertExpectations(t)
}

func TestMockParserExpressionNil(t *testing.T) {
	stream := &mockLexer{}
	obj := &MockParser{}
	obj.On("Expression", stream, mock.Anything).Return(nil, assert.AnError)

	result, err := obj.Expression(stream)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockParserExpressionNotNil(t *testing.T) {
	expected := &mockNode{}
	stream := &mockLexer{}
	obj := &MockParser{}
	obj.On("Expression", stream, mock.Anything).Return(expected, assert.AnError)

	result, err := obj.Expression(stream)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockParserStatementNil(t *testing.T) {
	stream := &mockLexer{}
	obj := &MockParser{}
	obj.On("Statement", stream, mock.Anything).Return(nil, assert.AnError)

	result, err := obj.Statement(stream)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockParserStatementNotNil(t *testing.T) {
	expected := &mockNode{}
	stream := &mockLexer{}
	obj := &MockParser{}
	obj.On("Statement", stream, mock.Anything).Return(expected, assert.AnError)

	result, err := obj.Statement(stream)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockParserStatementsNil(t *testing.T) {
	stream := &mockLexer{}
	obj := &MockParser{}
	obj.On("Statements", stream, mock.Anything).Return(nil, assert.AnError)

	result, err := obj.Statements(stream)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockParserStatementsNotNil(t *testing.T) {
	stream := &mockLexer{}
	obj := &MockParser{}
	obj.On("Statements", stream, mock.Anything).Return([]Node{}, assert.AnError)

	result, err := obj.Statements(stream)

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, []Node{}, result)
	obj.AssertExpectations(t)
}

func TestParserImplementsParser(t *testing.T) {
	assert.Implements(t, (*Parser)(nil), &parser{})
}

func TestNew(t *testing.T) {
	tab := Table{
		"ent": Entry{},
	}

	result := New(tab)

	assert.Equal(t, &parser{
		table: tab,
	}, result)
}

func TestParserTable(t *testing.T) {
	obj := &parser{
		table: Table{
			"ent": Entry{},
		},
	}

	result := obj.Table()

	assert.Equal(t, Table{
		"ent": Entry{},
	}, result)
}

func TestParserExpression(t *testing.T) {
	obj := &parser{
		table: Table{
			"ent": Entry{},
		},
	}
	stream := &mockLexer{}
	options := []Option{
		func(s State) {},
		func(s State) {},
	}
	state := &MockState{}
	node := &mockNode{}
	state.On("Expression", 0).Return(node, assert.AnError)
	defer patcher.SetVar(&newState, func(p Parser, str lexer.ILexer, options []Option) State {
		assert.Same(t, obj, p)
		assert.Same(t, stream, str)
		assert.Len(t, options, 2)
		return state
	}).Install().Restore()

	result, err := obj.Expression(stream, options...)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, node, result)
	state.AssertExpectations(t)
}

func TestParserStatement(t *testing.T) {
	obj := &parser{
		table: Table{
			"ent": Entry{},
		},
	}
	stream := &mockLexer{}
	options := []Option{
		func(s State) {},
		func(s State) {},
	}
	state := &MockState{}
	node := &mockNode{}
	state.On("Statement").Return(node, assert.AnError)
	defer patcher.SetVar(&newState, func(p Parser, str lexer.ILexer, options []Option) State {
		assert.Same(t, obj, p)
		assert.Same(t, stream, str)
		assert.Len(t, options, 2)
		return state
	}).Install().Restore()

	result, err := obj.Statement(stream, options...)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, node, result)
	state.AssertExpectations(t)
}

func TestParserStatementsBase(t *testing.T) {
	obj := &parser{
		table: Table{
			"ent": Entry{},
		},
	}
	stream := &mockLexer{}
	options := []Option{
		func(s State) {},
		func(s State) {},
	}
	state := &MockState{}
	nodes := []Node{&mockNode{}, &mockNode{}, &mockNode{}}
	for i, node := range nodes {
		node.(*mockNode).On("dummy", i) // make distinct
		state.On("Statement").Return(node, nil).Once()
	}
	state.On("Statement").Return(nil, nil)
	defer patcher.SetVar(&newState, func(p Parser, str lexer.ILexer, options []Option) State {
		assert.Same(t, obj, p)
		assert.Same(t, stream, str)
		assert.Len(t, options, 2)
		return state
	}).Install().Restore()

	result, err := obj.Statements(stream, options...)

	assert.NoError(t, err)
	assert.Equal(t, nodes, result)
	state.AssertExpectations(t)
}

func TestParserStatementsError(t *testing.T) {
	obj := &parser{
		table: Table{
			"ent": Entry{},
		},
	}
	stream := &mockLexer{}
	options := []Option{
		func(s State) {},
		func(s State) {},
	}
	state := &MockState{}
	state.On("Statement").Return(nil, assert.AnError)
	defer patcher.SetVar(&newState, func(p Parser, str lexer.ILexer, options []Option) State {
		assert.Same(t, obj, p)
		assert.Same(t, stream, str)
		assert.Len(t, options, 2)
		return state
	}).Install().Restore()

	result, err := obj.Statements(stream, options...)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	state.AssertExpectations(t)
}
