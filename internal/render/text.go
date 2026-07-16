package render

import (
	"fmt"
	"io"
    "sort"

	"github.com/ecce-machina/machina-trace/internal/diff"
	"github.com/ecce-machina/machina-trace/internal/features"
)

func WriteDiffText(w io.Writer, deltas []diff.CounterDelta) {
	for _, d := range deltas {
		names := make([]string, 0, len(d.Deltas))
		for name, delta := range d.Deltas {
			if delta == 0 && d.Rates[name] == 0 {
				continue
			}
			names = append(names, name)
		}
		sort.Strings(names)

		if len(names) == 0 {
			continue
		}

		fmt.Fprintf(w, "%s", d.Collector)
		if d.Object != "" {
			fmt.Fprintf(w, " %s", d.Object)
		}
		fmt.Fprintf(w, " interval=%.2fs\n", d.IntervalSec)

		for _, name := range names {
			fmt.Fprintf(
				w,
				" %s: delta=%d rate=%.2f/sec\n",
				name,
				d.Deltas[name],
				d.Rates[name],
			)
		}
	}
}

func WriteLustreClientIOFeaturesText(w io.Writer, deltas []diff.CounterDelta) {
	for _, d := range deltas {
		f, ok := features.FromLustreLLiteDelta(d)
		if !ok {
			continue
		}

		if f.ReadOpsPerSec == 0 &&
			f.WriteOpsPerSec == 0 &&
			f.ReadBytesPerSec == 0 &&
			f.WriteBytesPerSec == 0 &&
			f.FsyncsPerSec == 0 &&
			f.OpenOpsPerSec == 0 &&
			f.CloseOpsPerSec == 0 &&
			f.ReadUsecsPerSec == 0 &&
			f.WriteUsecsPerSec == 0 &&
			f.FsyncUsecsPerSec == 0 {
			continue
		}

		fmt.Fprintf(
			w,
			"lustre_client_io_features %s interval=%.2fs\n",
			f.Object,
			f.IntervalSec,
		)
		fmt.Fprintf(w, " read_ops_per_sec: %.2f\n", f.ReadOpsPerSec)
		fmt.Fprintf(w, " write_ops_per_sec: %.2f\n", f.WriteOpsPerSec)
		fmt.Fprintf(w, " read_bytes_per_sec: %.2f\n", f.ReadBytesPerSec)
		fmt.Fprintf(w, " write_bytes_per_sec: %.2f\n", f.WriteBytesPerSec)
		fmt.Fprintf(w, " fsyncs_per_sec: %.2f\n", f.FsyncsPerSec)
		fmt.Fprintf(w, " opens_per_sec: %.2f\n", f.OpenOpsPerSec)
		fmt.Fprintf(w, " closes_per_sec: %.2f\n", f.CloseOpsPerSec)
		fmt.Fprintf(w, " read_usecs_per_sec: %.2f\n", f.ReadUsecsPerSec)
		fmt.Fprintf(w, " write_usecs_per_sec: %.2f\n", f.WriteUsecsPerSec)
		fmt.Fprintf(w, " fsync_usecs_per_sec: %.2f\n", f.FsyncUsecsPerSec)
	}
}

func WriteLustreMetadataFeaturesText(w io.Writer, deltas []diff.CounterDelta) {
	for _, d := range deltas {
		if f, ok := features.FromLustreMDCMetadataDelta(d); ok {
			if f.CreatesPerSec == 0 &&
				f.ClosesPerSec == 0 &&
				f.GetattrsPerSec == 0 &&
				f.GetxattrsPerSec == 0 &&
				f.SetattrsPerSec == 0 &&
				f.RenamesPerSec == 0 &&
				f.UnlinksPerSec == 0 &&
				f.IntentLocksPerSec == 0 &&
				f.RevalidateLocksPerSec == 0 {
				continue
			}

			fmt.Fprintf(
				w,
				"lustre_metadata_features %s interval=%.2fs\n",
				f.Object,
				f.IntervalSec,
			)
			fmt.Fprintf(w, " creates_per_sec: %.2f\n", f.CreatesPerSec)
			fmt.Fprintf(w, " closes_per_sec: %.2f\n", f.ClosesPerSec)
			fmt.Fprintf(w, " getattrs_per_sec: %.2f\n", f.GetattrsPerSec)
			fmt.Fprintf(w, " getxattrs_per_sec: %.2f\n", f.GetxattrsPerSec)
			fmt.Fprintf(w, " setattrs_per_sec: %.2f\n", f.SetattrsPerSec)
			fmt.Fprintf(w, " renames_per_sec: %.2f\n", f.RenamesPerSec)
			fmt.Fprintf(w, " unlinks_per_sec: %.2f\n", f.UnlinksPerSec)
			fmt.Fprintf(w, " intent_locks_per_sec: %.2f\n", f.IntentLocksPerSec)
			fmt.Fprintf(
				w,
				" revalidate_locks_per_sec: %.2f\n",
				f.RevalidateLocksPerSec,
			)
			continue
		}

		f, ok := features.FromLustreMDCDelta(d)
		if !ok {
			continue
		}

		if f.RequestsPerSec == 0 &&
			f.RequestWaitUsecsPerSec == 0 &&
			f.LDLMEnqueuesPerSec == 0 &&
			f.LDLMCancelsPerSec == 0 &&
			f.MetadataClosesPerSec == 0 &&
			f.MetadataSyncsPerSec == 0 {
			continue
		}

		fmt.Fprintf(
			w,
			"lustre_mdc_rpc_features %s interval=%.2fs\n",
			f.Object,
			f.IntervalSec,
		)
		fmt.Fprintf(w, " requests_per_sec: %.2f\n", f.RequestsPerSec)
		fmt.Fprintf(
			w,
			" request_wait_usecs_per_sec: %.2f\n",
			f.RequestWaitUsecsPerSec,
		)
		fmt.Fprintf(w, " ldlm_enqueues_per_sec: %.2f\n", f.LDLMEnqueuesPerSec)
		fmt.Fprintf(w, " ldlm_cancels_per_sec: %.2f\n", f.LDLMCancelsPerSec)
		fmt.Fprintf(w, " metadata_closes_per_sec: %.2f\n", f.MetadataClosesPerSec)
		fmt.Fprintf(w, " metadata_syncs_per_sec: %.2f\n", f.MetadataSyncsPerSec)
	}
}

