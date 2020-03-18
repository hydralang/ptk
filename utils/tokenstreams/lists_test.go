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

	"github.com/hydralang/ptk/common"
)

func TestListTokenStreamImplementsTokenStream(t *testing.T) {
	assert.Implements(t, (*common.TokenStream)(nil), &listTokenStream{})
}

func TestNewListTokenStream(t *testing.T) {
	toks := []*common.Token{{}, {}, {}}

	result := NewListTokenStream(toks)

	assert.Equal(t, &listTokenStream{
		toks: toks,
	}, result)
}

func TestListTokenStreamNextUnstarted(t *testing.T) {
	toks := []*common.Token{{}, {}, {}}
	obj := &listTokenStream{
		toks: toks,
	}

	result := obj.Next()

	assert.Same(t, toks[0], result)
	assert.Equal(t, &listTokenStream{
		toks:    toks,
		started: true,
	}, obj)
}

func TestListTokenStreamNextStarted(t *testing.T) {
	toks := []*common.Token{{}, {}, {}}
	obj := &listTokenStream{
		toks:    toks,
		started: true,
	}

	result := obj.Next()

	assert.Same(t, toks[1], result)
	assert.Equal(t, &listTokenStream{
		toks:    toks,
		idx:     1,
		started: true,
	}, obj)
}

func TestListTokenStreamNextEnding(t *testing.T) {
	toks := []*common.Token{{}, {}, {}}
	obj := &listTokenStream{
		toks:    toks,
		idx:     2,
		started: true,
	}

	result := obj.Next()

	assert.Nil(t, result)
	assert.Equal(t, &listTokenStream{
		toks:    toks,
		idx:     2,
		started: true,
	}, obj)
}
