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

	"github.com/hydralang/ptk/internal"
	"github.com/hydralang/ptk/scanner"
)

type mockState struct {
	MockTokenStream
}

func (m *mockState) Lexer() Lexer {
	args := m.MethodCalled("Lexer")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Lexer)
	}

	return nil
}

func (m *mockState) CharStream() scanner.Scanner {
	args := m.MethodCalled("CharStream")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(scanner.Scanner)
	}

	return nil
}

func (m *mockState) AppState() interface{} {
	args := m.MethodCalled("AppState")

	return args.Get(0)
}

func (m *mockState) PushAppState(state interface{}) {
	m.MethodCalled("PushAppState", state)
}

func (m *mockState) PopAppState() interface{} {
	args := m.MethodCalled("PopAppState")

	return args.Get(0)
}

func (m *mockState) SetAppState(state interface{}) interface{} {
	args := m.MethodCalled("SetAppState", state)

	return args.Get(0)
}

func (m *mockState) Classifier() Classifier {
	args := m.MethodCalled("Classifier")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

func (m *mockState) PushClassifier(cls Classifier) {
	m.MethodCalled("PushClassifier", cls)
}

func (m *mockState) PopClassifier() Classifier {
	args := m.MethodCalled("PopClassifier")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

func (m *mockState) SetClassifier(cls Classifier) Classifier {
	args := m.MethodCalled("SetClassifier", cls)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Classifier)
	}

	return nil
}

func (m *mockState) Push(tok *Token) bool {
	args := m.MethodCalled("Push", tok)

	return args.Bool(0)
}

func TestStateImplementsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &state{})
}

func TestNewState(t *testing.T) {
	cls := &mockClassifier{}
	lexer := &mockLexer{}
	lexer.On("Classifier").Return(cls)
	src := &mockScanner{}
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
	bt := &mockBackTracker{}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Times(2)
	bt.On("Accept", 0).Once()
	obj := &state{
		bt:  bt,
		cls: internal.NewStack(),
	}
	rec1 := &mockRecognizer{}
	rec1.On("Recognize", obj, bt).Return(false)
	rec2 := &mockRecognizer{}
	rec2.On("Recognize", obj, bt).Return(true)
	rec3 := &mockRecognizer{}
	cls := &mockClassifier{}
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
	bt := &mockBackTracker{}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Times(4)
	bt.On("Accept", 0).Once()
	obj := &state{
		bt:  bt,
		cls: internal.NewStack(),
	}
	rec1 := &mockRecognizer{}
	rec1.On("Recognize", obj, bt).Return(false)
	rec2 := &mockRecognizer{}
	rec2.On("Recognize", obj, bt).Return(false)
	rec3 := &mockRecognizer{}
	rec3.On("Recognize", obj, bt).Return(false)
	cls := &mockClassifier{}
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
	bt := &mockBackTracker{}
	bt.On("SetMax", TrackAll).Once()
	bt.On("BackTrack").Once()
	bt.On("Accept", 0).Once()
	obj := &state{
		bt:  bt,
		cls: internal.NewStack(),
	}
	cls := &mockClassifier{}
	cls.On("Classify", obj, bt).Return([]Recognizer{})
	cls.On("Error", obj, bt)
	obj.cls.Push(cls)

	obj.next()

	cls.AssertExpectations(t)
	bt.AssertExpectations(t)
}

type fakeClassifier struct {
	tok *Token
}

func (f *fakeClassifier) Classify(state State, str IBackTracker) []Recognizer {
	return []Recognizer{}
}

func (f *fakeClassifier) Error(state State, str IBackTracker) {
	state.Push(f.tok)
}

