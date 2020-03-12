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

import "github.com/hydralang/ptk/common"

// UnaryOperator is a Node implementation that describes the use of a
// unary operator, e.g., "~".
type UnaryOperator struct {
	Loc common.Location // The location of the expression
	Op  *common.Token   // The unary operator
	Exp common.Node     // The expression acted upon
}

// UnaryFactory is a factory function that may be passed to Prefix,
// and which constructs a UnaryOperator node.
func UnaryFactory(op *common.Token, exp common.Node) (common.Node, error) {
	obj := &UnaryOperator{
		Op:  op,
		Exp: exp,
	}

	// Set up the location data
	expLoc := exp.Location()
	if op.Loc != nil && expLoc != nil {
		var err error
		obj.Loc, err = op.Loc.ThruEnd(expLoc)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// Location returns the node's location range.
func (u *UnaryOperator) Location() common.Location {
	return u.Loc
}

// Children returns a list of child nodes.
func (u *UnaryOperator) Children() []common.Node {
	return []common.Node{common.NewAnnotatedNode(u.Exp, "Exp")}
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (u *UnaryOperator) String() string {
	return u.Op.String()
}

// BinaryOperator is a Node implementation that describes the use of a
// binary operator, e.g., "~".
type BinaryOperator struct {
	Loc common.Location // The location of the expression
	Op  *common.Token   // The unary operator
	L   common.Node     // The left-hand side expression
	R   common.Node     // The right-hand side expression
}

// BinaryFactory is a factory function that may be passed to Prefix,
// and which constructs a BinaryOperator node.
func BinaryFactory(l, r common.Node, op *common.Token) (common.Node, error) {
	obj := &BinaryOperator{
		Op: op,
		L:  l,
		R:  r,
	}

	// Set up the location data
	lLoc := l.Location()
	rLoc := r.Location()
	if lLoc != nil && rLoc != nil {
		var err error
		obj.Loc, err = lLoc.ThruEnd(rLoc)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// Location returns the node's location range.
func (b *BinaryOperator) Location() common.Location {
	return b.Loc
}

// Children returns a list of child nodes.
func (b *BinaryOperator) Children() []common.Node {
	return []common.Node{
		common.NewAnnotatedNode(b.L, "L"),
		common.NewAnnotatedNode(b.R, "R"),
	}
}

// String returns a string describing the node.  This should include
// the location range that encompasses all of the node's tokens.
func (b *BinaryOperator) String() string {
	return b.Op.String()
}
