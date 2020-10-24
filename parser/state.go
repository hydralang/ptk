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
	"github.com/hydralang/ptk/internal"
	"github.com/hydralang/ptk/lexer"
)

// State represents the parser state.  This is passed to the parsing
// functions and should contain all the data they require to perform
// their operation.
type State interface {
	// Parser returns the parser object.
	Parser() Parser

	// AppState returns the current application state.
	AppState() interface{}

	// PushAppState allows pushing an alternative application
	// state onto the application state stack.  Use this, paired
	// with PopAppState, when your grammar needs different state
	// in particular sections of the file.
	PushAppState(state interface{})

	// PopAppState allows popping an application state pushed with
	// PushAppState off the application state stack.  Use this,
	// paired with PushAppState, when your grammar needs different
	// state in particular sections of the file.
	PopAppState() interface{}

	// SetAppState allows changing the current application state
	// on the fly.  Its action is similar to a PopAppState
	// followed by a PushAppState, so the number of entries in the
	// application state stack remains the same.
	SetAppState(state interface{}) interface{}

	// Table returns the parser table currently in use.
	Table() Table

	// PushTable allows pushing an alternative parser table onto
	// the parser table stack.  Use this, paired with PopTable,
	// when your grammar has different rules in particular
	// sections of the file, e.g., like PHP embedded in a web page
	// template.
	PushTable(tab Table)

	// PopTable allows popping a table pushed with PushTable off
	// the parser table stack.  Use this, paired with PushTable,
	// when your grammar has different rules in particular
	// sections of the file, e.g., like PHP embedded in a web page
	// template.
	PopTable() Table

	// SetTable allows changing the current table on the fly.  Its
	// action is similar to a PopTable followed by a PushTable, so
	// the number of entries in the parser table stack remains the
	// same.
	SetTable(tab Table) Table

	// Stream returns the token stream currently in use.
	Stream() lexer.TokenStream

	// PushStream allows pushing an alternative token stream onto
	// the token stream stack.  Use this, paired with PopStream,
	// when your grammar allows inclusion of alternate files.
	// Note that if the current token stream returns a nil, an
	// implicit PopStream will be performed.
	PushStream(ts lexer.TokenStream)

	// PopStream allows popping a token stream off the token
	// stream stack.  Use this, paired with PushStream, when your
	// grammar allows inclusion of alternate files.  Note that
	// PopStream is called implicitly if the token stream returns
	// a nil.
	PopStream() lexer.TokenStream

	// SetStream allows changing the current token stream on the
	// fly.  Its action is similar to a PopStream followed by a
	// PushStream, so the number of entries in the token stream
	// stack remains the same.
	SetStream(ts lexer.TokenStream) lexer.TokenStream

	// Token returns the token currently being processed.  It will
	// be nil if NextToken has not been called, or if NextToken
	// returned a nil value.
	Token() *lexer.Token

	// NextToken returns the next token to be processed.
	NextToken() *lexer.Token

	// MoreTokens returns a boolean true if there are more tokens
	// available, that is, if NextToken will not return nil.
	MoreTokens() bool

	// PushToken pushes a token back.  This token will end up
	// being the next token returned when NextToken is called.
	// Note that the token returned by Token is not changed.
	PushToken(tok *lexer.Token)

	// Expression parses a sub-expression.  This is the core,
	// workhorse function of expression parsing utilizing the
	// Pratt technique; many token parsing functions end up
	// calling this function recursively.  It should be called
	// with a "right binding power", which is used in operator
	// precedence calculations.
	Expression(rbp int) (common.Node, error)

	// Statement parses a single statement.  This is the core,
	// workhorse function for statement parsing utilizing the
	// Pratt technique.
	Statement() (common.Node, error)
}

// MockState is a mock implementation of the State interface.
type MockState struct {
	mock.Mock
}

// Parser returns the parser object.
func (m *MockState) Parser() Parser {
	args := m.MethodCalled("Parser")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Parser)
	}

	return nil
}

// AppState returns the current application state.
func (m *MockState) AppState() interface{} {
	args := m.MethodCalled("AppState")

	return args.Get(0)
}

