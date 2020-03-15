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

import "github.com/klmitch/kent"

// Option is a lexer option that may be passed to the Lex method.
type Option func(state State)

// AppState is an option allowing an application state to be set when
// lexing a character stream.
func AppState(state interface{}) Option {
	return func(s State) {
		s.PushAppState(state)
	}
}

// ScannerOption is a scanner option that may be passed to the
// NewScanner function.
type ScannerOption func(s *scanner)

// LineEndings is a scanner option that may be used to set the
// preferred line ending style.  A line ending style is an instance of
// LineStyle that controls how the scanner recognizes newlines.  The
// scanner always converts line endings into single newlines.
func LineEndings(ls LineStyle) ScannerOption {
	return func(s *scanner) {
		s.ls = ls
	}
}

// TabStop is a scanner option that may be used to specify a different
// tab stop setting than the default of 8.
func TabStop(ts int) ScannerOption {
	return func(s *scanner) {
		s.ts = ts
	}
}

// Reporter is a scanner option that may be used to specify a
// kent.Reporter to use for reporting encoding errors.  By default,
// encoding errors are reported to standard error.
func Reporter(rep kent.Reporter) ScannerOption {
	return func(s *scanner) {
		s.rep = rep
	}
}
