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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationErrorImplementsError(t *testing.T) {
	assert.Implements(t, (*error)(nil), &locationError{})
}

func TestLocationErrorBase(t *testing.T) {
	loc := &MockLocation{}

	result := LocationError(loc, assert.AnError)

	assert.Equal(t, &locationError{
		loc: loc,
		err: assert.AnError,
	}, result)
}

func TestLocationErrorNoLocation(t *testing.T) {
	result := LocationError(nil, assert.AnError)

	assert.Same(t, assert.AnError, result)
}

func TestLocationErrorError(t *testing.T) {
	loc := &MockLocation{}
	loc.On("String").Return("location")
	obj := &locationError{
		loc: loc,
		err: assert.AnError,
	}

	result := obj.Error()

	assert.Equal(t, fmt.Sprintf("location: %s", assert.AnError), result)
	loc.AssertExpectations(t)
}

func TestLocationErrorUnwrap(t *testing.T) {
	obj := &locationError{
		err: assert.AnError,
	}

	result := obj.Unwrap()

	assert.Same(t, assert.AnError, result)
}

type errorWrapper struct {
	err error
}

func (e *errorWrapper) Error() string {
	return ""
}

func (e *errorWrapper) Unwrap() error {
	return e.err
}

func TestLocationOfBase(t *testing.T) {
	loc := &MockLocation{}
	err := &errorWrapper{
		err: &locationError{
			loc: loc,
		},
	}

	result := LocationOf(err)

	assert.Same(t, loc, result)
}

func TestLocationOfNoLocation(t *testing.T) {
	err := &errorWrapper{
		err: assert.AnError,
	}

	result := LocationOf(err)

	assert.Nil(t, result)
}
