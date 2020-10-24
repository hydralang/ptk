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

	"github.com/hydralang/ptk/scanner"
)

type mockScanner struct {
	mock.Mock
}

func (m *mockScanner) Next() (scanner.Char, error) {
	args := m.MethodCalled("Next")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(scanner.Char), args.Error(1)
	}

	return scanner.Char{}, args.Error(1)
}

type mockBackTracker struct {
	mockScanner
}

func (m *mockBackTracker) More() bool {
	args := m.MethodCalled("More")

	return args.Bool(0)
}

func (m *mockBackTracker) SetMax(max int) {
	m.MethodCalled("SetMax", max)
}

func (m *mockBackTracker) Accept(leave int) {
	m.MethodCalled("Accept", leave)
}

func (m *mockBackTracker) Len() int {
	args := m.MethodCalled("Len")

	return args.Int(0)
}

func (m *mockBackTracker) Pos() int {
	args := m.MethodCalled("Pos")

	return args.Int(0)
}

func (m *mockBackTracker) BackTrack() {
	m.MethodCalled("BackTrack")
}

func TestBackTrackerImplementsBackTracker(t *testing.T) {
	assert.Implements(t, (*BackTracker)(nil), &backTracker{})
}

func TestNewBackTracker(t *testing.T) {
	src := &mockScanner{}

	result := NewBackTracker(src, 42)

	assert.Equal(t, &backTracker{
		src:   src,
		max:   42,
		saved: &list.List{},
		last: btElem{
			ch: scanner.Char{Rune: scanner.EOF},
		},
	}, result)
}

func TestBackTrackerNextBase(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(scanner.Char{Rune: 't'}, assert.AnError)
	obj := &backTracker{
		src:   src,
		max:   TrackAll,
		saved: &list.List{},
		last: btElem{
			ch: scanner.Char{Rune: 'b'},
		},
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, scanner.Char{Rune: 't'}, result)
	assert.Same(t, src, obj.src)
	assert.Equal(t, 1, obj.saved.Len())
	assert.Equal(t, btElem{
		ch:  scanner.Char{Rune: 't'},
		err: assert.AnError,
	}, obj.saved.Back().Value)
	assert.Nil(t, obj.next)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, btElem{
		ch: scanner.Char{Rune: 'b'},
	}, obj.last)
}

func TestBackTrackerNextTrackNone(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(scanner.Char{Rune: 't'}, assert.AnError)
	obj := &backTracker{
		src:   src,
		max:   0,
		saved: &list.List{},
		last: btElem{
			ch: scanner.Char{Rune: 'b'},
		},
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, scanner.Char{Rune: 't'}, result)
	assert.Same(t, src, obj.src)
	assert.Equal(t, 0, obj.saved.Len())
	assert.Nil(t, obj.next)
	assert.Equal(t, 0, obj.pos)
	assert.Equal(t, btElem{
		ch: scanner.Char{Rune: 'b'},
	}, obj.last)
}

func TestBackTrackerNextNoTrim(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(scanner.Char{Rune: 't'}, assert.AnError)
	obj := &backTracker{
		src:   src,
		max:   4,
		saved: &list.List{},
		last: btElem{
			ch: scanner.Char{Rune: 'b'},
		},
		pos: 3,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, scanner.Char{Rune: 't'}, result)
	assert.Same(t, src, obj.src)
	assert.Equal(t, 4, obj.saved.Len())
	assert.Equal(t, btElem{
		ch:  scanner.Char{Rune: 't'},
		err: assert.AnError,
	}, obj.saved.Back().Value)
	assert.Nil(t, obj.next)
	assert.Equal(t, 4, obj.pos)
	assert.Equal(t, btElem{
		ch: scanner.Char{Rune: 'b'},
	}, obj.last)
}

func TestBackTrackerNextWithTrim(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(scanner.Char{Rune: 't'}, assert.AnError)
	obj := &backTracker{
		src:   src,
		max:   3,
		saved: &list.List{},
		last: btElem{
			ch: scanner.Char{Rune: 'b'},
		},
		pos: 3,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, scanner.Char{Rune: 't'}, result)
	assert.Same(t, src, obj.src)
	assert.Equal(t, 3, obj.saved.Len())
	assert.Equal(t, btElem{
		ch:  scanner.Char{Rune: 't'},
		err: assert.AnError,
	}, obj.saved.Back().Value)
	assert.Nil(t, obj.next)
	assert.Equal(t, 3, obj.pos)
	assert.Equal(t, btElem{
		ch: scanner.Char{Rune: 'b'},
	}, obj.last)
}

