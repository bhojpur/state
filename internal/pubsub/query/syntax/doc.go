package syntax

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// It defines a scanner and parser for the Bhojpur State event filter
// query language. A query selects events by their types and attribute values.
//
// Grammar
//
// The grammar of the query language is defined by the following EBNF:
//
//   query      = conditions EOF
//   conditions = condition {"AND" condition}
//   condition  = tag comparison
//   comparison = equal / order / contains / "EXISTS"
//   equal      = "=" (date / number / time / value)
//   order      = cmp (date / number / time)
//   contains   = "CONTAINS" value
//   cmp        = "<" / "<=" / ">" / ">="
//
// The lexical terms are defined here using RE2 regular expression notation:
//
//   // The name of an event attribute (type.value)
//   tag    = #'\w+(\.\w+)*'
//
//   // A datestamp (YYYY-MM-DD)
//   date   = #'DATE \d{4}-\d{2}-\d{2}'
//
//   // A number with optional fractional parts (0, 10, 3.25)
//   number = #'\d+(\.\d+)?'
//
//   // An RFC3339 timestamp (2021-11-23T22:04:19-09:00)
//   time   = #'TIME \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}([-+]\d{2}:\d{2}|Z)'
//
//   // A quoted literal string value ('a b c')
//   value  = #'\'[^\']*\''
//
