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

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/lexer"
)

func TestExpectedTypesNil(t *testing.T) {
	result := expectedTypes(nil)

	assert.Equal(t, "", result)
}

func TestExpectedTypesEmpty(t *testing.T) {
	result := expectedTypes([]string{})

	assert.Equal(t, "", result)
}

func TestExpectedTypesOne(t *testing.T) {
	result := expectedTypes([]string{"one"})

	assert.Equal(t, "; expected token of type \"one\"", result)
}

func TestExpectedTypesTwo(t *testing.T) {
	result := expectedTypes([]string{"one", "two"})

	assert.Equal(t, "; expected tokens of type \"one\" or \"two\"", result)
}

func TestExpectedTypesThree(t *testing.T) {
	result := expectedTypes([]string{"one", "two", "three"})

	assert.Equal(t, "; expected tokens of type \"one\", \"two\", or \"three\"", result)
}

func TestUnknownTokenTypeNoTypes(t *testing.T) {
	tok := &lexer.Token{
		Type: "type",
	}

	result := UnknownTokenType(tok)

	assert.EqualError(t, result, "Unknown token type \"type\"")
}

func TestUnknownTokenTypeWithTypes(t *testing.T) {
	tok := &lexer.Token{
		Type: "type",
	}

	result := UnknownTokenType(tok, "one")

	assert.EqualError(t, result, "Unknown token type \"type\"; expected token of type \"one\"")
}

func TestExpectedTokenNoTypes(t *testing.T) {
	result := ExpectedToken()

	assert.EqualError(t, result, "No tokens available, but expected more")
}

func TestExpectedTokenWithTypes(t *testing.T) {
	result := ExpectedToken("one")

	assert.EqualError(t, result, "No tokens available; expected token of type \"one\"")
}

func TestUnexpectedTokenNoTypes(t *testing.T) {
	tok := &lexer.Token{
		Type: "type",
	}

	result := UnexpectedToken(tok)

	assert.EqualError(t, result, "Unexpected token of type \"type\"")
}

func TestUnexpectedTokenWithTypes(t *testing.T) {
	tok := &lexer.Token{
		Type: "type",
	}

	result := UnexpectedToken(tok, "one")

	assert.EqualError(t, result, "Unexpected token of type \"type\"; expected token of type \"one\"")
}
