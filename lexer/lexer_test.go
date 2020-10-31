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
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLexerImplementsILexer(t *testing.T) {
	assert.Implements(t, (*ILexer)(nil), &Lexer{})
}

func TestNewBase(t *testing.T) {
	src := &mockScanner{}
	state := &mockState{}

	result := New(src, state)

	require.NotNil(t, result.Scanner)
	assert.Same(t, src, result.Scanner.(*BackTracker).Src)
	assert.Same(t, state, result.State)
	assert.Equal(t, &list.List{}, result.toks)
}

func TestNewWithBackTracker(t *testing.T) {
	src := &mockBackTracker{}
	state := &mockState{}

	result := New(src, state)

	require.NotNil(t, result.Scanner)
	assert.Same(t, src, result.Scanner)
	assert.Same(t, state, result.State)
	assert.Equal(t, &list.List{}, result.toks)
}

func TestLexerNextInternalBase(t *testing.T) {
	bt := &mockBackTracker{}
	state := &mockState{}
	obj := &Lexer{
		Scanner: bt,
		State:   state,
		toks:    &list.List{},
	}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Times(2)
	bt.On("Accept", 0).Once()
	rec1 := &mockRecognizer{}
	rec1.On("Recognize", obj).Return(false)
	rec2 := &mockRecognizer{}
	rec2.On("Recognize", obj).Return(true)
	rec3 := &mockRecognizer{}
	cls := &mockClassifier{}
	cls.On("Classify", obj).Return([]Recognizer{rec1, rec2, rec3})
	state.On("Classifier").Return(cls)

	obj.next()

	bt.AssertExpectations(t)
	rec1.AssertExpectations(t)
	rec2.AssertExpectations(t)
	rec3.AssertExpectations(t)
	cls.AssertExpectations(t)
	state.AssertExpectations(t)
}

func TestLexerNextInternalUnrecognized(t *testing.T) {
	bt := &mockBackTracker{}
	state := &mockState{}
	obj := &Lexer{
		Scanner: bt,
		State:   state,
		toks:    &list.List{},
	}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Times(4)
	bt.On("Accept", 0).Once()
	rec1 := &mockRecognizer{}
	rec1.On("Recognize", obj).Return(false)
	rec2 := &mockRecognizer{}
	rec2.On("Recognize", obj).Return(false)
	rec3 := &mockRecognizer{}
	rec3.On("Recognize", obj).Return(false)
	cls := &mockClassifier{}
	cls.On("Classify", obj).Return([]Recognizer{rec1, rec2, rec3})
	cls.On("Error", obj)
	state.On("Classifier").Return(cls)

	obj.next()

	bt.AssertExpectations(t)
	rec1.AssertExpectations(t)
	rec2.AssertExpectations(t)
	rec3.AssertExpectations(t)
	cls.AssertExpectations(t)
	state.AssertExpectations(t)
}

func TestLexerNextInternalUnclassified(t *testing.T) {
	bt := &mockBackTracker{}
	state := &mockState{}
	obj := &Lexer{
		Scanner: bt,
		State:   state,
		toks:    &list.List{},
	}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Once()
	bt.On("Accept", 0).Once()
	cls := &mockClassifier{}
	cls.On("Classify", obj).Return([]Recognizer{})
	cls.On("Error", obj)
	state.On("Classifier").Return(cls)

	obj.next()

	bt.AssertExpectations(t)
	cls.AssertExpectations(t)
	state.AssertExpectations(t)
}

func TestLexerNextQueued(t *testing.T) {
	tok := &Token{}
	bt := &mockBackTracker{}
	state := &mockState{}
	obj := &Lexer{
		Scanner: bt,
		State:   state,
		toks:    &list.List{},
	}
	obj.toks.PushBack(tok)

	result := obj.Next()

	assert.Same(t, tok, result)
	assert.Equal(t, 0, obj.toks.Len())
	bt.AssertExpectations(t)
	state.AssertExpectations(t)
}

func TestLexerNextLex(t *testing.T) {
	tok := &Token{}
	bt := &mockBackTracker{}
	state := &mockState{}
	obj := &Lexer{
		Scanner: bt,
		State:   state,
		toks:    &list.List{},
	}
	bt.On("More").Return(true)
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Once()
	bt.On("Accept", 0).Once()
	cls := &mockClassifier{}
	cls.On("Classify", obj).Return([]Recognizer{}).Once()
	cls.On("Error", obj).Run(func(args mock.Arguments) {
		obj.toks.PushBack(tok)
	}).Once()
	state.On("Classifier").Return(cls)

	result := obj.Next()

	assert.Same(t, tok, result)
	assert.Equal(t, 0, obj.toks.Len())
	bt.AssertExpectations(t)
	state.AssertExpectations(t)
	cls.AssertExpectations(t)
}

func TestLexerNextClosed(t *testing.T) {
	bt := &mockBackTracker{}
	state := &mockState{}
	obj := &Lexer{
		Scanner: bt,
		State:   state,
		toks:    &list.List{},
	}
	bt.On("More").Return(false)

	result := obj.Next()

	assert.Nil(t, result)
	assert.Equal(t, 0, obj.toks.Len())
	bt.AssertExpectations(t)
	state.AssertExpectations(t)
}

func TestLexerPush(t *testing.T) {
	tok := &Token{}
	obj := &Lexer{
		toks: &list.List{},
	}

	result := obj.Push(tok)

	assert.True(t, result)
	assert.Equal(t, 1, obj.toks.Len())
	assert.Same(t, tok, obj.toks.Front().Value)
}
