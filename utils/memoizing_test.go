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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/common"
)

func TestMemoizingCharStreamImplementsCharStream(t *testing.T) {
	assert.Implements(t, (*common.CharStream)(nil), &memoizingCharStream{})
}

func TestNewMemoizingCharStream(t *testing.T) {
	src := &common.MockCharStream{}

	result := NewMemoizingCharStream(src)

	assert.Equal(t, &memoizingCharStream{
		chars: []common.Char{},
		src:   src,
	}, result)
}

func TestMemoizingCharStreamNextBase(t *testing.T) {
	src := &common.MockCharStream{}
	src.On("Next").Return(common.Char{Rune: 't'}, nil)
	obj := &memoizingCharStream{
		chars: []common.Char{},
		src:   src,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: 't'}, result)
	assert.Equal(t, []common.Char{
		{Rune: 't'},
	}, obj.chars)
	assert.Equal(t, 0, obj.idx)
	assert.False(t, obj.replay)
	src.AssertExpectations(t)
}

func TestMemoizingCharStreamNextEOF(t *testing.T) {
	src := &common.MockCharStream{}
	src.On("Next").Return(common.Char{Rune: common.EOF}, nil)
	obj := &memoizingCharStream{
		chars: []common.Char{},
		src:   src,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: common.EOF}, result)
	assert.Equal(t, []common.Char{
		{Rune: common.EOF},
	}, obj.chars)
	assert.Equal(t, 0, obj.idx)
	assert.True(t, obj.replay)
	src.AssertExpectations(t)
}

func TestMemoizingCharStreamNextError(t *testing.T) {
	src := &common.MockCharStream{}
	src.On("Next").Return(common.Char{Rune: common.EOF}, assert.AnError)
	obj := &memoizingCharStream{
		chars: []common.Char{},
		src:   src,
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, common.Char{Rune: common.EOF}, result)
	assert.Equal(t, []common.Char{
		{Rune: common.EOF},
	}, obj.chars)
	assert.Equal(t, 0, obj.idx)
	assert.True(t, obj.replay)
	src.AssertExpectations(t)
}

func TestMemoizingCharStreamReplay(t *testing.T) {
	src := &common.MockCharStream{}
	obj := &memoizingCharStream{
		chars: []common.Char{
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
	assert.Equal(t, common.Char{Rune: 't'}, result)
	assert.Equal(t, []common.Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}, obj.chars)
	assert.Equal(t, 1, obj.idx)
	assert.True(t, obj.replay)
	src.AssertExpectations(t)
}

func TestMemoizingCharStreamReplayWrap(t *testing.T) {
	src := &common.MockCharStream{}
	obj := &memoizingCharStream{
		chars: []common.Char{
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
	assert.Equal(t, common.Char{Rune: 't'}, result)
	assert.Equal(t, []common.Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}, obj.chars)
	assert.Equal(t, 0, obj.idx)
	assert.True(t, obj.replay)
	src.AssertExpectations(t)
}
