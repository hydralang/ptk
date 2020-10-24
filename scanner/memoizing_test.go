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

package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoizingScannerImplementsScanner(t *testing.T) {
	assert.Implements(t, (*Scanner)(nil), &MemoizingScanner{})
}

func TestNewMemoizingScanner(t *testing.T) {
	src := &mockScanner{}

	result := NewMemoizingScanner(src)

	assert.Equal(t, &MemoizingScanner{
		chars: []Char{},
		src:   src,
	}, result)
}

func TestMemoizingScannerNextBase(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(Char{Rune: 't'}, nil)
	obj := &MemoizingScanner{
		chars: []Char{},
		src:   src,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 't'}, result)
	assert.Equal(t, []Char{
		{Rune: 't'},
	}, obj.chars)
	assert.Equal(t, 0, obj.idx)
	assert.False(t, obj.replay)
	src.AssertExpectations(t)
}

func TestMemoizingScannerNextEOF(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(Char{Rune: EOF}, nil)
	obj := &MemoizingScanner{
		chars: []Char{},
		src:   src,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: EOF}, result)
	assert.Equal(t, []Char{
		{Rune: EOF},
	}, obj.chars)
	assert.Equal(t, 0, obj.idx)
	assert.True(t, obj.replay)
	src.AssertExpectations(t)
}

func TestMemoizingScannerNextError(t *testing.T) {
	src := &mockScanner{}
	src.On("Next").Return(Char{Rune: EOF}, assert.AnError)
	obj := &MemoizingScanner{
		chars: []Char{},
		src:   src,
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, Char{Rune: EOF}, result)
	assert.Equal(t, []Char{
		{Rune: EOF},
	}, obj.chars)
	assert.Equal(t, 0, obj.idx)
	assert.True(t, obj.replay)
	src.AssertExpectations(t)
}

func TestMemoizingScannerReplay(t *testing.T) {
	src := &mockScanner{}
	obj := &MemoizingScanner{
		chars: []Char{
			{Rune: 't'},
			{Rune: 'e'},
			{Rune: 's'},
			{Rune: 't'},
		},
		src:    src,
		replay: true,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 't'}, result)
	assert.Equal(t, []Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}, obj.chars)
	assert.Equal(t, 1, obj.idx)
	assert.True(t, obj.replay)
	src.AssertExpectations(t)
}

func TestMemoizingScannerReplayWrap(t *testing.T) {
	src := &mockScanner{}
	obj := &MemoizingScanner{
		chars: []Char{
			{Rune: 't'},
			{Rune: 'e'},
			{Rune: 's'},
			{Rune: 't'},
		},
		idx:    3,
		src:    src,
		replay: true,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 't'}, result)
	assert.Equal(t, []Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}, obj.chars)
	assert.Equal(t, 0, obj.idx)
	assert.True(t, obj.replay)
	src.AssertExpectations(t)
}
