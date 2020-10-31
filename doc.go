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

// Package ptk is a toolkit for creating language parsers using a
// top-down recursive descent technique first described by Vaughan
// Pratt.  This technique is based on the concept of creating the
// functions around token types, rather than around productions, and
// yields a more efficient recursive descent parser that doesn't
// require recursing as deeply.  The Pratt technique is particularly
// effective for parsing expressions, where operator precedence is
// important; instead of a different function for each operator class,
// the Pratt technique enables handling the situation with a single
// function.
//
// Creating a parser using ptk involves three main components--a
// scanner, a lexer, and a parser--each represented by a different
// subpackage.  The scanner takes a source representation and breaks
// it down into single characters that are paired their location
// within the source.  The lexer then takes these characters and,
// using user-provided classifiers and recognizers, groups them into
// tokens; a token is comprised of a token type, a location, and an
// optional semantic meaning (for instance, the characters "12345"
// would likely have a semantic meaning consisting of the numerical
// value represented by that character sequence).  Finally, the parser
// takes a sequence of tokens and, with a user-provided parse table,
// constructs an abstract syntax tree.  This AST is the final result
// of the parsing, and the caller may process it in any necessary
// fashion.
//
// The scanner subpackage contains two main scanners--a FileScanner,
// which actually takes an io.Reader instance, and an argument
// scanner, which takes a list of character strings and treats them as
// a single source separated by an optional separator that defaults to
// whitespace.  In addition, there are several composite scanners,
// such as the "memoizing" scanner and the chaining scanner; the
// argument scanner is actually built by composing these composite
// scanners around the arguments.  A ListScanner is also provided to
// simplify testing code.  This subpackage also contains the Location
// interface and types, which represent a range of the source
// containing a single character or group of characters.
//
// The lexer subpackage contains a Lexer type, which is initialized by
// passing in a user-created scanner and a lexer.State; the state is
// used by the lexer to find a Classifier.  The lexer calls the
// Classifier.Classify method to get a list of Recognizer instances
// that may be able to handle the input, which the lexer then calls in
// turn until one of the recognizers successfully recognizes the
// input.  Recognizers call the Lexer.Push method to register any
// tokens they extract from the input; but a recognizer need not push
// any tokens.  (For instance, in most languages, whitespace is simply
// discarded, so a recognizer for whitespace would not produce any
// tokens.)  A lexer.BaseState type is provided that has the basic
// functionality, but applications are free to extend or reimplement
// the functionality.
//
// The parser subpackage contains a Parser type, which is initialized
// by passing in a user-created lexer and a parser.State; the state is
// used by the parser to find a Table, which maps token types produced
// by the lexer to an Entry, which itself contains an integer binding
// Power, along with three callback functions related to three
// contexts in which a token may appear.  As with the lexer, a default
// simple implementation of parser.State is provided, named
// parser.BaseState; applications may extend or reimplement the
// functionality as needed.  The parser subpackage also contains some
// utilities for building an AST; specifically, the UnaryOperator and
// BinaryOperator types are implementations of parser.Node, and the
// Prefix, Infix, and InfixR functions can be used to construct the
// callback functions for some common casses, while Literal is a
// parser.ExprFirst function for literals.
//
// A fourth subpackage is provided: the parser/visualize subpackage
// contains the Visualize function, which allows rendering an abstract
// syntax tree into a visual representation.  This may be used while
// developing a parser to ensure that it is producing the correct AST
// for a given input.
package ptk
