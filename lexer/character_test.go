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

package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/common"
)

func TestMockCharStreamImplementsCharStream(t *testing.T) {
	assert.Implements(t, (*CharStream)(nil), &MockCharStream{})
}

func TestMockCharStreamNextNil(t *testing.T) {
	obj := &MockCharStream{}
	obj.On("Next").Return(nil, assert.AnError)

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, common.Char{}, result)
}

func TestMockCharStreamNextNotNil(t *testing.T) {
	obj := &MockCharStream{}
	obj.On("Next").Return(common.Char{Rune: 'c'}, assert.AnError)

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, common.Char{Rune: 'c'}, result)
}
