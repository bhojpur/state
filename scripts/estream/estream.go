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

// It is a manual testing tool for polling the event stream of a running
// Bhojpur State consensus node.

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/bhojpur/state/pkg/rpc/client/eventstream"
	rpcclient "github.com/bhojpur/state/pkg/rpc/client/http"
	"github.com/bhojpur/state/pkg/rpc/coretypes"
)

var (
	query      = flag.String("query", "", "Filter query")
	batchSize  = flag.Int("batch", 0, "Batch size")
	resumeFrom = flag.String("resume", "", "Resume cursor")
	numItems   = flag.Int("count", 0, "Number of items to read (0 to stream)")
	waitTime   = flag.Duration("poll", 0, "Long poll interval")
	rpcAddr    = flag.String("addr", "http://localhost:26657", "RPC service address")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %[1]s [options]

Connect to the Bhojpur State node whose RPC service is at -addr, and poll for events
matching the specified -query. If no query is given, all events are fetched.
The resulting event data are written to stdout as JSON.

Use -resume to pick up polling from a previously-reported event cursor.
Use -count to stop polling after a certain number of events has been reported.
Use -batch to override the default request batch size.
Use -poll to override the default long-polling interval.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	cli, err := rpcclient.New(*rpcAddr)
	if err != nil {
		log.Fatalf("RPC client: %v", err)
	}
	stream := eventstream.New(cli, *query, &eventstream.StreamOptions{
		BatchSize:  *batchSize,
		ResumeFrom: *resumeFrom,
		WaitTime:   *waitTime,
	})

	// Shut down cleanly on SIGINT.  Don't attempt clean shutdown for other
	// fatal signals.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var nr int
	if err := stream.Run(ctx, func(itm *coretypes.EventItem) error {
		nr++
		bits, err := json.Marshal(itm)
		if err != nil {
			return err
		}
		fmt.Println(string(bits))
		if *numItems > 0 && nr >= *numItems {
			return eventstream.ErrStopRunning
		}
		return nil
	}); err != nil {
		log.Fatalf("Stream failed: %v", err)
	}
}
