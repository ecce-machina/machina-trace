package features

import "github.com/ecce-machina/machina-trace/internal/diff"

type NetworkFeatures struct {
	Object                string
	IntervalSec           float64
	ReceiveBytesPerSec    float64
	TransmitBytesPerSec   float64
	ReceivePacketsPerSec  float64
	TransmitPacketsPerSec float64
	ReceiveErrorsPerSec   float64
	TransmitErrorsPerSec  float64
	ReceiveDropsPerSec    float64
	TransmitDropsPerSec   float64
}

func FromNetdevDelta(d diff.CounterDelta) (NetworkFeatures, bool) {
	if d.Collector != "proc_net_dev" {
		return NetworkFeatures{}, false
	}

	return NetworkFeatures{
		Object:                d.Object,
		IntervalSec:           d.IntervalSec,
		ReceiveBytesPerSec:    d.Rates["receive_bytes"],
		TransmitBytesPerSec:   d.Rates["transmit_bytes"],
		ReceivePacketsPerSec:  d.Rates["receive_packets"],
		TransmitPacketsPerSec: d.Rates["transmit_packets"],
		ReceiveErrorsPerSec:   d.Rates["receive_errors"],
		TransmitErrorsPerSec:  d.Rates["transmit_errors"],
		ReceiveDropsPerSec:    d.Rates["receive_dropped"],
		TransmitDropsPerSec:   d.Rates["transmit_dropped"],
	}, true
}
