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

// FileOption is an option that may be passed to the NewFileScanner
// function.
type FileOption interface {
	// fileApply applies the option to FileScanner.
	fileApply(s *FileScanner)
}

// argOptions is a set of options for NewArgumentScanner.
type argOptions struct {
	joiner string       // The string to use to join arguments
	opts   []FileOption // Options for the scanners
}

// ArgOption is an option that may be passed to the NewArgumentScanner
// function.
type ArgOption interface {
	// argApply applies the option to the ArgumentScanner.
	argApply(o *argOptions)
}

// lineEndings is the type that stores the line ending style to use.
type lineEndings struct {
	ls LineStyle // The line ending style to use
}

// fileApply applies the option to FileScanner.
func (o lineEndings) fileApply(s *FileScanner) {
	s.ls = o.ls
}

// LineEndings is a file scanner option that may be used to set the
// preferred line ending style.  A line ending style is an instance of
// LineStyle that controls how the scanner recognizes newlines.  The
// scanner always converts line endings into single newlines.
func LineEndings(ls LineStyle) FileOption {
	return lineEndings{ls: ls}
}

// TabStop is a file scanner option that specifies the tab stop to
// apply.  The default tab stop is 8.
type TabStop int

// fileApply applies the option to FileScanner.
func (o TabStop) fileApply(s *FileScanner) {
	s.ts = int(o)
}

// EncodingErrorOption is the type that stores the encoding error
// handler that the file scanner should use.
type EncodingErrorOption struct {
	enc EncodingErrorHandler // The encoding error handler to use
}

// fileApply applies the option to FileScanner.
func (eeo EncodingErrorOption) fileApply(s *FileScanner) {
	s.enc = eeo.enc
}

// argApply applies the option to argOptions.
func (eeo EncodingErrorOption) argApply(o *argOptions) {
	o.opts = append(o.opts, eeo)
}

// EncodingError is an option that may be passed to either
// NewFileScanner or NewArgumentScanner.  It is used to specify an
// EncodingErrorHandler to use to handle encoding errors.
func EncodingError(enc EncodingErrorHandler) EncodingErrorOption {
	return EncodingErrorOption{enc: enc}
}

// ArgJoiner is an argument option that specifies the string that
// should logically be expected between each argument.  By default,
// this is a single space (" ").
type ArgJoiner string

// argApply applies the option to argOptions.
func (opt ArgJoiner) argApply(o *argOptions) {
	o.joiner = string(opt)
}
