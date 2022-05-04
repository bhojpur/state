package proxy

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
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	MetricsSubsystem = "abci_connection"
)

// Metrics contains the prometheus metrics exposed by the proxy package.
type Metrics struct {
	MethodTiming metrics.Histogram
}

// PrometheusMetrics constructs a Metrics instance that collects metrics samples.
// The resulting metrics will be prefixed with namespace and labeled with the
// defaultLabelsAndValues. defaultLabelsAndValues must be a list of string pairs
// where the first of each pair is the label and the second is the value.
func PrometheusMetrics(namespace string, defaultLabelsAndValues ...string) *Metrics {
	defaultLabels := []string{}
	for i := 0; i < len(defaultLabelsAndValues); i += 2 {
		defaultLabels = append(defaultLabels, defaultLabelsAndValues[i])
	}
	return &Metrics{
		MethodTiming: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "method_timing",
			Help:      "ABCI Method Timing",
			Buckets:   []float64{.0001, .0004, .002, .009, .02, .1, .65, 2, 6, 25},
		}, append(defaultLabels, []string{"method", "type"}...)).With(defaultLabelsAndValues...),
	}
}

// NopMetrics constructs a Metrics instance that discards all samples and is suitable
// for testing.
func NopMetrics() *Metrics {
	return &Metrics{
		MethodTiming: discard.NewHistogram(),
	}
}
