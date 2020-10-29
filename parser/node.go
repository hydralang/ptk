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
	"fmt"

	"github.com/hydralang/ptk/lexer"
	"github.com/hydralang/ptk/scanner"
)

// Node describes one node in an abstract syntax tree.
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

// TokenNode is an implementation of Node that wraps lexer.Token.  It
// provides a default implementation of Children that returns an empty
// list of child nodes.
type TokenNode struct {
	lexer.Token
}

// Children returns a list of child nodes.
func (tn *TokenNode) Children() []Node {
	return []Node{}
}

// AnnotatedNode is a wrapper for Node that implements Node.  The
// Location and String calls are proxied through, but the String
// method includes a specified annotation.  This is used to allow
// attaching annotations to the string representations of nodes for
// the purposes of visualizing the AST.
type AnnotatedNode struct {
	Node       Node   // The wrapped node
	Annotation string // The annotation text
}

// NewAnnotatedNode returns a new AnnotatedNode wrapping a given node
// with the specified annotation.
func NewAnnotatedNode(node Node, annotation string) *AnnotatedNode {
	return &AnnotatedNode{
		Node:       node,
		Annotation: annotation,
	}
}

// Location returns the node's location range.
func (an *AnnotatedNode) Location() scanner.Location {
	return an.Node.Location()
}

// Children returns a list of child nodes.
func (an *AnnotatedNode) Children() []Node {
	return an.Node.Children()
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (an *AnnotatedNode) String() string {
	return fmt.Sprintf("%s: %s", an.Annotation, an.Node)
}

// UnaryOperator is a Node implementation that describes the use of a
// unary operator, e.g., "~".
type UnaryOperator struct {
	Loc scanner.Location // The location of the expression
	Op  *lexer.Token     // The unary operator
	Exp Node             // The expression acted upon
}

// UnaryFactory is a factory function that may be passed to Prefix,
// and which constructs a UnaryOperator node.
func UnaryFactory(p IParser, s State, lex IPushBackLexer, op *lexer.Token, exp Node) (Node, error) {
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
func (u *UnaryOperator) Children() []Node {
	return []Node{NewAnnotatedNode(u.Exp, "Exp")}
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
	L   Node             // The left-hand side expression
	R   Node             // The right-hand side expression
}

// BinaryFactory is a factory function that may be passed to Infix or
// InfixR, and which constructs a BinaryOperator node.
func BinaryFactory(p IParser, s State, lex IPushBackLexer, l, r Node, op *lexer.Token) (Node, error) {
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
func (b *BinaryOperator) Children() []Node {
	return []Node{
		NewAnnotatedNode(b.L, "L"),
		NewAnnotatedNode(b.R, "R"),
	}
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (b *BinaryOperator) String() string {
	return b.Op.String()
}

// literal is a ExprFirst function for literal tokens.  Its
// implementation is trivial: it simply returns the token wrapped in a
// TokenNode.
func literal(p IParser, s State, lex IPushBackLexer, pow int, tok *lexer.Token) (Node, error) {
	return &TokenNode{Token: *tok}, nil
}

// Literal is an ExprFirst function for literal tokens.  It may be
// used directly to initialize a First field in an Entry.  The
// implementation of Literal is trivial: it simply returns the token
// wrapped in a TokenNode.
var Literal = ExprFirst(literal)

// Prefix constructs an ExprFirst function for prefix operators, e.g.,
// "+" or "-" when used directly before a token--that is, not in a
// binary operator context.  For example, "+123" or "-12".  The Prefix
// function should be passed a "factory" function that constructs a
// Node; this function will be called with the token representing the
// operator and the expression to the right of the operator.  It
// should also be called with a binding power; typically, this binding
// power will be higher than the binding power for the same operator
// as a binary operator, which is why it is separate from the Entry.
func Prefix(factory func(p IParser, s State, lex IPushBackLexer, op *lexer.Token, exp Node) (Node, error), power int) ExprFirst {
	return func(p IParser, s State, lex IPushBackLexer, pow int, tok *lexer.Token) (Node, error) {
		// Get the sub-expression on the right
		exp, err := p.Expression(power)
		if err != nil {
			return nil, err
		}

		// Construct the node and return it
		return factory(p, s, lex, tok, exp)
	}
}

// Infix constructs an ExprNext function for infix-style binary
// operators, e.g., "+", "*", etc.  These operators are
// left-associative; that is, an expression like "1 + 2 + 3" is
// equivalent to "(1 + 2) + 3".  The Infix function should be passed a
// "factory" function that constructs a Node; this function will be
// called with the left and right nodes and the token representing the
// operator.
func Infix(factory func(p IParser, s State, lex IPushBackLexer, l, r Node, op *lexer.Token) (Node, error)) ExprNext {
	return func(p IParser, s State, lex IPushBackLexer, pow int, l Node, tok *lexer.Token) (Node, error) {
		// Get the sub-expression on the right
		r, err := p.Expression(pow)
		if err != nil {
			return nil, err
		}

		// Construct the node and return it
		return factory(p, s, lex, l, r, tok)
	}
}

// InfixR is identical to Infix, with the exception that it is used
// for right-associative operators, e.g., "**".  In this case, an
// expression like "1 ** 2 ** 3" is equivalent to "1 ** (2 ** 3)".
// The InfixR should be passed a factory function, which will be
// called with the left and right nodes and the token representing the
// operator.
func InfixR(factory func(p IParser, s State, lex IPushBackLexer, l, r Node, op *lexer.Token) (Node, error)) ExprNext {
	return func(p IParser, s State, lex IPushBackLexer, pow int, l Node, tok *lexer.Token) (Node, error) {
		// Get the sub-expression on the right
		r, err := p.Expression(pow - 1)
		if err != nil {
			return nil, err
		}

		// Construct the node and return it
		return factory(p, s, lex, l, r, tok)
	}
}
