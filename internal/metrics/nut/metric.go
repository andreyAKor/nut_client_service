package nut

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/andreyAKor/nut_parser/internal/http/clients/nut"
)

var metrics = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "nut_parser",
	Name:      "ups_variables",
	Help:      "Variables of UPS list",
	Buckets:   []float64{0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
}, []string{"ups", "variable", "value"})

type Metric struct {
	interval  time.Duration
	nutClient *nut.Client
}

func New(interval string, nutClient *nut.Client) (*Metric, error) {
	intervalDur, err := time.ParseDuration(interval)
	if err != nil {
		return nil, errors.Wrapf(err, "interval parsing fail (%s)", interval)
	}

	return &Metric{
		interval:  intervalDur,
		nutClient: nutClient,
	}, nil
}

func (m *Metric) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.interval)

	go func() {
		select {
		case <-ctx.Done():
			ticker.Stop()
		}
	}()

	for _ = range ticker.C {
		list, err := m.nutClient.GetUPSList()
		if err != nil {
			return errors.Wrap(err, "get UPS list fail")
		}

		for _, ups := range list {
			for _, v := range ups.Variables {
				metrics.WithLabelValues(ups.Name, v.Name, fmt.Sprintf("%v", v.Value))
			}
		}
	}

	return nil
}
