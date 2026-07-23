package workload

import (
	"fmt"

	"github.com/ecce-machina/machina-trace/internal/features"
)

// BuildProfile translates measured workload features into a semantic
// description of the filesystem load.
//
// Each profile dimension is assessed independently. Dimensions that cannot be
// supported by the currently available features are explicitly reported as
// unknown.
func BuildProfile(f features.WorkloadFeatures) Profile {
	return Profile{
		ReadWriteBalance:  assessReadWriteBalance(f),
		IOSize:            assessIOSize(f),
		AccessPattern:     assessAccessPattern(f),
		DataIntensity:     assessDataIntensity(f),
		MetadataIntensity: assessMetadataIntensity(f),
		NamespaceBehavior: assessNamespaceBehavior(f),
		CacheBehavior:     assessCacheBehavior(f),
	}
}

func assessReadWriteBalance(
	f features.WorkloadFeatures,
) Assessment[Balance] {
	readBytes := f.ReadBytesPerSec
	writeBytes := f.WriteBytesPerSec
	totalBytes := readBytes + writeBytes

	if totalBytes <= 0 {
		return Assessment[Balance]{
			Value:      BalanceUnknown,
			Confidence: 0,
			Evidence: []string{
				"no read or write throughput observed",
				evidenceFloat("read_bytes_per_sec", readBytes),
				evidenceFloat("write_bytes_per_sec", writeBytes),
			},
		}
	}

	readFraction := readBytes / totalBytes
	writeFraction := writeBytes / totalBytes

	var value Balance

	switch {
	case readFraction >= 0.95:
		value = MostlyRead

	case readFraction >= 0.65:
		value = ReadHeavy

	case writeFraction >= 0.95:
		value = MostlyWrite

	case writeFraction >= 0.65:
		value = WriteHeavy

	default:
		value = Balanced
	}

	return Assessment[Balance]{
		Value:      value,
		Confidence: balanceConfidence(value, readFraction, writeFraction),
		Evidence: []string{
			evidenceFloat("read_bytes_per_sec", readBytes),
			evidenceFloat("write_bytes_per_sec", writeBytes),
			evidenceFloat("read_fraction", readFraction),
			evidenceFloat("write_fraction", writeFraction),
		},
	}
}

func assessIOSize(
	f features.WorkloadFeatures,
) Assessment[IOSize] {
	readOps := f.ReadOpsPerSec
	writeOps := f.WriteOpsPerSec
	totalOps := readOps + writeOps

	readBytes := f.ReadBytesPerSec
	writeBytes := f.WriteBytesPerSec
	totalBytes := readBytes + writeBytes

	if totalOps <= 0 {
		return Assessment[IOSize]{
			Value:      SizeUnknown,
			Confidence: 0,
			Evidence: []string{
				"no read or write operations observed",
				evidenceFloat("read_ops_per_sec", readOps),
				evidenceFloat("write_ops_per_sec", writeOps),
			},
		}
	}

	averageSize := totalBytes / totalOps
	value := classifyIOSize(averageSize)

	return Assessment[IOSize]{
		Value:      value,
		Confidence: 0.5,
		Evidence: []string{
			evidenceFloat("average_io_size_bytes", averageSize),
			evidenceFloat("data_bytes_per_sec", totalBytes),
			evidenceFloat("data_ops_per_sec", totalOps),
		},
	}
}

func assessAccessPattern(
	f features.WorkloadFeatures,
) Assessment[AccessPattern] {
	return Assessment[AccessPattern]{
		Value:      AccessUnknown,
		Confidence: 0,
		Evidence: []string{
			"no offset-locality or seek-distance signal is available",
		},
	}
}

func assessDataIntensity(
	f features.WorkloadFeatures,
) Assessment[Level] {
	bytesPerSec := f.ReadBytesPerSec + f.WriteBytesPerSec
	opsPerSec := f.ReadOpsPerSec + f.WriteOpsPerSec

	value := classifyDataIntensity(bytesPerSec)

	confidence := 0.5
	if value == LevelNone {
		confidence = 1
	}

	return Assessment[Level]{
		Value:      value,
		Confidence: confidence,
		Evidence: []string{
			evidenceFloat("data_bytes_per_sec", bytesPerSec),
			evidenceFloat("data_ops_per_sec", opsPerSec),
		},
	}
}

