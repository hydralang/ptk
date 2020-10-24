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

package tokenstreams

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/lexer"
)

func TestChanTokenStreamImplementsTokenStream(t *testing.T) {
	assert.Implements(t, (*lexer.TokenStream)(nil), &ChanTokenStream{})
}

func TestNewChanTokenStream(t *testing.T) {
	result := NewChanTokenStream()

	assert.NotNil(t, result.Chan)
}

func TestChanTokenStreamNextOpen(t *testing.T) {
	tok := &lexer.Token{}
	obj := NewChanTokenStream()
	obj.Chan <- tok

	result := obj.Next()

	assert.Same(t, tok, result)
}

func TestChanTokenStreamNextClosed(t *testing.T) {
	obj := NewChanTokenStream()
	close(obj.Chan)

	result := obj.Next()

	assert.Nil(t, result)
}

func TestChanTokenStreamPushBase(t *testing.T) {
	obj := NewChanTokenStream()
	tok := &lexer.Token{}

	result := obj.Push(tok)

	assert.True(t, result)
	assert.Same(t, tok, <-obj.Chan)
}

func TestChanTokenStreamPushDone(t *testing.T) {
	obj := NewChanTokenStream()
	close(obj.Chan)
	tok := &lexer.Token{}

	result := obj.Push(tok)

	assert.False(t, result)
}

func TestChanTokenStreamDone(t *testing.T) {
	obj := NewChanTokenStream()

	obj.Done()

	_, ok := <-obj.Chan
	assert.False(t, ok)
}
