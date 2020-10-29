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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockState struct {
	mock.Mock
}

func (m *mockState) Table() Table {
	args := m.MethodCalled("Table")

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Table)
	}

	return nil
}

func TestBaseStateImplementsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &BaseState{})
}

func TestBaseStateTable(t *testing.T) {
	obj := &BaseState{
		Tab: Table{},
	}

	result := obj.Table()

	assert.Equal(t, Table{}, result)
}
