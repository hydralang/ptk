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

func TestListCharStreamImplementsCharStream(t *testing.T) {
	assert.Implements(t, (*common.CharStream)(nil), &listCharStream{})
}

func TestNewListCharStream(t *testing.T) {
	chars := []common.Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}

	result := NewListCharStream(chars, assert.AnError)

	assert.Equal(t, &listCharStream{
		chars: chars,
		err:   assert.AnError,
	}, result)
}

func TestListCharStreamNextFirst(t *testing.T) {
	chars := []common.Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}
	obj := &listCharStream{
		chars: chars,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: 't'}, result)
	assert.Equal(t, &listCharStream{
		chars: chars,
		pos:   1,
	}, obj)
}

func TestListCharStreamNextLast(t *testing.T) {
	chars := []common.Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}
	obj := &listCharStream{
		chars: chars,
		pos:   3,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: 't'}, result)
	assert.Equal(t, &listCharStream{
		chars: chars,
		pos:   4,
	}, obj)
}

func TestListCharStreamNextAgain(t *testing.T) {
	chars := []common.Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}
	obj := &listCharStream{
		chars: chars,
		pos:   4,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: 't'}, result)
	assert.Equal(t, &listCharStream{
		chars: chars,
		pos:   4,
	}, obj)
}

func TestListCharStreamNextEOF(t *testing.T) {
	chars := []common.Char{
		{Rune: common.EOF},
	}
	obj := &listCharStream{
		chars: chars,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: common.EOF}, result)
	assert.Equal(t, &listCharStream{
		chars: chars,
		pos:   1,
	}, obj)
}

func TestListCharStreamNextEOFWithError(t *testing.T) {
	chars := []common.Char{
		{Rune: common.EOF},
	}
	obj := &listCharStream{
		chars: chars,
		err:   assert.AnError,
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, common.Char{Rune: common.EOF}, result)
	assert.Equal(t, &listCharStream{
		chars: chars,
		pos:   1,
	}, obj)
}
