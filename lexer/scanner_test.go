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

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/klmitch/kent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hydralang/ptk/common"
)

func TestScannerImplementsCharStream(t *testing.T) {
	assert.Implements(t, (*CharStream)(nil), &scanner{})
}

func TestNewScanner(t *testing.T) {
	src := &bytes.Buffer{}
	loc := &common.MockLocation{}
	var opt1Called *scanner
	var opt2Called *scanner
	options := []ScannerOption{
		func(s *scanner) {
			opt1Called = s
		},
		func(s *scanner) {
			opt2Called = s
		},
	}

	result := NewScanner(src, loc, options...)

	assert.Equal(t, &scanner{
		src:   src,
		ts:    DefaultTabStop,
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ls:    UnknownLineStyle,
		saved: sentinel,
		loc:   loc,
		rep:   kent.Root(),
	}, result)
	assert.Same(t, result, opt1Called)
	assert.Same(t, result, opt2Called)
}

func TestScannerNextCharBufferedASCII(t *testing.T) {
	obj := &scanner{
		src: &bytes.Buffer{},
		buf: [scanBuf + 1]byte{'b', 'u', 'f', utf8.RuneSelf},
		end: 3,
	}

	c := obj.next()

	assert.Equal(t, 'b', c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{'b', 'u', 'f', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, 3, obj.end)
	assert.Nil(t, obj.err)
}

func TestScannerNextCharBufferedMultibyte(t *testing.T) {
	obj := &scanner{
		src: &bytes.Buffer{},
		buf: [scanBuf + 1]byte{195, 177, 'i', 'n', 'o', utf8.RuneSelf},
		end: 5,
	}

	c := obj.next()

	assert.Equal(t, '\xf1', c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{195, 177, 'i', 'n', 'o', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 2, obj.pos)
	assert.Equal(t, 5, obj.end)
	assert.Nil(t, obj.err)
}

func TestScannerNextCharBufferedBadChar(t *testing.T) {
	chrLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", utf8.RuneError, DefaultTabStop).Return(chrLoc)
	rep := &kent.MockReporter{}
	rep.On("Report", common.LocationError(chrLoc, ErrBadEncoding))
	obj := &scanner{
		src: &bytes.Buffer{},
		buf: [scanBuf + 1]byte{128, 'b', 'u', 'f', utf8.RuneSelf},
		end: 4,
		ts:  DefaultTabStop,
		loc: loc,
		rep: rep,
	}

	c := obj.next()

	assert.Equal(t, utf8.RuneError, c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{128, 'b', 'u', 'f', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, 4, obj.end)
	assert.Nil(t, obj.err)
	loc.AssertExpectations(t)
	rep.AssertExpectations(t)
}

func TestScannerNextCharEOF(t *testing.T) {
	obj := &scanner{
		buf: [scanBuf + 1]byte{utf8.RuneSelf},
	}

	c := obj.next()

	assert.Equal(t, common.EOF, c)
	assert.Nil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 0, obj.pos)
	assert.Equal(t, 0, obj.end)
	assert.Nil(t, obj.err)
}

func TestScannerNextCharError(t *testing.T) {
	obj := &scanner{
		buf: [scanBuf + 1]byte{utf8.RuneSelf},
		err: assert.AnError,
	}

	c := obj.next()

	assert.Equal(t, errRune, c)
	assert.Nil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 0, obj.pos)
	assert.Equal(t, 0, obj.end)
	assert.Same(t, assert.AnError, obj.err)
}

func TestScannerNextCharFromEmpty(t *testing.T) {
	obj := &scanner{
		src: bytes.NewBufferString("test"),
		buf: [scanBuf + 1]byte{utf8.RuneSelf},
	}

	c := obj.next()

	assert.Equal(t, 't', c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{'t', 'e', 's', 't', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, 4, obj.end)
	assert.Nil(t, obj.err)
}

func TestScannerNextCharNotEmpty(t *testing.T) {
	obj := &scanner{
		src: bytes.NewBufferString("test"),
		buf: [scanBuf + 1]byte{'b', 'u', 'f', utf8.RuneSelf},
		pos: 3,
		end: 3,
	}

	c := obj.next()

	assert.Equal(t, 't', c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{'t', 'e', 's', 't', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, 4, obj.end)
	assert.Nil(t, obj.err)
}

func TestScannerNextCharSplitMultibyte(t *testing.T) {
	obj := &scanner{
		src: bytes.NewBuffer([]byte{177, 'i', 'n', 'o'}),
		buf: [scanBuf + 1]byte{'b', 'u', 'f', 195, utf8.RuneSelf},
		pos: 3,
		end: 4,
	}

	c := obj.next()

	assert.Equal(t, '\xf1', c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{195, 177, 'i', 'n', 'o', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 2, obj.pos)
	assert.Equal(t, 5, obj.end)
	assert.Nil(t, obj.err)
}

func TestScannerNextCharReadEOF(t *testing.T) {
	obj := &scanner{
		src: &bytes.Buffer{},
		buf: [scanBuf + 1]byte{'b', 'u', 'f', utf8.RuneSelf},
		pos: 3,
		end: 3,
	}

	c := obj.next()

	assert.Equal(t, common.EOF, c)
	assert.Nil(t, obj.src)
	assert.Nil(t, obj.err)
}

type mockReader struct {
	mock.Mock
}

func (r *mockReader) Read(b []byte) (int, error) {
	args := r.MethodCalled("Read")

	data := args.Get(0).([]byte)
	copy(b, data)

	return args.Int(1), args.Error(2)
}

func TestScannerNextCharReadErrorDelayed(t *testing.T) {
	src := &mockReader{}
	src.On("Read").Return([]byte{'t', 'e', 's', 't'}, 4, assert.AnError)
	obj := &scanner{
		src: src,
		buf: [scanBuf + 1]byte{'b', 'u', 'f', utf8.RuneSelf},
		pos: 3,
		end: 3,
	}

	c := obj.next()

	assert.Equal(t, 't', c)
	assert.Nil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{'t', 'e', 's', 't', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, 4, obj.end)
	assert.Same(t, assert.AnError, obj.err)
}

func TestScannerNextCharReadErrorImmediate(t *testing.T) {
	src := &mockReader{}
	src.On("Read").Return([]byte{}, 0, assert.AnError)
	obj := &scanner{
		src: src,
		buf: [scanBuf + 1]byte{'b', 'u', 'f', utf8.RuneSelf},
		pos: 3,
		end: 3,
	}

	c := obj.next()

	assert.Equal(t, errRune, c)
	assert.Nil(t, obj.src)
	assert.Same(t, assert.AnError, obj.err)
}

func TestScannerNextBase(t *testing.T) {
	nextLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", 't', DefaultTabStop).Return(nextLoc)
	ls := &MockLineStyle{}
	obj := &scanner{
		src:   bytes.NewBufferString("test"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 't',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{'t', 'e', 's', 't', utf8.RuneSelf}, obj.buf)
	assert.Same(t, ls, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
}

func TestScannerNextSaved(t *testing.T) {
	nextLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", 's', DefaultTabStop).Return(nextLoc)
	ls := &MockLineStyle{}
	obj := &scanner{
		src:   bytes.NewBufferString("test"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: 's',
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: 's',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{utf8.RuneSelf}, obj.buf)
	assert.Same(t, ls, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
}

func TestScannerNextError(t *testing.T) {
	nextLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", common.EOF, DefaultTabStop).Return(nextLoc)
	ls := &MockLineStyle{}
	obj := &scanner{
		src:   bytes.NewBufferString("test"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		err:   assert.AnError,
		ts:    DefaultTabStop,
		ls:    ls,
		saved: errRune,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, common.Char{
		Rune: common.EOF,
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{utf8.RuneSelf}, obj.buf)
	assert.Same(t, ls, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
}

func TestScannerNextDisNewline(t *testing.T) {
	nextLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", '\n', DefaultTabStop).Return(nextLoc)
	newLs := &MockLineStyle{}
	ls := &MockLineStyle{}
	ls.On("Handle", []rune{'\n'}).Return(LineDisNewline, newLs)
	obj := &scanner{
		src:   bytes.NewBufferString("\n"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: '\n',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{'\n', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, 1, obj.end)
	assert.Same(t, newLs, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
	ls.AssertExpectations(t)
}

func TestScannerNextDisSpace(t *testing.T) {
	nextLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", ' ', DefaultTabStop).Return(nextLoc)
	newLs := &MockLineStyle{}
	ls := &MockLineStyle{}
	ls.On("Handle", []rune{'\n'}).Return(LineDisSpace, newLs)
	obj := &scanner{
		src:   bytes.NewBufferString("\n"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: ' ',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{'\n', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, 1, obj.end)
	assert.Same(t, newLs, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
	ls.AssertExpectations(t)
}

func TestScannerNextDisMoreNewline(t *testing.T) {
	nextLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", '\n', DefaultTabStop).Return(nextLoc)
	newLs := &MockLineStyle{}
	ls := &MockLineStyle{}
	ls.On("Handle", []rune{'\r'}).Return(LineDisMore, newLs)
	ls.On("Handle", []rune{'\r', '\r'}).Return(LineDisNewline, newLs)
	obj := &scanner{
		src:   bytes.NewBufferString("\r\r"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: '\n',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{'\r', '\r', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 2, obj.pos)
	assert.Equal(t, 2, obj.end)
	assert.Same(t, newLs, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
	ls.AssertExpectations(t)
}

func TestScannerNextDisMoreNewlineSave(t *testing.T) {
	nextLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", '\n', DefaultTabStop).Return(nextLoc)
	newLs := &MockLineStyle{}
	ls := &MockLineStyle{}
	ls.On("Handle", []rune{'\r'}).Return(LineDisMore, newLs)
	ls.On("Handle", []rune{'\r', '\r'}).Return(LineDisNewlineSave, newLs)
	obj := &scanner{
		src:   bytes.NewBufferString("\r\r"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: '\n',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{'\r', '\r', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 2, obj.pos)
	assert.Equal(t, 2, obj.end)
	assert.Same(t, newLs, obj.ls)
	assert.Equal(t, '\r', obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
	ls.AssertExpectations(t)
}

func TestScannerNextDisMoreSpace(t *testing.T) {
	nextLoc := &common.MockLocation{}
	loc := &common.MockLocation{}
	loc.On("Incr", ' ', DefaultTabStop).Return(nextLoc)
	newLs := &MockLineStyle{}
	ls := &MockLineStyle{}
	ls.On("Handle", []rune{'\r'}).Return(LineDisMore, newLs)
	ls.On("Handle", []rune{'\r', '\r'}).Return(LineDisSpace, newLs)
	obj := &scanner{
		src:   bytes.NewBufferString("\r\r"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, common.Char{
		Rune: ' ',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{'\r', '\r', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 2, obj.pos)
	assert.Equal(t, 2, obj.end)
	assert.Same(t, newLs, obj.ls)
	assert.Equal(t, '\r', obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
	ls.AssertExpectations(t)
}
