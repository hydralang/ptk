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

	"github.com/hydralang/ptk/common"
)

// Patch points to enable testing functions below in isolation.
var (
	newState func(Table, common.TokenStream) State = NewState
)

// Parser represents the actual parser.  This is passed to the parsing
// functions along with the parsing state.
type Parser interface {
	// Expression parses a single expression from the specified
	// token stream.
	Expression(stream common.TokenStream) (common.Node, error)

	// Statement parses a single statement from the specified
	// token stream.
	Statement(stream common.TokenStream) (common.Node, error)

	// Statements parses all statements from the specified token
	// stream.  It is essentially equivalent to running Statement
	// in a loop until all tokens are exhausted.
	Statements(stream common.TokenStream) ([]common.Node, error)
}

// MockParser is a mock implementation of the Parser interface.
type MockParser struct {
	mock.Mock
}

// Expression parses a single expression from the specified token
// stream.
func (m *MockParser) Expression(stream common.TokenStream) (common.Node, error) {
	args := m.MethodCalled("Expression", stream)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(common.Node), args.Error(1)
	}

	return nil, args.Error(1)
}

// Statement parses a single statement from the specified token
// stream.
func (m *MockParser) Statement(stream common.TokenStream) (common.Node, error) {
	args := m.MethodCalled("Statement", stream)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(common.Node), args.Error(1)
	}

	return nil, args.Error(1)
}

// Statements parses all statements from the specified token stream.
// It is essentially equivalent to running Statement in a loop until
// all tokens are exhausted.
func (m *MockParser) Statements(stream common.TokenStream) ([]common.Node, error) {
	args := m.MethodCalled("Statements", stream)

	if tmp := args.Get(0); tmp != nil {
		return tmp.([]common.Node), args.Error(1)
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

// Expression parses a single expression from the specified token
// stream.
func (p *parser) Expression(stream common.TokenStream) (common.Node, error) {
	// Construct a state
	s := newState(p.table, stream)

	// Parse an expression
	return s.Expression(p, 0)
}

// Statement parses a single statement from the specified token
// stream.
func (p *parser) Statement(stream common.TokenStream) (common.Node, error) {
	// Construct a state
	s := newState(p.table, stream)

	// Parse a statement
	return s.Statement(p)
}

// Statements parses all statements from the specified token stream.
// It is essentially equivalent to running Statement in a loop until
// all tokens are exhausted.
func (p *parser) Statements(stream common.TokenStream) ([]common.Node, error) {
	// Construct a state
	s := newState(p.table, stream)

	// Parse as many statements as possible
	var err error
	var node common.Node
	nodes := []common.Node{}
	for node, err = s.Statement(p); node != nil; node, err = s.Statement(p) {
		// Add the new node to the list
		nodes = append(nodes, node)
	}
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
