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

	"github.com/hydralang/ptk/common"
	"github.com/hydralang/ptk/lexer"
)

func TestEntryCallFirstNil(t *testing.T) {
	s := &MockState{}
	tok := &lexer.Token{}
	ent := Entry{}

	result, err := ent.callFirst(s, tok)

	assert.True(t, errors.Is(err, ErrUnexpectedToken))
	assert.Nil(t, result)
}

func TestEntryCallFirstNotNil(t *testing.T) {
	s := &MockState{}
	tok := &lexer.Token{}
	node := &common.MockNode{}
	ent := Entry{
		Power: 42,
		First: ExprFirst(func(es State, pow int, eTok *lexer.Token) (common.Node, error) {
			assert.Same(t, s, es)
			assert.Equal(t, 42, pow)
			assert.Same(t, tok, eTok)
			return node, assert.AnError
		}),
	}

	result, err := ent.callFirst(s, tok)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, node, result)
}

func TestEntryCallNextNil(t *testing.T) {
	s := &MockState{}
	l := &common.MockNode{}
	tok := &lexer.Token{}
	ent := Entry{}

	result, err := ent.callNext(s, l, tok)

	assert.True(t, errors.Is(err, ErrUnexpectedToken))
	assert.Nil(t, result)
}

func TestEntryCallNextNotNil(t *testing.T) {
	s := &MockState{}
	l := &common.MockNode{}
	tok := &lexer.Token{}
	node := &common.MockNode{}
	ent := Entry{
		Power: 42,
		Next: ExprNext(func(es State, pow int, el common.Node, eTok *lexer.Token) (common.Node, error) {
			assert.Same(t, s, es)
			assert.Equal(t, 42, pow)
			assert.Same(t, l, el)
			assert.Same(t, tok, eTok)
			return node, assert.AnError
		}),
	}

	result, err := ent.callNext(s, l, tok)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, node, result)
}

func TestEntryCallStmtNil(t *testing.T) {
	s := &MockState{}
	tok := &lexer.Token{}
	ent := Entry{}

	result, err := ent.callStmt(s, tok)

	assert.True(t, errors.Is(err, ErrUnexpectedToken))
	assert.Nil(t, result)
}

func TestEntryCallStmtNotNil(t *testing.T) {
	s := &MockState{}
	tok := &lexer.Token{}
	node := &common.MockNode{}
	ent := Entry{
		Stmt: Statement(func(es State, eTok *lexer.Token) (common.Node, error) {
			assert.Same(t, s, es)
			assert.Same(t, tok, eTok)
			return node, assert.AnError
		}),
	}

	result, err := ent.callStmt(s, tok)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, node, result)
}
