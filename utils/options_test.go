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
	"testing"

	"github.com/klmitch/kent"
	"github.com/stretchr/testify/assert"
)

func TestLineEndings(t *testing.T) {
	ls := &MockLineStyle{}
	s := &scanner{}

	opt := LineEndings(ls)
	opt(s)

	assert.Same(t, ls, s.ls)
}

func TestTabStop(t *testing.T) {
	s := &scanner{}

	opt := TabStop(42)
	opt(s)

	assert.Equal(t, 42, s.ts)
}

func TestReporter(t *testing.T) {
	rep := &kent.MockReporter{}
	s := &scanner{}

	opt := Reporter(rep)
	opt(s)

	assert.Same(t, rep, s.rep)
}
