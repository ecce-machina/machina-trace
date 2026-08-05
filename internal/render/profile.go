package render

import (
	"fmt"
	"io"

	"github.com/ecce-machina/machina-trace/internal/workload"
)

func WriteWorkloadProfileText(w io.Writer, p workload.Profile) {
	fmt.Fprintln(w, "workload_profile")

	writeAssessment(
		w,
		"read_write_balance",
		balanceString(p.ReadWriteBalance.Value),
		p.ReadWriteBalance.Confidence,
	)
	writeAssessment(
		w,
		"io_size",
		ioSizeString(p.IOSize.Value),
		p.IOSize.Confidence,
	)
	writeAssessment(
		w,
		"access_pattern",
		accessPatternString(p.AccessPattern.Value),
		p.AccessPattern.Confidence,
	)
	writeAssessment(
		w,
		"data_intensity",
		p.DataIntensity.Value.String(),
		p.DataIntensity.Confidence,
	)
	writeAssessment(
		w,
		"metadata_intensity",
		p.MetadataIntensity.Value.String(),
		p.MetadataIntensity.Confidence,
	)
	writeAssessment(
		w,
		"namespace_behavior",
		namespaceProfileString(p.NamespaceBehavior.Value),
		p.NamespaceBehavior.Confidence,
	)
	writeAssessment(
		w,
		"cache_behavior",
		cacheProfileString(p.CacheBehavior.Value),
		p.CacheBehavior.Confidence,
	)
}

func writeAssessment(
	w io.Writer,
	name string,
	value string,
	confidence float64,
) {
	fmt.Fprintf(
		w,
		"  %s: %s confidence=%.2f\n",
		name,
		value,
		confidence,
	)
}

func balanceString(v workload.Balance) string {
	switch v {
	case workload.MostlyRead:
		return "mostly_read"
	case workload.ReadHeavy:
		return "read_heavy"
	case workload.Balanced:
		return "balanced"
	case workload.WriteHeavy:
		return "write_heavy"
	case workload.MostlyWrite:
		return "mostly_write"
	default:
		return "unknown"
	}
}

func ioSizeString(v workload.IOSize) string {
	switch v {
	case workload.SizeTiny:
		return "tiny"
	case workload.SizeSmall:
		return "small"
	case workload.SizeMedium:
		return "medium"
	case workload.SizeLarge:
		return "large"
	case workload.SizeHuge:
		return "huge"
	default:
		return "unknown"
	}
}

func accessPatternString(v workload.AccessPattern) string {
	switch v {
	case workload.Sequential:
		return "sequential"
	case workload.Random:
		return "random"
	case workload.Mixed:
		return "mixed"
	default:
		return "unknown"
	}
}

func namespaceProfileString(v workload.NamespaceProfile) string {
	switch v {
	case workload.NamespaceIdle:
		return "idle"
	case workload.LookupHeavy:
		return "lookup_heavy"
	case workload.CreateHeavy:
		return "create_heavy"
	case workload.DeleteHeavy:
		return "delete_heavy"
	case workload.RenameHeavy:
		return "rename_heavy"
	case workload.MixedNamespace:
		return "mixed"
	default:
		return "unknown"
	}
}

func cacheProfileString(v workload.CacheProfile) string {
	switch v {
	case workload.MostlyCached:
		return "mostly_cached"
	case workload.MixedCache:
		return "mixed"
	case workload.MostlyMisses:
		return "mostly_misses"
	default:
		return "unknown"
	}
}
