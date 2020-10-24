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

package scanner

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFileLocationImplementsLocation(t *testing.T) {
	assert.Implements(t, (*Location)(nil), &FileLocation{})
}

func TestFileLocationString0Columns(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 2},
	}

	result := loc.String()

	assert.Equal(t, "file:3:2", result)
}

func TestFileLocationString1Column(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.String()

	assert.Equal(t, "file:3:2", result)
}

func TestFileLocationString2Columns(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 4},
	}

	result := loc.String()

	assert.Equal(t, "file:3:2-4", result)
}

func TestFileLocationString2Lines(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{4, 2},
	}

	result := loc.String()

	assert.Equal(t, "file:3:2-4:2", result)
}

func TestFileLocationThruBase(t *testing.T) {
	loc1 := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}
	loc2 := FileLocation{
		File: "file",
		B:    FilePos{3, 5},
		E:    FilePos{3, 6},
	}

	result, err := loc1.Thru(loc2)

	assert.NoError(t, err)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 5},
	}, result)
}

func TestFileLocationThruSplitFile(t *testing.T) {
	loc1 := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}
	loc2 := FileLocation{
		File: "other",
		B:    FilePos{3, 5},
		E:    FilePos{3, 6},
	}

	result, err := loc1.Thru(loc2)

	assert.Same(t, ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestFileLocationThruSplitLocation(t *testing.T) {
	loc1 := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}
	loc2 := &mockLocation{}

	result, err := loc1.Thru(loc2)

	assert.Same(t, ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestFileLocationThruEndBase(t *testing.T) {
	loc1 := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}
	loc2 := FileLocation{
		File: "file",
		B:    FilePos{3, 5},
		E:    FilePos{3, 6},
	}

	result, err := loc1.ThruEnd(loc2)

	assert.NoError(t, err)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 6},
	}, result)
}

func TestFileLocationThruEndSplitFile(t *testing.T) {
	loc1 := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}
	loc2 := FileLocation{
		File: "other",
		B:    FilePos{3, 5},
		E:    FilePos{3, 6},
	}

	result, err := loc1.ThruEnd(loc2)

	assert.Same(t, ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestFileLocationThruEndSplitLocation(t *testing.T) {
	loc1 := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}
	loc2 := &mockLocation{}

	result, err := loc1.ThruEnd(loc2)

	assert.Same(t, ErrSplitLocation, err)
	assert.Nil(t, result)
}

func TestFileLocationAdvanceColumn(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.advance(FilePos{C: 2})

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 3},
		E:    FilePos{3, 5},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}, loc)
}

func TestFileLocationAdvanceLine(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.advance(FilePos{L: 1, C: 2})

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 3},
		E:    FilePos{4, 3},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}, loc)
}

func TestFileLocationIncrBase(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.Incr('c', 8)

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 3},
		E:    FilePos{3, 4},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}, loc)
}

func TestFileLocationIncrEOF(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.Incr(EOF, 8)

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 3},
		E:    FilePos{3, 3},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}, loc)
}

func TestFileLocationIncrNewline(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.Incr('\n', 8)

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 3},
		E:    FilePos{4, 1},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}, loc)
}

func TestFileLocationIncrTab8(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.Incr('\t', 8)

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 3},
		E:    FilePos{3, 9},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}, loc)
}

func TestFileLocationIncrTab4(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.Incr('\t', 4)

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 3},
		E:    FilePos{3, 5},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}, loc)
}

func TestFileLocationIncrFormFeedMidLine(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}

	result := loc.Incr('\f', 8)

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 3},
		E:    FilePos{3, 4},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 2},
		E:    FilePos{3, 3},
	}, loc)
}

func TestFileLocationIncrFormFeedBeginningOfLine(t *testing.T) {
	loc := FileLocation{
		File: "file",
		B:    FilePos{3, 1},
		E:    FilePos{3, 2},
	}

	result := loc.Incr('\f', 8)

	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 1},
		E:    FilePos{3, 2},
	}, result)
	assert.Equal(t, FileLocation{
		File: "file",
		B:    FilePos{3, 1},
		E:    FilePos{3, 2},
	}, loc)
}

func TestFileScannerImplementsScanner(t *testing.T) {
	assert.Implements(t, (*Scanner)(nil), &FileScanner{})
}

