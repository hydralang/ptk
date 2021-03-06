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

// NewAsyncLexer wraps another lexer and uses the ChanLexer to allow
// running that other lexer in a separate goroutine.
func NewAsyncLexer(ts ILexer) ILexer {
	// Construct the ChanLexer
	obj := NewChanLexer()

	// Run the other lexer in a goroutine and push all its tokens
	go func() {
		for tok := ts.Next(); tok != nil; tok = ts.Next() {
			obj.Push(tok)
		}
		obj.Done()
	}()

	return obj
}
