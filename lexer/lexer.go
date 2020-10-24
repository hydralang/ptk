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

import (
	"container/list"

	"github.com/hydralang/ptk/scanner"
)

// ILexer presents a stream of tokens.  The basic lexer does not
// provide token pushback.
type ILexer interface {
	// Next returns the next token.  At the end of the lexer, a
	// nil should be returned.
	Next() *Token
}

// Lexer is an implementation of ILexer.
type Lexer struct {
	Scanner IBackTracker // The character source, wrapped in a BackTracker
	State   State        // The state of the lexer
	toks    *list.List   // List of tokens to produce
}

// New constructs a new Lexer using the provided source and state.
func New(src scanner.Scanner, state State) *Lexer {
	return &Lexer{
		Scanner: NewBackTracker(src, TrackAll),
		State:   state,
		toks:    &list.List{},
	}
}

// next is the actual implementation of the lexer.  This is the
// routine that calls the Classify, Recognize, and Error methods
// provided by the state.
func (l *Lexer) next() {
	// Reset the backtracker
	l.Scanner.SetMax(TrackAll)

	// Classify the contents
	for _, rec := range l.State.Classifier().Classify(l, l.State, l.Scanner) {
		l.Scanner.BackTrack()
		if rec.Recognize(l, l.State, l.Scanner) {
			l.Scanner.Accept(0)
			return
		}
	}

	// None of the recognizers recognized the contents
	l.Scanner.BackTrack()
	l.State.Classifier().Error(l, l.State, l.Scanner)
	l.Scanner.Accept(0)
}

// Next returns the next token.  At the end of the lexer, a nil should
// be returned.
func (l *Lexer) Next() *Token {
	// Loop until we have a token or all characters have been
	// processed
	for l.toks.Len() <= 0 {
		if !l.Scanner.More() {
			return nil
		}

		l.next()
	}

	// Return a token off the token queue
	defer func() {
		l.toks.Remove(l.toks.Front())
	}()
	return l.toks.Front().Value.(*Token)
}

// Push pushes a token onto the list of tokens to be returned by the
// lexer.  Recognizers should call this method with the token or
// tokens that they recognize from the input.
func (l *Lexer) Push(tok *Token) bool {
	l.toks.PushBack(tok)
	return true
}
