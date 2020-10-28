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

package parser

import (
	"github.com/hydralang/ptk/lexer"
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

// TokenNode is an implementation of Node that wraps lexer.Token.
type TokenNode struct {
	lexer.Token
}

// Children returns a list of child nodes.
func (tn *TokenNode) Children() []Node {
	return []Node{}
}
