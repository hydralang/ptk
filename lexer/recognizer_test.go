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
)

func TestMockRecognizerImplementsRecognizer(t *testing.T) {
	assert.Implements(t, (*Recognizer)(nil), &MockRecognizer{})
}

func TestMockRecognizerRecognize(t *testing.T) {
	s := &MockState{}
	str := &BackTracker{}
	obj := &MockRecognizer{}
	obj.On("Recognize", s, str).Return(true)

	result := obj.Recognize(s, str)

	assert.True(t, result)
	obj.AssertExpectations(t)
}
