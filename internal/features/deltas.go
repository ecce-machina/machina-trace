package features

import "github.com/ecce-machina/machina-trace/internal/diff"

// WorkloadFeaturesFromDeltas derives workload features from a set of counter
// deltas belonging to one observation.
//
// The caller is responsible for ensuring that the deltas represent the same
// node and observation interval.
func WorkloadFeaturesFromDeltas(
	deltas []diff.CounterDelta,
) WorkloadFeatures {
	var clientIO LustreClientIOFeatures
	var metadata LustreMetadataFeatures
	var mdc LustreMDCFeatures

	for _, delta := range deltas {
		if f, ok := FromLustreLLiteDelta(delta); ok {
			clientIO.IntervalSec = f.IntervalSec
			clientIO.ReadOpsPerSec += f.ReadOpsPerSec
			clientIO.WriteOpsPerSec += f.WriteOpsPerSec
			clientIO.ReadBytesPerSec += f.ReadBytesPerSec
			clientIO.WriteBytesPerSec += f.WriteBytesPerSec
		}

		if f, ok := FromLustreMDCMetadataDelta(delta); ok {
			metadata.IntervalSec = f.IntervalSec
			metadata.CreatesPerSec += f.CreatesPerSec
			metadata.ClosesPerSec += f.ClosesPerSec
			metadata.GetattrsPerSec += f.GetattrsPerSec
			metadata.GetxattrsPerSec += f.GetxattrsPerSec
			metadata.SetattrsPerSec += f.SetattrsPerSec
			metadata.RenamesPerSec += f.RenamesPerSec
			metadata.UnlinksPerSec += f.UnlinksPerSec
			metadata.IntentLocksPerSec += f.IntentLocksPerSec
			metadata.RevalidateLocksPerSec += f.RevalidateLocksPerSec
		}

		if f, ok := FromLustreMDCDelta(delta); ok {
			mdc.IntervalSec = f.IntervalSec
			mdc.RequestsPerSec += f.RequestsPerSec
			mdc.LDLMEnqueuesPerSec += f.LDLMEnqueuesPerSec
			mdc.LDLMCancelsPerSec += f.LDLMCancelsPerSec
		}
	}

	return DeriveWorkloadFeatures(clientIO, metadata, mdc)
}
