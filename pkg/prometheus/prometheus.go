package prometheus

import (
	"github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Counter struct {
	vec        *prometheus.CounterVec
	labelNames []string
}

func NewCounter(vec *prometheus.CounterVec) *Counter {
	return &Counter{vec: vec}
}

func NewCounterForm(opts prometheus.CounterOpts, labelNames []string) *Counter {
	vec := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(vec)
	return NewCounter(vec)
}

func (c *Counter) With(labelValue ...string) metrics.Counter {
	return &Counter{
		vec:        c.vec,
		labelNames: append(c.labelNames, labelValue...),
	}
}

func (c *Counter) Add(delta float64) {
	c.vec.With(makeLabels(c.labelNames...)).Add(delta)
}

// todo
func makeLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}

type Summary struct {
	svec        *prometheus.SummaryVec
	labelValues []string
}

func NewSummary(vec *prometheus.SummaryVec) *Summary {
	return &Summary{
		svec: vec,
	}
}

func NewSummaryForm(opts prometheus.SummaryOpts, labelsName []string) *Summary {
	sv := prometheus.NewSummaryVec(opts, labelsName)
	prometheus.MustRegister(sv)
	return NewSummary(sv)
}

func (s *Summary) With(labelValues ...string) metrics.Histogram {
	return &Summary{
		svec:        s.svec,
		labelValues: append(s.labelValues, labelValues...),
	}
}

func (s *Summary) Observe(value float64) {
	s.svec.With(makeLabels(s.labelValues...)).Observe(value)
}
