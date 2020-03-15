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
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/common"
)

// CharStream presents a stream of characters.  The basic character
// stream does not provide backtracking or character push-back.  The
// Scanner is one implementation of CharStream.
type CharStream interface {
	// Next returns the next character from the stream as a Char,
	// which will include the character's location.  If an error
	// was encountered, that will also be returned.
	Next() (common.Char, error)
}

// MockCharStream is a mock implementation of the CharStream
// interface.
type MockCharStream struct {
	mock.Mock
}

// Next returns the next character from the stream as a Char, which
// will include the character's location.  If an error was
// encountered, that will also be returned.
func (m *MockCharStream) Next() (common.Char, error) {
	args := m.MethodCalled("Next")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(common.Char), args.Error(1)
	}

	return common.Char{}, args.Error(1)
}
