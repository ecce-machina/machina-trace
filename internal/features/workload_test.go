package features

import (
	"math"
	"testing"
)

func TestDeriveWorkloadFeatures(t *testing.T) {
	io := LustreClientIOFeatures{
		IntervalSec:      10,
		ReadOpsPerSec:    20,
		WriteOpsPerSec:   80,
		ReadBytesPerSec:  20 * 4096,
		WriteBytesPerSec: 80 * 16384,
	}

	metadata := LustreMetadataFeatures{
		CreatesPerSec:         10,
		ClosesPerSec:          8,
		GetattrsPerSec:        30,
		GetxattrsPerSec:       2,
		SetattrsPerSec:        3,
		RenamesPerSec:         4,
		UnlinksPerSec:         5,
		IntentLocksPerSec:     12,
		RevalidateLocksPerSec: 7,
	}

	mdc := LustreMDCFeatures{
		RequestsPerSec:     50,
		LDLMEnqueuesPerSec: 6,
		LDLMCancelsPerSec:  2,
	}

	got := DeriveWorkloadFeatures(io, metadata, mdc)

	assertFloatEqual(t, "AverageReadSizeBytes", got.AverageReadSizeBytes, 4096)
	assertFloatEqual(t, "AverageWriteSizeBytes", got.AverageWriteSizeBytes, 16384)
	assertFloatEqual(t, "ReadFraction", got.ReadFraction, 0.2)
	assertFloatEqual(t, "WriteFraction", got.WriteFraction, 0.8)

	assertFloatEqual(t, "MetadataOpsPerSec", got.MetadataOpsPerSec, 62)
	assertFloatEqual(t, "LockOpsPerSec", got.LockOpsPerSec, 27)
	assertFloatEqual(t, "MDCRPCsPerSec", got.MDCRPCsPerSec, 50)
}

func TestDeriveWorkloadFeaturesHandlesZeroIO(t *testing.T) {
	got := DeriveWorkloadFeatures(
		LustreClientIOFeatures{},
		LustreMetadataFeatures{},
		LustreMDCFeatures{},
	)

	assertFloatEqual(t, "AverageReadSizeBytes", got.AverageReadSizeBytes, 0)
	assertFloatEqual(t, "AverageWriteSizeBytes", got.AverageWriteSizeBytes, 0)
	assertFloatEqual(t, "ReadFraction", got.ReadFraction, 0)
	assertFloatEqual(t, "WriteFraction", got.WriteFraction, 0)
}

func assertFloatEqual(t *testing.T, name string, got, want float64) {
	t.Helper()

	const tolerance = 0.000001

	if math.Abs(got-want) > tolerance {
		t.Fatalf("%s = %f, want %f", name, got, want)
	}
}
