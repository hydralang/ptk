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

func TestMockNodeImplementsNode(t *testing.T) {
	assert.Implements(t, (*Node)(nil), &MockNode{})
}

func TestMockNodeLocationNil(t *testing.T) {
	obj := &MockNode{}
	obj.On("Location").Return(nil)

	result := obj.Location()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockNodeLocationNotNil(t *testing.T) {
	expected := &MockLocation{}
	obj := &MockNode{}
	obj.On("Location").Return(expected)

	result := obj.Location()

	assert.Same(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockNodeChildrenNil(t *testing.T) {
	obj := &MockNode{}
	obj.On("Children").Return(nil)

	result := obj.Children()

	assert.Nil(t, result)
	obj.AssertExpectations(t)
}

func TestMockNodeChildrenNotNil(t *testing.T) {
	expected := []Node{&MockNode{}, &MockNode{}}
	obj := &MockNode{}
	obj.On("Children").Return(expected)

	result := obj.Children()

	assert.Equal(t, expected, result)
	obj.AssertExpectations(t)
}

func TestMockNodeString(t *testing.T) {
	obj := &MockNode{}
	obj.On("String").Return("string")

	result := obj.String()

	assert.Equal(t, "string", result)
	obj.AssertExpectations(t)
}

func TestAnnotatedNodeImplementsNode(t *testing.T) {
	assert.Implements(t, (*Node)(nil), &AnnotatedNode{})
}

func TestNewAnnotatedNode(t *testing.T) {
	node := &MockNode{}

	result := NewAnnotatedNode(node, "annotation")

	assert.Equal(t, &AnnotatedNode{
		node: node,
		ann:  "annotation",
	}, result)
}

func TestAnnotatedNodeLocation(t *testing.T) {
	loc := &MockLocation{}
	node := &MockNode{}
	node.On("Location").Return(loc)
	obj := &AnnotatedNode{
		node: node,
	}

	result := obj.Location()

	assert.Same(t, loc, result)
	node.AssertExpectations(t)
}

func TestAnnotatedNodeChildren(t *testing.T) {
	children := []Node{&MockNode{}, &MockNode{}, &MockNode{}}
	node := &MockNode{}
	node.On("Children").Return(children)
	obj := &AnnotatedNode{
		node: node,
	}

	result := obj.Children()

	assert.Same(t, children[0], result[0])
	assert.Same(t, children[1], result[1])
	assert.Same(t, children[2], result[2])
	node.AssertExpectations(t)
}

func TestAnnotatedNodeString(t *testing.T) {
	node := &MockNode{}
	node.On("String").Return("mock node")
	obj := &AnnotatedNode{
		node: node,
		ann:  "annotation",
	}

	result := obj.String()

	assert.Equal(t, "annotation: mock node", result)
	node.AssertExpectations(t)
}

func TestAnnotatedNodeUnwrap(t *testing.T) {
	node := &MockNode{}
	obj := &AnnotatedNode{
		node: node,
	}

	result := obj.Unwrap()

	assert.Same(t, node, result)
}
