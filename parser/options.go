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

// Option is a parse option that may be passed to one of the Parser
// methods.
type Option func(state State)

// AppState is an option allowing an application state to be set when
// parsing an expression or statements.
func AppState(state interface{}) Option {
	return func(s State) {
		s.PushAppState(state)
	}
}