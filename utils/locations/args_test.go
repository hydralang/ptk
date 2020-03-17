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

package locations

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/common"
)

func TestArgLocationImplementsLocation(t *testing.T) {
	assert.Implements(t, (*common.Location)(nil), &ArgLocation{})
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
	loc2 := &common.MockLocation{}

	result, err := loc1.Thru(loc2)

	assert.Same(t, common.ErrSplitLocation, err)
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
	loc2 := &common.MockLocation{}

	result, err := loc1.ThruEnd(loc2)

	assert.Same(t, common.ErrSplitLocation, err)
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

	result := loc.Incr(common.EOF, 8)

	assert.Equal(t, ArgLocation{
		B: ArgPos{1, 5},
		E: ArgPos{1, 5},
	}, result)
	assert.Equal(t, ArgLocation{
		B: ArgPos{1, 3},
		E: ArgPos{1, 5},
	}, loc)
}
