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
	"github.com/hydralang/ptk/lexer"
	"github.com/hydralang/ptk/parser"
)

// literal is a ExprFirst function for literal tokens.  Its
// implementation is trivial: it simply returns the token, as Token
// implements Node.
func literal(p parser.IParser, s parser.State, lex parser.IPushBackLexer, pow int, tok *lexer.Token) (parser.Node, error) {
	return &parser.TokenNode{Token: *tok}, nil
}

// Literal is an ExprFirst function for literal tokens.  It may be
// used directly to initialize a First field in an Entry.  The
// implementation of Literal is trivial: it simply returns the token,
// as Token implements Node.
var Literal = parser.ExprFirst(literal)

// Prefix constructs an ExprFirst function for prefix operators, e.g.,
// "+" or "-" when used directly before a token--that is, not in a
// binary operator context.  For example, "+123" or "-12".  The Prefix
// function should be passed a "factory" function that constructs a
// Node; this function will be called with the token representing the
// operator and the expression to the right of the operator.  It
// should also be called with a binding power; typically, this binding
// power will be higher than the binding power for the same operator
// as a binary operator, which is why it is separate from the Entry.
func Prefix(factory func(p parser.IParser, s parser.State, lex parser.IPushBackLexer, op *lexer.Token, exp parser.Node) (parser.Node, error), power int) parser.ExprFirst {
	return func(p parser.IParser, s parser.State, lex parser.IPushBackLexer, pow int, tok *lexer.Token) (parser.Node, error) {
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
func Infix(factory func(p parser.IParser, s parser.State, lex parser.IPushBackLexer, l, r parser.Node, op *lexer.Token) (parser.Node, error)) parser.ExprNext {
	return func(p parser.IParser, s parser.State, lex parser.IPushBackLexer, pow int, l parser.Node, tok *lexer.Token) (parser.Node, error) {
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
func InfixR(factory func(p parser.IParser, s parser.State, lex parser.IPushBackLexer, l, r parser.Node, op *lexer.Token) (parser.Node, error)) parser.ExprNext {
	return func(p parser.IParser, s parser.State, lex parser.IPushBackLexer, pow int, l parser.Node, tok *lexer.Token) (parser.Node, error) {
		// Get the sub-expression on the right
		r, err := p.Expression(pow - 1)
		if err != nil {
			return nil, err
		}

		// Construct the node and return it
		return factory(p, s, lex, l, r, tok)
	}
}