// PushAppState allows pushing an alternative application state onto
// the application state stack.  Use this, paired with PopAppState,
// when your grammar needs different state in particular sections of
// the file.
func (m *MockState) PushAppState(state interface{}) {
	m.MethodCalled("PushAppState", state)
}

// PopAppState allows popping an application state pushed with
// PushAppState off the application state stack.  Use this, paired
// with PushAppState, when your grammar needs different state in
// particular sections of the file.
func (m *MockState) PopAppState() interface{} {
	args := m.MethodCalled("PopAppState")

	return args.Get(0)
}

// SetAppState allows changing the current application state on the
// fly.  Its action is similar to a PopAppState followed by a
// PushAppState, so the number of entries in the application state
// stack remains the same.
func (m *MockState) SetAppState(state interface{}) interface{} {
	args := m.MethodCalled("SetAppState", state)

	return args.Get(0)
}

// Table returns the parser table currently in use.
func (m *MockState) Table() Table {
	args := m.MethodCalled("Table")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Table)
	}

	return nil
}

// PushTable allows pushing an alternative parser table onto the
// parser table stack.  Use this, paired with PopTable, when your
// grammar has different rules in particular sections of the file,
// e.g., like PHP embedded in a web page template.
func (m *MockState) PushTable(tab Table) {
	m.MethodCalled("PushTable", tab)
}

// PopTable allows popping a table pushed with PushTable off the
// parser table stack.  Use this, paired with PushTable, when your
// grammar has different rules in particular sections of the file,
// e.g., like PHP embedded in a web page template.
func (m *MockState) PopTable() Table {
	args := m.MethodCalled("PopTable")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Table)
	}

	return nil
}

// SetTable allows changing the current table on the fly.  Its action
// is similar to a PopTable followed by a PushTable, so the number of
// entries in the parser table stack remains the same.
func (m *MockState) SetTable(tab Table) Table {
	args := m.MethodCalled("SetTable", tab)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Table)
	}

	return nil
}

// Stream returns the token stream currently in use.
func (m *MockState) Stream() lexer.TokenStream {
	args := m.MethodCalled("Stream")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(lexer.TokenStream)
	}

	return nil
}

// PushStream allows pushing an alternative token stream onto the
// token stream stack.  Use this, paired with PopStream, when your
// grammar allows inclusion of alternate files.  Note that if the
// current token stream returns a nil, an implicit PopStream will be
// performed.
func (m *MockState) PushStream(ts lexer.TokenStream) {
	m.MethodCalled("PushStream", ts)
}

// PopStream allows popping a token stream off the token stream stack.
// Use this, paired with PushStream, when your grammar allows
// inclusion of alternate files.  Note that PopStream is called
// implicitly if the token stream returns a nil.
func (m *MockState) PopStream() lexer.TokenStream {
	args := m.MethodCalled("PopStream")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(lexer.TokenStream)
	}

	return nil
}

// SetStream allows changing the current token stream on the fly.  Its
// action is similar to a PopStream followed by a PushStream, so the
// number of entries in the token stream stack remains the same.
func (m *MockState) SetStream(ts lexer.TokenStream) lexer.TokenStream {
	args := m.MethodCalled("SetStream", ts)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(lexer.TokenStream)
	}

	return nil
}

// Token returns the token currently being processed.  It will be nil
// if NextToken has not been called, or if NextToken returned a nil
// value.
func (m *MockState) Token() *lexer.Token {
	args := m.MethodCalled("Token")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(*lexer.Token)
	}

	return nil
}

// NextToken returns the next token to be processed.
func (m *MockState) NextToken() *lexer.Token {
	args := m.MethodCalled("NextToken")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(*lexer.Token)
	}

	return nil
}

// MoreTokens returns a boolean true if there are more tokens
// available, that is, if NextToken will not return nil.
func (m *MockState) MoreTokens() bool {
	args := m.MethodCalled("MoreTokens")

	return args.Bool(0)
}

// PushToken pushes a token back.  This token will end up being the
// next token returned when NextToken is called.  Note that the token
// returned by Token is not changed.
func (m *MockState) PushToken(tok *lexer.Token) {
	m.MethodCalled("PushToken", tok)
}

