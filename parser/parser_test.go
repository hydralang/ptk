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

	"github.com/hydralang/ptk/common"
)

func TestMockParserImplementsParser(t *testing.T) {
	assert.Implements(t, (*Parser)(nil), &MockParser{})
}

func TestMockParserExpressionNil(t *testing.T) {
	stream := &common.MockTokenStream{}
	obj := &MockParser{}
	obj.On("Expression", stream).Return(nil, assert.AnError)

	result, err := obj.Expression(stream)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockParserExpressionNotNil(t *testing.T) {
	expected := &common.MockNode{}
	stream := &common.MockTokenStream{}
	obj := &MockParser{}
	obj.On("Expression", stream).Return(expected, assert.AnError)

	result, err := obj.Expression(stream)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockParserStatementNil(t *testing.T) {
	stream := &common.MockTokenStream{}
	obj := &MockParser{}
	obj.On("Statement", stream).Return(nil, assert.AnError)

	result, err := obj.Statement(stream)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockParserStatementNotNil(t *testing.T) {
	expected := &common.MockNode{}
	stream := &common.MockTokenStream{}
	obj := &MockParser{}
	obj.On("Statement", stream).Return(expected, assert.AnError)

	result, err := obj.Statement(stream)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockParserStatementsNil(t *testing.T) {
	stream := &common.MockTokenStream{}
	obj := &MockParser{}
	obj.On("Statements", stream).Return(nil, assert.AnError)

	result, err := obj.Statements(stream)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockParserStatementsNotNil(t *testing.T) {
	stream := &common.MockTokenStream{}
	obj := &MockParser{}
	obj.On("Statements", stream).Return([]common.Node{}, assert.AnError)

	result, err := obj.Statements(stream)

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, []common.Node{}, result)
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

func TestParserExpression(t *testing.T) {
	obj := &parser{
		table: Table{
			"ent": Entry{},
		},
	}
	stream := &common.MockTokenStream{}
	state := &MockState{}
	node := &common.MockNode{}
	state.On("Expression", obj, 0).Return(node, assert.AnError)
	defer patcher.SetVar(&newState, func(tab Table, str common.TokenStream) State {
		assert.Equal(t, obj.table, tab)
		assert.Same(t, stream, str)
		return state
	}).Install().Restore()

	result, err := obj.Expression(stream)

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
	stream := &common.MockTokenStream{}
	state := &MockState{}
	node := &common.MockNode{}
	state.On("Statement", obj).Return(node, assert.AnError)
	defer patcher.SetVar(&newState, func(tab Table, str common.TokenStream) State {
		assert.Equal(t, obj.table, tab)
		assert.Same(t, stream, str)
		return state
	}).Install().Restore()

	result, err := obj.Statement(stream)

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
	stream := &common.MockTokenStream{}
	state := &MockState{}
	nodes := []common.Node{&common.MockNode{}, &common.MockNode{}, &common.MockNode{}}
	for i, node := range nodes {
		node.(*common.MockNode).On("dummy", i) // make distinct
		state.On("Statement", obj).Return(node, nil).Once()
	}
	state.On("Statement", obj).Return(nil, nil)
	defer patcher.SetVar(&newState, func(tab Table, str common.TokenStream) State {
		assert.Equal(t, obj.table, tab)
		assert.Same(t, stream, str)
		return state
	}).Install().Restore()

	result, err := obj.Statements(stream)

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
	stream := &common.MockTokenStream{}
	state := &MockState{}
	state.On("Statement", obj).Return(nil, assert.AnError)
	defer patcher.SetVar(&newState, func(tab Table, str common.TokenStream) State {
		assert.Equal(t, obj.table, tab)
		assert.Same(t, stream, str)
		return state
	}).Install().Restore()

	result, err := obj.Statements(stream)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	state.AssertExpectations(t)
}
