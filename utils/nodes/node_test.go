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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/lexer"
	"github.com/hydralang/ptk/parser"
	"github.com/hydralang/ptk/scanner"
)

type mockLocation struct {
	mock.Mock
}

func (m *mockLocation) String() string {
	args := m.MethodCalled("String")

	return args.String(0)
}

func (m *mockLocation) Thru(other scanner.Location) (scanner.Location, error) {
	args := m.MethodCalled("Thru", other)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(scanner.Location), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockLocation) ThruEnd(other scanner.Location) (scanner.Location, error) {
	args := m.MethodCalled("ThruEnd", other)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(scanner.Location), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockLocation) Incr(c rune, tabstop int) scanner.Location {
	args := m.MethodCalled("Incr", c, tabstop)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(scanner.Location)
	}

	return nil
}

type mockNode struct {
	mock.Mock
}

func (m *mockNode) Location() scanner.Location {
	args := m.MethodCalled("Location")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(scanner.Location)
	}

	return nil
}

func (m *mockNode) Children() []parser.Node {
	args := m.MethodCalled("Children")

	if tmp := args.Get(0); tmp != nil {
		return tmp.([]parser.Node)
	}

	return nil
}

func (m *mockNode) String() string {
	args := m.MethodCalled("String")

	return args.String(0)
}

func TestAnnotatedNodeImplementsNode(t *testing.T) {
	assert.Implements(t, (*parser.Node)(nil), &AnnotatedNode{})
}

func TestNewAnnotatedNode(t *testing.T) {
	node := &mockNode{}

	result := NewAnnotatedNode(node, "annotation")

	assert.Equal(t, &AnnotatedNode{
		node: node,
		ann:  "annotation",
	}, result)
}

func TestAnnotatedNodeLocation(t *testing.T) {
	loc := &mockLocation{}
	node := &mockNode{}
	node.On("Location").Return(loc)
	obj := &AnnotatedNode{
		node: node,
	}

	result := obj.Location()

	assert.Same(t, loc, result)
	node.AssertExpectations(t)
}

func TestAnnotatedNodeChildren(t *testing.T) {
	children := []parser.Node{&mockNode{}, &mockNode{}, &mockNode{}}
	node := &mockNode{}
	node.On("Children").Return(children)
	obj := &AnnotatedNode{
		node: node,
	}

	result := obj.Children()

	assert.Same(t, children[0], result[0])
	assert.Same(t, children[1], result[1])
	assert.Same(t, children[2], result[2])
	node.AssertExpectations(t)
}

func TestAnnotatedNodeString(t *testing.T) {
	node := &mockNode{}
	node.On("String").Return("mock node")
	obj := &AnnotatedNode{
		node: node,
		ann:  "annotation",
	}

	result := obj.String()

	assert.Equal(t, "annotation: mock node", result)
	node.AssertExpectations(t)
}

func TestAnnotatedNodeUnwrap(t *testing.T) {
	node := &mockNode{}
	obj := &AnnotatedNode{
		node: node,
	}

	result := obj.Unwrap()

	assert.Same(t, node, result)
}

func TestUnaryOperatorImplementsNode(t *testing.T) {
	assert.Implements(t, (*parser.Node)(nil), &UnaryOperator{})
}

func TestUnaryFactoryBase(t *testing.T) {
	s := &mockState{}
	op := &lexer.Token{}
	exp := &mockNode{}
	exp.On("Location").Return(nil)

	result, err := UnaryFactory(s, op, exp)

	assert.NoError(t, err)
	assert.Equal(t, &UnaryOperator{
		Op:  op,
		Exp: exp,
	}, result)
}

func TestUnaryFactoryLocation(t *testing.T) {
	s := &mockState{}
	finalLoc := &mockLocation{}
	opLoc := &mockLocation{}
	expLoc := &mockLocation{}
	opLoc.On("ThruEnd", expLoc).Return(finalLoc, nil)
	op := &lexer.Token{
		Loc: opLoc,
	}
	exp := &mockNode{}
	exp.On("Location").Return(expLoc)

	result, err := UnaryFactory(s, op, exp)

	assert.NoError(t, err)
	assert.Equal(t, &UnaryOperator{
		Loc: finalLoc,
		Op:  op,
		Exp: exp,
	}, result)
	assert.Same(t, finalLoc, result.(*UnaryOperator).Loc)
	opLoc.AssertExpectations(t)
	exp.AssertExpectations(t)
}

