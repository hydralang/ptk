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
