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

// State represents the state of the lexer.  This is an interface; an
// implementation must be provided by the user.  A base implementation
// is available as BaseState.
type State interface {
	// Classifier must return the classifier to use.  It is safe
	// for the application to return different Classifier
	// implementations depending on the lexer state.
	Classifier() Classifier
}

// BaseState is a basic implementation of the State interface.  It
// assumes a fixed Classifier for the lifetime of the lexer's
// operation.
type BaseState struct {
	Cls Classifier // The classifier for the lex
}

// Classifier must return the classifier to use.  It is safe for the
// application to return different Classifier implementations
// depending on the lexer state.
func (bs *BaseState) Classifier() Classifier {
	return bs.Cls
}
