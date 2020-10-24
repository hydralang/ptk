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

	"github.com/stretchr/testify/assert"
)

func TestListLexerImplementsLexer(t *testing.T) {
	assert.Implements(t, (*ILexer)(nil), &ListLexer{})
}

func TestNewListLexer(t *testing.T) {
	toks := []*Token{{}, {}, {}}

	result := NewListLexer(toks)

	assert.Equal(t, &ListLexer{
		toks: toks,
	}, result)
}

func TestListLexerNextUnstarted(t *testing.T) {
	toks := []*Token{{}, {}, {}}
	obj := &ListLexer{
		toks: toks,
	}

	result := obj.Next()

	assert.Same(t, toks[0], result)
	assert.Equal(t, &ListLexer{
		toks:    toks,
		started: true,
	}, obj)
}

func TestListLexerNextStarted(t *testing.T) {
	toks := []*Token{{}, {}, {}}
	obj := &ListLexer{
		toks:    toks,
		started: true,
	}

	result := obj.Next()

	assert.Same(t, toks[1], result)
	assert.Equal(t, &ListLexer{
		toks:    toks,
		idx:     1,
		started: true,
	}, obj)
}

func TestListLexerNextEnding(t *testing.T) {
	toks := []*Token{{}, {}, {}}
	obj := &ListLexer{
		toks:    toks,
		idx:     2,
		started: true,
	}

	result := obj.Next()

	assert.Nil(t, result)
	assert.Equal(t, &ListLexer{
		toks:    toks,
		idx:     2,
		started: true,
	}, obj)
}
