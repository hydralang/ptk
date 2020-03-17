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
	"github.com/stretchr/testify/require"

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
	cs := &common.MockCharStream{}
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

func TestStateImplementsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &state{})
}

func TestNewState(t *testing.T) {
	cls := &MockClassifier{}
	lexer := &MockLexer{}
	lexer.On("Classifier").Return(cls)
	src := &common.MockCharStream{}
	var opt1Called State
	var opt2Called State
	options := []Option{
		func(s State) {
			opt1Called = s
		},
		func(s State) {
			opt2Called = s
		},
	}

	result := NewState(lexer, src, options)

	assert.Same(t, result, opt1Called)
	assert.Same(t, result, opt2Called)
	state, ok := result.(*state)
	require.True(t, ok)
	assert.Same(t, lexer, state.lexer)
	assert.Same(t, src, state.src)
	assert.Equal(t, NewBackTracker(src, TrackAll), state.bt)
	assert.Equal(t, 0, state.appState.Len())
	assert.Equal(t, 1, state.cls.Len())
	assert.Same(t, cls, state.cls.Get())
	assert.Equal(t, 0, state.toks.Len())
}

func TestNextInternalBase(t *testing.T) {
	bt := &MockBackTracker{}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Times(2)
	bt.On("Accept", 0).Once()
	obj := &state{
		bt:  bt,
		cls: common.NewStack(),
	}
	rec1 := &MockRecognizer{}
	rec1.On("Recognize", obj, bt).Return(false)
	rec2 := &MockRecognizer{}
	rec2.On("Recognize", obj, bt).Return(true)
	rec3 := &MockRecognizer{}
	cls := &MockClassifier{}
	cls.On("Classify", obj, bt).Return([]Recognizer{rec1, rec2, rec3})
	obj.cls.Push(cls)

	obj.next()

	cls.AssertExpectations(t)
	rec1.AssertExpectations(t)
	rec2.AssertExpectations(t)
	rec3.AssertExpectations(t)
	bt.AssertExpectations(t)
}

func TestNextInternalUnrecognized(t *testing.T) {
	bt := &MockBackTracker{}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Times(4)
	bt.On("Accept", 0).Once()
	obj := &state{
		bt:  bt,
		cls: common.NewStack(),
	}
	rec1 := &MockRecognizer{}
	rec1.On("Recognize", obj, bt).Return(false)
	rec2 := &MockRecognizer{}
	rec2.On("Recognize", obj, bt).Return(false)
	rec3 := &MockRecognizer{}
	rec3.On("Recognize", obj, bt).Return(false)
	cls := &MockClassifier{}
	cls.On("Classify", obj, bt).Return([]Recognizer{rec1, rec2, rec3})
	cls.On("Error", obj, bt)
	obj.cls.Push(cls)

	obj.next()

	cls.AssertExpectations(t)
	rec1.AssertExpectations(t)
	rec2.AssertExpectations(t)
	rec3.AssertExpectations(t)
	bt.AssertExpectations(t)
}

func TestNextInternalUnclassified(t *testing.T) {
	bt := &MockBackTracker{}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Once()
	bt.On("Accept", 0).Once()
	obj := &state{
		bt:  bt,
		cls: common.NewStack(),
	}
	cls := &MockClassifier{}
	cls.On("Classify", obj, bt).Return([]Recognizer{})
	cls.On("Error", obj, bt)
	obj.cls.Push(cls)

	obj.next()

	cls.AssertExpectations(t)
	bt.AssertExpectations(t)
}

type fakeClassifier struct {
	tok *common.Token
}

func (f *fakeClassifier) Classify(state State, str BackTracker) []Recognizer {
	return []Recognizer{}
}

func (f *fakeClassifier) Error(state State, str BackTracker) {
	state.Push(f.tok)
}

func TestStateNextQueued(t *testing.T) {
	tok := &common.Token{}
	bt := &MockBackTracker{}
	obj := &state{
		bt:   bt,
		cls:  common.NewStack(),
		toks: &list.List{},
	}
	obj.cls.Push(&fakeClassifier{})
	obj.toks.PushBack(tok)

	result := obj.Next()

	assert.Same(t, tok, result)
	assert.Equal(t, 0, obj.toks.Len())
	bt.AssertExpectations(t)
}

func TestStateNextLex(t *testing.T) {
	tok := &common.Token{}
	src := &common.MockCharStream{}
	bt := &MockBackTracker{}
	bt.On("More").Return(true)
	bt.On("SetMax", TrackAll)
	bt.On("BackTrack")
	bt.On("Accept", 0)
	obj := &state{
		src:  src,
		bt:   bt,
		cls:  common.NewStack(),
		toks: &list.List{},
	}
	obj.cls.Push(&fakeClassifier{tok: tok})

	result := obj.Next()

	assert.Same(t, tok, result)
	assert.Equal(t, 0, obj.toks.Len())
	bt.AssertExpectations(t)
}

