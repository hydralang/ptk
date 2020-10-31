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

	"github.com/hydralang/ptk/lexer"
)

func TestEntryCallFirstNil(t *testing.T) {
	p := &Parser{}
	tok := &lexer.Token{}
	ent := Entry{}

	result, err := ent.callFirst(p, tok)

	assert.True(t, errors.Is(err, ErrUnexpectedToken))
	assert.Nil(t, result)
}

func TestEntryCallFirstNotNil(t *testing.T) {
	p := &Parser{}
	tok := &lexer.Token{}
	node := &mockNode{}
	ent := Entry{
		Power: 42,
		First: ExprFirst(func(ep *Parser, pow int, eTok *lexer.Token) (Node, error) {
			assert.Same(t, p, ep)
			assert.Equal(t, 42, pow)
			assert.Same(t, tok, eTok)
			return node, assert.AnError
		}),
	}

	result, err := ent.callFirst(p, tok)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, node, result)
}

func TestEntryCallNextNil(t *testing.T) {
	p := &Parser{}
	l := &mockNode{}
	tok := &lexer.Token{}
	ent := Entry{}

	result, err := ent.callNext(p, l, tok)

	assert.True(t, errors.Is(err, ErrUnexpectedToken))
	assert.Nil(t, result)
}

func TestEntryCallNextNotNil(t *testing.T) {
	p := &Parser{}
	l := &mockNode{}
	tok := &lexer.Token{}
	node := &mockNode{}
	ent := Entry{
		Power: 42,
		Next: ExprNext(func(ep *Parser, pow int, el Node, eTok *lexer.Token) (Node, error) {
			assert.Same(t, p, ep)
			assert.Equal(t, 42, pow)
			assert.Same(t, l, el)
			assert.Same(t, tok, eTok)
			return node, assert.AnError
		}),
	}

	result, err := ent.callNext(p, l, tok)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, node, result)
}

func TestEntryCallStmtNil(t *testing.T) {
	p := &Parser{}
	tok := &lexer.Token{}
	ent := Entry{}

	result, err := ent.callStmt(p, tok)

	assert.True(t, errors.Is(err, ErrUnexpectedToken))
	assert.Nil(t, result)
}

func TestEntryCallStmtNotNil(t *testing.T) {
	p := &Parser{}
	tok := &lexer.Token{}
	node := &mockNode{}
	ent := Entry{
		Stmt: Statement(func(ep *Parser, eTok *lexer.Token) (Node, error) {
			assert.Same(t, p, ep)
			assert.Same(t, tok, eTok)
			return node, assert.AnError
		}),
	}

	result, err := ent.callStmt(p, tok)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, node, result)
}
