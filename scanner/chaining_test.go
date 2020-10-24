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

func TestChainingScannerImplementsScanner(t *testing.T) {
	assert.Implements(t, (*Scanner)(nil), &ChainingScanner{})
}

func TestNewChainingScanner(t *testing.T) {
	streams := []Scanner{
		&mockScanner{},
		&mockScanner{},
		&mockScanner{},
	}

	result := NewChainingScanner(streams)

	assert.Equal(t, &ChainingScanner{
		streams: streams,
		last: Char{
			Rune: EOF,
		},
	}, result)
}

func TestChainingScannerNextBase(t *testing.T) {
	str1 := &mockScanner{}
	str1.On("Next").Return(Char{Rune: 'c'}, nil)
	obj := &ChainingScanner{
		streams: []Scanner{str1},
		last: Char{
			Rune: EOF,
		},
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 'c'}, result)
	assert.Equal(t, 0, obj.idx)
	assert.Equal(t, Char{Rune: EOF}, obj.last)
	str1.AssertExpectations(t)
}

func TestChainingScannerNextNextStream(t *testing.T) {
	loc := &mockLocation{}
	str1 := &mockScanner{}
	str1.On("Next").Return(Char{Rune: EOF, Loc: loc}, nil)
	str2 := &mockScanner{}
	str2.On("Next").Return(Char{Rune: 'c'}, nil)
	obj := &ChainingScanner{
		streams: []Scanner{str1, str2},
		last: Char{
			Rune: EOF,
		},
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 'c'}, result)
	assert.Equal(t, 1, obj.idx)
	assert.Equal(t, Char{Rune: EOF, Loc: loc}, obj.last)
	str1.AssertExpectations(t)
	str2.AssertExpectations(t)
}

func TestChainingScannerNextLastStream(t *testing.T) {
	loc := &mockLocation{}
	str1 := &mockScanner{}
	str1.On("Next").Return(Char{Rune: EOF, Loc: loc}, nil)
	obj := &ChainingScanner{
		streams: []Scanner{str1},
		last: Char{
			Rune: EOF,
		},
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: EOF, Loc: loc}, result)
	assert.Equal(t, 1, obj.idx)
	assert.Equal(t, Char{Rune: EOF, Loc: loc}, obj.last)
	str1.AssertExpectations(t)
}

func TestChainingScannerNextEmptyStreams(t *testing.T) {
	loc1 := &mockLocation{}
	loc1.On("dummy", 1)
	str1 := &mockScanner{}
	str1.On("Next").Return(Char{Rune: EOF, Loc: loc1}, nil)
	loc2 := &mockLocation{}
	loc2.On("dummy", 2)
	str2 := &mockScanner{}
	str2.On("Next").Return(Char{Rune: EOF, Loc: loc2}, nil)
	str3 := &mockScanner{}
	str3.On("Next").Return(Char{Rune: 'c'}, nil)
	obj := &ChainingScanner{
		streams: []Scanner{str1, str2, str3},
		last: Char{
			Rune: EOF,
		},
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 'c'}, result)
	assert.Equal(t, 2, obj.idx)
	assert.Equal(t, Char{Rune: EOF, Loc: loc2}, obj.last)
	str1.AssertExpectations(t)
	str2.AssertExpectations(t)
	str3.AssertExpectations(t)
}

func TestChainingScannerNextError(t *testing.T) {
	loc := &mockLocation{}
	str1 := &mockScanner{}
	str1.On("Next").Return(Char{Rune: EOF, Loc: loc}, assert.AnError)
	str2 := &mockScanner{}
	obj := &ChainingScanner{
		streams: []Scanner{str1, str2},
		last: Char{
			Rune: EOF,
		},
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, Char{Rune: EOF, Loc: loc}, result)
	assert.Equal(t, 0, obj.idx)
	assert.Equal(t, Char{Rune: EOF}, obj.last)
	str1.AssertExpectations(t)
	str2.AssertExpectations(t)
}