// Expression parses a sub-expression.  This is the core, workhorse
// function of expression parsing utilizing the Pratt technique; many
// token parsing functions end up calling this function recursively.
// It should be called with a "right binding power", which is used in
// operator precedence calculations.
func (m *MockState) Expression(rbp int) (common.Node, error) {
	args := m.MethodCalled("Expression", rbp)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(common.Node), args.Error(1)
	}

	return nil, args.Error(1)
}

// Statement parses a single statement.  This is the core, workhorse
// function for statement parsing utilizing the Pratt technique.
func (m *MockState) Statement() (common.Node, error) {
	args := m.MethodCalled("Statement")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(common.Node), args.Error(1)
	}

	return nil, args.Error(1)
}

// state is an implementation of State.
type state struct {
	parser   Parser         // The parser being used
	appState internal.Stack // Stack for application state
	table    internal.Stack // Stack for tables
	stream   internal.Stack // Stack for token streams
	tokens   internal.Stack // Stack of pushed-back tokens
	tok      *lexer.Token   // Last returned token
}

// NewState constructs and returns a new state, with the specified
// table and stream.
func NewState(parser Parser, stream lexer.TokenStream, options []Option) State {
	obj := &state{
		parser:   parser,
		appState: internal.NewStack(),
		table:    internal.NewStack(),
		stream:   internal.NewStack(),
		tokens:   internal.NewStack(),
	}

	// Push the initial table and stream
	obj.table.Push(parser.Table())
	obj.stream.Push(stream)

	// Apply options
	for _, opt := range options {
		opt(obj)
	}

	return obj
}

// Parser returns the parser object.
func (s *state) Parser() Parser {
	return s.parser
}

// AppState returns the current application state.
func (s *state) AppState() interface{} {
	return s.appState.Get()
}

// PushAppState allows pushing an alternative application state onto
// the application state stack.  Use this, paired with PopAppState,
// when your grammar needs different state in particular sections of
// the file.
func (s *state) PushAppState(state interface{}) {
	s.appState.Push(state)
}

// PopAppState allows popping an application state pushed with
// PushAppState off the application state stack.  Use this, paired
// with PushAppState, when your grammar needs different state in
// particular sections of the file.
func (s *state) PopAppState() interface{} {
	return s.appState.Pop()
}

// SetAppState allows changing the current application state on the
// fly.  Its action is similar to a PopAppState followed by a
// PushAppState, so the number of entries in the application state
// stack remains the same.
func (s *state) SetAppState(state interface{}) interface{} {
	return s.appState.Set(state)
}

// Table returns the parser table currently in use.
func (s *state) Table() Table {
	if tmp := s.table.Get(); tmp != nil {
		return tmp.(Table)
	}

	return nil
}

// PushTable allows pushing an alternative parser table onto the
// parser table stack.  Use this, paired with PopTable, when your
// grammar has different rules in particular sections of the file,
// e.g., like PHP embedded in a web page template.
func (s *state) PushTable(tab Table) {
	s.table.Push(tab)
}

// PopTable allows popping a table pushed with PushTable off the
// parser table stack.  Use this, paired with PushTable, when your
// grammar has different rules in particular sections of the file,
// e.g., like PHP embedded in a web page template.
func (s *state) PopTable() Table {
	if tmp := s.table.Pop(); tmp != nil {
		return tmp.(Table)
	}

	return nil
}

// SetTable allows changing the current table on the fly.  Its action
// is similar to a PopTable followed by a PushTable, so the number of
// entries in the parser table stack remains the same.
func (s *state) SetTable(tab Table) Table {
	if tmp := s.table.Set(tab); tmp != nil {
		return tmp.(Table)
	}

	return nil
}

// Stream returns the token stream currently in use.
func (s *state) Stream() lexer.TokenStream {
	if tmp := s.stream.Get(); tmp != nil {
		return tmp.(lexer.TokenStream)
	}

	return nil
}

// PushStream allows pushing an alternative token stream onto the
// token stream stack.  Use this, paired with PopStream, when your
// grammar allows inclusion of alternate files.  Note that if the
// current token stream returns a nil, an implicit PopStream will be
// performed.
func (s *state) PushStream(ts lexer.TokenStream) {
	s.stream.Push(ts)
}

// PopStream allows popping a token stream off the token stream stack.
// Use this, paired with PushStream, when your grammar allows
// inclusion of alternate files.  Note that PopStream is called
// implicitly if the token stream returns a nil.
func (s *state) PopStream() lexer.TokenStream {
	if tmp := s.stream.Pop(); tmp != nil {
		return tmp.(lexer.TokenStream)
	}

	return nil
}

