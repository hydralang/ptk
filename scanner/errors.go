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
	"errors"
	"fmt"
)

// Simple errors that may be generated within the package.
var (
	ErrSplitLocation = errors.New("Attempt to range file location through an incompatible location")
	ErrBadEncoding   = errors.New("Invalid UTF-8 encoding")
)

// EncodingErrorHandler is an interface for an encoding error handler.
// A scanner will call the Handle method of the error handler, which
// must return either nil or an error (which may be the same error).
type EncodingErrorHandler interface {
	// Handle handles the reported encoding error.  If it returns
	// non-nil, the scanner will report an error.
	Handle(e error) error
}

// locationError is an implementation of error that wraps an error and
// includes a location.
type locationError struct {
	loc Location // The location of the error
	err error    // The error that occurred there
}

// LocationError wraps an error and includes a location.
func LocationError(loc Location, err error) error {
	var le *locationError
	if loc == nil || errors.As(err, &le) {
		return err
	}

	return &locationError{
		loc: loc,
		err: err,
	}
}

// Error returns the error message for a locationError.  This
// implementation prefixes the error message with the location.
func (le *locationError) Error() string {
	return fmt.Sprintf("%s: %s", le.loc, le.err)
}

// Unwrap allows unwrapping the locationError to retrieve the
// underlying error.
func (le *locationError) Unwrap() error {
	return le.err
}

// LocationOf attempts to retrieve the location of an error.  If the
// location is not available, it returns nil.
func LocationOf(err error) Location {
	// Look for a locationError in the error chain
	var le *locationError
	if errors.As(err, &le) {
		return le.loc
	}

	return nil
}
