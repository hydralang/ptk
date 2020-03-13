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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockLocationImplementsLocation(t *testing.T) {
	assert.Implements(t, (*Location)(nil), &MockLocation{})
}

func TestMockLocationString(t *testing.T) {
	obj := &MockLocation{}
	obj.On("String").Return("location")

	result := obj.String()

	assert.Equal(t, "location", result)
	obj.AssertExpectations(t)
}

func TestMockLocationThruNil(t *testing.T) {
	other := &MockLocation{}
	obj := &MockLocation{}
	obj.On("Thru", other).Return(nil, assert.AnError)

	result, err := obj.Thru(other)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockLocationThruNotNil(t *testing.T) {
	expected := &MockLocation{}
	other := &MockLocation{}
	obj := &MockLocation{}
	obj.On("Thru", other).Return(expected, assert.AnError)

	result, err := obj.Thru(other)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockLocationThruEndNil(t *testing.T) {
	other := &MockLocation{}
	obj := &MockLocation{}
	obj.On("ThruEnd", other).Return(nil, assert.AnError)

	result, err := obj.ThruEnd(other)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockLocationThruEndNotNil(t *testing.T) {
	expected := &MockLocation{}
	other := &MockLocation{}
	obj := &MockLocation{}
	obj.On("ThruEnd", other).Return(expected, assert.AnError)

	result, err := obj.ThruEnd(other)

	assert.Same(t, assert.AnError, err)
	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}
