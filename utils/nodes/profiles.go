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

package nodes

import (
	"fmt"
	"io"

	"github.com/hydralang/ptk/common"
)

// Profile represents a visual profile for visualizing the tree.
type Profile struct {
	Start  rune // Joiner character to start a tree at the root
	Last   rune // Joiner character for last child
	Branch rune // Joiner character for a branch
	Skip   rune // Joiner character indicating a node is not a branch
	Into   rune // Joiner character leading into a node
}

// Render renders a node into the buffer with the appropriate tree
// indicator characters drawn from the profile.  It returns the prefix
// that should be used when rendering children of this node.
func (p *Profile) Render(buf io.Writer, prefix string, node common.Node, last bool) string {
	// Pick the joiner and prepare the next prefix
	var joiner rune
	var nextPrefix string
	if prefix == "" {
		joiner = p.Start
		nextPrefix = "   "
	} else if last {
		joiner = p.Last
		nextPrefix = fmt.Sprintf("%s   ", prefix)
	} else {
		joiner = p.Branch
		nextPrefix = fmt.Sprintf("%s%c  ", prefix, p.Skip)
	}

	// Render the node into the buffer
	fmt.Fprintf(buf, "%s%c%c %s\n", prefix, joiner, p.Into, node)

	return nextPrefix
}

// Pre-configured visual profiles for trees.
var (
	ASCII = Profile{
		Start:  '-',
		Last:   '`',
		Branch: '+',
		Skip:   '|',
		Into:   '-',
	}
	Rounded = Profile{
		Start:  '\u2500',
		Last:   '\u2570',
		Branch: '\u251c',
		Skip:   '\u2502',
		Into:   '\u2500',
	}
	Square = Profile{
		Start:  '\u2500',
		Last:   '\u2514',
		Branch: '\u251c',
		Skip:   '\u2502',
		Into:   '\u2500',
	}
)
