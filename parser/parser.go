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

import "github.com/hydralang/ptk/lexer"

// IParser is the interface implemented by the Parser.
type IParser interface {
	// Expression parses a single expression from the token stream
	// provided by the lexer.  The method will be called with a
	// "right binding power", which should be 0 for consumers of
	// the parser, but will be non-zero when called recursively.
	Expression(rbp int) (Node, error)

	// Statement parses a single statement from the token stream
	// provided by the lexer.
	Statement() (Node, error)
}

// Parser is an implementation of IParser.
type Parser struct {
	Lexer IPushBackLexer // The lexer providing the tokens
	State State          // The state of the parser
}

// New constructs a new Parser using the provided lexer and state.
func New(l lexer.ILexer, state State) *Parser {
	// Wrap the lexer to allow for push-back
	var ok bool
	var pbl IPushBackLexer
	if pbl, ok = l.(IPushBackLexer); !ok {
		pbl = NewPushBackLexer(l)
	}

	return &Parser{
		Lexer: pbl,
		State: state,
	}
}

// Expression parses a single expression from the token stream
// provided by the lexer.  The method will be called with a "right
// binding power", which should be 0 for consumers of the parser, but
// will be non-zero when called recursively.
func (p *Parser) Expression(rbp int) (Node, error) {
	// Get a token from the lexer
	tok := p.Lexer.Next()
	if tok == nil {
		return nil, ExpectedToken()
	}

	// Get the table entry for it
	ent, ok := p.State.Table()[tok.Type]
	if !ok {
		return nil, UnknownTokenType(tok)
	}

	// Process the token
	node, err := ent.callFirst(p, p.State, p.Lexer, tok)
	if err != nil {
		return nil, err
	}

	// Handle subsequent tokens
	for tok = p.Lexer.Next(); tok != nil; tok = p.Lexer.Next() {
		// Get the table entry for the token
		ent, ok = p.State.Table()[tok.Type]
		if !ok {
			return nil, UnknownTokenType(tok)
		}

		// Check the binding power of the token
		if rbp >= ent.Power {
			p.Lexer.PushBack(tok)
			break
		}

		// Process the token
		node, err = ent.callNext(p, p.State, p.Lexer, node, tok)
		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

// Statement parses a single statement from the token stream provided
// by the lexer.
func (p *Parser) Statement() (Node, error) {
	// Get a token from the lexer
	tok := p.Lexer.Next()
	if tok == nil {
		// No statement
		return nil, nil
	}

	// Get the table entry for it
	ent, ok := p.State.Table()[tok.Type]
	if !ok {
		return nil, UnknownTokenType(tok)
	}

	// Process the token and return the result
	return ent.callStmt(p, p.State, p.Lexer, tok)
}
