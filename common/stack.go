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
	"container/list"

	"github.com/stretchr/testify/mock"
)

// Stack is a generic stack.
type Stack interface {
	// Len returns the number of items on the stack.
	Len() int

	// Get returns the current head of the stack.
	Get() interface{}

	// Push pushes a new item onto the stack.
	Push(item interface{})

	// Pop pops an item off the stack and returns it.
	Pop() interface{}

	// Set replaces the current item at the top of the stack with
	// a new item.  It is equivalent to calling Pop, then Push.
	// It returns the old item at the top of the stack.
	Set(item interface{}) interface{}
}

// MockStack is a mock implementation of the Stack interface.
type MockStack struct {
	mock.Mock
}

// Len returns the number of items on the stack.
func (m *MockStack) Len() int {
	args := m.MethodCalled("Len")

	return args.Int(0)
}

// Get returns the current head of the stack.
func (m *MockStack) Get() interface{} {
	args := m.MethodCalled("Get")

	return args.Get(0)
}

// Push pushes a new item onto the stack.
func (m *MockStack) Push(item interface{}) {
	m.MethodCalled("Push", item)
}

// Pop pops an item off the stack and returns it.
func (m *MockStack) Pop() interface{} {
	args := m.MethodCalled("Pop")

	return args.Get(0)
}

// Set replaces the current item at the top of the stack with a new
// item.  It is equivalent to calling Pop, then Push.  It returns the
// old item at the top of the stack.
func (m *MockStack) Set(item interface{}) interface{} {
	args := m.MethodCalled("Set", item)

	return args.Get(0)
}

// stack is an implementation of the Stack interface, based on
// list.List.
type stack struct {
	list.List
}

// NewStack returns a new, empty stack.
func NewStack() Stack {
	return &stack{}
}

// Get returns the current head of the stack.
func (s *stack) Get() interface{} {
	if s.Len() <= 0 {
		return nil
	}

	return s.Front().Value
}

// Push pushes a new item onto the stack.
func (s *stack) Push(item interface{}) {
	s.PushFront(item)
}

// Pop pops an item off the stack and returns it.
func (s *stack) Pop() interface{} {
	if s.Len() <= 0 {
		return nil
	}

	defer func() {
		s.Remove(s.Front())
	}()

	return s.Get()
}

// Set replaces the current item at the top of the stack with a new
// item.  It is equivalent to calling Pop, then Push.  It returns the
// old item at the top of the stack.
func (s *stack) Set(item interface{}) (old interface{}) {
	if s.Len() <= 0 {
		s.Push(item)
		return nil
	}

	old = s.Front().Value
	s.Front().Value = item
	return
}
