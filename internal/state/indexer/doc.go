package indexer

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
It defines Bhojpur State's block and transaction event indexing logic.

Bhojpur State supports two primary means of block and transaction event indexing:

1. A key-value sink via an embedded database with a proprietary query language.
2. A Postgres-based sink.

An ABCI application can emit events during block and transaction execution in the form

   <abci.Event.Type>.<abci.EventAttributeKey>=<abci.EventAttributeValue>

for example "transfer.amount=10000".

An operator can enable one or both of the supported indexing sinks via the
'tx-index.indexer' Bhojpur State configuration.

Example:

	[tx-index]
	indexer = ["kv", "psql"]

If an operator wants to completely disable indexing, they may simply just provide
the "null" sink option in the configuration. All other sinks will be ignored if
"null" is provided.

If indexing is enabled, the indexer.Service will iterate over all enabled sinks
and invoke block and transaction indexing via the appropriate IndexBlockEvents
and IndexTxEvents methods.

Note, the "kv" sink is considered deprecated and its query functionality is very
limited, but does allow users to directly query for block and transaction events
against Bhojpur State's RPC. Instead, operators are encouraged to use the "psql"
indexing sink when more complex queries are required and for reliability purposes
as PostgreSQL can scale.

Prior to starting Bhojpur State with the "psql" indexing sink enabled, operators
must ensure the following:

1. The "psql" indexing sink is provided in Bhojpur State's configuration.
2. A 'tx-index.psql-conn' value is provided that contains the PostgreSQL connection URI.
3. The block and transaction event schemas have been created in the PostgreSQL database.

Bhojpur State provides the block and transaction event schemas in the following
path: state/indexer/sink/psql/schema.sql

To create the schema in a PostgreSQL database, perform the schema query
manually or invoke schema creation via the CLI:

	$ psql <flags> -f state/indexer/sink/psql/schema.sql

The "psql" indexing sink prohibits queries via RPC. When using a PostgreSQL sink,
queries can and should be made directly against the database using SQL.

The following are some example SQL queries against the database schema:

* Query for all transaction events for a given transaction hash:

	SELECT * FROM tx_events WHERE hash = '3E7D1F...';

* Query for all transaction events for a given block height:

	SELECT * FROM tx_events WHERE height = 25;

* Query for transaction events that have a given type (i.e. value wildcard):

	SELECT * FROM tx_events WHERE key LIKE '%transfer.recipient%';

Note that if a complete abci.TxResult is needed, you will need to join "tx_events" with
"tx_results" via a foreign key, to obtain contains the raw protobuf-encoded abci.TxResult.
*/
