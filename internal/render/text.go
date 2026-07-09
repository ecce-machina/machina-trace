package render

import (
	"fmt"
	"io"

	"github.com/ecce-machina/machina-trace/internal/diff"
	"github.com/ecce-machina/machina-trace/internal/features"
)

func WriteDiffText(w io.Writer, deltas []diff.CounterDelta) {
	for _, d := range deltas {
		printedHeader := false

		for name, rate := range d.Rates {
			if rate == 0 {
				continue
			}

			if !printedHeader {
				fmt.Fprintf(w, "%s", d.Collector)
				if d.Object != "" {
					fmt.Fprintf(w, " %s", d.Object)
				}
				fmt.Fprintf(w, " interval=%.2fs\n", d.IntervalSec)
				printedHeader = true
			}

			fmt.Fprintf(w, "  %s_per_sec: %.2f\n", name, rate)
		}
	}
}

func WriteDiskFeaturesText(w io.Writer, deltas []diff.CounterDelta) {
	for _, d := range deltas {
		f, ok := features.FromDiskstatsDelta(d)
		if !ok {
			continue
		}

		if f.ReadIOPS == 0 && f.WriteIOPS == 0 && f.ReadBytesPerSec == 0 && f.WriteBytesPerSec == 0 {
			continue
		}

		fmt.Fprintf(w, "disk_features %s interval=%.2fs\n", f.Object, f.IntervalSec)
		fmt.Fprintf(w, "  read_iops: %.2f\n", f.ReadIOPS)
		fmt.Fprintf(w, "  write_iops: %.2f\n", f.WriteIOPS)
		fmt.Fprintf(w, "  read_bytes_per_sec: %.2f\n", f.ReadBytesPerSec)
		fmt.Fprintf(w, "  write_bytes_per_sec: %.2f\n", f.WriteBytesPerSec)
		fmt.Fprintf(w, "  avg_read_size_bytes: %.2f\n", f.AvgReadSizeBytes)
		fmt.Fprintf(w, "  avg_write_size_bytes: %.2f\n", f.AvgWriteSizeBytes)
	}
}
