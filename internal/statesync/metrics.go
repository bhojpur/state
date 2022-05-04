package statesync

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
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this package.
	MetricsSubsystem = "statesync"
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
	TotalSnapshots      metrics.Counter
	ChunkProcessAvgTime metrics.Gauge
	SnapshotHeight      metrics.Gauge
	SnapshotChunk       metrics.Counter
	SnapshotChunkTotal  metrics.Gauge
	BackFilledBlocks    metrics.Counter
	BackFillBlocksTotal metrics.Gauge
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
		TotalSnapshots: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "total_snapshots",
			Help:      "The total number of snapshots discovered.",
		}, labels).With(labelsAndValues...),
		ChunkProcessAvgTime: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "chunk_process_avg_time",
			Help:      "The average processing time per chunk.",
		}, labels).With(labelsAndValues...),
		SnapshotHeight: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "snapshot_height",
			Help:      "The height of the current snapshot the has been processed.",
		}, labels).With(labelsAndValues...),
		SnapshotChunk: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "snapshot_chunk",
			Help:      "The current number of chunks that have been processed.",
		}, labels).With(labelsAndValues...),
		SnapshotChunkTotal: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "snapshot_chunks_total",
			Help:      "The total number of chunks in the current snapshot.",
		}, labels).With(labelsAndValues...),
		BackFilledBlocks: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "backfilled_blocks",
			Help:      "The current number of blocks that have been back-filled.",
		}, labels).With(labelsAndValues...),
		BackFillBlocksTotal: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "backfilled_blocks_total",
			Help:      "The total number of blocks that need to be back-filled.",
		}, labels).With(labelsAndValues...),
	}
}

// NopMetrics returns no-op Metrics.
func NopMetrics() *Metrics {
	return &Metrics{
		TotalSnapshots:      discard.NewCounter(),
		ChunkProcessAvgTime: discard.NewGauge(),
		SnapshotHeight:      discard.NewGauge(),
		SnapshotChunk:       discard.NewCounter(),
		SnapshotChunkTotal:  discard.NewGauge(),
		BackFilledBlocks:    discard.NewCounter(),
		BackFillBlocksTotal: discard.NewGauge(),
	}
}