func WriteLustreOSTFeaturesText(w io.Writer, deltas []diff.CounterDelta) {
	for _, d := range deltas {
		f, ok := features.FromLustreOSCDelta(d)
		if !ok {
			continue
		}

		if f.ReadRPCsPerSec == 0 &&
			f.WriteRPCsPerSec == 0 &&
			f.ReadBytesPerSec == 0 &&
			f.WriteBytesPerSec == 0 &&
			f.RequestsPerSec == 0 &&
			f.RequestWaitUsecsPerSec == 0 &&
			f.ExtentEnqueuesPerSec == 0 &&
			f.SyncsPerSec == 0 {
			continue
		}

		fmt.Fprintf(
			w,
			"lustre_ost_features %s interval=%.2fs\n",
			f.Object,
			f.IntervalSec,
		)
		fmt.Fprintf(w, " read_rpcs_per_sec: %.2f\n", f.ReadRPCsPerSec)
		fmt.Fprintf(w, " write_rpcs_per_sec: %.2f\n", f.WriteRPCsPerSec)
		fmt.Fprintf(w, " read_bytes_per_sec: %.2f\n", f.ReadBytesPerSec)
		fmt.Fprintf(w, " write_bytes_per_sec: %.2f\n", f.WriteBytesPerSec)
		fmt.Fprintf(w, " requests_per_sec: %.2f\n", f.RequestsPerSec)
		fmt.Fprintf(
			w,
			" request_wait_usecs_per_sec: %.2f\n",
			f.RequestWaitUsecsPerSec,
		)
		fmt.Fprintf(
			w,
			" extent_enqueues_per_sec: %.2f\n",
			f.ExtentEnqueuesPerSec,
		)
		fmt.Fprintf(w, " syncs_per_sec: %.2f\n", f.SyncsPerSec)
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

func WriteNetworkFeaturesText(w io.Writer, deltas []diff.CounterDelta) {
	for _, d := range deltas {
		f, ok := features.FromNetdevDelta(d)
		if !ok {
			continue
		}

		if f.ReceiveBytesPerSec == 0 &&
			f.TransmitBytesPerSec == 0 &&
			f.ReceivePacketsPerSec == 0 &&
			f.TransmitPacketsPerSec == 0 &&
			f.ReceiveErrorsPerSec == 0 &&
			f.TransmitErrorsPerSec == 0 &&
			f.ReceiveDropsPerSec == 0 &&
			f.TransmitDropsPerSec == 0 {
			continue
		}

		fmt.Fprintf(w, "network_features %s interval=%.2fs\n", f.Object, f.IntervalSec)
		fmt.Fprintf(w, "  receive_bytes_per_sec: %.2f\n", f.ReceiveBytesPerSec)
		fmt.Fprintf(w, "  transmit_bytes_per_sec: %.2f\n", f.TransmitBytesPerSec)
		fmt.Fprintf(w, "  receive_packets_per_sec: %.2f\n", f.ReceivePacketsPerSec)
		fmt.Fprintf(w, "  transmit_packets_per_sec: %.2f\n", f.TransmitPacketsPerSec)
		fmt.Fprintf(w, "  receive_errors_per_sec: %.2f\n", f.ReceiveErrorsPerSec)
		fmt.Fprintf(w, "  transmit_errors_per_sec: %.2f\n", f.TransmitErrorsPerSec)
		fmt.Fprintf(w, "  receive_drops_per_sec: %.2f\n", f.ReceiveDropsPerSec)
		fmt.Fprintf(w, "  transmit_drops_per_sec: %.2f\n", f.TransmitDropsPerSec)
	}
}

