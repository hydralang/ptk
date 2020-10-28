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
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/lexer"
)

// Patch points to enable testing functions below in isolation.
var (
	newState func(Parser, lexer.ILexer, []Option) State = NewState
)

// Option is a parse option that may be passed to one of the Parser
// methods.
type Option func(state State)

// AppState is an option allowing an application state to be set when
// parsing an expression or statements.
func AppState(state interface{}) Option {
	return func(s State) {
		s.PushAppState(state)
	}
}

// Parser represents the actual parser.  This is passed to the parsing
// functions along with the parsing state.
type Parser interface {
	// Table returns the default table that will be used to
	// initialize the state.
	Table() Table

	// Expression parses a single expression from the specified
	// token stream.
	Expression(stream lexer.ILexer, options ...Option) (Node, error)

	// Statement parses a single statement from the specified
	// token stream.
	Statement(stream lexer.ILexer, options ...Option) (Node, error)

	// Statements parses all statements from the specified token
	// stream.  It is essentially equivalent to running Statement
	// in a loop until all tokens are exhausted.
	Statements(stream lexer.ILexer, options ...Option) ([]Node, error)
}

// MockParser is a mock implementation of the Parser interface.
type MockParser struct {
	mock.Mock
}

// Table returns the default table that will be used to initialize the
// state.
func (m *MockParser) Table() Table {
	args := m.MethodCalled("Table")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Table)
	}

	return nil
}

// Expression parses a single expression from the specified token
// stream.
func (m *MockParser) Expression(stream lexer.ILexer, options ...Option) (Node, error) {
	args := m.MethodCalled("Expression", stream, options)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Node), args.Error(1)
	}

	return nil, args.Error(1)
}

// Statement parses a single statement from the specified token
// stream.
func (m *MockParser) Statement(stream lexer.ILexer, options ...Option) (Node, error) {
	args := m.MethodCalled("Statement", stream, options)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Node), args.Error(1)
	}

	return nil, args.Error(1)
}

// Statements parses all statements from the specified token stream.
// It is essentially equivalent to running Statement in a loop until
// all tokens are exhausted.
func (m *MockParser) Statements(stream lexer.ILexer, options ...Option) ([]Node, error) {
	args := m.MethodCalled("Statements", stream, options)

	if tmp := args.Get(0); tmp != nil {
		return tmp.([]Node), args.Error(1)
	}

	return nil, args.Error(1)
}

// parser is an implementation of Parser.
type parser struct {
	table Table // The initial table to use when constructing a state
}

// New constructs a new parser, with the specified table.
func New(table Table) Parser {
	return &parser{
		table: table,
	}
}

// Table returns the default table that will be used to initialize the
// state.
func (p *parser) Table() Table {
	return p.table
}

// Expression parses a single expression from the specified token
// stream.
func (p *parser) Expression(stream lexer.ILexer, options ...Option) (Node, error) {
	// Construct a state
	s := newState(p, stream, options)

	// Parse an expression
	return s.Expression(0)
}

// Statement parses a single statement from the specified token
// stream.
func (p *parser) Statement(stream lexer.ILexer, options ...Option) (Node, error) {
	// Construct a state
	s := newState(p, stream, options)

	// Parse a statement
	return s.Statement()
}

// Statements parses all statements from the specified token stream.
// It is essentially equivalent to running Statement in a loop until
// all tokens are exhausted.
func (p *parser) Statements(stream lexer.ILexer, options ...Option) ([]Node, error) {
	// Construct a state
	s := newState(p, stream, options)

	// Parse as many statements as possible
	var err error
	var node Node
	nodes := []Node{}
	for node, err = s.Statement(); node != nil; node, err = s.Statement() {
		// Add the new node to the list
		nodes = append(nodes, node)
	}
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
