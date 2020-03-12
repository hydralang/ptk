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

	"github.com/stretchr/testify/mock"
)

// Node describes one node in an abstract syntax tree.  Note that it
// is deliberate that Token implements Node.
type Node interface {
	// Location returns the node's location range.
	Location() Location

	// Children returns a list of child nodes.
	Children() []Node

	// String returns a string describing the node.  This should
	// include the location range that encompasses all of the
	// node's tokens.
	String() string
}

// MockNode is a mock implementation of the Node interface.
type MockNode struct {
	mock.Mock
}

// Location returns the node's location range.
func (m *MockNode) Location() Location {
	args := m.MethodCalled("Location")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Location)
	}

	return nil
}

// Children returns a list of child nodes.
func (m *MockNode) Children() []Node {
	args := m.MethodCalled("Children")

	if tmp := args.Get(0); tmp != nil {
		return tmp.([]Node)
	}

	return nil
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (m *MockNode) String() string {
	args := m.MethodCalled("String")

	return args.String(0)
}

// AnnotatedNode is a wrapper for Node that implements Node.  The
// Location and String calls are proxied through, and there is an
// Unwrap call to retrieve the wrapped Node, but the String method
// includes a specified annotation.  This is used to allow attaching
// annotations to the string representations of nodes for the purposes
// of visualizing the AST.
type AnnotatedNode struct {
	node Node   // The wrapped node
	ann  string // The annotation text
}

// NewAnnotatedNode returns a new AnnotatedNode wrapping a given node
// with the specified annotation.
func NewAnnotatedNode(node Node, annotation string) *AnnotatedNode {
	return &AnnotatedNode{
		node: node,
		ann:  annotation,
	}
}

// Location returns the node's location range.
func (an *AnnotatedNode) Location() Location {
	return an.node.Location()
}

// Children returns a list of child nodes.
func (an *AnnotatedNode) Children() []Node {
	return an.node.Children()
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (an *AnnotatedNode) String() string {
	return fmt.Sprintf("%s: %s", an.ann, an.node)
}

// Unwrap returns the underlying node.  This may be used when the
// underlying node contains data or other methods that are not
// otherwise accessible.
func (an *AnnotatedNode) Unwrap() Node {
	return an.node
}
