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
	"bytes"

	"github.com/klmitch/kent"

	"github.com/hydralang/ptk/common"
)

// argOptions is a set of options for NewArgumentCharStream.
type argOptions struct {
	joiner string          // The string to use to join arguments
	opts   []ScannerOption // Options for the scanners
}

// ArgumentOption is an option that may be passed to the
// NewArgumentCharStream function.
type ArgumentOption func(*argOptions)

// ArgumentJoiner specifies an alternate argument joiner; this
// character sequence appears between each argument.  The default is a
// single space; pass this as an option to specify a different joiner.
func ArgumentJoiner(joiner string) ArgumentOption {
	return func(o *argOptions) {
		o.joiner = joiner
	}
}

// ArgumentReporter specifies a kent.Reporter to use for reporting
// encoding errors.  By default, encoding errors are not reported.
func ArgumentReporter(rep kent.Reporter) ArgumentOption {
	return func(o *argOptions) {
		o.opts = append(o.opts, Reporter(rep))
	}
}

// NewArgumentCharStream constructs and returns a common.Charstream
// implementation that returns characters drawn from a provided list
// of argument strings.  This is intended for use with arguments taken
// from the command line, but could be useful in other contexts as
// well.  The strings are logically joined by spaces; to use a
// different joiner, pass that as an option.
func NewArgumentCharStream(args []string, options ...ArgumentOption) common.CharStream {
	// Process the options
	opts := &argOptions{
		joiner: " ",
		opts:   []ScannerOption{LineEndings(NoLineStyle)},
	}
	for _, opt := range options {
		opt(opts)
	}

	// Construct the joiner character stream
	loc := ArgLocation{
		B: ArgPos{0, 1},
		E: ArgPos{0, 1},
	}
	joiner := NewMemoizingCharStream(NewScanner(bytes.NewBufferString(opts.joiner), loc, opts.opts...))

	// Construct a list of character streams
	streams := []common.CharStream{}
	for i, arg := range args {
		if i != 0 {
			streams = append(streams, joiner)
		}

		loc := ArgLocation{
			B: ArgPos{i + 1, 1},
			E: ArgPos{i + 1, 1},
		}
		streams = append(streams, NewScanner(bytes.NewBufferString(arg), loc, opts.opts...))
	}

	return NewChainingCharStream(streams)
}
