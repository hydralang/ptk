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
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/common"
)

func TestMockLexerImplementsLexer(t *testing.T) {
	assert.Implements(t, (*Lexer)(nil), &MockLexer{})
}

func TestMockLexerClassifierNil(t *testing.T) {
	obj := &MockLexer{}
	obj.On("Classifier").Return(nil)

	result := obj.Classifier()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockLexerClassifierNotNil(t *testing.T) {
	cls := &MockClassifier{}
	obj := &MockLexer{}
	obj.On("Classifier").Return(cls)

	result := obj.Classifier()

	assert.Equal(t, cls, result)
	obj.AssertExpectations(t)
}

func TestMockLexerLexNil(t *testing.T) {
	cs := &MockCharStream{}
	obj := &MockLexer{}
	obj.On("Lex", cs, mock.Anything).Return(nil)

	result := obj.Lex(cs)

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockLexerLexNotNil(t *testing.T) {
	stream := &common.MockTokenStream{}
	cs := &MockCharStream{}
	obj := &MockLexer{}
	obj.On("Lex", cs, mock.Anything).Return(stream)

	result := obj.Lex(cs)

	assert.Same(t, stream, result)
	obj.AssertExpectations(t)
}
