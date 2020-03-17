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

package charstreams

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/common"
)

func TestChainingCharStreamImplementsCharStream(t *testing.T) {
	assert.Implements(t, (*common.CharStream)(nil), &chainingCharStream{})
}

func TestNewChainingCharStream(t *testing.T) {
	streams := []common.CharStream{
		&common.MockCharStream{},
		&common.MockCharStream{},
		&common.MockCharStream{},
	}

	result := NewChainingCharStream(streams)

	assert.Equal(t, &chainingCharStream{
		streams: streams,
		last: common.Char{
			Rune: common.EOF,
		},
	}, result)
}

func TestChainingCharStreamNextBase(t *testing.T) {
	str1 := &common.MockCharStream{}
	str1.On("Next").Return(common.Char{Rune: 'c'}, nil)
	obj := &chainingCharStream{
		streams: []common.CharStream{str1},
		last: common.Char{
			Rune: common.EOF,
		},
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: 'c'}, result)
	assert.Equal(t, 0, obj.idx)
	assert.Equal(t, common.Char{Rune: common.EOF}, obj.last)
	str1.AssertExpectations(t)
}

func TestChainingCharStreamNextNextStream(t *testing.T) {
	loc := &common.MockLocation{}
	str1 := &common.MockCharStream{}
	str1.On("Next").Return(common.Char{Rune: common.EOF, Loc: loc}, nil)
	str2 := &common.MockCharStream{}
	str2.On("Next").Return(common.Char{Rune: 'c'}, nil)
	obj := &chainingCharStream{
		streams: []common.CharStream{str1, str2},
		last: common.Char{
			Rune: common.EOF,
		},
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: 'c'}, result)
	assert.Equal(t, 1, obj.idx)
	assert.Equal(t, common.Char{Rune: common.EOF, Loc: loc}, obj.last)
	str1.AssertExpectations(t)
	str2.AssertExpectations(t)
}

func TestChainingCharStreamNextLastStream(t *testing.T) {
	loc := &common.MockLocation{}
	str1 := &common.MockCharStream{}
	str1.On("Next").Return(common.Char{Rune: common.EOF, Loc: loc}, nil)
	obj := &chainingCharStream{
		streams: []common.CharStream{str1},
		last: common.Char{
			Rune: common.EOF,
		},
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: common.EOF, Loc: loc}, result)
	assert.Equal(t, 1, obj.idx)
	assert.Equal(t, common.Char{Rune: common.EOF, Loc: loc}, obj.last)
	str1.AssertExpectations(t)
}

func TestChainingCharStreamNextEmptyStreams(t *testing.T) {
	loc1 := &common.MockLocation{}
	loc1.On("dummy", 1)
	str1 := &common.MockCharStream{}
	str1.On("Next").Return(common.Char{Rune: common.EOF, Loc: loc1}, nil)
	loc2 := &common.MockLocation{}
	loc2.On("dummy", 2)
	str2 := &common.MockCharStream{}
	str2.On("Next").Return(common.Char{Rune: common.EOF, Loc: loc2}, nil)
	str3 := &common.MockCharStream{}
	str3.On("Next").Return(common.Char{Rune: 'c'}, nil)
	obj := &chainingCharStream{
		streams: []common.CharStream{str1, str2, str3},
		last: common.Char{
			Rune: common.EOF,
		},
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: 'c'}, result)
	assert.Equal(t, 2, obj.idx)
	assert.Equal(t, common.Char{Rune: common.EOF, Loc: loc2}, obj.last)
	str1.AssertExpectations(t)
	str2.AssertExpectations(t)
	str3.AssertExpectations(t)
}

func TestChainingCharStreamNextError(t *testing.T) {
	loc := &common.MockLocation{}
	str1 := &common.MockCharStream{}
	str1.On("Next").Return(common.Char{Rune: common.EOF, Loc: loc}, assert.AnError)
	str2 := &common.MockCharStream{}
	obj := &chainingCharStream{
		streams: []common.CharStream{str1, str2},
		last: common.Char{
			Rune: common.EOF,
		},
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, common.Char{Rune: common.EOF, Loc: loc}, result)
	assert.Equal(t, 0, obj.idx)
	assert.Equal(t, common.Char{Rune: common.EOF}, obj.last)
	str1.AssertExpectations(t)
	str2.AssertExpectations(t)
}
