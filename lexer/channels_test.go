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

func TestChanLexerImplementsLexer(t *testing.T) {
	assert.Implements(t, (*ILexer)(nil), &ChanLexer{})
}

func TestNewChanLexer(t *testing.T) {
	result := NewChanLexer()

	assert.NotNil(t, result.Chan)
}

func TestChanLexerNextOpen(t *testing.T) {
	tok := &Token{}
	obj := NewChanLexer()
	obj.Chan <- tok

	result := obj.Next()

	assert.Same(t, tok, result)
}

func TestChanLexerNextClosed(t *testing.T) {
	obj := NewChanLexer()
	close(obj.Chan)

	result := obj.Next()

	assert.Nil(t, result)
}

func TestChanLexerPushBase(t *testing.T) {
	obj := NewChanLexer()
	tok := &Token{}

	result := obj.Push(tok)

	assert.True(t, result)
	assert.Same(t, tok, <-obj.Chan)
}

func TestChanLexerPushDone(t *testing.T) {
	obj := NewChanLexer()
	close(obj.Chan)
	tok := &Token{}

	result := obj.Push(tok)

	assert.False(t, result)
}

func TestChanLexerDone(t *testing.T) {
	obj := NewChanLexer()

	obj.Done()

	_, ok := <-obj.Chan
	assert.False(t, ok)
}
