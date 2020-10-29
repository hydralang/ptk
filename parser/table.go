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
)

// ExprFirst functions are called to process the first token in an
// expression.  Functions of this type are typically declared on
// literal tokens or prefix operators.
type ExprFirst func(parser IParser, state State, lex IPushBackLexer, power int, tok *lexer.Token) (Node, error)

// ExprNext functions are called to process the subsequent tokens in
// an expression.  Functions of this type are typically declared with
// a "left binding power" (a measure of how tightly an operator binds
// to its operands), and are used with binary operators.
type ExprNext func(parser IParser, state State, lex IPushBackLexer, power int, left Node, tok *lexer.Token) (Node, error)

// Statement functions are called to process a statement.  They're
// called with the first token of the statement, and should read in
// additional tokens.
type Statement func(parser IParser, state State, lex IPushBackLexer, tok *lexer.Token) (Node, error)

// Entry is an entry in the parser table.  The Pratt technique is
// table driven, based on the token type; objects of this type contain
// a single entry from the table.
type Entry struct {
	Power int       // The binding power of the token type
	First ExprFirst // The function to call for an initial token
	Next  ExprNext  // The function to call for the next token
	Stmt  Statement // The function to call for statement tokens
}

// callFirst is a helper to call the declared ExprFirst function.
func (e Entry) callFirst(p IParser, s State, lex IPushBackLexer, tok *lexer.Token) (Node, error) {
	if e.First == nil {
		return nil, UnexpectedToken(tok)
	}

	return e.First(p, s, lex, e.Power, tok)
}

// callNext is a helper to call the declared ExprNext function.
func (e Entry) callNext(p IParser, s State, lex IPushBackLexer, l Node, tok *lexer.Token) (Node, error) {
	if e.Next == nil {
		return nil, UnexpectedToken(tok)
	}

	return e.Next(p, s, lex, e.Power, l, tok)
}

// callStmt is a helper to call the declared Statement function.
func (e Entry) callStmt(p IParser, s State, lex IPushBackLexer, tok *lexer.Token) (Node, error) {
	if e.Stmt == nil {
		return nil, UnexpectedToken(tok)
	}

	return e.Stmt(p, s, lex, tok)
}

// Table is a table of entries by their token type.  The Pratt
// technique is table driven, based on the token type; objects of this
// type contain the table.
type Table map[string]Entry