func assessMetadataIntensity(
	f features.WorkloadFeatures,
) Assessment[Level] {
	metadataOps := f.MetadataOpsPerSec
	value := classifyRate(metadataOps)

	confidence := 0.5
	if value == LevelNone {
		confidence = 1
	}

	return Assessment[Level]{
		Value:      value,
		Confidence: confidence,
		Evidence: []string{
			evidenceFloat("metadata_ops_per_sec", metadataOps),
			evidenceFloat("mdc_rpcs_per_sec", f.MDCRPCsPerSec),
			evidenceFloat("lock_ops_per_sec", f.LockOpsPerSec),
		},
	}
}

func assessNamespaceBehavior(
	f features.WorkloadFeatures,
) Assessment[NamespaceProfile] {
	return Assessment[NamespaceProfile]{
		Value:      NamespaceUnknown,
		Confidence: 0,
		Evidence: []string{
			"individual lookup, create, unlink, and rename rates are unavailable",
			evidenceFloat("metadata_ops_per_sec", f.MetadataOpsPerSec),
			evidenceFloat("mdc_rpcs_per_sec", f.MDCRPCsPerSec),
			evidenceFloat("lock_ops_per_sec", f.LockOpsPerSec),
		},
	}
}

func assessCacheBehavior(
	f features.WorkloadFeatures,
) Assessment[CacheProfile] {
	return Assessment[CacheProfile]{
		Value:      CacheUnknown,
		Confidence: 0,
		Evidence: []string{
			"application-visible operations cannot yet be compared with backend RPC activity",
		},
	}
}

func classifyIOSize(bytesPerOp float64) IOSize {
	switch {
	case bytesPerOp <= 0:
		return SizeUnknown

	case bytesPerOp < 4*1024:
		return SizeTiny

	case bytesPerOp < 64*1024:
		return SizeSmall

	case bytesPerOp < 512*1024:
		return SizeMedium

	case bytesPerOp < 2*1024*1024:
		return SizeLarge

	default:
		return SizeHuge
	}
}

func classifyDataIntensity(bytesPerSec float64) Level {
	switch {
	case bytesPerSec <= 0:
		return LevelNone

	case bytesPerSec < 4*1024:
		return LevelVeryLow

	case bytesPerSec < 1024*1024:
		return LevelLow

	case bytesPerSec < 64*1024*1024:
		return LevelMedium

	case bytesPerSec < 1024*1024*1024:
		return LevelHigh

	default:
		return LevelVeryHigh
	}
}

func classifyRate(opsPerSec float64) Level {
	switch {
	case opsPerSec <= 0:
		return LevelNone

	case opsPerSec < 1:
		return LevelVeryLow

	case opsPerSec < 10:
		return LevelLow

	case opsPerSec < 100:
		return LevelMedium

	case opsPerSec < 1000:
		return LevelHigh

	default:
		return LevelVeryHigh
	}
}

func balanceConfidence(
	value Balance,
	readFraction float64,
	writeFraction float64,
) float64 {
	switch value {
	case MostlyRead:
		return readFraction

	case ReadHeavy:
		return readFraction

	case MostlyWrite:
		return writeFraction

	case WriteHeavy:
		return writeFraction

	case Balanced:
		// Balanced covers read fractions from 0.35 through 0.65. Confidence is
		// highest at exactly 0.5 and declines toward either boundary.
		distance := abs(readFraction - 0.5)
		return clamp01(1 - distance/0.15)

	default:
		return 0
	}
}

func evidenceFloat(name string, value float64) string {
	return fmt.Sprintf("%s=%.2f", name, value)
}

func clamp01(value float64) float64 {
	switch {
	case value < 0:
		return 0
	case value > 1:
		return 1
	default:
		return value
	}
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}

	return value
}
