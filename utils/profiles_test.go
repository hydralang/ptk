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

package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/common"
)

func TestRenderRoot(t *testing.T) {
	node := &common.MockNode{}
	node.On("String").Return("node")
	buf := &bytes.Buffer{}

	result := ASCII.Render(buf, "", node, false)

	assert.Equal(t, "   ", result)
	assert.Equal(t, "-- node\n", buf.String())
	node.AssertExpectations(t)
}

func TestRenderLast(t *testing.T) {
	node := &common.MockNode{}
	node.On("String").Return("node")
	buf := &bytes.Buffer{}

	result := ASCII.Render(buf, " ", node, true)

	assert.Equal(t, "    ", result)
	assert.Equal(t, " `- node\n", buf.String())
	node.AssertExpectations(t)
}

func TestRenderBranch(t *testing.T) {
	node := &common.MockNode{}
	node.On("String").Return("node")
	buf := &bytes.Buffer{}

	result := ASCII.Render(buf, " ", node, false)

	assert.Equal(t, " |  ", result)
	assert.Equal(t, " +- node\n", buf.String())
	node.AssertExpectations(t)
}
