package features

type WorkloadFeatures struct {
	IntervalSec float64

	ReadOpsPerSec    float64
	WriteOpsPerSec   float64
	ReadBytesPerSec  float64
	WriteBytesPerSec float64

	AverageReadSizeBytes  float64
	AverageWriteSizeBytes float64
	ReadFraction          float64
	WriteFraction         float64

	MetadataOpsPerSec float64
	LockOpsPerSec     float64
	MDCRPCsPerSec     float64
}

func DeriveWorkloadFeatures(
	io LustreClientIOFeatures,
	metadata LustreMetadataFeatures,
	mdc LustreMDCFeatures,
) WorkloadFeatures {
	totalIOOps := io.ReadOpsPerSec + io.WriteOpsPerSec

	metadataOps := metadata.CreatesPerSec +
		metadata.ClosesPerSec +
		metadata.GetattrsPerSec +
		metadata.GetxattrsPerSec +
		metadata.SetattrsPerSec +
		metadata.RenamesPerSec +
		metadata.UnlinksPerSec

	lockOps := metadata.IntentLocksPerSec +
		metadata.RevalidateLocksPerSec +
		mdc.LDLMEnqueuesPerSec +
		mdc.LDLMCancelsPerSec

	return WorkloadFeatures{
		IntervalSec: io.IntervalSec,

		ReadOpsPerSec:    io.ReadOpsPerSec,
		WriteOpsPerSec:   io.WriteOpsPerSec,
		ReadBytesPerSec:  io.ReadBytesPerSec,
		WriteBytesPerSec: io.WriteBytesPerSec,

		AverageReadSizeBytes: safeDivide(
			io.ReadBytesPerSec,
			io.ReadOpsPerSec,
		),
		AverageWriteSizeBytes: safeDivide(
			io.WriteBytesPerSec,
			io.WriteOpsPerSec,
		),
		ReadFraction: safeDivide(
			io.ReadOpsPerSec,
			totalIOOps,
		),
		WriteFraction: safeDivide(
			io.WriteOpsPerSec,
			totalIOOps,
		),

		MetadataOpsPerSec: metadataOps,
		LockOpsPerSec:     lockOps,
		MDCRPCsPerSec:     mdc.RequestsPerSec,
	}
}

func safeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}
