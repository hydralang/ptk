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
	"bytes"
	"errors"
	"fmt"

	"github.com/hydralang/ptk/common"
)

// Simple errors that may be generated within the package.
var (
	ErrExpectedToken    = errors.New("No tokens available")
	ErrUnknownTokenType = errors.New("Unknown token type")
	ErrUnexpectedToken  = errors.New("Unexpected token")
	ErrNoTable          = errors.New("Programming error: Parse table missing")
)

// expectedTypes takes a list of token types and generates a string
// describing the expected types.  It is a helper for the error
// creation functions below.
func expectedTypes(types []string) string {
	// Are there any types?
	if types == nil || len(types) <= 0 {
		return ""
	}

	// Handle the case of a single token type
	if len(types) == 1 {
		return fmt.Sprintf("; expected token of type %q", types[0])
	} else if len(types) == 2 {
		return fmt.Sprintf("; expected tokens of type %q or %q", types[0], types[1])
	}

	// Handle the general case
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "; expected tokens of type ")
	for i, tokType := range types {
		switch i {
		case 0:
			fmt.Fprintf(buf, "%q", tokType)
		case len(types) - 1:
			fmt.Fprintf(buf, ", or %q", tokType)
		default:
			fmt.Fprintf(buf, ", %q", tokType)
		}
	}
	return buf.String()
}

// ExpectedToken constructs and returns an ErrExpectedToken.
func ExpectedToken(types ...string) error {
	tail := expectedTypes(types)
	if tail == "" {
		return fmt.Errorf("%w, but expected more", ErrExpectedToken)
	}
	return fmt.Errorf("%w%s", ErrExpectedToken, tail)
}

// UnknownTokenType constructs and returns an ErrUnknownTokenType.
func UnknownTokenType(tok *common.Token, types ...string) error {
	return common.LocationError(tok.Loc, fmt.Errorf("%w %q%s", ErrUnknownTokenType, tok.Type, expectedTypes(types)))
}

// UnexpectedToken cunstructs and returns an ErrUnexpectedToken.
func UnexpectedToken(tok *common.Token, types ...string) error {
	return common.LocationError(tok.Loc, fmt.Errorf("%w of type %q%s", ErrUnexpectedToken, tok.Type, expectedTypes(types)))
}