func TestUnaryFactoryLocationError(t *testing.T) {
	s := &mockState{}
	opLoc := &mockLocation{}
	expLoc := &mockLocation{}
	opLoc.On("ThruEnd", expLoc).Return(nil, assert.AnError)
	op := &lexer.Token{
		Loc: opLoc,
	}
	exp := &mockNode{}
	exp.On("Location").Return(expLoc)

	result, err := UnaryFactory(s, op, exp)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	opLoc.AssertExpectations(t)
	exp.AssertExpectations(t)
}

func TestUnaryOperatorLocation(t *testing.T) {
	loc := &mockLocation{}
	obj := &UnaryOperator{
		Loc: loc,
	}

	result := obj.Location()

	assert.Same(t, loc, result)
}

func TestUnaryOperatorChildren(t *testing.T) {
	exp := &mockNode{}
	obj := &UnaryOperator{
		Exp: exp,
	}

	result := obj.Children()

	assert.Equal(t, []parser.Node{
		NewAnnotatedNode(exp, "Exp"),
	}, result)
}

func TestUnaryOperatorString(t *testing.T) {
	obj := &UnaryOperator{
		Op: &lexer.Token{Type: "op"},
	}

	result := obj.String()

	assert.Equal(t, "<op> token", result)
}

func TestBinaryOperatorImplementsNode(t *testing.T) {
	assert.Implements(t, (*parser.Node)(nil), &BinaryOperator{})
}

func TestBinaryFactoryBase(t *testing.T) {
	s := &mockState{}
	op := &lexer.Token{}
	l := &mockNode{}
	l.On("Location").Return(nil)
	r := &mockNode{}
	r.On("Location").Return(nil)

	result, err := BinaryFactory(s, l, r, op)

	assert.NoError(t, err)
	assert.Equal(t, &BinaryOperator{
		Op: op,
		L:  l,
		R:  r,
	}, result)
	assert.Same(t, l, result.(*BinaryOperator).L)
	assert.Same(t, r, result.(*BinaryOperator).R)
}

func TestBinaryFactoryLocation(t *testing.T) {
	s := &mockState{}
	finalLoc := &mockLocation{}
	lLoc := &mockLocation{}
	rLoc := &mockLocation{}
	lLoc.On("ThruEnd", rLoc).Return(finalLoc, nil)
	op := &lexer.Token{}
	l := &mockNode{}
	l.On("Location").Return(lLoc)
	r := &mockNode{}
	r.On("Location").Return(rLoc)

	result, err := BinaryFactory(s, l, r, op)

	assert.NoError(t, err)
	assert.Equal(t, &BinaryOperator{
		Loc: finalLoc,
		Op:  op,
		L:   l,
		R:   r,
	}, result)
	assert.Same(t, finalLoc, result.(*BinaryOperator).Loc)
	assert.Same(t, l, result.(*BinaryOperator).L)
	assert.Same(t, r, result.(*BinaryOperator).R)
	lLoc.AssertExpectations(t)
	l.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestBinaryFactoryLocationError(t *testing.T) {
	s := &mockState{}
	lLoc := &mockLocation{}
	rLoc := &mockLocation{}
	lLoc.On("ThruEnd", rLoc).Return(nil, assert.AnError)
	op := &lexer.Token{}
	l := &mockNode{}
	l.On("Location").Return(lLoc)
	r := &mockNode{}
	r.On("Location").Return(rLoc)

	result, err := BinaryFactory(s, l, r, op)

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	lLoc.AssertExpectations(t)
	l.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestBinaryOperatorLocation(t *testing.T) {
	loc := &mockLocation{}
	obj := &BinaryOperator{
		Loc: loc,
	}

	result := obj.Location()

	assert.Same(t, loc, result)
}

func TestBinaryOperatorChildren(t *testing.T) {
	l := &mockNode{}
	l.On("dummy", "left")
	r := &mockNode{}
	r.On("dummy", "right")
	obj := &BinaryOperator{
		L: l,
		R: r,
	}

	result := obj.Children()

	assert.Equal(t, []parser.Node{
		NewAnnotatedNode(l, "L"),
		NewAnnotatedNode(r, "R"),
	}, result)
}

func TestBinaryOperatorString(t *testing.T) {
	obj := &BinaryOperator{
		Op: &lexer.Token{Type: "op"},
	}

	result := obj.String()

	assert.Equal(t, "<op> token", result)
}
