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

func TestListScannerImplementsScanner(t *testing.T) {
	assert.Implements(t, (*Scanner)(nil), &ListScanner{})
}

func TestNewListScanner(t *testing.T) {
	chars := []Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}

	result := NewListScanner(chars, assert.AnError)

	assert.Equal(t, &ListScanner{
		chars: chars,
		err:   assert.AnError,
	}, result)
}

func TestListScannerNextFirst(t *testing.T) {
	chars := []Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}
	obj := &ListScanner{
		chars: chars,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 't'}, result)
	assert.Equal(t, &ListScanner{
		chars: chars,
		pos:   1,
	}, obj)
}

func TestListScannerNextLast(t *testing.T) {
	chars := []Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}
	obj := &ListScanner{
		chars: chars,
		pos:   3,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 't'}, result)
	assert.Equal(t, &ListScanner{
		chars: chars,
		pos:   4,
	}, obj)
}

func TestListScannerNextAgain(t *testing.T) {
	chars := []Char{
		{Rune: 't'},
		{Rune: 'e'},
		{Rune: 's'},
		{Rune: 't'},
	}
	obj := &ListScanner{
		chars: chars,
		pos:   4,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: 't'}, result)
	assert.Equal(t, &ListScanner{
		chars: chars,
		pos:   4,
	}, obj)
}

func TestListScannerNextEOF(t *testing.T) {
	chars := []Char{
		{Rune: EOF},
	}
	obj := &ListScanner{
		chars: chars,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: EOF}, result)
	assert.Equal(t, &ListScanner{
		chars: chars,
		pos:   1,
	}, obj)
}

func TestListScannerNextEOFWithError(t *testing.T) {
	chars := []Char{
		{Rune: EOF},
	}
	obj := &ListScanner{
		chars: chars,
		err:   assert.AnError,
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, Char{Rune: EOF}, result)
	assert.Equal(t, &ListScanner{
		chars: chars,
		pos:   1,
	}, obj)
}
