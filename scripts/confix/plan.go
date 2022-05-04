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

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/creachadair/tomledit"
	"github.com/creachadair/tomledit/parser"
	"github.com/creachadair/tomledit/transform"
)

// The plan is the sequence of transformation steps that should be applied, in
// the given order, to convert a configuration file to be compatible with the
// current version of the config grammar.
//
// Transformation steps are specific to the target config version.  For this
// reason, you must exercise caution when backporting changes to this script
// into older releases.
var plan = transform.Plan{
	{
		Desc: "Rename everything from snake_case to kebab-case",
		T:    transform.SnakeToKebab(),
	},
	{
		// [fastsync]  renamed
		// [blocksync] removed
		Desc: "Remove [fastsync] and [blocksync] sections",
		T: transform.Func(func(_ context.Context, doc *tomledit.Document) error {
			doc.First("fast-sync").Remove()
			transform.FindTable(doc, "fastsync").Remove()
			transform.FindTable(doc, "blocksync").Remove()
			return nil
		}),
		ErrorOK: true,
	},
	{
		Desc: `Add top-level mode setting (default "full")`,
		T: transform.EnsureKey(nil, &parser.KeyValue{
			Block: parser.Comments{"Mode of Node: full | validator | seed"},
			Name:  parser.Key{"mode"},
			Value: parser.MustValue(`"full"`),
		}),
		ErrorOK: true,
	},
	{
		Desc: "Remove gRPC settings from the [rpc] section",
		T: transform.Func(func(_ context.Context, doc *tomledit.Document) error {
			doc.First("rpc", "grpc-laddr").Remove()
			doc.First("rpc", "grpc-max-open-connections").Remove()
			return nil
		}),
	},
	{
		Desc: "Remove per-node consensus timeouts (converted to consensus parameters)",
		T: transform.Remove(
			parser.Key{"consensus", "skip-timeout-commit"},
			parser.Key{"consensus", "timeout-commit"},
			parser.Key{"consensus", "timeout-precommit"},
			parser.Key{"consensus", "timeout-precommit-delta"},
			parser.Key{"consensus", "timeout-prevote"},
			parser.Key{"consensus", "timeout-prevote-delta"},
			parser.Key{"consensus", "timeout-propose"},
			parser.Key{"consensus", "timeout-propose-delta"},
		),
		ErrorOK: true,
	},
	{
		// Removed wal-dir
		// Removed version
		Desc: "Remove vestigial mempool.wal-dir settings",
		T: transform.Remove(
			parser.Key{"mempool", "wal-dir"},
			parser.Key{"mempool", "version"},
		),
		ErrorOK: true,
	},
	{
		Desc: "Add new [p2p] queue-type setting",
		T: transform.EnsureKey(parser.Key{"p2p"}, &parser.KeyValue{
			Block: parser.Comments{"Select the p2p internal queue"},
			Name:  parser.Key{"queue-type"},
			Value: parser.MustValue(`"priority"`),
		}),
		ErrorOK: true,
	},
	{
		Desc: "Add [p2p] connection count and rate limit settings",
		T: transform.Func(func(_ context.Context, doc *tomledit.Document) error {
			tab := transform.FindTable(doc, "p2p")
			if tab == nil {
				return errors.New("p2p table not found")
			}
			transform.InsertMapping(tab.Section, &parser.KeyValue{
				Block: parser.Comments{"Maximum number of connections (inbound and outbound)."},
				Name:  parser.Key{"max-connections"},
				Value: parser.MustValue("64"),
			}, false)
			transform.InsertMapping(tab.Section, &parser.KeyValue{
				Block: parser.Comments{
					"Rate limits the number of incoming connection attempts per IP address.",
				},
				Name:  parser.Key{"max-incoming-connection-attempts"},
				Value: parser.MustValue("100"),
			}, false)
			return nil
		}),
	},
	{
		// Added "chunk-fetchers"
		// This value was backported
		// Renamed to "fetchers"
		Desc: "Rename statesync.chunk-fetchers to statesync.fetchers",
		T: transform.Func(func(ctx context.Context, doc *tomledit.Document) error {
			// If the key already exists, rename it preserving its value.
			if found := doc.First("statesync", "chunk-fetchers"); found != nil {
				found.KeyValue.Name = parser.Key{"fetchers"}
				return nil
			}

			// Otherwise, add it.
			return transform.EnsureKey(parser.Key{"statesync"}, &parser.KeyValue{
				Block: parser.Comments{
					"The number of concurrent chunk and block fetchers to run (default: 4).",
				},
				Name:  parser.Key{"fetchers"},
				Value: parser.MustValue("4"),
			})(ctx, doc)
		}),
	},
	{
		// Backported
		Desc: "Add statesync.use-p2p setting",
		T: transform.EnsureKey(parser.Key{"statesync"}, &parser.KeyValue{
			Block: parser.Comments{
				"# State sync uses light client verification to verify state. This can be done either through the",
				"# P2P layer or RPC layer. Set this to true to use the P2P layer. If false (default), RPC layer",
				"# will be used.",
			},
			Name:  parser.Key{"use-p2p"},
			Value: parser.MustValue("false"),
		}),
	},
	{
		Desc: "Move priv-validator settings under [priv-validator]",
		T: transform.Func(func(_ context.Context, doc *tomledit.Document) error {
			const pvPrefix = "priv-validator-"

			var found []*tomledit.Entry
			doc.Global.Scan(func(key parser.Key, e *tomledit.Entry) bool {
				if len(key) == 1 && strings.HasPrefix(key[0], pvPrefix) {
					found = append(found, e)
				}
				return true
			})
			if len(found) == 0 {
				return nil // nothing to do
			}

			// Now that we know we have work to do, find the target table.
			var sec *tomledit.Section
			if dst := transform.FindTable(doc, "priv-validator"); dst == nil {
				// If the table doesn't exist, create it. Old config files
				// probably will not have it, so plug in the comment too.
				sec = &tomledit.Section{
					Heading: &parser.Heading{
						Block: parser.Comments{
							"#######################################################",
							"###       Priv Validator Configuration              ###",
							"#######################################################",
						},
						Name: parser.Key{"priv-validator"},
					},
				}
				doc.Sections = append(doc.Sections, sec)
			} else {
				sec = dst.Section
			}

			for _, e := range found {
				e.Remove()
				e.Name = parser.Key{strings.TrimPrefix(e.Name[0], pvPrefix)}
				sec.Items = append(sec.Items, e.KeyValue)
			}
			return nil
		}),
	},
	{
		Desc: "Convert tx-index.indexer from a string to a list of strings",
		T: transform.Func(func(ctx context.Context, doc *tomledit.Document) error {
			idx := doc.First("tx-index", "indexer")
			if idx == nil {
				// No previous indexer setting: Default to ["null"]
				return transform.EnsureKey(parser.Key{"tx-index"}, &parser.KeyValue{
					Block: parser.Comments{"The backend database list to back the indexer."},
					Name:  parser.Key{"indexer"},
					Value: parser.MustValue(`["null"]`),
				})(ctx, doc)
			}

			// use an array of strings.
			switch idx.KeyValue.Value.X.(type) {
			case parser.Array:
				// OK, this is already up-to-date.
				return nil
			case parser.Token:
				// Wrap the value in a single-element array.
				idx.KeyValue.Value.X = parser.Array{idx.KeyValue.Value}
				return nil
			}
			return fmt.Errorf("unrecognized value: %v", idx.KeyValue)
		}),
	},
}
