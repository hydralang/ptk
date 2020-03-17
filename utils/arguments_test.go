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

	"github.com/klmitch/kent"
	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/common"
)

func TestArgumentJoiner(t *testing.T) {
	obj := &argOptions{
		joiner: " ",
	}

	opt := ArgumentJoiner("|")
	opt(obj)

	assert.Equal(t, &argOptions{
		joiner: "|",
	}, obj)
}

func TestArgumentReporter(t *testing.T) {
	obj := &argOptions{
		opts: []ScannerOption{},
	}

	opt := ArgumentReporter(&kent.MockReporter{})
	opt(obj)

	assert.Len(t, obj.opts, 1)
}

func TestNewArgumentCharStreamEmptyList(t *testing.T) {
	var opt1Called *argOptions
	var opt2Called *argOptions
	options := []ArgumentOption{
		func(o *argOptions) {
			opt1Called = o
		},
		func(o *argOptions) {
			opt2Called = o
		},
	}

	result := NewArgumentCharStream([]string{}, options...)

	assert.NotNil(t, result)
	ch, err := result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{Rune: common.EOF}, ch)
	assert.NotNil(t, opt1Called)
	assert.Same(t, opt1Called, opt2Called)
	assert.Equal(t, " ", opt1Called.joiner)
	assert.Len(t, opt1Called.opts, 1)
}

func TestNewArgumentCharStream1ElemList(t *testing.T) {
	var opt1Called *argOptions
	var opt2Called *argOptions
	options := []ArgumentOption{
		func(o *argOptions) {
			opt1Called = o
		},
		func(o *argOptions) {
			opt2Called = o
		},
	}

	result := NewArgumentCharStream([]string{"one"}, options...)

	assert.NotNil(t, result)
	ch, err := result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'o',
		Loc: ArgLocation{
			B: ArgPos{1, 1},
			E: ArgPos{1, 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'n',
		Loc: ArgLocation{
			B: ArgPos{1, 2},
			E: ArgPos{1, 3},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'e',
		Loc: ArgLocation{
			B: ArgPos{1, 3},
			E: ArgPos{1, 4},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: common.EOF,
		Loc: ArgLocation{
			B: ArgPos{1, 4},
			E: ArgPos{1, 4},
		},
	}, ch)
	assert.NotNil(t, opt1Called)
	assert.Same(t, opt1Called, opt2Called)
	assert.Equal(t, " ", opt1Called.joiner)
	assert.Len(t, opt1Called.opts, 1)
}

func TestNewArgumentCharStream3ElemList(t *testing.T) {
	var opt1Called *argOptions
	var opt2Called *argOptions
	options := []ArgumentOption{
		func(o *argOptions) {
			opt1Called = o
		},
		func(o *argOptions) {
			opt2Called = o
		},
	}

	result := NewArgumentCharStream([]string{"o", "n", "e"}, options...)

	assert.NotNil(t, result)
	ch, err := result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'o',
		Loc: ArgLocation{
			B: ArgPos{1, 1},
			E: ArgPos{1, 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: ' ',
		Loc: ArgLocation{
			B: ArgPos{0, 1},
			E: ArgPos{0, 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'n',
		Loc: ArgLocation{
			B: ArgPos{2, 1},
			E: ArgPos{2, 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: ' ',
		Loc: ArgLocation{
			B: ArgPos{0, 1},
			E: ArgPos{0, 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'e',
		Loc: ArgLocation{
			B: ArgPos{3, 1},
			E: ArgPos{3, 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: common.EOF,
		Loc: ArgLocation{
			B: ArgPos{3, 2},
			E: ArgPos{3, 2},
		},
	}, ch)
	assert.NotNil(t, opt1Called)
	assert.Same(t, opt1Called, opt2Called)
	assert.Equal(t, " ", opt1Called.joiner)
	assert.Len(t, opt1Called.opts, 1)
}
