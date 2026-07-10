package features

import "github.com/ecce-machina/machina-trace/internal/diff"

type NetworkStats struct {
	RXBytesPerSec   float64
	TXBytesPerSec   float64
	RXPacketsPerSec float64
	TXPacketsPerSec float64
	RXErrorsPerSec  float64
	TXErrorsPerSec  float64
	RXDropsPerSec   float64
	TXDropsPerSec   float64
}

func NetworkFeatures(d diff.CounterDelta) (NetworkStats, bool) {
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

