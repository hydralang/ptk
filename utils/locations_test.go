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

func TestFileLocationImplementsLocation(t *testing.T) {
	assert.Implements(t, (*common.Location)(nil), &FileLocation{})
}

func TestFileLocationAdvanceColumn(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}

	loc.Advance(FilePos{C: 2})

	assert.Equal(t, FilePos{L: 3, C: 3}, loc.B)
	assert.Equal(t, FilePos{L: 3, C: 5}, loc.E)
}

func TestFileLocationAdvanceLine(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}

	loc.Advance(FilePos{L: 1, C: 2})

	assert.Equal(t, FilePos{L: 3, C: 3}, loc.B)
	assert.Equal(t, FilePos{L: 4, C: 3}, loc.E)
}

func TestFileLocationAdvanceTab8(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}

	loc.AdvanceTab(8)

	assert.Equal(t, FilePos{L: 3, C: 3}, loc.B)
	assert.Equal(t, FilePos{L: 3, C: 9}, loc.E)
}

func TestFileLocationAdvanceTab4(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B: FilePos{
			L: 3,
			C: 2,
		},
		E: FilePos{
			L: 3,
			C: 3,
		},
	}

	loc.AdvanceTab(4)

	assert.Equal(t, FilePos{L: 3, C: 3}, loc.B)
	assert.Equal(t, FilePos{L: 3, C: 5}, loc.E)
}

func TestFileLocationString0Columns(t *testing.T) {
	loc := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 2}}

	result := loc.String()

	assert.Equal(t, "file:3:2", result)
}

func TestFileLocationString1Column(t *testing.T) {
	loc := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}

	result := loc.String()

	assert.Equal(t, "file:3:2", result)
}

func TestFileLocationString2Columns(t *testing.T) {
	loc := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 4}}

	result := loc.String()

	assert.Equal(t, "file:3:2-4", result)
}

func TestFileLocationString2Lines(t *testing.T) {
	loc := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{4, 2}}

	result := loc.String()

	assert.Equal(t, "file:3:2-4:2", result)
}

func TestFileLocationThruBase(t *testing.T) {
	loc1 := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := FileLocation{File: "file", B: FilePos{3, 5}, E: FilePos{3, 6}}

	result, err := loc1.Thru(loc2)

	assert.NoError(t, err)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 5},
	}, result)
}

func TestFileLocationThruSplitFile(t *testing.T) {
	loc1 := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := FileLocation{File: "other", B: FilePos{3, 5}, E: FilePos{3, 6}}

	result, err := loc1.Thru(loc2)

	assert.Same(t, common.ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestFileLocationThruSplitLocation(t *testing.T) {
	loc1 := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := &common.MockLocation{}

	result, err := loc1.Thru(loc2)

	assert.Same(t, common.ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestFileLocationThruEndBase(t *testing.T) {
	loc1 := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := FileLocation{File: "file", B: FilePos{3, 5}, E: FilePos{3, 6}}

	result, err := loc1.ThruEnd(loc2)

	assert.NoError(t, err)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 6},
	}, result)
}

func TestFileLocationThruEndSplitFile(t *testing.T) {
	loc1 := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := FileLocation{File: "other", B: FilePos{3, 5}, E: FilePos{3, 6}}

	result, err := loc1.ThruEnd(loc2)

	assert.Same(t, common.ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestFileLocationThruEndSplitLocation(t *testing.T) {
	loc1 := FileLocation{File: "file", B: FilePos{3, 2}, E: FilePos{3, 3}}
	loc2 := &common.MockLocation{}

	result, err := loc1.ThruEnd(loc2)

	assert.Same(t, common.ErrSplitLocation, err)
	assert.Nil(t, result)
}
