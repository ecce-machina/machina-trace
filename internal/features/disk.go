package features

import "github.com/ecce-machina/machina-trace/internal/diff"

const SectorSizeBytes = 512

type DiskIOFeatures struct {
	Object            string
	IntervalSec       float64
	ReadIOPS          float64
	WriteIOPS         float64
	ReadBytesPerSec   float64
	WriteBytesPerSec  float64
	AvgReadSizeBytes  float64
	AvgWriteSizeBytes float64
}

func FromDiskstatsDelta(d diff.CounterDelta) (DiskIOFeatures, bool) {
	if d.Collector != "proc_diskstats" {
		return DiskIOFeatures{}, false
	}

	readIOPS := d.Rates["reads_completed"]
	writeIOPS := d.Rates["writes_completed"]
	readSectorsPerSec := d.Rates["sectors_read"]
	writeSectorsPerSec := d.Rates["sectors_written"]

	readBytesPerSec := readSectorsPerSec * SectorSizeBytes
	writeBytesPerSec := writeSectorsPerSec * SectorSizeBytes

	var avgReadSize float64
	if readIOPS > 0 {
		avgReadSize = readBytesPerSec / readIOPS
	}

	var avgWriteSize float64
	if writeIOPS > 0 {
		avgWriteSize = writeBytesPerSec / writeIOPS
	}

	return DiskIOFeatures{
		Object:            d.Object,
		IntervalSec:       d.IntervalSec,
		ReadIOPS:          readIOPS,
		WriteIOPS:         writeIOPS,
		ReadBytesPerSec:   readBytesPerSec,
		WriteBytesPerSec:  writeBytesPerSec,
		AvgReadSizeBytes:  avgReadSize,
		AvgWriteSizeBytes: avgWriteSize,
	}, true
}
