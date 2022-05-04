package sink

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
	"errors"
	"strings"

	"github.com/bhojpur/state/internal/state/indexer"
	"github.com/bhojpur/state/internal/state/indexer/sink/kv"
	"github.com/bhojpur/state/internal/state/indexer/sink/null"
	"github.com/bhojpur/state/internal/state/indexer/sink/psql"
	"github.com/bhojpur/state/pkg/config"
)

// EventSinksFromConfig constructs a slice of indexer.EventSink using the provided
// configuration.
func EventSinksFromConfig(cfg *config.Config, dbProvider config.DBProvider, chainID string) ([]indexer.EventSink, error) {
	if len(cfg.TxIndex.Indexer) == 0 {
		return []indexer.EventSink{null.NewEventSink()}, nil
	}

	// check for duplicated sinks
	sinks := map[string]struct{}{}
	for _, s := range cfg.TxIndex.Indexer {
		sl := strings.ToLower(s)
		if _, ok := sinks[sl]; ok {
			return nil, errors.New("found duplicated sinks, please check the tx-index section in the config.toml")
		}
		sinks[sl] = struct{}{}
	}
	eventSinks := []indexer.EventSink{}
	for k := range sinks {
		switch indexer.EventSinkType(k) {
		case indexer.NULL:
			// When we see null in the config, the eventsinks will be reset with the
			// nullEventSink.
			return []indexer.EventSink{null.NewEventSink()}, nil

		case indexer.KV:
			store, err := dbProvider(&config.DBContext{ID: "tx_index", Config: cfg})
			if err != nil {
				return nil, err
			}

			eventSinks = append(eventSinks, kv.NewEventSink(store))

		case indexer.PSQL:
			conn := cfg.TxIndex.PsqlConn
			if conn == "" {
				return nil, errors.New("the psql connection settings cannot be empty")
			}

			es, err := psql.NewEventSink(conn, chainID)
			if err != nil {
				return nil, err
			}
			eventSinks = append(eventSinks, es)
		default:
			return nil, errors.New("unsupported event sink type")
		}
	}
	return eventSinks, nil

}
