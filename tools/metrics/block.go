package metrics

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type LastBlockMetric struct {
	id   prometheus.Gauge
	time prometheus.Gauge
}

// Constructor for prometheus metrics
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

// Update last block for prometheus metric
func (m *LastBlockMetric) OnNewBlock(block blocks.Resource) {
	blockTime, err := time.Parse(time.RFC3339, block.Timestamp)
	helpers.CheckErr(err)

	m.id.Set(float64(block.ID))
	m.time.Set(float64(blockTime.Unix()))
}