func TestNewFileScanner(t *testing.T) {
	src := &bytes.Buffer{}
	loc := &mockLocation{}
	opt1 := &mockFileOption{}
	opt1.On("fileApply", mock.Anything)
	opt2 := &mockFileOption{}
	opt2.On("fileApply", mock.Anything)

	result := NewFileScanner(src, loc, opt1, opt2)

	assert.Equal(t, &FileScanner{
		src:   src,
		ts:    DefaultTabStop,
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ls:    UnknownLineStyle,
		saved: sentinel,
		loc:   loc,
	}, result)
	opt1.AssertExpectations(t)
	opt2.AssertExpectations(t)
}

func TestFileScannerNextCharBufferedASCII(t *testing.T) {
	obj := &FileScanner{
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

func TestFileScannerNextCharBufferedMultibyte(t *testing.T) {
	obj := &FileScanner{
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

func TestFileScannerNextCharBufferedBadCharBase(t *testing.T) {
	chrLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", utf8.RuneError, DefaultTabStop).Return(chrLoc)
	obj := &FileScanner{
		src: &bytes.Buffer{},
		buf: [scanBuf + 1]byte{128, 'b', 'u', 'f', utf8.RuneSelf},
		end: 4,
		ts:  DefaultTabStop,
		loc: loc,
	}

	c := obj.next()

	assert.Equal(t, errRune, c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{128, 'b', 'u', 'f', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 0, obj.pos)
	assert.Equal(t, 4, obj.end)
	assert.Equal(t, LocationError(chrLoc, ErrBadEncoding), obj.err)
	loc.AssertExpectations(t)
}

func TestFileScannerNextCharBufferedBadCharHandled(t *testing.T) {
	chrLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", utf8.RuneError, DefaultTabStop).Return(chrLoc)
	handler := &mockEncodingErrorHandler{}
	handler.On("Handle", LocationError(chrLoc, ErrBadEncoding)).Return(nil)
	obj := &FileScanner{
		src: &bytes.Buffer{},
		buf: [scanBuf + 1]byte{128, 'b', 'u', 'f', utf8.RuneSelf},
		end: 4,
		ts:  DefaultTabStop,
		loc: loc,
		enc: handler,
	}

	c := obj.next()

	assert.Equal(t, utf8.RuneError, c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{128, 'b', 'u', 'f', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 1, obj.pos)
	assert.Equal(t, 4, obj.end)
	assert.Nil(t, obj.err)
	loc.AssertExpectations(t)
	handler.AssertExpectations(t)
}

func TestFileScannerNextCharBufferedBadCharChanged(t *testing.T) {
	chrLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", utf8.RuneError, DefaultTabStop).Return(chrLoc)
	handler := &mockEncodingErrorHandler{}
	handler.On("Handle", LocationError(chrLoc, ErrBadEncoding)).Return(assert.AnError)
	obj := &FileScanner{
		src: &bytes.Buffer{},
		buf: [scanBuf + 1]byte{128, 'b', 'u', 'f', utf8.RuneSelf},
		end: 4,
		ts:  DefaultTabStop,
		loc: loc,
		enc: handler,
	}

	c := obj.next()

	assert.Equal(t, errRune, c)
	assert.NotNil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{128, 'b', 'u', 'f', utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 0, obj.pos)
	assert.Equal(t, 4, obj.end)
	assert.Same(t, assert.AnError, obj.err)
	loc.AssertExpectations(t)
	handler.AssertExpectations(t)
}

func TestFileScannerNextCharEOF(t *testing.T) {
	obj := &FileScanner{
		buf: [scanBuf + 1]byte{utf8.RuneSelf},
	}

	c := obj.next()

	assert.Equal(t, EOF, c)
	assert.Nil(t, obj.src)
	assert.Equal(t, [scanBuf + 1]byte{utf8.RuneSelf}, obj.buf)
	assert.Equal(t, 0, obj.pos)
	assert.Equal(t, 0, obj.end)
	assert.Nil(t, obj.err)
}

func TestFileScannerNextCharError(t *testing.T) {
	obj := &FileScanner{
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

func TestFileScannerNextCharFromEmpty(t *testing.T) {
	obj := &FileScanner{
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

func TestFileScannerNextCharNotEmpty(t *testing.T) {
	obj := &FileScanner{
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

func TestFileScannerNextCharSplitMultibyte(t *testing.T) {
	obj := &FileScanner{
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

func TestFileScannerNextCharReadEOF(t *testing.T) {
	obj := &FileScanner{
		src: &bytes.Buffer{},
		buf: [scanBuf + 1]byte{'b', 'u', 'f', utf8.RuneSelf},
		pos: 3,
		end: 3,
	}

	c := obj.next()

	assert.Equal(t, EOF, c)
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

func TestFileScannerNextCharReadErrorDelayed(t *testing.T) {
	src := &mockReader{}
	src.On("Read").Return([]byte{'t', 'e', 's', 't'}, 4, assert.AnError)
	obj := &FileScanner{
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

func TestFileScannerNextCharReadErrorImmediate(t *testing.T) {
	src := &mockReader{}
	src.On("Read").Return([]byte{}, 0, assert.AnError)
	obj := &FileScanner{
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

func TestFileScannerNextBase(t *testing.T) {
	nextLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", 't', DefaultTabStop).Return(nextLoc)
	ls := &mockLineStyle{}
	obj := &FileScanner{
		src:   bytes.NewBufferString("test"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: 't',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{'t', 'e', 's', 't', utf8.RuneSelf}, obj.buf)
	assert.Same(t, ls, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
}

func TestFileScannerNextSaved(t *testing.T) {
	nextLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", 's', DefaultTabStop).Return(nextLoc)
	ls := &mockLineStyle{}
	obj := &FileScanner{
		src:   bytes.NewBufferString("test"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: 's',
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{
		Rune: 's',
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{utf8.RuneSelf}, obj.buf)
	assert.Same(t, ls, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
}

func TestFileScannerNextError(t *testing.T) {
	nextLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", EOF, DefaultTabStop).Return(nextLoc)
	ls := &mockLineStyle{}
	obj := &FileScanner{
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
	assert.Equal(t, Char{
		Rune: EOF,
		Loc:  nextLoc,
	}, result)
	assert.Equal(t, [scanBuf + 1]byte{utf8.RuneSelf}, obj.buf)
	assert.Same(t, ls, obj.ls)
	assert.Equal(t, sentinel, obj.saved)
	assert.Same(t, nextLoc, obj.loc)
	loc.AssertExpectations(t)
}

func TestFileScannerNextDisNewline(t *testing.T) {
	nextLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", '\n', DefaultTabStop).Return(nextLoc)
	newLs := &mockLineStyle{}
	ls := &mockLineStyle{}
	ls.On("Handle", []rune{'\n'}).Return(LineDisNewline, newLs)
	obj := &FileScanner{
		src:   bytes.NewBufferString("\n"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{
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

func TestFileScannerNextDisSpace(t *testing.T) {
	nextLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", ' ', DefaultTabStop).Return(nextLoc)
	newLs := &mockLineStyle{}
	ls := &mockLineStyle{}
	ls.On("Handle", []rune{'\n'}).Return(LineDisSpace, newLs)
	obj := &FileScanner{
		src:   bytes.NewBufferString("\n"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{
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

func TestFileScannerNextDisMoreNewline(t *testing.T) {
	nextLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", '\n', DefaultTabStop).Return(nextLoc)
	newLs := &mockLineStyle{}
	ls := &mockLineStyle{}
	ls.On("Handle", []rune{'\r'}).Return(LineDisMore, newLs)
	ls.On("Handle", []rune{'\r', '\r'}).Return(LineDisNewline, newLs)
	obj := &FileScanner{
		src:   bytes.NewBufferString("\r\r"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{
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

func TestFileScannerNextDisMoreNewlineSave(t *testing.T) {
	nextLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", '\n', DefaultTabStop).Return(nextLoc)
	newLs := &mockLineStyle{}
	ls := &mockLineStyle{}
	ls.On("Handle", []rune{'\r'}).Return(LineDisMore, newLs)
	ls.On("Handle", []rune{'\r', '\r'}).Return(LineDisNewlineSave, newLs)
	obj := &FileScanner{
		src:   bytes.NewBufferString("\r\r"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{
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

func TestFileScannerNextDisMoreSpace(t *testing.T) {
	nextLoc := &mockLocation{}
	loc := &mockLocation{}
	loc.On("Incr", ' ', DefaultTabStop).Return(nextLoc)
	newLs := &mockLineStyle{}
	ls := &mockLineStyle{}
	ls.On("Handle", []rune{'\r'}).Return(LineDisMore, newLs)
	ls.On("Handle", []rune{'\r', '\r'}).Return(LineDisSpace, newLs)
	obj := &FileScanner{
		src:   bytes.NewBufferString("\r\r"),
		buf:   [scanBuf + 1]byte{utf8.RuneSelf},
		ts:    DefaultTabStop,
		ls:    ls,
		saved: sentinel,
		loc:   loc,
	}

	result, err := obj.Next()

	assert.NoError(t, err)
	assert.Equal(t, Char{
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
