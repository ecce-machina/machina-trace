package workload

import (
	"testing"

	"github.com/ecce-machina/machina-trace/internal/features"
)

func TestBuildProfileNoActivity(t *testing.T) {
	f := features.WorkloadFeatures{}

	got := BuildProfile(f)

	if got.ReadWriteBalance.Value != BalanceUnknown {
		t.Errorf(
			"ReadWriteBalance.Value = %v, want %v",
			got.ReadWriteBalance.Value,
			BalanceUnknown,
		)
	}

	if got.IOSize.Value != SizeUnknown {
		t.Errorf(
			"IOSize.Value = %v, want %v",
			got.IOSize.Value,
			SizeUnknown,
		)
	}

	if got.DataIntensity.Value != LevelNone {
		t.Errorf(
			"DataIntensity.Value = %v, want %v",
			got.DataIntensity.Value,
			LevelNone,
		)
	}

	if got.MetadataIntensity.Value != LevelNone {
		t.Errorf(
			"MetadataIntensity.Value = %v, want %v",
			got.MetadataIntensity.Value,
			LevelNone,
		)
	}

	if got.AccessPattern.Value != AccessUnknown {
		t.Errorf(
			"AccessPattern.Value = %v, want %v",
			got.AccessPattern.Value,
			AccessUnknown,
		)
	}

	if got.NamespaceBehavior.Value != NamespaceUnknown {
		t.Errorf(
			"NamespaceBehavior.Value = %v, want %v",
			got.NamespaceBehavior.Value,
			NamespaceUnknown,
		)
	}

	if got.CacheBehavior.Value != CacheUnknown {
		t.Errorf(
			"CacheBehavior.Value = %v, want %v",
			got.CacheBehavior.Value,
			CacheUnknown,
		)
	}
}

func TestBuildProfileLargeWriteWorkload(t *testing.T) {
	f := features.WorkloadFeatures{
		WriteBytesPerSec: 100 * 1024 * 1024,
		WriteOpsPerSec:   100,
	}

	got := BuildProfile(f)

	if got.ReadWriteBalance.Value != MostlyWrite {
		t.Errorf(
			"ReadWriteBalance.Value = %v, want %v",
			got.ReadWriteBalance.Value,
			MostlyWrite,
		)
	}

	if got.IOSize.Value != SizeLarge {
		t.Errorf(
			"IOSize.Value = %v, want %v",
			got.IOSize.Value,
			SizeLarge,
		)
	}

	if got.DataIntensity.Value != LevelHigh {
		t.Errorf(
			"DataIntensity.Value = %v, want %v",
			got.DataIntensity.Value,
			LevelHigh,
		)
	}
}

func TestBuildProfileBalancedIO(t *testing.T) {
	f := features.WorkloadFeatures{
		ReadBytesPerSec:  32 * 1024 * 1024,
		WriteBytesPerSec: 32 * 1024 * 1024,
		ReadOpsPerSec:    512,
		WriteOpsPerSec:   512,
	}

	got := BuildProfile(f)

	if got.ReadWriteBalance.Value != Balanced {
		t.Errorf(
			"ReadWriteBalance.Value = %v, want %v",
			got.ReadWriteBalance.Value,
			Balanced,
		)
	}

	if got.IOSize.Value != SizeMedium {
		t.Errorf(
			"IOSize.Value = %v, want %v",
			got.IOSize.Value,
			SizeMedium,
		)
	}

	if got.DataIntensity.Value != LevelHigh {
		t.Errorf(
			"DataIntensity.Value = %v, want %v",
			got.DataIntensity.Value,
			LevelHigh,
		)
	}
}

func TestBuildProfileMetadataIntensity(t *testing.T) {
	tests := []struct {
		name      string
		opsPerSec float64
		want      Level
	}{
		{
			name:      "none",
			opsPerSec: 0,
			want:      LevelNone,
		},
		{
			name:      "very low",
			opsPerSec: 0.5,
			want:      LevelVeryLow,
		},
		{
			name:      "low",
			opsPerSec: 5,
			want:      LevelLow,
		},
		{
			name:      "medium",
			opsPerSec: 50,
			want:      LevelMedium,
		},
		{
			name:      "high",
			opsPerSec: 500,
			want:      LevelHigh,
		},
		{
			name:      "very high",
			opsPerSec: 1500,
			want:      LevelVeryHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := features.WorkloadFeatures{
				MetadataOpsPerSec: tt.opsPerSec,
			}

			got := BuildProfile(f)

			if got.MetadataIntensity.Value != tt.want {
				t.Errorf(
					"MetadataIntensity.Value = %v, want %v",
					got.MetadataIntensity.Value,
					tt.want,
				)
			}
		})
	}
}
