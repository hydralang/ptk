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

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockStackImplementsStack(t *testing.T) {
	assert.Implements(t, (*Stack)(nil), &MockStack{})
}

func TestMockStackLen(t *testing.T) {
	obj := &MockStack{}
	obj.On("Len").Return(42)

	result := obj.Len()

	assert.Equal(t, 42, result)
	obj.AssertExpectations(t)
}

func TestMockStackGet(t *testing.T) {
	obj := &MockStack{}
	obj.On("Get").Return("item")

	result := obj.Get()

	assert.Equal(t, "item", result)
	obj.AssertExpectations(t)
}

func TestMockStackPush(t *testing.T) {
	obj := &MockStack{}
	obj.On("Push", "item")

	obj.Push("item")

	obj.AssertExpectations(t)
}

func TestMockStackPop(t *testing.T) {
	obj := &MockStack{}
	obj.On("Pop").Return("item")

	result := obj.Pop()

	assert.Equal(t, "item", result)
	obj.AssertExpectations(t)
}

func TestMockStackSet(t *testing.T) {
	obj := &MockStack{}
	obj.On("Set", "new").Return("old")

	result := obj.Set("new")

	assert.Equal(t, "old", result)
	obj.AssertExpectations(t)
}

func TestStackImplementsStack(t *testing.T) {
	assert.Implements(t, (*Stack)(nil), &stack{})
}

func TestNewStack(t *testing.T) {
	result := NewStack()

	assert.Equal(t, &stack{}, result)
}

func TestStackGetBase(t *testing.T) {
	s := &stack{}
	s.PushFront("item")

	result := s.Get()

	assert.Equal(t, "item", result)
	assert.Equal(t, 1, s.Len())
}

func TestStackGetEmpty(t *testing.T) {
	s := &stack{}

	result := s.Get()

	assert.Nil(t, result)
	assert.Equal(t, 0, s.Len())
}

func TestStackPush(t *testing.T) {
	s := &stack{}

	s.Push("item")

	require.Equal(t, 1, s.Len())
	assert.Equal(t, "item", s.Front().Value)
}

func TestStackPopBase(t *testing.T) {
	s := &stack{}
	s.PushFront("item")

	result := s.Pop()

	assert.Equal(t, "item", result)
	assert.Equal(t, 0, s.Len())
}

func TestStackPopEmpty(t *testing.T) {
	s := &stack{}

	result := s.Pop()

	assert.Nil(t, result)
	assert.Equal(t, 0, s.Len())
}

func TestStackSetBase(t *testing.T) {
	s := &stack{}
	s.PushFront("item")

	result := s.Set("other")

	assert.Equal(t, "item", result)
	require.Equal(t, 1, s.Len())
	assert.Equal(t, "other", s.Front().Value)
}

func TestStackSetEmpty(t *testing.T) {
	s := &stack{}

	result := s.Set("other")

	assert.Nil(t, result)
	require.Equal(t, 1, s.Len())
	assert.Equal(t, "other", s.Front().Value)
}