func TestStateNextClosed(t *testing.T) {
	bt := &MockBackTracker{}
	bt.On("More").Return(false)
	obj := &state{
		bt:   bt,
		cls:  common.NewStack(),
		toks: &list.List{},
	}
	obj.cls.Push(&fakeClassifier{})

	result := obj.Next()

	assert.Nil(t, result)
	assert.Equal(t, 0, obj.toks.Len())
	bt.AssertExpectations(t)
}

func TestStateLexer(t *testing.T) {
	l := &MockLexer{}
	obj := &state{
		lexer: l,
	}

	result := obj.Lexer()

	assert.Same(t, l, result)
}

func TestStateCharStream(t *testing.T) {
	src := &common.MockCharStream{}
	obj := &state{
		src: src,
	}

	result := obj.CharStream()

	assert.Same(t, src, result)
}

func TestStateAppState(t *testing.T) {
	appStack := &common.MockStack{}
	appStack.On("Get").Return("state")
	obj := &state{
		appState: appStack,
	}

	result := obj.AppState()

	assert.Equal(t, "state", result)
	appStack.AssertExpectations(t)
}

func TestStatePushAppState(t *testing.T) {
	appStack := &common.MockStack{}
	appStack.On("Push", "state")
	obj := &state{
		appState: appStack,
	}

	obj.PushAppState("state")

	appStack.AssertExpectations(t)
}

func TestStatePopAppState(t *testing.T) {
	appStack := &common.MockStack{}
	appStack.On("Pop").Return("state")
	obj := &state{
		appState: appStack,
	}

	result := obj.PopAppState()

	assert.Equal(t, "state", result)
	appStack.AssertExpectations(t)
}

func TestStateSetAppState(t *testing.T) {
	appStack := &common.MockStack{}
	appStack.On("Set", "new").Return("old")
	obj := &state{
		appState: appStack,
	}

	result := obj.SetAppState("new")

	assert.Equal(t, "old", result)
	appStack.AssertExpectations(t)
}

// XXX

func TestStateClassifierBase(t *testing.T) {
	cls := &MockClassifier{}
	clsStack := &common.MockStack{}
	clsStack.On("Get").Return(cls)
	obj := &state{
		cls: clsStack,
	}

	result := obj.Classifier()

	assert.Same(t, cls, result)
	clsStack.AssertExpectations(t)
}

func TestStateClassifierNil(t *testing.T) {
	clsStack := &common.MockStack{}
	clsStack.On("Get").Return(nil)
	obj := &state{
		cls: clsStack,
	}

	result := obj.Classifier()

	assert.Nil(t, result)
	clsStack.AssertExpectations(t)
}

func TestStatePushClassifier(t *testing.T) {
	cls := &MockClassifier{}
	clsStack := &common.MockStack{}
	clsStack.On("Push", cls)
	obj := &state{
		cls: clsStack,
	}

	obj.PushClassifier(cls)

	clsStack.AssertExpectations(t)
}

func TestStatePopClassifierBase(t *testing.T) {
	cls := &MockClassifier{}
	clsStack := &common.MockStack{}
	clsStack.On("Pop").Return(cls)
	obj := &state{
		cls: clsStack,
	}

	result := obj.PopClassifier()

	assert.Same(t, cls, result)
	clsStack.AssertExpectations(t)
}

func TestStatePopClassifierNil(t *testing.T) {
	clsStack := &common.MockStack{}
	clsStack.On("Pop").Return(nil)
	obj := &state{
		cls: clsStack,
	}

	result := obj.PopClassifier()

	assert.Nil(t, result)
	clsStack.AssertExpectations(t)
}

func TestStateSetClassifierBase(t *testing.T) {
	cls := &MockClassifier{}
	newCls := &MockClassifier{}
	clsStack := &common.MockStack{}
	clsStack.On("Set", newCls).Return(cls)
	obj := &state{
		cls: clsStack,
	}

	result := obj.SetClassifier(newCls)

	assert.Same(t, cls, result)
	clsStack.AssertExpectations(t)
}

func TestStateSetClassifierNil(t *testing.T) {
	newCls := &MockClassifier{}
	clsStack := &common.MockStack{}
	clsStack.On("Set", newCls).Return(nil)
	obj := &state{
		cls: clsStack,
	}

	result := obj.SetClassifier(newCls)

	assert.Nil(t, result)
	clsStack.AssertExpectations(t)
}

func TestStatePush(t *testing.T) {
	tok := &common.Token{}
	obj := &state{
		toks: &list.List{},
	}

	result := obj.Push(tok)

	assert.True(t, result)
	assert.Equal(t, 1, obj.toks.Len())
	assert.Same(t, tok, obj.toks.Front().Value)
}
