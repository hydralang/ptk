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
	"github.com/stretchr/testify/mock"
)

func TestArgLocationImplementsLocation(t *testing.T) {
	assert.Implements(t, (*Location)(nil), &ArgLocation{})
}

func TestArgLocationString0Columns(t *testing.T) {
	loc := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 3},
	}

	result := loc.String()

	assert.Equal(t, "1:3", result)
}

func TestArgLocationString1Columns(t *testing.T) {
	loc := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 4},
	}

	result := loc.String()

	assert.Equal(t, "1:3", result)
}

func TestArgLocationString2Columns(t *testing.T) {
	loc := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}

	result := loc.String()

	assert.Equal(t, "1:3-5", result)
}

func TestArgLocationString2Arguments(t *testing.T) {
	loc := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{2, 1},
	}

	result := loc.String()

	assert.Equal(t, "1:3-2:1", result)
}

func TestArgLocationThruBase(t *testing.T) {
	loc1 := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}
	loc2 := ArgLocation{
		B: ArgPos{2, 2},
		E: ArgPos{2, 4},
	}

	result, err := loc1.Thru(loc2)

	assert.NoError(t, err)
	assert.Equal(t, ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{2, 2},
	}, result)
}

func TestArgLocationThruSplitLocation(t *testing.T) {
	loc1 := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}
	loc2 := &mockLocation{}

	result, err := loc1.Thru(loc2)

	assert.Same(t, ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestArgLocationThruEndBase(t *testing.T) {
	loc1 := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}
	loc2 := ArgLocation{
		B: ArgPos{2, 2},
		E: ArgPos{2, 4},
	}

	result, err := loc1.ThruEnd(loc2)

	assert.NoError(t, err)
	assert.Equal(t, ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{2, 4},
	}, result)
}

func TestArgLocationThruEndSplitLocation(t *testing.T) {
	loc1 := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}
	loc2 := &mockLocation{}

	result, err := loc1.ThruEnd(loc2)

	assert.Same(t, ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestArgLocationIncrBase(t *testing.T) {
	loc := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}

	result := loc.Incr('t', 8)

	assert.Equal(t, ArgLocation{
		B: ArgPos{1, 5},
		E: ArgPos{1, 6},
	}, result)
	assert.Equal(t, ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}, loc)
}

func TestArgLocationIncrEOF(t *testing.T) {
	loc := ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}

	result := loc.Incr(EOF, 8)

	assert.Equal(t, ArgLocation{
		B: ArgPos{1, 5},
		E: ArgPos{1, 5},
	}, result)
	assert.Equal(t, ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}, loc)
}

func TestNewArgumentScannerEmptyList(t *testing.T) {
	opt1 := &mockArgOption{}
	opt1.On("argApply", mock.Anything)
	opt2 := &mockArgOption{}
	opt2.On("argApply", mock.Anything)

	result := NewArgumentScanner([]string{}, opt1, opt2)

	assert.NotNil(t, result)
	ch, err := result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{Rune: EOF}, ch)
	opt1.AssertExpectations(t)
	opt2.AssertExpectations(t)
}

func TestNewArgumentScanner1ElemList(t *testing.T) {
	opt1 := &mockArgOption{}
	opt1.On("argApply", mock.Anything)
	opt2 := &mockArgOption{}
	opt2.On("argApply", mock.Anything)

	result := NewArgumentScanner([]string{"one"}, opt1, opt2)

	assert.NotNil(t, result)
	ch, err := result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: 'o',
		Loc: ArgLocation{
			B: ArgPos{I: 1, C: 1},
			E: ArgPos{I: 1, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: 'n',
		Loc: ArgLocation{
			B: ArgPos{I: 1, C: 2},
			E: ArgPos{I: 1, C: 3},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: 'e',
		Loc: ArgLocation{
			B: ArgPos{I: 1, C: 3},
			E: ArgPos{I: 1, C: 4},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: EOF,
		Loc: ArgLocation{
			B: ArgPos{I: 1, C: 4},
			E: ArgPos{I: 1, C: 4},
		},
	}, ch)
	opt1.AssertExpectations(t)
	opt2.AssertExpectations(t)
}

func TestNewArgumentScanner3ElemList(t *testing.T) {
	opt1 := &mockArgOption{}
	opt1.On("argApply", mock.Anything)
	opt2 := &mockArgOption{}
	opt2.On("argApply", mock.Anything)

	result := NewArgumentScanner([]string{"o", "n", "e"}, opt1, opt2)

	assert.NotNil(t, result)
	ch, err := result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: 'o',
		Loc: ArgLocation{
			B: ArgPos{I: 1, C: 1},
			E: ArgPos{I: 1, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: ' ',
		Loc: ArgLocation{
			B: ArgPos{I: 0, C: 1},
			E: ArgPos{I: 0, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: 'n',
		Loc: ArgLocation{
			B: ArgPos{I: 2, C: 1},
			E: ArgPos{I: 2, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: ' ',
		Loc: ArgLocation{
			B: ArgPos{I: 0, C: 1},
			E: ArgPos{I: 0, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: 'e',
		Loc: ArgLocation{
			B: ArgPos{I: 3, C: 1},
			E: ArgPos{I: 3, C: 2},
		},
	}, ch)
	ch, err = result.Next()
	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: EOF,
		Loc: ArgLocation{
			B: ArgPos{I: 3, C: 2},
			E: ArgPos{I: 3, C: 2},
		},
	}, ch)
	opt1.AssertExpectations(t)
	opt2.AssertExpectations(t)
}
