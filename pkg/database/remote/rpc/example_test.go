package rpc_test

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

import (
	"bytes"
	"context"
	"log"

	v1 "github.com/bhojpur/state/pkg/api/v1/database"
	grpcdb "github.com/bhojpur/state/pkg/database/remote/rpc"
)

func Example() {
	addr := ":8998"
	cert := "server.crt"
	key := "server.key"
	go func() {
		if err := grpcdb.ListenAndServe(addr, cert, key); err != nil {
			log.Fatalf("BindServer: %v", err)
		}
	}()

	client, err := grpcdb.NewClient(addr, cert)
	if err != nil {
		log.Fatalf("Failed to create grpcDB client: %v", err)
	}

	ctx := context.Background()
	// 1. Initialize the DB
	in := &v1.Init{
		Type: "leveldb",
		Name: "grpc-uno-test",
		Dir:  ".",
	}
	if _, err := client.Init(ctx, in); err != nil {
		log.Fatalf("Init error: %v", err)
	}

	// 2. Now it can be used!
	query1 := &v1.Entity{Key: []byte("Project"), Value: []byte("Bslibs-on-gRPC")}
	if _, err := client.SetSync(ctx, query1); err != nil {
		log.Fatalf("SetSync err: %v", err)
	}

	query2 := &v1.Entity{Key: []byte("Project")}
	read, err := client.Get(ctx, query2)
	if err != nil {
		log.Fatalf("Get err: %v", err)
	}
	if g, w := read.Value, []byte("Bslibs-on-gRPC"); !bytes.Equal(g, w) {
		log.Fatalf("got= (%q ==> % X)\nwant=(%q ==> % X)", g, g, w, w)
	}
}
