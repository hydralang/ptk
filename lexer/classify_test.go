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

func TestMockClassifierImplementsClassifier(t *testing.T) {
	assert.Implements(t, (*Classifier)(nil), &MockClassifier{})
}

func TestMockClassifierClassifyNil(t *testing.T) {
	s := &MockState{}
	str := &common.MockBackTracker{}
	obj := &MockClassifier{}
	obj.On("Classify", s, str).Return(nil)

	result := obj.Classify(s, str)

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockClassifierClassifyNotNil(t *testing.T) {
	recs := []Recognizer{&MockRecognizer{}, &MockRecognizer{}}
	s := &MockState{}
	str := &common.MockBackTracker{}
	obj := &MockClassifier{}
	obj.On("Classify", s, str).Return(recs)

	result := obj.Classify(s, str)

	assert.Equal(t, recs, result)
	obj.AssertExpectations(t)
}

func TestMockClassifierError(t *testing.T) {
	s := &MockState{}
	str := &common.MockBackTracker{}
	obj := &MockClassifier{}
	obj.On("Error", s, str)

	obj.Error(s, str)

	obj.AssertExpectations(t)
}
