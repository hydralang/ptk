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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hydralang/ptk/common"
)

func TestVisualizeInner(t *testing.T) {
	children := []common.Node{&common.MockNode{}, &common.MockNode{}, &common.MockNode{}}
	for i, child := range children {
		n := child.(*common.MockNode)
		n.On("String").Return(fmt.Sprintf("child%d", i))
		n.On("Children").Return([]common.Node{})
	}
	node := &common.MockNode{}
	node.On("String").Return("node")
	node.On("Children").Return(children)
	buf := &bytes.Buffer{}

	visualize(buf, ASCII, node, "", false)

	assert.Equal(t, "-- node\n   +- child0\n   +- child1\n   `- child2\n", buf.String())
}

func TestVisualize(t *testing.T) {
	children := []common.Node{&common.MockNode{}, &common.MockNode{}, &common.MockNode{}}
	for i, child := range children {
		n := child.(*common.MockNode)
		n.On("String").Return(fmt.Sprintf("child%d", i))
		n.On("Children").Return([]common.Node{})
	}
	node := &common.MockNode{}
	node.On("String").Return("node")
	node.On("Children").Return(children)

	result := Visualize(ASCII, node)

	assert.Equal(t, "-- node\n   +- child0\n   +- child1\n   `- child2\n", result)
}