func TestBackTrackerNextSaveEOF(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(scanner.Char{Rune: scanner.EOF}, assert.AnError)
	obj := &backTracker{
		src:   src,
		max:   TrackAll,
		saved: &list.List{},
		last: btElem{
			ch: scanner.Char{Rune: 'b'},
		},
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, scanner.Char{Rune: scanner.EOF}, result)
	assert.Nil(t, obj.src)
	assert.Equal(t, 1, obj.saved.Len())
	assert.Equal(t, btElem{
		ch:  scanner.Char{Rune: scanner.EOF},
		err: assert.AnError,
	}, obj.saved.Back().Value)
	assert.Nil(t, obj.next)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, btElem{
		ch: scanner.Char{Rune: scanner.EOF},
	}, obj.last)
}

func TestBackTrackerNextBackTracked(t *testing.T) {
	src := &mockScanner{}
	obj := &backTracker{
		src:   src,
		max:   TrackAll,
		saved: &list.List{},
		last: btElem{
			ch: scanner.Char{Rune: 'b'},
		},
		pos: 0,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}, err: assert.AnError})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.next = obj.saved.Front()

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, scanner.Char{Rune: 't'}, result)
	assert.Same(t, src, obj.src)
	assert.Equal(t, 4, obj.saved.Len())
	assert.Equal(t, btElem{
		ch: scanner.Char{Rune: 't'},
	}, obj.saved.Back().Value)
	assert.Same(t, obj.saved.Front().Next(), obj.next)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, btElem{
		ch: scanner.Char{Rune: 'b'},
	}, obj.last)
}

func TestBackTrackerNextExtension(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(scanner.Char{Rune: 't'}, assert.AnError)
	obj := &backTracker{
		max:   TrackAll,
		saved: &list.List{},
		last: btElem{
			ch: scanner.Char{Rune: 'b'},
		},
		pos: 42,
	}

	result, err := obj.Next()

	assert.Nil(t, err)
	assert.Equal(t, scanner.Char{Rune: 'b'}, result)
	assert.Nil(t, obj.src)
	assert.Equal(t, 0, obj.saved.Len())
	assert.Nil(t, obj.next)
	assert.Equal(t, 42, obj.pos)
	assert.Equal(t, btElem{
		ch: scanner.Char{Rune: 'b'},
	}, obj.last)
}

func TestBackTrackerMoreBackTracked(t *testing.T) {
	obj := &backTracker{
		saved: &list.List{},
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.next = obj.saved.Front()

	result := obj.More()

	assert.True(t, result)
}

func TestBackTrackerMoreHaveSrc(t *testing.T) {
	src := &mockScanner{}
	obj := &backTracker{
		src: src,
	}

	result := obj.More()

	assert.True(t, result)
}

func TestBackTrackerMoreNoMore(t *testing.T) {
	obj := &backTracker{}

	result := obj.More()

	assert.False(t, result)
}

func TestBackTrackerSetMaxBase(t *testing.T) {
	obj := &backTracker{
		max:   42,
		saved: &list.List{},
		pos:   4,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	obj.SetMax(TrackAll)

	assert.Equal(t, TrackAll, obj.max)
	assert.Equal(t, 4, obj.saved.Len())
	assert.Equal(t, 4, obj.pos)
}

func TestBackTrackerSetMax0(t *testing.T) {
	obj := &backTracker{
		max:   42,
		saved: &list.List{},
		pos:   4,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	obj.SetMax(0)

	assert.Equal(t, 0, obj.max)
	assert.Equal(t, 0, obj.saved.Len())
	assert.Equal(t, 0, obj.pos)
}

func TestBackTrackerSetMaxIncrease(t *testing.T) {
	obj := &backTracker{
		max:   3,
		saved: &list.List{},
		pos:   4,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	obj.SetMax(4)

	assert.Equal(t, 4, obj.max)
	assert.Equal(t, 4, obj.saved.Len())
	assert.Equal(t, 4, obj.pos)
}

func TestBackTrackerSetMaxDecrease(t *testing.T) {
	obj := &backTracker{
		max:   4,
		saved: &list.List{},
		pos:   4,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	obj.SetMax(2)

	assert.Equal(t, 2, obj.max)
	assert.Equal(t, 2, obj.saved.Len())
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 's'}}, obj.saved.Front().Value)
	assert.Equal(t, 2, obj.pos)
}

func TestBackTrackerAcceptUnsaved(t *testing.T) {
	obj := &backTracker{
		max:   0,
		saved: &list.List{},
		pos:   0,
	}

	obj.Accept(0)

	assert.Equal(t, 0, obj.saved.Len())
	assert.Equal(t, 0, obj.pos)
}

func TestBackTrackerAccept0Current(t *testing.T) {
	obj := &backTracker{
		max:   TrackAll,
		saved: &list.List{},
		pos:   4,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	obj.Accept(0)

	assert.Equal(t, 0, obj.saved.Len())
	assert.Equal(t, 0, obj.pos)
}

func TestBackTrackerAccept2Current(t *testing.T) {
	obj := &backTracker{
		max:   TrackAll,
		saved: &list.List{},
		pos:   4,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	obj.Accept(2)

	assert.Equal(t, 2, obj.saved.Len())
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 's'}}, obj.saved.Front().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 't'}}, obj.saved.Front().Next().Value)
	assert.Equal(t, 2, obj.pos)
}

