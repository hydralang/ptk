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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFileOption struct {
	mock.Mock
}

func (m *mockFileOption) fileApply(s *FileScanner) {
	m.MethodCalled("fileApply", s)
}

type mockArgOption struct {
	mock.Mock
}

func (m *mockArgOption) argApply(o *argOptions) {
	m.MethodCalled("argApply", o)
}

func TestLineEndingsImplementsFileOption(t *testing.T) {
	assert.Implements(t, (*FileOption)(nil), lineEndings{})
}

func TestLineEndingsFileApply(t *testing.T) {
	ls := &mockLineStyle{}
	s := &FileScanner{}
	obj := lineEndings{ls: ls}

	obj.fileApply(s)

	assert.Same(t, ls, s.ls)
}

func TestLineEndings(t *testing.T) {
	ls := &mockLineStyle{}

	result := LineEndings(ls)

	assert.Equal(t, lineEndings{ls: ls}, result)
}

func TestTabStopImplementsFileOption(t *testing.T) {
	assert.Implements(t, (*FileOption)(nil), TabStop(0))
}

func TestTabStopFileApply(t *testing.T) {
	s := &FileScanner{}
	obj := TabStop(42)

	obj.fileApply(s)

	assert.Equal(t, 42, s.ts)
}

func TestEncodingErrorOptionImplementsFileOption(t *testing.T) {
	assert.Implements(t, (*FileOption)(nil), EncodingErrorOption{})
}

func TestEncodingErrorOptionImplementsArgOption(t *testing.T) {
	assert.Implements(t, (*ArgOption)(nil), EncodingErrorOption{})
}

func TestEncodingErrorOptionFileApply(t *testing.T) {
	enc := &mockEncodingErrorHandler{}
	s := &FileScanner{}
	obj := EncodingErrorOption{enc: enc}

	obj.fileApply(s)

	assert.Same(t, enc, s.enc)
}

func TestEncodingErrorOptionArgApply(t *testing.T) {
	enc := &mockEncodingErrorHandler{}
	o := &argOptions{}
	obj := EncodingErrorOption{enc: enc}

	obj.argApply(o)

	assert.Equal(t, []FileOption{obj}, o.opts)
}

func TestEncodingError(t *testing.T) {
	enc := &mockEncodingErrorHandler{}

	result := EncodingError(enc)

	assert.Equal(t, EncodingErrorOption{enc: enc}, result)
}

func TestArgJoinerImplementsArgOption(t *testing.T) {
	assert.Implements(t, (*ArgOption)(nil), ArgJoiner(""))
}

func TestArgJoinerArgApply(t *testing.T) {
	o := &argOptions{}
	obj := ArgJoiner("|")

	obj.argApply(o)

	assert.Equal(t, &argOptions{
		joiner: "|",
	}, o)
}
