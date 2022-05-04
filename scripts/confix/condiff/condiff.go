package main

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

// It performs a keyspace diff on two TOML documents.

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/creachadair/tomledit"
	"github.com/creachadair/tomledit/parser"
	"github.com/creachadair/tomledit/transform"
)

var (
	doDesnake = flag.Bool("desnake", false, "Convert snake_case to kebab-case before comparing")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %[1]s [options] f1 f2

Diff the keyspaces of the TOML documents in files f1 and f2.
The output prints one line per key that differs:

   -S name    -- section exists in f1 but not f2
   +S name    -- section exists in f2 but not f1
   -M name    -- mapping exists in f1 but not f2
   +M name    -- mapping exists in f2 but not f1

Comments, order, and values are ignored for comparison purposes.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatalf("Usage: %[1]s <lhs> <rhs>", filepath.Base(os.Args[0]))
	}
	lhs := mustParse(flag.Arg(0))
	rhs := mustParse(flag.Arg(1))
	if *doDesnake {
		log.Printf("Converting all names from snake_case to kebab-case")
		fix := transform.SnakeToKebab()
		_ = fix(context.Background(), lhs)
		_ = fix(context.Background(), rhs)
	}
	diffDocs(os.Stdout, lhs, rhs)
}

func mustParse(path string) *tomledit.Document {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Opening TOML input: %v", err)
	}
	defer f.Close()
	doc, err := tomledit.Parse(f)
	if err != nil {
		log.Fatalf("Parsing %q: %v", path, err)
	}
	return doc
}

func allKeys(s *tomledit.Section) []string {
	var keys []string
	s.Scan(func(key parser.Key, _ *tomledit.Entry) bool {
		keys = append(keys, key.String())
		return true
	})
	return keys
}

const (
	delSection = "-S"
	delMapping = "-M"
	addSection = "+S"
	addMapping = "+M"

	delMapSep = "\n" + delMapping + " "
	addMapSep = "\n" + addMapping + " "
)

func diffDocs(w io.Writer, lhs, rhs *tomledit.Document) {
	diffSections(w, lhs.Global, rhs.Global)
	lsec, rsec := lhs.Sections, rhs.Sections
	transform.SortSectionsByName(lsec)
	transform.SortSectionsByName(rsec)

	i, j := 0, 0
	for i < len(lsec) && j < len(rsec) {
		if lsec[i].Name.Before(rsec[j].Name) {
			fmt.Fprintln(w, delSection, lsec[i].Name)
			fmt.Fprintln(w, delMapping, strings.Join(allKeys(lsec[i]), delMapSep))
			i++
		} else if rsec[j].Name.Before(lsec[i].Name) {
			fmt.Fprintln(w, addSection, rsec[j].Name)
			fmt.Fprintln(w, addMapping, strings.Join(allKeys(rsec[j]), addMapSep))
			j++
		} else {
			diffSections(w, lsec[i], rsec[j])
			i++
			j++
		}
	}
	for ; i < len(lsec); i++ {
		fmt.Fprintln(w, delSection, lsec[i].Name)
		fmt.Fprintln(w, delMapping, strings.Join(allKeys(lsec[i]), delMapSep))
	}
	for ; j < len(rsec); j++ {
		fmt.Fprintln(w, addSection, rsec[j].Name)
		fmt.Fprintln(w, addMapping, strings.Join(allKeys(rsec[j]), addMapSep))
	}
}

func diffSections(w io.Writer, lhs, rhs *tomledit.Section) {
	diffKeys(w, allKeys(lhs), allKeys(rhs))
}

func diffKeys(w io.Writer, lhs, rhs []string) {
	sort.Strings(lhs)
	sort.Strings(rhs)

	i, j := 0, 0
	for i < len(lhs) && j < len(rhs) {
		if lhs[i] < rhs[j] {
			fmt.Fprintln(w, delMapping, lhs[i])
			i++
		} else if lhs[i] > rhs[j] {
			fmt.Fprintln(w, addMapping, rhs[j])
			j++
		} else {
			i++
			j++
		}
	}
	for ; i < len(lhs); i++ {
		fmt.Fprintln(w, delMapping, lhs[i])
	}
	for ; j < len(rhs); j++ {
		fmt.Fprintln(w, addMapping, rhs[j])
	}
}