func TestBackTrackerAccept10Current(t *testing.T) {
	obj := &backTracker{
		max:   TrackAll,
		saved: &list.List{},
		pos:   4,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	obj.Accept(10)

	assert.Equal(t, 4, obj.saved.Len())
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 't'}}, obj.saved.Front().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 'e'}}, obj.saved.Front().Next().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 's'}}, obj.saved.Front().Next().Next().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 't'}}, obj.saved.Front().Next().Next().Next().Value)
	assert.Equal(t, 4, obj.pos)
}

func TestBackTrackerAccept0Point(t *testing.T) {
	obj := &backTracker{
		max:   TrackAll,
		saved: &list.List{},
		pos:   3,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.next = obj.saved.Back()

	obj.Accept(0)

	assert.Equal(t, 1, obj.saved.Len())
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 't'}}, obj.saved.Front().Value)
	assert.Equal(t, 0, obj.pos)
}

func TestBackTrackerAccept2Point(t *testing.T) {
	obj := &backTracker{
		max:   TrackAll,
		saved: &list.List{},
		pos:   3,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.next = obj.saved.Back()

	obj.Accept(2)

	assert.Equal(t, 3, obj.saved.Len())
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 'e'}}, obj.saved.Front().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 's'}}, obj.saved.Front().Next().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 't'}}, obj.saved.Front().Next().Next().Value)
	assert.Equal(t, 2, obj.pos)
}

func TestBackTrackerAccept10Point(t *testing.T) {
	obj := &backTracker{
		max:   TrackAll,
		saved: &list.List{},
		pos:   3,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.next = obj.saved.Back()

	obj.Accept(10)

	assert.Equal(t, 4, obj.saved.Len())
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 't'}}, obj.saved.Front().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 'e'}}, obj.saved.Front().Next().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 's'}}, obj.saved.Front().Next().Next().Value)
	assert.Equal(t, btElem{ch: scanner.Char{Rune: 't'}}, obj.saved.Front().Next().Next().Next().Value)
	assert.Equal(t, 3, obj.pos)
}

func TestBackTrackerLen(t *testing.T) {
	obj := &backTracker{
		saved: &list.List{},
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	result := obj.Len()

	assert.Equal(t, 4, result)
}

func TestBackTrackerPos(t *testing.T) {
	obj := &backTracker{
		pos: 42,
	}

	result := obj.Pos()

	assert.Equal(t, 41, result)
}

func TestBackTrackerBackTrack(t *testing.T) {
	obj := &backTracker{
		saved: &list.List{},
		pos:   4,
	}
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 'e'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 's'}})
	obj.saved.PushBack(btElem{ch: scanner.Char{Rune: 't'}})

	obj.BackTrack()

	assert.Same(t, obj.saved.Front(), obj.next)
	assert.Equal(t, 0, obj.pos)
}
