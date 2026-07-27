package aggregate

import "github.com/ecce-machina/machina-trace/internal/features"

// FilesystemFeatures preserves the per-node feature sets while also exposing
// the filesystem-wide aggregate derived from them.
type FilesystemFeatures struct {
	Nodes     map[string]features.WorkloadFeatures
	Aggregate features.WorkloadFeatures
}

// AggregateWorkloadFeatures derives filesystem-wide workload features from
// per-node workload features.
//
// Rates are additive across nodes. Derived values such as average I/O size and
// read/write fractions are recomputed from the aggregate measurements.
func AggregateWorkloadFeatures(
	nodes map[string]features.WorkloadFeatures,
) FilesystemFeatures {
	perNode := make(map[string]features.WorkloadFeatures, len(nodes))
	aggregate := features.WorkloadFeatures{}

	for node, nodeFeatures := range nodes {
		perNode[node] = nodeFeatures

		aggregate.ReadOpsPerSec += nodeFeatures.ReadOpsPerSec
		aggregate.WriteOpsPerSec += nodeFeatures.WriteOpsPerSec

		aggregate.ReadBytesPerSec += nodeFeatures.ReadBytesPerSec
		aggregate.WriteBytesPerSec += nodeFeatures.WriteBytesPerSec

		aggregate.MetadataOpsPerSec += nodeFeatures.MetadataOpsPerSec
		aggregate.LockOpsPerSec += nodeFeatures.LockOpsPerSec
		aggregate.MDCRPCsPerSec += nodeFeatures.MDCRPCsPerSec

		if aggregate.IntervalSec == 0 {
			aggregate.IntervalSec = nodeFeatures.IntervalSec
		}
	}

	totalIOOps := aggregate.ReadOpsPerSec + aggregate.WriteOpsPerSec

	aggregate.AverageReadSizeBytes = divide(
		aggregate.ReadBytesPerSec,
		aggregate.ReadOpsPerSec,
	)
	aggregate.AverageWriteSizeBytes = divide(
		aggregate.WriteBytesPerSec,
		aggregate.WriteOpsPerSec,
	)
	aggregate.ReadFraction = divide(
		aggregate.ReadOpsPerSec,
		totalIOOps,
	)
	aggregate.WriteFraction = divide(
		aggregate.WriteOpsPerSec,
		totalIOOps,
	)

	return FilesystemFeatures{
		Nodes:     perNode,
		Aggregate: aggregate,
	}
}

func divide(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}
