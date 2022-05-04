package inspect

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
It provides a tool for investigating the state of a failed Bhojpur State node.

This package provides the Inspector type. The Inspector type runs a subset of the Bhojpur State
RPC endpoints that are useful for debugging issues with Bhojpur State consensus.

When a node running the Bhojpur State consensus engine detects an inconsistent consensus state,
the entire node will crash. The Bhojpur State consensus engine cannot run in this
inconsistent state so the node will not be able to start up again.

The RPC endpoints provided by the Inspector type allow for a node operator to inspect
the block store and state store to better understand what may have caused the inconsistent state.


The Inspector type's lifecycle is controlled by a context.Context
  ins := inspect.NewFromConfig(rpcConfig)
  ctx, cancelFunc:= context.WithCancel(context.Background())

  // Run blocks until the Inspector server is shut down.
  go ins.Run(ctx)
  ...

  // calling the cancel function will stop the running inspect server
  cancelFunc()

Inspector serves its RPC endpoints on the address configured in the RPC configuration

  rpcConfig.ListenAddress = "tcp://127.0.0.1:26657"
  ins := inspect.NewFromConfig(rpcConfig)
  go ins.Run(ctx)

The list of available RPC endpoints can then be viewed by navigating to
http://127.0.0.1:26657/ in the web browser.
*/
