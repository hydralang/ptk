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

type mockLineStyle struct {
	mock.Mock
}

func (m *mockLineStyle) Handle(chs []rune) (LineDis, LineStyle) {
	args := m.MethodCalled("Handle", chs)

	if tmp := args.Get(1); tmp != nil {
		return args.Get(0).(LineDis), tmp.(LineStyle)
	}

	return args.Get(0).(LineDis), nil
}

func TestUNIXLineStyleImplementsLineStyle(t *testing.T) {
	assert.Implements(t, (*LineStyle)(nil), UNIXLineStyle)
}

func TestUNIXLineStyleHandleCR(t *testing.T) {
	dis, next := UNIXLineStyle.Handle([]rune{'\r'})

	assert.Equal(t, LineDisSpace, dis)
	assert.Same(t, UNIXLineStyle, next)
}

func TestUNIXLineStyleHandleNL(t *testing.T) {
	dis, next := UNIXLineStyle.Handle([]rune{'\n'})

	assert.Equal(t, LineDisNewline, dis)
	assert.Same(t, UNIXLineStyle, next)
}

func TestMacLineStyleImplementsLineStyle(t *testing.T) {
	assert.Implements(t, (*LineStyle)(nil), MacLineStyle)
}

func TestMacLineStyleHandleCR(t *testing.T) {
	dis, next := MacLineStyle.Handle([]rune{'\r'})

	assert.Equal(t, LineDisNewline, dis)
	assert.Same(t, MacLineStyle, next)
}

func TestMacLineStyleHandleNL(t *testing.T) {
	dis, next := MacLineStyle.Handle([]rune{'\n'})

	assert.Equal(t, LineDisSpace, dis)
	assert.Same(t, MacLineStyle, next)
}

func TestDOSLineStyleImplementsLineStyle(t *testing.T) {
	assert.Implements(t, (*LineStyle)(nil), DOSLineStyle)
}

func TestDOSLineStyleHandleCR(t *testing.T) {
	dis, next := DOSLineStyle.Handle([]rune{'\r'})

	assert.Equal(t, LineDisMore, dis)
	assert.Same(t, DOSLineStyle, next)
}

func TestDOSLineStyleHandleNL(t *testing.T) {
	dis, next := DOSLineStyle.Handle([]rune{'\n'})

	assert.Equal(t, LineDisNewline, dis)
	assert.Same(t, DOSLineStyle, next)
}

func TestDOSLineStyleHandleCRCR(t *testing.T) {
	dis, next := DOSLineStyle.Handle([]rune{'\r', '\r'})

	assert.Equal(t, LineDisSpace, dis)
	assert.Same(t, DOSLineStyle, next)
}

func TestDOSLineStyleHandleCRNL(t *testing.T) {
	dis, next := DOSLineStyle.Handle([]rune{'\r', '\n'})

	assert.Equal(t, LineDisNewline, dis)
	assert.Same(t, DOSLineStyle, next)
}

func TestDOSLineStyleHandleCREOF(t *testing.T) {
	dis, next := DOSLineStyle.Handle([]rune{'\r', EOF})

	assert.Equal(t, LineDisSpace, dis)
	assert.Same(t, DOSLineStyle, next)
}

func TestUnknownLineStyleImplementsLineStyle(t *testing.T) {
	assert.Implements(t, (*LineStyle)(nil), UnknownLineStyle)
}

func TestUnknownLineStyleHandleCR(t *testing.T) {
	dis, next := UnknownLineStyle.Handle([]rune{'\r'})

	assert.Equal(t, LineDisMore, dis)
	assert.Same(t, UnknownLineStyle, next)
}

func TestUnknownLineStyleHandleNL(t *testing.T) {
	dis, next := UnknownLineStyle.Handle([]rune{'\n'})

	assert.Equal(t, LineDisNewline, dis)
	assert.Same(t, UNIXLineStyle, next)
}

func TestUnknownLineStyleHandleCRCR(t *testing.T) {
	dis, next := UnknownLineStyle.Handle([]rune{'\r', '\r'})

	assert.Equal(t, LineDisNewlineSave, dis)
	assert.Same(t, MacLineStyle, next)
}

func TestUnknownLineStyleHandleCRNL(t *testing.T) {
	dis, next := UnknownLineStyle.Handle([]rune{'\r', '\n'})

	assert.Equal(t, LineDisNewline, dis)
	assert.Same(t, DOSLineStyle, next)
}

func TestUnknownLineStyleHandleCREOF(t *testing.T) {
	dis, next := UnknownLineStyle.Handle([]rune{'\r', EOF})

	assert.Equal(t, LineDisNewlineSave, dis)
	assert.Same(t, MacLineStyle, next)
}

func TestNoLineStyleImplementsLineStyle(t *testing.T) {
	assert.Implements(t, (*LineStyle)(nil), NoLineStyle)
}

func TestNoLineStyleHandleCR(t *testing.T) {
	dis, next := NoLineStyle.Handle([]rune{'\r'})

	assert.Equal(t, LineDisSpace, dis)
	assert.Same(t, NoLineStyle, next)
}

func TestNoLineStyleHandleNL(t *testing.T) {
	dis, next := NoLineStyle.Handle([]rune{'\n'})

	assert.Equal(t, LineDisSpace, dis)
	assert.Same(t, NoLineStyle, next)
}