func TestStateNextQueued(t *testing.T) {
	tok := &Token{}
	bt := &mockBackTracker{}
	obj := &state{
		bt:   bt,
		cls:  internal.NewStack(),
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
	tok := &Token{}
	src := &mockScanner{}
	bt := &mockBackTracker{}
	bt.On("More").Return(true)
	bt.On("SetMax", TrackAll)
	bt.On("BackTrack")
	bt.On("Accept", 0)
	obj := &state{
		src:  src,
		bt:   bt,
		cls:  internal.NewStack(),
		toks: &list.List{},
	}
	obj.cls.Push(&fakeClassifier{tok: tok})

	result := obj.Next()

	assert.Same(t, tok, result)
	assert.Equal(t, 0, obj.toks.Len())
	bt.AssertExpectations(t)
}

func TestStateNextClosed(t *testing.T) {
	bt := &mockBackTracker{}
	bt.On("More").Return(false)
	obj := &state{
		bt:   bt,
		cls:  internal.NewStack(),
		toks: &list.List{},
	}
	obj.cls.Push(&fakeClassifier{})

	result := obj.Next()

	assert.Nil(t, result)
	assert.Equal(t, 0, obj.toks.Len())
	bt.AssertExpectations(t)
}

func TestStateLexer(t *testing.T) {
	l := &mockLexer{}
	obj := &state{
		lexer: l,
	}

	result := obj.Lexer()

	assert.Same(t, l, result)
}

func TestStateCharStream(t *testing.T) {
	src := &mockScanner{}
	obj := &state{
		src: src,
	}

	result := obj.CharStream()

	assert.Same(t, src, result)
}

func TestStateAppState(t *testing.T) {
	appStack := &internal.MockStack{}
	appStack.On("Get").Return("state")
	obj := &state{
		appState: appStack,
	}

	result := obj.AppState()

	assert.Equal(t, "state", result)
	appStack.AssertExpectations(t)
}

func TestStatePushAppState(t *testing.T) {
	appStack := &internal.MockStack{}
	appStack.On("Push", "state")
	obj := &state{
		appState: appStack,
	}

	obj.PushAppState("state")

	appStack.AssertExpectations(t)
}

func TestStatePopAppState(t *testing.T) {
	appStack := &internal.MockStack{}
	appStack.On("Pop").Return("state")
	obj := &state{
		appState: appStack,
	}

	result := obj.PopAppState()

	assert.Equal(t, "state", result)
	appStack.AssertExpectations(t)
}

func TestStateSetAppState(t *testing.T) {
	appStack := &internal.MockStack{}
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
	cls := &mockClassifier{}
	clsStack := &internal.MockStack{}
	clsStack.On("Get").Return(cls)
	obj := &state{
		cls: clsStack,
	}

	result := obj.Classifier()

	assert.Same(t, cls, result)
	clsStack.AssertExpectations(t)
}

func TestStateClassifierNil(t *testing.T) {
	clsStack := &internal.MockStack{}
	clsStack.On("Get").Return(nil)
	obj := &state{
		cls: clsStack,
	}

	result := obj.Classifier()

	assert.Nil(t, result)
	clsStack.AssertExpectations(t)
}

func TestStatePushClassifier(t *testing.T) {
	cls := &mockClassifier{}
	clsStack := &internal.MockStack{}
	clsStack.On("Push", cls)
	obj := &state{
		cls: clsStack,
	}

	obj.PushClassifier(cls)

	clsStack.AssertExpectations(t)
}

func TestStatePopClassifierBase(t *testing.T) {
	cls := &mockClassifier{}
	clsStack := &internal.MockStack{}
	clsStack.On("Pop").Return(cls)
	obj := &state{
		cls: clsStack,
	}

	result := obj.PopClassifier()

	assert.Same(t, cls, result)
	clsStack.AssertExpectations(t)
}

func TestStatePopClassifierNil(t *testing.T) {
	clsStack := &internal.MockStack{}
	clsStack.On("Pop").Return(nil)
	obj := &state{
		cls: clsStack,
	}

	result := obj.PopClassifier()

	assert.Nil(t, result)
	clsStack.AssertExpectations(t)
}

func TestStateSetClassifierBase(t *testing.T) {
	cls := &mockClassifier{}
	newCls := &mockClassifier{}
	clsStack := &internal.MockStack{}
	clsStack.On("Set", newCls).Return(cls)
	obj := &state{
		cls: clsStack,
	}

	result := obj.SetClassifier(newCls)

	assert.Same(t, cls, result)
	clsStack.AssertExpectations(t)
}

func TestStateSetClassifierNil(t *testing.T) {
	newCls := &mockClassifier{}
	clsStack := &internal.MockStack{}
	clsStack.On("Set", newCls).Return(nil)
	obj := &state{
		cls: clsStack,
	}

	result := obj.SetClassifier(newCls)

	assert.Nil(t, result)
	clsStack.AssertExpectations(t)
}

func TestStatePush(t *testing.T) {
	tok := &Token{}
	obj := &state{
		toks: &list.List{},
	}

	result := obj.Push(tok)

	assert.True(t, result)
	assert.Equal(t, 1, obj.toks.Len())
	assert.Same(t, tok, obj.toks.Front().Value)
}
