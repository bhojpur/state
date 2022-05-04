package eventlog

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
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// gauge is the subset of the Prometheus gauge interface used here.
type gauge interface {
	Set(float64)
}

// Metrics define the metrics exported by the eventlog package.
type Metrics struct {
	numItemsGauge gauge
}

// discard is a no-op implementation of the gauge interface.
type discard struct{}

func (discard) Set(float64) {}

const eventlogSubsystem = "eventlog"

// PrometheusMetrics returns a collection of eventlog metrics for Prometheus.
func PrometheusMetrics(ns string, fields ...string) *Metrics {
	var labels []string
	for i := 0; i < len(fields); i += 2 {
		labels = append(labels, fields[i])
	}
	return &Metrics{
		numItemsGauge: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: ns,
			Subsystem: eventlogSubsystem,
			Name:      "num_items",
			Help:      "Number of items currently resident in the event log.",
		}, labels).With(fields...),
	}
}
