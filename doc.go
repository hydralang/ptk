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

// Package ptk is a toolkit for building parsers using the Vaughan
// Pratt top-down recursive-descent technique.  This technique is
// based on the concept of creating the functions of recursive descent
// around tokens, rather than around productions, and yields a more
// efficient parser.  It's particularly effective for parsing
// expressions, where operator precedence is important; with a
// straight recursive-descent technique, a different function has to
// be created (and recursed through) for each class of operators,
// whereas the same parser in the Pratt construction requires
// essentially a single function.
//
// The parser framework itself lives in the parser subdirectory.  The
// Pratt technique is a table-driven technique, and developers
// utilizing this toolkit need to create a Table, keyed by the token
// type (a simple string), and mapping to an Entry.  Entries have a
// "binding power" (assigned to the Power field of Entry), where
// higher values result in tighter binding--so, for instance, the "*"
// multiplication operator binds tighter than "+", so should have a
// higher binding power.  In addition, an entry has three optional
// functions: ExprFirst (assigned to the First field); ExprNext
// (assigned to the Next field); and Statement (assigned to the Stmt
// field).  The semantics of these functions will be described below.
//
// The user of a parser built using this toolkit will create a new
// Parser object by calling parser.New and passing it a Table.  An
// expression may then be parsed by calling the Parser.Expression
// method, and passing it a TokenStream.  The Parser.Expression method
// constructs a parser.State object with the Table and the
// TokenStream, and then proceeds to call the State.Expression method,
// which performs the heavy lifting.  That method begins by looking up
// the Entry corresponding to the first token; calling the First
// function defined there, passing it the parser, state, power, and
// token; then running the Next function on subsequent tokens until
// reading a token whose binding power is less than the power the
// specific instance of State.Expression was called with.  Typically,
// the First function would get called on a literal token, which it
// would simply return as-is; the Next function would then typically
// call the State.Expression function recursively with the operator's
// binding power (which is passed to Next) to retrieve the right-hand
// side expression.
//
// For parsing statements, the Parser.Statement method would be used,
// again passing it a TokenStream.  The Parser.Statement method is
// much simpler, simply calling the Stmt function defined in the Entry
// associated with the first token's token type.  Subsequent calls
// would extract subsequent statements, so a Parser.Statements
// function (added "s") can be used to retrieve a list of statements.
//
// The Token and TokenStream types are defined in the common
// subdirectory.  A Token is a simple structure, containing the token
// type (a string), and three optional fields: the Location of the
// token, any Value associated with the token, and the original Text
// of the token.  The application's lexical analyzer should produce
// instances of this struct and communicate them to the parser through
// a TokenStream instance; TokenStream is an interface with a Next
// method, which returns subsequent Token instances.  (Several
// implementations of TokenStream are provided in the utils package,
// but the lexical analyzer could also be designed to implement the
// TokenStream interface.)
//
// The Location type is also defined in common, as is LocationError.
// The former is an interface with three methods, the most important
// being String; the latter allows wrapping an error with the location
// of the error.  There is also an implementation of Location,
// FileLocation, in the utils package; this Location implementation
// should be suitable for most uses.
//
// The final important type in common is the Node type, which is the
// return type for all the major functions involved in the parser
// framework.  Typically, the Node would be a set of data structures
// which form an abstract syntax tree, but the use of an interface
// here should allow alternatives depending on the application.  An
// AnnotatedNode is also available in the utils package; this allows
// wrapping arbitrary Node instances with additional annotation, which
// may be useful for the visualization utilities.
//
// The utils directory contains utilities that may be useful to
// developers working with ptk.  Besides the implementations mentioned
// above, it contains UnaryOperator and BinaryOperator, which are
// implementations of Node that can be used with unary and binary
// operators, as the name implies.  It also provides the Literal
// function, which may be used for literal tokens and can be assigned
// to Entry.First; the Prefix function, which generates a closure
// assignable to Entry.First for tokens representing unary operators;
// and the Infix and InfixR functions, which generate closures
// assignable to Entry.Next for tokens representing left-associative
// and right-associative binary operators, respectively.
//
// The utils directory also contains the Visualize function, which can
// be used to create a string dump of an abstract syntax tree; the
// exact look of the string is controlled by the Profile, with ASCII,
// Rounded, and Square available.
package ptk