// SetStream allows changing the current token stream on the fly.  Its
// action is similar to a PopStream followed by a PushStream, so the
// number of entries in the token stream stack remains the same.
func (s *state) SetStream(ts lexer.TokenStream) lexer.TokenStream {
	if tmp := s.stream.Set(ts); tmp != nil {
		return tmp.(lexer.TokenStream)
	}

	return nil
}

// Token returns the token currently being processed.  It will be nil
// if NextToken has not been called, or if NextToken returned a nil
// value.
func (s *state) Token() *lexer.Token {
	return s.tok
}

// getToken is a helper that retrieves a token from the token stream.
func (s *state) getToken() *lexer.Token {
	// Loop while we have a stream
	for tmp := s.stream.Get(); tmp != nil; tmp = s.stream.Get() {
		stream := tmp.(lexer.TokenStream)

		// Get the next token from it
		tok := stream.Next()
		if tok == nil {
			// Stream exhausted, pop it off and try again
			s.stream.Pop()
			continue
		}

		// OK, we have a token
		return tok
	}

	return nil
}

// NextToken returns the next token to be processed.
func (s *state) NextToken() *lexer.Token {
	// Are there pushed-back tokens?
	if s.tokens.Len() > 0 {
		s.tok = s.tokens.Pop().(*lexer.Token)
	} else {
		// Get a token from the stream(s)
		s.tok = s.getToken()
	}

	return s.tok
}

// MoreTokens returns a boolean true if there are more tokens
// available, that is, if NextToken will not return nil.
func (s *state) MoreTokens() bool {
	// Check the pushed-back token stack first
	if s.tokens.Len() > 0 {
		return true
	}

	// Get a token from the stream(s)
	tok := s.getToken()

	// Did we get one?
	if tok == nil {
		return false
	}

	// Push the token back so it'll be returned next
	s.PushToken(tok)
	return true
}

// PushToken pushes a token back.  This token will end up being the
// next token returned when NextToken is called.  Note that the token
// returned by Token is not changed.
func (s *state) PushToken(tok *lexer.Token) {
	s.tokens.Push(tok)
}

// getEntry is a helper that retrieves the table entry for a
// particular token.
func (s *state) getEntry(tok *lexer.Token) (Entry, error) {
	// Get the table to use
	tab := s.Table()
	if tab == nil {
		return Entry{}, ErrNoTable
	}

	// Look up the entry
	ent, ok := tab[tok.Type]
	if !ok {
		return Entry{}, UnknownTokenType(tok)
	}

	return ent, nil
}

// Expression parses a sub-expression.  This is the core, workhorse
// function of expression parsing utilizing the Pratt technique; many
// token parsing functions end up calling this function recursively.
// It should be called with a "right binding power", which is used in
// operator precedence calculations.
func (s *state) Expression(rbp int) (common.Node, error) {
	// Get a token from the state
	tok := s.NextToken()
	if tok == nil {
		return nil, ExpectedToken()
	}

	// Get the table entry for it
	ent, err := s.getEntry(tok)
	if err != nil {
		return nil, err
	}

	// Process the token
	node, err := ent.callFirst(s, tok)
	if err != nil {
		return nil, err
	}

	// Handle subsequent tokens
	for tok = s.NextToken(); tok != nil; tok = s.NextToken() {
		// Get the table entry for the token
		ent, err = s.getEntry(tok)
		if err != nil {
			return nil, err
		}

		// Check the binding power of the token
		if rbp >= ent.Power {
			s.PushToken(tok)
			break
		}

		// Process the token
		node, err = ent.callNext(s, node, tok)
		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

// Statement parses a single statement.  This is the core, workhorse
// function for statement parsing utilizing the Pratt technique.
func (s *state) Statement() (common.Node, error) {
	// Get a token from the state
	tok := s.NextToken()
	if tok == nil {
		// No statement
		return nil, nil
	}

	// Get the table entry for it
	ent, err := s.getEntry(tok)
	if err != nil {
		return nil, err
	}

	// Process the token and return the result
	return ent.callStmt(s, tok)
}
