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

package visualize

import (
	"bytes"
	"io"

	"github.com/hydralang/ptk/parser"
)

// visualize is a recursive function that renders a tree node.
func visualize(buf io.Writer, prof Profile, node parser.Node, prefix string, last bool) {
	// Render this node
	prefix = prof.Render(buf, prefix, node, last)

	// Recurse to each of its children
	children := node.Children()
	for i, child := range children {
		visualize(buf, prof, child, prefix, i == len(children)-1)
	}
}

// Visualize returns a string containing a visualization of the tree
// rooted at Node.
func Visualize(prof Profile, node parser.Node) string {
	buf := &bytes.Buffer{}
	visualize(buf, prof, node, "", false)
	return buf.String()
}
