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

	"github.com/klmitch/kent"
	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/common"
	"github.com/hydralang/ptk/utils/locations"
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
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 1, C: 1},
			E: locations.ArgPos{I: 1, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'n',
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 1, C: 2},
			E: locations.ArgPos{I: 1, C: 3},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'e',
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 1, C: 3},
			E: locations.ArgPos{I: 1, C: 4},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: common.EOF,
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 1, C: 4},
			E: locations.ArgPos{I: 1, C: 4},
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
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 1, C: 1},
			E: locations.ArgPos{I: 1, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: ' ',
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 0, C: 1},
			E: locations.ArgPos{I: 0, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'n',
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 2, C: 1},
			E: locations.ArgPos{I: 2, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: ' ',
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 0, C: 1},
			E: locations.ArgPos{I: 0, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 'e',
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 3, C: 1},
			E: locations.ArgPos{I: 3, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: common.EOF,
		Loc: locations.ArgLocation{
			B: locations.ArgPos{I: 3, C: 2},
			E: locations.ArgPos{I: 3, C: 2},
		},
	}, ch)
	assert.NotNil(t, opt1Called)
	assert.Same(t, opt1Called, opt2Called)
	assert.Equal(t, " ", opt1Called.joiner)
	assert.Len(t, opt1Called.opts, 1)
}
