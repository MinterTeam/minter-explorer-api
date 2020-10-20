package metrics

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/centrifugal/centrifuge-go"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type LastBlockMetric struct {
	id   prometheus.Gauge
	time prometheus.Gauge
}

func NewLastBlockMetric() *LastBlockMetric {
	prometheusLastBlockIdMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "explorer_last_block_id",
		},
	)

	prometheusLastBlockTimeMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "explorer_last_block_time",
		},
	)

	prometheus.MustRegister(prometheusLastBlockIdMetric)
	prometheus.MustRegister(prometheusLastBlockTimeMetric)

	return &LastBlockMetric{
		id:   prometheusLastBlockIdMetric,
		time: prometheusLastBlockTimeMetric,
	}
}

func (m *LastBlockMetric) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	var block blocks.Resource
	err := json.Unmarshal(e.Data, &block)
	helpers.CheckErr(err)

	blockTime, err := time.Parse(time.RFC3339, block.Timestamp)
	helpers.CheckErr(err)

	m.id.Set(float64(block.ID))
	m.time.Set(float64(blockTime.Unix()))
}
