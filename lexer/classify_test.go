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

import "github.com/stretchr/testify/mock"

type mockClassifier struct {
	mock.Mock
}

func (m *mockClassifier) Classify(lexer *Lexer) []Recognizer {
	args := m.MethodCalled("Classify", lexer)

	if tmp := args.Get(0); tmp != nil {
		return tmp.([]Recognizer)
	}

	return nil
}

func (m *mockClassifier) Error(lexer *Lexer) {
	m.MethodCalled("Error", lexer)
}
