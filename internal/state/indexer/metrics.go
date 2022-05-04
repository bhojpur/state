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

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"

	prometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// MetricsSubsystem is a the subsystem label for the indexer package.
const MetricsSubsystem = "indexer"

// Metrics contains metrics exposed by this package.
type Metrics struct {
	// Latency for indexing block events.
	BlockEventsSeconds metrics.Histogram

	// Latency for indexing transaction events.
	TxEventsSeconds metrics.Histogram

	// Number of complete blocks indexed.
	BlocksIndexed metrics.Counter

	// Number of transactions indexed.
	TransactionsIndexed metrics.Counter
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
// Optionally, labels can be provided along with their values ("foo",
// "fooValue").
func PrometheusMetrics(namespace string, labelsAndValues ...string) *Metrics {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}
	return &Metrics{
		BlockEventsSeconds: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "block_events_seconds",
			Help:      "Latency for indexing block events.",
		}, labels).With(labelsAndValues...),
		TxEventsSeconds: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "tx_events_seconds",
			Help:      "Latency for indexing transaction events.",
		}, labels).With(labelsAndValues...),
		BlocksIndexed: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "blocks_indexed",
			Help:      "Number of complete blocks indexed.",
		}, labels).With(labelsAndValues...),
		TransactionsIndexed: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "transactions_indexed",
			Help:      "Number of transactions indexed.",
		}, labels).With(labelsAndValues...),
	}
}

// NopMetrics returns an indexer metrics stub that discards all samples.
func NopMetrics() *Metrics {
	return &Metrics{
		BlockEventsSeconds:  discard.NewHistogram(),
		TxEventsSeconds:     discard.NewHistogram(),
		BlocksIndexed:       discard.NewCounter(),
		TransactionsIndexed: discard.NewCounter(),
	}
}
