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

package nodes

import (
	"fmt"

	"github.com/hydralang/ptk/lexer"
	"github.com/hydralang/ptk/parser"
	"github.com/hydralang/ptk/scanner"
)

// AnnotatedNode is a wrapper for Node that implements Node.  The
// Location and String calls are proxied through, and there is an
// Unwrap call to retrieve the wrapped Node, but the String method
// includes a specified annotation.  This is used to allow attaching
// annotations to the string representations of nodes for the purposes
// of visualizing the AST.
type AnnotatedNode struct {
	node parser.Node // The wrapped node
	ann  string      // The annotation text
}

// NewAnnotatedNode returns a new AnnotatedNode wrapping a given node
// with the specified annotation.
func NewAnnotatedNode(node parser.Node, annotation string) *AnnotatedNode {
	return &AnnotatedNode{
		node: node,
		ann:  annotation,
	}
}

// Location returns the node's location range.
func (an *AnnotatedNode) Location() scanner.Location {
	return an.node.Location()
}

// Children returns a list of child nodes.
func (an *AnnotatedNode) Children() []parser.Node {
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
func (an *AnnotatedNode) Unwrap() parser.Node {
	return an.node
}

// UnaryOperator is a Node implementation that describes the use of a
// unary operator, e.g., "~".
type UnaryOperator struct {
	Loc scanner.Location // The location of the expression
	Op  *lexer.Token     // The unary operator
	Exp parser.Node      // The expression acted upon
}

// UnaryFactory is a factory function that may be passed to Prefix,
// and which constructs a UnaryOperator node.
func UnaryFactory(s parser.State, op *lexer.Token, exp parser.Node) (parser.Node, error) {
	obj := &UnaryOperator{
		Op:  op,
		Exp: exp,
	}

	// Set up the location data
	expLoc := exp.Location()
	if op.Loc != nil && expLoc != nil {
		var err error
		obj.Loc, err = op.Loc.ThruEnd(expLoc)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// Location returns the node's location range.
func (u *UnaryOperator) Location() scanner.Location {
	return u.Loc
}

// Children returns a list of child nodes.
func (u *UnaryOperator) Children() []parser.Node {
	return []parser.Node{NewAnnotatedNode(u.Exp, "Exp")}
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (u *UnaryOperator) String() string {
	return u.Op.String()
}

// BinaryOperator is a Node implementation that describes the use of a
// binary operator, e.g., "*".
type BinaryOperator struct {
	Loc scanner.Location // The location of the expression
	Op  *lexer.Token     // The unary operator
	L   parser.Node      // The left-hand side expression
	R   parser.Node      // The right-hand side expression
}

// BinaryFactory is a factory function that may be passed to Infix or
// InfixR, and which constructs a BinaryOperator node.
func BinaryFactory(s parser.State, l, r parser.Node, op *lexer.Token) (parser.Node, error) {
	obj := &BinaryOperator{
		Op: op,
		L:  l,
		R:  r,
	}

	// Set up the location data
	lLoc := l.Location()
	rLoc := r.Location()
	if lLoc != nil && rLoc != nil {
		var err error
		obj.Loc, err = lLoc.ThruEnd(rLoc)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// Location returns the node's location range.
func (b *BinaryOperator) Location() scanner.Location {
	return b.Loc
}

// Children returns a list of child nodes.
func (b *BinaryOperator) Children() []parser.Node {
	return []parser.Node{
		NewAnnotatedNode(b.L, "L"),
		NewAnnotatedNode(b.R, "R"),
	}
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (b *BinaryOperator) String() string {
	return b.Op.String()
}
