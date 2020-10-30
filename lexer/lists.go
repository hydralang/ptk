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

package lexer

// ListLexer is an implementation of Lexer that is initialized with a
// list of tokens, and simply returns the tokens in sequence.
type ListLexer struct {
	toks    []*Token // The list of tokens
	idx     int      // The index of the current token to return
	started bool     // A boolean indicating whether the iterator has started
}

// NewListLexer returns a Lexer that retrieves its tokens from a list
// passed to the function.  This actually uses a ChanLexer under the
// covers.
func NewListLexer(toks []*Token) *ListLexer {
	// return obj
	return &ListLexer{
		toks: toks,
	}
}

// Next returns the next token.  At the end of the lexer, a nil should
// be returned.
func (lts *ListLexer) Next() *Token {
	// Check the state
	switch {
	case !lts.started: // Need to start?
		lts.started = true
		if lts.idx >= len(lts.toks)-1 {
			return nil
		}

	case lts.idx >= len(lts.toks)-1:
		return nil

	default:
		lts.idx++
	}

	// Return the indexed token
	return lts.toks[lts.idx]
}
