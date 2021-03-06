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
	"unicode"
)

// Special runes used to signal end-of-file or other conditions in the
// scanners.
const (
	EOF      rune = unicode.MaxRune + iota + 1 // Signals end-of-file
	errRune                                    // Signals an error
	sentinel                                   // Signals saved is unset
)

// Char represents a character retrieved from the source input stream.
// It bundles together a rune and a location.
type Char struct {
	Rune rune     // The rune read from the source
	Loc  Location // The location of the rune within the stream
}
