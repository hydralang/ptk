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

package tokenstreams

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/lexer"
)

func TestNewAsyncTokenStream(t *testing.T) {
	toks := []*lexer.Token{{}, {}, {}}
	ts := NewListTokenStream(toks)

	result := NewAsyncTokenStream(ts)

	i := 0
	for tok := result.Next(); tok != nil; tok = result.Next() {
		assert.Same(t, toks[i], tok)
		i++
	}
	assert.Equal(t, len(toks), i)
}
