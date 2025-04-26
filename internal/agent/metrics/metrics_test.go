package metrics

import (
	"testing"

	"github.com/echo9et/alerting/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestMetricsRuntime_Update(t *testing.T) {
	tests := []struct {
		name  string
		err   string
		mType string
	}{
		{
			name:  "PollCount",
			err:   "PollCount should be updated",
			mType: entities.Counter,
		},
		{
			name:  "RandomValue",
			err:   "RandomValue should be updated",
			mType: entities.Gauge,
		},
		{
			name:  "GCCPUFraction",
			err:   "GCCPUFraction should be updated",
			mType: entities.Gauge,
		},
	}

	metrics := NewMetricsRuntime()
	metrics.Update()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mType == entities.Counter {
				assert.NotZero(t, metrics.data.Counters[tt.name], tt.err)
			} else if tt.mType == entities.Gauge {
				assert.NotZero(t, metrics.data.Gauges[tt.name], tt.err)
			}
		})
	}
}

func TestMetricsRuntime_ToJSON(t *testing.T) {
	tests := []struct {
		name  string
		err   string
		mType string
	}{
		{
			name:  "PollCount",
			err:   "Delta for PollCount should not be nil",
			mType: entities.Counter,
		},
		{
			name:  "RandomValue",
			err:   "Value for RandomValue should not be nil",
			mType: entities.Gauge,
		},
		{
			name:  "GCCPUFraction",
			err:   "Value for GCCPUFraction should not be nil",
			mType: entities.Gauge,
		},
	}

	metrics := NewMetricsRuntime()
	metrics.Update()
	jsonMetrics := metrics.ToJSON()
	for _, tt := range tests {
		for _, metric := range jsonMetrics {
			if metric.ID == tt.name && metric.MType == tt.mType {
				if tt.mType == entities.Counter {
					assert.NotNil(t, metric.Delta, tt.err)
				} else if tt.mType == entities.Gauge {
					assert.NotNil(t, metric.Value, tt.err)
				}
			}
		}
	}
}

func TestMetricsMem_Update(t *testing.T) {
	tests := []struct {
		name  string
		err   string
		mType string
	}{
		{
			name:  "TotalMemory",
			err:   "TotalMemory should be updated",
			mType: entities.Counter,
		},
		{
			name:  "FreeMemory",
			err:   "FreeMemory should be updated",
			mType: entities.Counter,
		},
	}

	metrics := NewMetricsMem()
	metrics.Update()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotZero(t, metrics.data.Counters[tt.name], tt.err)
		})
	}
}

func TestMetricsMem_ToJSON(t *testing.T) {
	tests := []struct {
		name  string
		err   string
		mType string
	}{
		{
			name:  "TotalMemory",
			err:   "TotalMemory should be updated",
			mType: entities.Counter,
		},
		{
			name:  "FreeMemory",
			err:   "FreeMemory should be updated",
			mType: entities.Counter,
		},
	}

	metrics := NewMetricsMem()
	metrics.Update()
	jsonMetrics := metrics.ToJSON()
	for _, tt := range tests {
		for _, metric := range jsonMetrics {
			if metric.ID == tt.name && metric.MType == tt.mType {
				if tt.mType == entities.Counter {
					assert.NotNil(t, metric.Delta, tt.err)
				} else if tt.mType == entities.Gauge {
					assert.NotNil(t, metric.Value, tt.err)
				}
			}
		}
	}
}
