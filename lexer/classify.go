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

// Classifier represents a character classification tool.  A
// classifier has a Classify method that takes the lexer, the state,
// and a backtracker, and returns a list of recognizers, which the
// lexer then runs in order until one of them succeeds.
type Classifier interface {
	// Classify takes a lexer, a state, and a backtracking scanner
	// and determines one or more recognizers to extract a token
	// or a set of tokens from the lexer input.
	Classify(lexer *Lexer) []Recognizer

	// Error is called by the lexer if all recognizers returned by
	// Classify return without success.
	Error(lexer *Lexer)
}
