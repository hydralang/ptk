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

package common

import (
	"github.com/stretchr/testify/mock"
)

// Location is an interface for location data.  Each token and node
// should have attached location data that reports its location.  This
// aids in finding the location of errors.
type Location interface {
	// String constructs a string representation of the location.
	String() string

	// Thru creates a new Location that ranges from the beginning
	// of this location to the beginning of another Location.
	Thru(other Location) (Location, error)

	// ThruEnd is similar to Thru, except that it creates a new
	// Location that ranges from the beginning of this location to
	// the ending of another location.
	ThruEnd(other Location) (Location, error)
}

// MockLocation is a mock for Location
type MockLocation struct {
	mock.Mock
}

// String constructs a string representation of the location.
func (m *MockLocation) String() string {
	args := m.MethodCalled("String")

	return args.String(0)
}

// Thru creates a new Location that ranges from the beginning of this
// location to the beginning of another Location.
func (m *MockLocation) Thru(other Location) (Location, error) {
	args := m.MethodCalled("Thru", other)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Location), args.Error(1)
	}

	return nil, args.Error(1)
}

// ThruEnd is similar to Thru, except that it creates a new Location
// that ranges from the beginning of this location to the ending of
// another location.
func (m *MockLocation) ThruEnd(other Location) (Location, error) {
	args := m.MethodCalled("ThruEnd", other)

	if tmp := args.Get(0); tmp != nil {
		return tmp.(Location), args.Error(1)
	}

	return nil, args.Error(1)
}
