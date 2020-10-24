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
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/scanner"
)

// Node describes one node in an abstract syntax tree.  Note that it
// is deliberate that Token implements Node.
type Node interface {
	// Location returns the node's location range.
	Location() scanner.Location

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
func (m *MockNode) Location() scanner.Location {
	args := m.MethodCalled("Location")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(scanner.Location)
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
