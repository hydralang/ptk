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
	"bytes"
	"fmt"

	"github.com/hydralang/ptk/scanner"
)

// Token represents a single token emitted by the lexical analyzer.  A
// token has an associated symbol, a location, and optionally the
// original text and a semantic value.
type Token struct {
	Type  string           // The type of token
	Loc   scanner.Location // The location of the token
	Value interface{}      // The semantic value of the token; optional
	Text  string           // The original text of the token; optional
}

// Location returns the node's location range.
func (t *Token) Location() scanner.Location {
	return t.Loc
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (t *Token) String() string {
	buf := &bytes.Buffer{}

	// Prefix with location
	if t.Loc != nil {
		fmt.Fprintf(buf, "%s: ", t.Loc)
	}

	// Add the prefix
	fmt.Fprintf(buf, "<%s> token", t.Type)

	// Add the semantic value, if present
	if t.Value != nil {
		fmt.Fprintf(buf, ": %v", t.Value)
	}

	return buf.String()
}
