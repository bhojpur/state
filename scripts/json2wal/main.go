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

/*
	json2wal converts JSON file to binary WAL file.

	Usage:
			json2wal <path-to-JSON>  <path-to-wal>
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bhojpur/state/internal/consensus"
	"github.com/bhojpur/state/pkg/types"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "missing arguments: Usage:json2wal <path-to-JSON>  <path-to-wal>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(fmt.Errorf("failed to open WAL file: %w", err))
	}
	defer f.Close()

	walFile, err := os.OpenFile(os.Args[2], os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Errorf("failed to open WAL file: %w", err))
	}
	defer walFile.Close()

	// the length of bhojpur/wal/MsgInfo in the wal.json may exceed the defaultBufSize(4096) of bufio
	// because of the byte array in BlockPart
	// leading to unmarshal error: unexpected end of JSON input
	br := bufio.NewReaderSize(f, int(2*types.BlockPartSizeBytes))
	dec := consensus.NewWALEncoder(walFile)

	for {
		msgJSON, _, err := br.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(fmt.Errorf("failed to read file: %w", err))
		}
		// ignore the ENDHEIGHT in json.File
		if strings.HasPrefix(string(msgJSON), "ENDHEIGHT") {
			continue
		}

		var msg consensus.TimedWALMessage
		err = json.Unmarshal(msgJSON, &msg)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal json: %w", err))
		}

		err = dec.Encode(&msg)
		if err != nil {
			panic(fmt.Errorf("failed to encode msg: %w", err))
		}
	}
}
