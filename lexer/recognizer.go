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

// Recognizer describes a recognizer.  A recognizer is an object
// returned by the Classify method of a Classifier; its Recognize
// method will be passed the lexer, the state, and a backtracker, and
// it should read input from the backtracker until it has a complete
// lexeme (think "word" in your grammar).  Assuming that lexeme is a
// valid token (a comment or a run of whitespace would not be), the
// Recognize method should then use Lexer.Push to push one or more
// tokens.
type Recognizer interface {
	// Recognize is called by the lexer on the objects returned by
	// the Classifier.  Each will be called in turn until one of
	// the methods returns a boolean true value.  If no recognizer
	// returns true, or if the Classifier returns an empty list,
	// then the Error recognizer will be called, if one is
	// declared, after which the character will be discarded.  The
	// Recognize method will be called with the lexer, the state,
	// and a backtracking scanner.
	Recognize(lexer *Lexer) bool
}
