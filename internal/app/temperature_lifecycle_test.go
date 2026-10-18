package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func TestQuickSettingTemperatureLifecycle(t *testing.T) {
	tempSensorGauge.Reset()
	previousTempSensors = map[string]prometheus.Labels{}
	registry := prometheus.NewRegistry()
	registry.MustRegister(tempSensorGauge)
	count := func() int {
		families, err := registry.Gather()
		if err != nil {
			t.Fatal(err)
		}
		n := 0
		for _, family := range families {
			n += len(family.Metric)
		}
		return n
	}
	// Reuse of an NVMe index must not retain the disk previously at that index.
	updatePrometheusSensors(nil, []TempSensor{{Key: "Nv00", Name: "Internal", Value: 33}, {Key: "TS0p", Name: "SSD", Value: 40}})
	if count() != 2 {
		t.Fatal("missing initial sensors")
	}
	updatePrometheusSensors(nil, []TempSensor{{Key: "Nv00", Name: "External", Value: 43}, {Key: "TS0p", Name: "SSD", Value: 41}})
	if count() != 2 {
		t.Fatal("old disk labels retained")
	}
	// Failure/removal must stop exporting the cached disk, preserving other sensors.
	updatePrometheusSensors(nil, []TempSensor{{Key: "TS0p", Name: "SSD", Value: 42}})
	if count() != 1 {
		t.Fatal("failed disk remained exported")
	}
	updatePrometheusSensors(nil, nil)
	if count() != 0 {
		t.Fatal("removed sensors remained exported")
	}
	updatePrometheusSensors(nil, []TempSensor{{Key: "Nv00", Name: "Internal", Value: 35}})
	if count() != 1 {
		t.Fatal("sensor did not recover")
	}
}
