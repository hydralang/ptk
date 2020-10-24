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

// Package scanner is a subpackage of ptk that contains various
// implementations of the Scanner, along with related support types
// and code such as Char and Location.  A scanner is an object that
// can be queried for the next character in a character stream, such
// as a file or a list of arguments; its return value is a Char, which
// packages together a Unicode rune and a location that may be used
// for error reporting.
package scanner
