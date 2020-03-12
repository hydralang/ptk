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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenImplementsNode(t *testing.T) {
	assert.Implements(t, (*Node)(nil), &Token{})
}

func TestTokenLocation(t *testing.T) {
	loc := &MockLocation{}
	obj := &Token{
		Loc: loc,
	}

	result := obj.Location()

	assert.Same(t, loc, result)
}

func TestTokenChildren(t *testing.T) {
	obj := &Token{}

	result := obj.Children()

	assert.Equal(t, []Node{}, result)
}

func TestTokenStringBase(t *testing.T) {
	loc := &MockLocation{}
	loc.On("String").Return("location")
	obj := &Token{
		Type: "type",
		Loc:  loc,
	}

	result := obj.String()

	assert.Equal(t, "location: <type> token", result)
}

func TestTokenStringNoLocation(t *testing.T) {
	obj := &Token{
		Type: "type",
	}

	result := obj.String()

	assert.Equal(t, "<type> token", result)
}

func TestTokenStringValue(t *testing.T) {
	loc := &MockLocation{}
	loc.On("String").Return("location")
	obj := &Token{
		Type:  "type",
		Loc:   loc,
		Value: 42,
	}

	result := obj.String()

	assert.Equal(t, "location: <type> token: 42", result)
}

func TestMockTokenStreamImplementsTokenStream(t *testing.T) {
	assert.Implements(t, (*TokenStream)(nil), &MockTokenStream{})
}

func TestMockTokenStreamNextNil(t *testing.T) {
	obj := &MockTokenStream{}
	obj.On("Next").Return(nil)

	result := obj.Next()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockTokenStreamNextNotNil(t *testing.T) {
	tok := &Token{}
	obj := &MockTokenStream{}
	obj.On("Next").Return(tok)

	result := obj.Next()

	assert.Same(t, tok, result)
	obj.AssertExpectations(t)
}

func TestNewChanTokenStream(t *testing.T) {
	result := NewChanTokenStream()

	assert.NotNil(t, result.Chan)
}

func TestChanTokenStreamNextOpen(t *testing.T) {
	tok := &Token{}
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
	tok := &Token{}

	result := obj.Push(tok)

	assert.True(t, result)
	assert.Same(t, tok, <-obj.Chan)
}

func TestChanTokenStreamPushDone(t *testing.T) {
	obj := NewChanTokenStream()
	close(obj.Chan)
	tok := &Token{}

	result := obj.Push(tok)

	assert.False(t, result)
}

func TestChanTokenStreamDone(t *testing.T) {
	obj := NewChanTokenStream()

	obj.Done()

	_, ok := <-obj.Chan
	assert.False(t, ok)
}

func TestNewListTokenStream(t *testing.T) {
	toks := []*Token{{}, {}, {}}

	result := NewListTokenStream(toks)

	i := 0
	for tok := result.Next(); tok != nil; tok = result.Next() {
		assert.Same(t, toks[i], tok)
		i++
	}
	assert.Equal(t, len(toks), i)
}
