package nut

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/nut_client_service/internal/http/clients/nut"
)

var metrics = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "nut_client_service",
	Name:      "ups_variables",
	Help:      "Variables of UPS list",
}, []string{"ups", "variable"})

var mappingUPSStatuses = map[string]int{
	"CAL":     0,
	"TRIM":    1,
	"BOOST":   2,
	"OL":      3,
	"OB":      4,
	"OVER":    5,
	"LB":      6,
	"RB":      7,
	"BYPASS":  8,
	"OFF":     9,
	"CHRG":    10,
	"DISCHRG": 11,
}

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
		list, err := m.nutClient.GetUPSList(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("get UPS list fail")
			continue
		}

		for _, ups := range list {
			for _, v := range ups.Variables {
				if v.Type == "INTEGER" || v.Type == "FLOAT_64" {
					value, err := strconv.ParseFloat(fmt.Sprintf("%v", v.Value), 64)
					if err != nil {
						log.Warn().Err(err).Msg("parse float64 of value fail")
						continue
					}

					metrics.WithLabelValues(ups.Name, v.Name).Set(value)
				} else {
					if v.Type == "STRING" {
						str, ok := v.Value.(string)
						if !ok {
							log.Warn().Err(err).Msg("type cast to string fail")
							continue
						}

						if v.Name == "ups.status" {
							metrics.WithLabelValues(ups.Name, v.Name).Set(float64(mapUpsStatus(str)))
						} else {
							// TODO: Parser another string values to float metrics representations
							//metrics.WithLabelValues(ups.Name, v.Name).Set(str)
						}
					}
				}
			}
		}
	}

	return nil
}

// mapUpsStatus Mapping string value of UPS status to int constants
func mapUpsStatus(value string) int {
	res, ok := mappingUPSStatuses[value]
	if !ok {
		return -1
	}

	return res
}
