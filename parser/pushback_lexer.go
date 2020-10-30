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
	"container/list"

	"github.com/hydralang/ptk/lexer"
)

// IPushBackLexer is an interface for a lexer supporting push-back of
// tokens.  A token that is pushed back will be returned by the next
// call to the objects Next method.
type IPushBackLexer interface {
	lexer.ILexer

	// PushBack pushes a token back into the lexer.  This token
	// will be returned on the next call to the Next method.
	PushBack(tok *lexer.Token)
}

// PushBackLexer is an implementation of lexer.ILexer that includes
// token push-back capability.  A PushBackLexer wraps another
// lexer.ILexer, but provides an additional method for pushing back
// tokens that will be returned later.
type PushBackLexer struct {
	Lexer lexer.ILexer // The source lexer
	toks  *list.List   // A list of pushed-back tokens
}

// NewPushBackLexer wraps another lexer in a PushBackLexer.
func NewPushBackLexer(l lexer.ILexer) *PushBackLexer {
	return &PushBackLexer{
		Lexer: l,
		toks:  &list.List{},
	}
}

// Next returns the next token.  At the end of the lexer, a nil should
// be returned.
func (pbl *PushBackLexer) Next() *lexer.Token {
	// If there are pushed back tokens, return them
	if pbl.toks.Len() > 0 {
		defer func() {
			pbl.toks.Remove(pbl.toks.Front())
		}()

		return pbl.toks.Front().Value.(*lexer.Token)
	}

	// OK, no pushed back tokens; see if there's a lexer
	var tok *lexer.Token
	if pbl.Lexer != nil {
		tok = pbl.Lexer.Next()
		if tok == nil {
			// Exhaused the lexer
			pbl.Lexer = nil
		}
	}

	return tok
}

// PushBack pushes a token back into the lexer.  This token will be
// returned on the next call to the Next method.
func (pbl *PushBackLexer) PushBack(tok *lexer.Token) {
	pbl.toks.PushFront(tok)
}
