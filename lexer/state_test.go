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

	"github.com/hydralang/ptk/common"
)

func TestMockStateImplementsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &MockState{})
}

func TestMockStateLexerNil(t *testing.T) {
	obj := &MockState{}
	obj.On("Lexer").Return(nil)

	result := obj.Lexer()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateLexerNotNil(t *testing.T) {
	lex := &MockLexer{}
	obj := &MockState{}
	obj.On("Lexer").Return(lex)

	result := obj.Lexer()

	assert.Same(t, lex, result)
	obj.AssertExpectations(t)
}

func TestMockStateCharStreamNil(t *testing.T) {
	obj := &MockState{}
	obj.On("CharStream").Return(nil)

	result := obj.CharStream()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateCharStreamNotNil(t *testing.T) {
	cs := &MockCharStream{}
	obj := &MockState{}
	obj.On("CharStream").Return(cs)

	result := obj.CharStream()

	assert.Same(t, cs, result)
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

func TestMockStateClassifierNil(t *testing.T) {
	obj := &MockState{}
	obj.On("Classifier").Return(nil)

	result := obj.Classifier()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateClassifierNotNil(t *testing.T) {
	expected := &MockClassifier{}
	obj := &MockState{}
	obj.On("Classifier").Return(expected)

	result := obj.Classifier()

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStatePushClassifier(t *testing.T) {
	cls := &MockClassifier{}
	obj := &MockState{}
	obj.On("PushClassifier", cls)

	obj.PushClassifier(cls)

	obj.AssertExpectations(t)
}

func TestMockStatePopClassifierNil(t *testing.T) {
	obj := &MockState{}
	obj.On("PopClassifier").Return(nil)

	result := obj.PopClassifier()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStatePopClassifierNotNil(t *testing.T) {
	expected := &MockClassifier{}
	obj := &MockState{}
	obj.On("PopClassifier").Return(expected)

	result := obj.PopClassifier()

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStateSetClassifierNil(t *testing.T) {
	cls := &MockClassifier{}
	obj := &MockState{}
	obj.On("SetClassifier", cls).Return(nil)

	result := obj.SetClassifier(cls)

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockStateSetClassifierNotNil(t *testing.T) {
	cls := &MockClassifier{}
	expected := &MockClassifier{}
	obj := &MockState{}
	obj.On("SetClassifier", cls).Return(expected)

	result := obj.SetClassifier(cls)

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockStatePush(t *testing.T) {
	tok := &common.Token{}
	obj := &MockState{}
	obj.On("Push", tok).Return(true)

	result := obj.Push(tok)

	assert.True(t, result)
	obj.AssertExpectations(t)
}
