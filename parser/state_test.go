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

	"github.com/hydralang/ptk/common"
	"github.com/hydralang/ptk/scanner"
)

func TestMockStateImplementsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &MockState{})
}

func TestMockStateParserNil(t *testing.T) {
	obj := &MockState{}
	obj.On("Parser").Return(nil)

	result := obj.Parser()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateParserNotNil(t *testing.T) {
	p := &MockParser{}
	obj := &MockState{}
	obj.On("Parser").Return(p)

	result := obj.Parser()

	assert.Same(t, p, result)
	obj.AssertExpectations(t)
}

func TestMockStateAppState(t *testing.T) {
	obj := &MockState{}
	obj.On("AppState").Return("state")

	result := obj.AppState()

	assert.Equal(t, "state", result)
	obj.AssertExpectations(t)
}

func TestMockStatePushAppState(t *testing.T) {
	obj := &MockState{}
	obj.On("PushAppState", "state")

	obj.PushAppState("state")

	obj.AssertExpectations(t)
}

func TestMockStatePopAppState(t *testing.T) {
	obj := &MockState{}
	obj.On("PopAppState").Return("state")

	result := obj.PopAppState()

	assert.Equal(t, "state", result)
	obj.AssertExpectations(t)
}

func TestMockStateSetAppState(t *testing.T) {
	obj := &MockState{}
	obj.On("SetAppState", "new").Return("old")

	result := obj.SetAppState("new")

	assert.Equal(t, "old", result)
	obj.AssertExpectations(t)
}

func TestMockStateTableNil(t *testing.T) {
	obj := &MockState{}
	obj.On("Table").Return(nil)

	result := obj.Table()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateTableNotNil(t *testing.T) {
	expected := Table{}
	obj := &MockState{}
	obj.On("Table").Return(expected)

	result := obj.Table()

	assert.Equal(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStatePushTable(t *testing.T) {
	tab := Table{
		"foo": {},
	}
	obj := &MockState{}
	obj.On("PushTable", tab)

	obj.PushTable(tab)

	obj.AssertExpectations(t)
}

func TestMockStatePopTableNil(t *testing.T) {
	obj := &MockState{}
	obj.On("PopTable").Return(nil)

	result := obj.PopTable()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStatePopTableNotNil(t *testing.T) {
	expected := Table{}
	obj := &MockState{}
	obj.On("PopTable").Return(expected)

	result := obj.PopTable()

	assert.Equal(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStateSetTableNil(t *testing.T) {
	tab := Table{
		"foo": Entry{},
	}
	obj := &MockState{}
	obj.On("SetTable", tab).Return(nil)

	result := obj.SetTable(tab)

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateSetTableNotNil(t *testing.T) {
	tab := Table{
		"foo": Entry{},
	}
	expected := Table{}
	obj := &MockState{}
	obj.On("SetTable", tab).Return(expected)

	result := obj.SetTable(tab)

	assert.Equal(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStateStreamNil(t *testing.T) {
	obj := &MockState{}
	obj.On("Stream").Return(nil)

	result := obj.Stream()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateStreamNotNil(t *testing.T) {
	expected := &common.MockTokenStream{}
	obj := &MockState{}
	obj.On("Stream").Return(expected)

	result := obj.Stream()

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStatePushStream(t *testing.T) {
	stream := &common.MockTokenStream{}
	obj := &MockState{}
	obj.On("PushStream", stream)

	obj.PushStream(stream)

	obj.AssertExpectations(t)
}

func TestMockStatePopStreamNil(t *testing.T) {
	obj := &MockState{}
	obj.On("PopStream").Return(nil)

	result := obj.PopStream()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStatePopStreamNotNil(t *testing.T) {
	expected := &common.MockTokenStream{}
	obj := &MockState{}
	obj.On("PopStream").Return(expected)

	result := obj.PopStream()

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStateSetStreamNil(t *testing.T) {
	stream := &common.MockTokenStream{}
	obj := &MockState{}
	obj.On("SetStream", stream).Return(nil)

	result := obj.SetStream(stream)

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateSetStreamNotNil(t *testing.T) {
	stream := &common.MockTokenStream{}
	expected := &common.MockTokenStream{}
	obj := &MockState{}
	obj.On("SetStream", stream).Return(expected)

	result := obj.SetStream(stream)

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStateTokenNil(t *testing.T) {
	obj := &MockState{}
	obj.On("Token").Return(nil)

	result := obj.Token()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateTokenNotNil(t *testing.T) {
	expected := &common.Token{}
	obj := &MockState{}
	obj.On("Token").Return(expected)

	result := obj.Token()

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStateNextTokenNil(t *testing.T) {
	obj := &MockState{}
	obj.On("NextToken").Return(nil)

	result := obj.NextToken()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateNextTokenNotNil(t *testing.T) {
	expected := &common.Token{}
	obj := &MockState{}
	obj.On("NextToken").Return(expected)

	result := obj.NextToken()

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStateMoreTokens(t *testing.T) {
	obj := &MockState{}
	obj.On("MoreTokens").Return(true)

	result := obj.MoreTokens()

	assert.True(t, result)
	obj.AssertExpectations(t)
}

func TestPushToken(t *testing.T) {
	tok := &common.Token{}
	obj := &MockState{}
	obj.On("PushToken", tok)

	obj.PushToken(tok)

	obj.AssertExpectations(t)
}

func TestMockStateExpressionNil(t *testing.T) {
	obj := &MockState{}
	obj.On("Expression", 42).Return(nil, assert.AnError)

	result, err := obj.Expression(42)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateExpressionNotNil(t *testing.T) {
	expected := &common.MockNode{}
	obj := &MockState{}
	obj.On("Expression", 42).Return(expected, assert.AnError)

	result, err := obj.Expression(42)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStateStatementNil(t *testing.T) {
	obj := &MockState{}
	obj.On("Statement").Return(nil, assert.AnError)

	result, err := obj.Statement()

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateStatementNotNil(t *testing.T) {
	expected := &common.MockNode{}
	obj := &MockState{}
	obj.On("Statement").Return(expected, assert.AnError)

	result, err := obj.Statement()

	assert.Same(t, assert.AnError, err)
	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestStateImplementsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &state{})
}

func TestNewState(t *testing.T) {
	tab := Table{
		"foo": Entry{},
	}
	parser := &MockParser{}
	parser.On("Table").Return(tab)
	stream := &common.MockTokenStream{}
	var opt1Called State
	var opt2Called State
	options := []Option{
		func(s State) {
			opt1Called = s
		},
		func(s State) {
			opt2Called = s
		},
	}

	result := NewState(parser, stream, options)

	assert.Same(t, result, opt1Called)
	assert.Same(t, result, opt2Called)
	state, ok := result.(*state)
	require.True(t, ok)
	assert.Same(t, parser, state.parser)
	assert.Equal(t, 0, state.appState.Len())
	assert.Equal(t, 1, state.table.Len())
	assert.Equal(t, tab, state.table.Get())
	assert.Equal(t, 1, state.stream.Len())
	assert.Same(t, stream, state.stream.Get())
	assert.Equal(t, 0, state.tokens.Len())
	assert.Nil(t, state.tok)
}

func TestStateParser(t *testing.T) {
	p := &MockParser{}
	obj := &state{
		parser: p,
	}

	result := obj.Parser()

	assert.Same(t, p, result)
}

func TestStateAppState(t *testing.T) {
	appStack := &common.MockStack{}
	appStack.On("Get").Return("state")
	obj := &state{
		appState: appStack,
	}

	result := obj.AppState()

	assert.Equal(t, "state", result)
	appStack.AssertExpectations(t)
}

func TestStatePushAppState(t *testing.T) {
	appStack := &common.MockStack{}
	appStack.On("Push", "state")
	obj := &state{
		appState: appStack,
	}

	obj.PushAppState("state")

	appStack.AssertExpectations(t)
}

func TestStatePopAppState(t *testing.T) {
	appStack := &common.MockStack{}
	appStack.On("Pop").Return("state")
	obj := &state{
		appState: appStack,
	}

	result := obj.PopAppState()

	assert.Equal(t, "state", result)
	appStack.AssertExpectations(t)
}

func TestStateSetAppState(t *testing.T) {
	appStack := &common.MockStack{}
	appStack.On("Set", "new").Return("old")
	obj := &state{
		appState: appStack,
	}

	result := obj.SetAppState("new")

	assert.Equal(t, "old", result)
	appStack.AssertExpectations(t)
}

func TestStateTableBase(t *testing.T) {
	tab := Table{
		"foo": Entry{},
	}
	tableStack := &common.MockStack{}
	tableStack.On("Get").Return(tab)
	obj := &state{
		table: tableStack,
	}

	result := obj.Table()

	assert.Equal(t, tab, result)
	tableStack.AssertExpectations(t)
}

func TestStateTableNil(t *testing.T) {
	tableStack := &common.MockStack{}
	tableStack.On("Get").Return(nil)
	obj := &state{
		table: tableStack,
	}

	result := obj.Table()

	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
}

func TestStatePushTable(t *testing.T) {
	tab := Table{
		"foo": Entry{},
	}
	tableStack := &common.MockStack{}
	tableStack.On("Push", tab)
	obj := &state{
		table: tableStack,
	}

	obj.PushTable(tab)

	tableStack.AssertExpectations(t)
}

func TestStatePopTableBase(t *testing.T) {
	tab := Table{
		"foo": Entry{},
	}
	tableStack := &common.MockStack{}
	tableStack.On("Pop").Return(tab)
	obj := &state{
		table: tableStack,
	}

	result := obj.PopTable()

	assert.Equal(t, tab, result)
	tableStack.AssertExpectations(t)
}

func TestStatePopTableNil(t *testing.T) {
	tableStack := &common.MockStack{}
	tableStack.On("Pop").Return(nil)
	obj := &state{
		table: tableStack,
	}

	result := obj.PopTable()

	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
}

func TestStateSetTableBase(t *testing.T) {
	tab := Table{
		"foo": Entry{},
	}
	newTab := Table{
		"bar": Entry{},
	}
	tableStack := &common.MockStack{}
	tableStack.On("Set", newTab).Return(tab)
	obj := &state{
		table: tableStack,
	}

	result := obj.SetTable(newTab)

	assert.Equal(t, tab, result)
	tableStack.AssertExpectations(t)
}

func TestStateSetTableNil(t *testing.T) {
	newTab := Table{
		"bar": Entry{},
	}
	tableStack := &common.MockStack{}
	tableStack.On("Set", newTab).Return(nil)
	obj := &state{
		table: tableStack,
	}

	result := obj.SetTable(newTab)

	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
}

func TestStateStreamBase(t *testing.T) {
	stream := &common.MockTokenStream{}
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(stream)
	obj := &state{
		stream: streamStack,
	}

	result := obj.Stream()

	assert.Same(t, stream, result)
	streamStack.AssertExpectations(t)
}

func TestStateStreamNil(t *testing.T) {
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(nil)
	obj := &state{
		stream: streamStack,
	}

	result := obj.Stream()

	assert.Nil(t, result)
	streamStack.AssertExpectations(t)
}

func TestStatePushStream(t *testing.T) {
	stream := &common.MockTokenStream{}
	streamStack := &common.MockStack{}
	streamStack.On("Push", stream)
	obj := &state{
		stream: streamStack,
	}

	obj.PushStream(stream)

	streamStack.AssertExpectations(t)
}

func TestStatePopStreamBase(t *testing.T) {
	stream := &common.MockTokenStream{}
	streamStack := &common.MockStack{}
	streamStack.On("Pop").Return(stream)
	obj := &state{
		stream: streamStack,
	}

	result := obj.PopStream()

	assert.Same(t, stream, result)
	streamStack.AssertExpectations(t)
}

func TestStatePopStreamNil(t *testing.T) {
	streamStack := &common.MockStack{}
	streamStack.On("Pop").Return(nil)
	obj := &state{
		stream: streamStack,
	}

	result := obj.PopStream()

	assert.Nil(t, result)
	streamStack.AssertExpectations(t)
}

func TestStateSetStreamBase(t *testing.T) {
	stream := &common.MockTokenStream{}
	newStream := &common.MockTokenStream{}
	streamStack := &common.MockStack{}
	streamStack.On("Set", newStream).Return(stream)
	obj := &state{
		stream: streamStack,
	}

	result := obj.SetStream(newStream)

	assert.Same(t, stream, result)
	streamStack.AssertExpectations(t)
}

func TestStateSetStreamNil(t *testing.T) {
	newStream := &common.MockTokenStream{}
	streamStack := &common.MockStack{}
	streamStack.On("Set", newStream).Return(nil)
	obj := &state{
		stream: streamStack,
	}

	result := obj.SetStream(newStream)

	assert.Nil(t, result)
	streamStack.AssertExpectations(t)
}

func TestStateToken(t *testing.T) {
	tok := &common.Token{}
	obj := &state{
		tok: tok,
	}

	result := obj.Token()

	assert.Same(t, tok, result)
}

func TestStateGetTokenBase(t *testing.T) {
	tok := &common.Token{}
	stream := &common.MockTokenStream{}
	stream.On("Next").Return(tok)
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(stream)
	obj := &state{
		stream: streamStack,
	}

	result := obj.getToken()

	assert.Same(t, tok, result)
	stream.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateGetTokenPopStream(t *testing.T) {
	tok := &common.Token{}
	stream1 := &common.MockTokenStream{}
	stream1.On("Next").Return(tok)
	stream2 := &common.MockTokenStream{}
	stream2.On("Next").Return(nil)
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(stream2).Once()
	streamStack.On("Pop").Return(stream2)
	streamStack.On("Get").Return(stream1)
	obj := &state{
		stream: streamStack,
	}

	result := obj.getToken()

	assert.Same(t, tok, result)
	stream1.AssertExpectations(t)
	stream2.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateGetTokenExhausted(t *testing.T) {
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(nil)
	obj := &state{
		stream: streamStack,
	}

	result := obj.getToken()

	assert.Nil(t, result)
	streamStack.AssertExpectations(t)
}

func TestStateNextTokenBase(t *testing.T) {
	tok := &common.Token{}
	stream := &common.MockTokenStream{}
	stream.On("Next").Return(tok)
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(stream)
	tokenStack := &common.MockStack{}
	tokenStack.On("Len").Return(0)
	obj := &state{
		stream: streamStack,
		tokens: tokenStack,
	}

	result := obj.NextToken()

	assert.Same(t, tok, result)
	assert.Same(t, tok, obj.tok)
	stream.AssertExpectations(t)
	streamStack.AssertExpectations(t)
	tokenStack.AssertExpectations(t)
}

func TestStateNextTokenPushed(t *testing.T) {
	tok := &common.Token{}
	stream := &common.MockTokenStream{}
	streamStack := &common.MockStack{}
	tokenStack := &common.MockStack{}
	tokenStack.On("Len").Return(1)
	tokenStack.On("Pop").Return(tok)
	obj := &state{
		stream: streamStack,
		tokens: tokenStack,
	}

	result := obj.NextToken()

	assert.Same(t, tok, result)
	assert.Same(t, tok, obj.tok)
	stream.AssertExpectations(t)
	streamStack.AssertExpectations(t)
	tokenStack.AssertExpectations(t)
}

func TestStateNextTokenExhausted(t *testing.T) {
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(nil)
	tokenStack := &common.MockStack{}
	tokenStack.On("Len").Return(0)
	obj := &state{
		stream: streamStack,
		tokens: tokenStack,
	}

	result := obj.NextToken()

	assert.Nil(t, result)
	assert.Nil(t, obj.tok)
	streamStack.AssertExpectations(t)
	tokenStack.AssertExpectations(t)
}

func TestStateMoreTokensBase(t *testing.T) {
	tok := &common.Token{}
	stream := &common.MockTokenStream{}
	stream.On("Next").Return(tok)
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(stream)
	tokenStack := &common.MockStack{}
	tokenStack.On("Len").Return(0)
	tokenStack.On("Push", tok)
	obj := &state{
		stream: streamStack,
		tokens: tokenStack,
	}

	result := obj.MoreTokens()

	assert.True(t, result)
	assert.Nil(t, obj.tok)
	stream.AssertExpectations(t)
	streamStack.AssertExpectations(t)
	tokenStack.AssertExpectations(t)
}

func TestStateMoreTokensPushed(t *testing.T) {
	stream := &common.MockTokenStream{}
	streamStack := &common.MockStack{}
	tokenStack := &common.MockStack{}
	tokenStack.On("Len").Return(1)
	obj := &state{
		stream: streamStack,
		tokens: tokenStack,
	}

	result := obj.MoreTokens()

	assert.True(t, result)
	assert.Nil(t, obj.tok)
	stream.AssertExpectations(t)
	streamStack.AssertExpectations(t)
	tokenStack.AssertExpectations(t)
}

func TestStateMoreTokensExhausted(t *testing.T) {
	streamStack := &common.MockStack{}
	streamStack.On("Get").Return(nil)
	tokenStack := &common.MockStack{}
	tokenStack.On("Len").Return(0)
	obj := &state{
		stream: streamStack,
		tokens: tokenStack,
	}

	result := obj.MoreTokens()

	assert.False(t, result)
	assert.Nil(t, obj.tok)
	streamStack.AssertExpectations(t)
	tokenStack.AssertExpectations(t)
}

func TestStatePushToken(t *testing.T) {
	tok := &common.Token{}
	tokenStream := &common.MockStack{}
	tokenStream.On("Push", tok)
	obj := &state{
		tokens: tokenStream,
	}

	obj.PushToken(tok)

	tokenStream.AssertExpectations(t)
}

func TestStateGetEntryBase(t *testing.T) {
	tok := &common.Token{Type: "type"}
	tab := Table{
		"type": Entry{Power: 42},
	}
	tableStack := &common.MockStack{}
	tableStack.On("Get").Return(tab)
	obj := &state{
		table: tableStack,
	}

	ent, err := obj.getEntry(tok)

	assert.NoError(t, err)
	assert.Equal(t, Entry{Power: 42}, ent)
	tableStack.AssertExpectations(t)
}

func TestStateGetEntryNoTable(t *testing.T) {
	tok := &common.Token{Type: "type"}
	tableStack := &common.MockStack{}
	tableStack.On("Get").Return(nil)
	obj := &state{
		table: tableStack,
	}

	ent, err := obj.getEntry(tok)

	assert.Same(t, ErrNoTable, err)
	assert.Equal(t, Entry{}, ent)
	tableStack.AssertExpectations(t)
}

func TestStateGetEntryNoEntry(t *testing.T) {
	tok := &common.Token{Type: "type"}
	tab := Table{}
	tableStack := &common.MockStack{}
	tableStack.On("Get").Return(tab)
	obj := &state{
		table: tableStack,
	}

	ent, err := obj.getEntry(tok)

	assert.True(t, errors.Is(err, ErrUnknownTokenType))
	assert.Equal(t, Entry{}, ent)
	tableStack.AssertExpectations(t)
}

type binaryNode struct {
	L  common.Node
	R  common.Node
	Op *common.Token
}

func (bn *binaryNode) Location() scanner.Location {
	return nil
}

func (bn *binaryNode) Children() []common.Node {
	return []common.Node{bn.L, bn.R}
}

func (bn *binaryNode) String() string {
	return bn.Op.String()
}

func TestStateExpressionBase(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}
	var first ExprFirst = func(s State, pow int, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		assert.Equal(t, 0, pow)
		return tok, nil
	}
	var next ExprNext = func(s State, pow int, l common.Node, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		assert.Equal(t, 10, pow)
		r, _ := s.Expression(pow)
		return &binaryNode{
			L:  l,
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
	tableStack.On("Get").Return(tab)
	streamStack.On("Get").Return(nil)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Expression(0)

	assert.NoError(t, err)
	assert.Equal(t, &binaryNode{
		Op: &common.Token{Type: "+"},
		L: &binaryNode{
			Op: &common.Token{Type: "+"},
			L:  &common.Token{Type: "n", Value: 1},
			R:  &common.Token{Type: "n", Value: 2},
		},
		R: &common.Token{Type: "n", Value: 3},
	}, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateExpressionPrecedence(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "*"},
		{Type: "n", Value: 3},
		{Type: "+"},
		{Type: "n", Value: 4},
	}
	var first ExprFirst = func(s State, pow int, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		assert.Equal(t, 0, pow)
		return tok, nil
	}
	var next ExprNext = func(s State, pow int, l common.Node, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		switch tok.Type {
		case "+":
			assert.Equal(t, 10, pow)
		case "*":
			assert.Equal(t, 20, pow)
		}
		r, _ := s.Expression(pow)
		return &binaryNode{
			L:  l,
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
	tableStack.On("Get").Return(tab)
	streamStack.On("Get").Return(nil)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Expression(0)

	assert.NoError(t, err)
	assert.Equal(t, &binaryNode{
		Op: &common.Token{Type: "+"},
		L: &binaryNode{
			Op: &common.Token{Type: "+"},
			L:  &common.Token{Type: "n", Value: 1},
			R: &binaryNode{
				Op: &common.Token{Type: "*"},
				L:  &common.Token{Type: "n", Value: 2},
				R:  &common.Token{Type: "n", Value: 3},
			},
		},
		R: &common.Token{Type: "n", Value: 4},
	}, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateExpressionNoTokens(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	streamStack.On("Get").Return(nil)

	result, err := obj.Expression(0)

	assert.True(t, errors.Is(err, ErrExpectedToken))
	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateExpressionFirstEntryMissing(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}
	tab := Table{}
	tableStack.On("Get").Return(tab)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Expression(0)

	assert.True(t, errors.Is(err, ErrUnknownTokenType))
	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateExpressionFirstFails(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}
	var first ExprFirst = func(s State, pow int, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		assert.Equal(t, 0, pow)
		return nil, assert.AnError
	}
	var next ExprNext = func(s State, pow int, l common.Node, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		assert.Equal(t, 10, pow)
		r, _ := s.Expression(pow)
		return &binaryNode{
			L:  l,
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
	tableStack.On("Get").Return(tab)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Expression(0)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateExpressionNextEntryMissing(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}
	var first ExprFirst = func(s State, pow int, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		assert.Equal(t, 0, pow)
		return tok, nil
	}
	tab := Table{
		"n": Entry{
			Power: 0,
			First: first,
		},
	}
	tableStack.On("Get").Return(tab)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Expression(0)

	assert.True(t, errors.Is(err, ErrUnknownTokenType))
	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateExpressionNextFails(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "n", Value: 1},
		{Type: "+"},
		{Type: "n", Value: 2},
		{Type: "+"},
		{Type: "n", Value: 3},
	}
	var first ExprFirst = func(s State, pow int, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		assert.Equal(t, 0, pow)
		return tok, nil
	}
	var next ExprNext = func(s State, pow int, l common.Node, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
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
	tableStack.On("Get").Return(tab)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Expression(0)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

type stmtNode struct {
	toks []*common.Token
}

func (sn *stmtNode) Location() scanner.Location {
	return nil
}

func (sn *stmtNode) Children() []common.Node {
	children := []common.Node{}
	for _, tok := range sn.toks[1:] {
		children = append(children, tok)
	}
	return children
}

func (sn *stmtNode) String() string {
	return sn.toks[0].String()
}

func TestStateStatementBase(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "stmt"},
		{Type: "kw", Value: "kw1"},
		{Type: "kw", Value: "kw2"},
		{Type: "kw", Value: "kw3"},
		{Type: "end"},
	}
	var stmt Statement = func(s State, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		node := &stmtNode{
			toks: []*common.Token{tok},
		}
		for nextTok := s.NextToken(); nextTok != nil; nextTok = s.NextToken() {
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
	tableStack.On("Get").Return(tab)
	streamStack.On("Get").Return(nil)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Statement()

	assert.NoError(t, err)
	assert.Equal(t, &stmtNode{
		toks: tokens,
	}, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateStatementNoTokens(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	streamStack.On("Get").Return(nil)

	result, err := obj.Statement()

	assert.NoError(t, err)
	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateStatementNoEntry(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "stmt"},
		{Type: "kw", Value: "kw1"},
		{Type: "kw", Value: "kw2"},
		{Type: "kw", Value: "kw3"},
		{Type: "end"},
	}
	tab := Table{}
	tableStack.On("Get").Return(tab)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Statement()

	assert.True(t, errors.Is(err, ErrUnknownTokenType))
	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}

func TestStateStatementStmtFailed(t *testing.T) {
	tableStack := &common.MockStack{}
	streamStack := &common.MockStack{}
	tokensStack := common.NewStack()
	obj := &state{
		table:  tableStack,
		stream: streamStack,
		tokens: tokensStack,
	}
	tokens := []*common.Token{
		{Type: "stmt"},
		{Type: "kw", Value: "kw1"},
		{Type: "kw", Value: "kw2"},
		{Type: "kw", Value: "kw3"},
		{Type: "end"},
	}
	var stmt Statement = func(s State, tok *common.Token) (common.Node, error) {
		assert.Same(t, obj, s)
		return nil, assert.AnError
	}
	tab := Table{
		"stmt": Entry{
			Stmt: stmt,
		},
	}
	tableStack.On("Get").Return(tab)
	for i := range tokens {
		tok := tokens[len(tokens)-i-1]
		tokensStack.Push(tok)
	}

	result, err := obj.Statement()

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	tableStack.AssertExpectations(t)
	streamStack.AssertExpectations(t)
}
