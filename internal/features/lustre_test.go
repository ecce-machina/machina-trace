package features

import (
	"reflect"
	"testing"

	"github.com/ecce-machina/machina-trace/internal/diff"
)

func TestFromLustreLLiteDelta(t *testing.T) {
	d := diff.CounterDelta{
		Collector:   "lustre_llite_stats",
		Object:      "lustrefs-ffff9bc043537800",
		IntervalSec: 2.5,
		Rates: map[string]float64{
			"read":        10,
			"write":       20,
			"read_bytes":  30,
			"write_bytes": 40,
			"fsync":       50,
			"open":        60,
			"close":       70,
			"read_usecs":  80,
			"write_usecs": 90,
			"fsync_usecs": 100,
		},
	}

	got, ok := FromLustreLLiteDelta(d)
	if !ok {
		t.Fatal("FromLustreLLiteDelta() returned false for lustre_llite_stats")
	}

	want := LustreClientIOFeatures{
		Object:           "lustrefs-ffff9bc043537800",
		IntervalSec:      2.5,
		ReadOpsPerSec:    10,
		WriteOpsPerSec:   20,
		ReadBytesPerSec:  30,
		WriteBytesPerSec: 40,
		FsyncsPerSec:     50,
		OpenOpsPerSec:    60,
		CloseOpsPerSec:   70,
		ReadUsecsPerSec:  80,
		WriteUsecsPerSec: 90,
		FsyncUsecsPerSec: 100,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"FromLustreLLiteDelta() mismatch:\ngot:  %+v\nwant: %+v",
			got,
			want,
		)
	}
}

func TestFromLustreMDCMetadataDelta(t *testing.T) {
	d := diff.CounterDelta{
		Collector:   "lustre_mdc_md_stats",
		Object:      "lustrefs-MDT0000-mdc-ffff9bc043537800",
		IntervalSec: 5,
		Rates: map[string]float64{
			"create":          1,
			"close":           2,
			"getattr":         3,
			"getxattr":        4,
			"setattr":         5,
			"rename":          6,
			"unlink":          7,
			"intent_lock":     8,
			"revalidate_lock": 9,
		},
	}

	got, ok := FromLustreMDCMetadataDelta(d)
	if !ok {
		t.Fatal(
			"FromLustreMDCMetadataDelta() returned false for lustre_mdc_md_stats",
		)
	}

	want := LustreMetadataFeatures{
		Object:                "lustrefs-MDT0000-mdc-ffff9bc043537800",
		IntervalSec:           5,
		CreatesPerSec:         1,
		ClosesPerSec:          2,
		GetattrsPerSec:        3,
		GetxattrsPerSec:       4,
		SetattrsPerSec:        5,
		RenamesPerSec:         6,
		UnlinksPerSec:         7,
		IntentLocksPerSec:     8,
		RevalidateLocksPerSec: 9,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"FromLustreMDCMetadataDelta() mismatch:\ngot:  %+v\nwant: %+v",
			got,
			want,
		)
	}
}

func TestFromLustreMDCDelta(t *testing.T) {
	d := diff.CounterDelta{
		Collector:   "lustre_mdc_stats",
		Object:      "lustrefs-MDT0000-mdc-ffff9bc043537800",
		IntervalSec: 10,
		Rates: map[string]float64{
			"req_waittime":       11,
			"req_waittime_usecs": 12,
			"ldlm_ibits_enqueue": 13,
			"ldlm_cancel":        14,
			"mds_close":          15,
			"mds_sync":           16,
		},
	}

	got, ok := FromLustreMDCDelta(d)
	if !ok {
		t.Fatal("FromLustreMDCDelta() returned false for lustre_mdc_stats")
	}

	want := LustreMDCFeatures{
		Object:                 "lustrefs-MDT0000-mdc-ffff9bc043537800",
		IntervalSec:            10,
		RequestsPerSec:         11,
		RequestWaitUsecsPerSec: 12,
		LDLMEnqueuesPerSec:     13,
		LDLMCancelsPerSec:      14,
		MetadataClosesPerSec:   15,
		MetadataSyncsPerSec:    16,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"FromLustreMDCDelta() mismatch:\ngot:  %+v\nwant: %+v",
			got,
			want,
		)
	}
}

func TestFromLustreOSCDelta(t *testing.T) {
	d := diff.CounterDelta{
		Collector:   "lustre_osc_stats",
		Object:      "lustrefs-OST0000-osc-ffff9bc043537800",
		IntervalSec: 1.5,
		Rates: map[string]float64{
			"ost_read":            21,
			"ost_write":           22,
			"read_bytes":          23,
			"write_bytes":         24,
			"req_waittime":        25,
			"req_waittime_usecs":  26,
			"ldlm_extent_enqueue": 27,
			"ost_sync":            28,
		},
	}

	got, ok := FromLustreOSCDelta(d)
	if !ok {
		t.Fatal("FromLustreOSCDelta() returned false for lustre_osc_stats")
	}

	want := LustreOSTFeatures{
		Object:                 "lustrefs-OST0000-osc-ffff9bc043537800",
		IntervalSec:            1.5,
		ReadRPCsPerSec:         21,
		WriteRPCsPerSec:        22,
		ReadBytesPerSec:        23,
		WriteBytesPerSec:       24,
		RequestsPerSec:         25,
		RequestWaitUsecsPerSec: 26,
		ExtentEnqueuesPerSec:   27,
		SyncsPerSec:            28,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"FromLustreOSCDelta() mismatch:\ngot:  %+v\nwant: %+v",
			got,
			want,
		)
	}
}

func TestLustreFeatureFunctionsRejectWrongCollector(t *testing.T) {
	d := diff.CounterDelta{
		Collector: "some_other_collector",
		Rates:     map[string]float64{},
	}

	tests := []struct {
		name string
		run  func(diff.CounterDelta) bool
	}{
		{
			name: "llite",
			run: func(d diff.CounterDelta) bool {
				_, ok := FromLustreLLiteDelta(d)
				return ok
			},
		},
		{
			name: "mdc metadata",
			run: func(d diff.CounterDelta) bool {
				_, ok := FromLustreMDCMetadataDelta(d)
				return ok
			},
		},
		{
			name: "mdc rpc",
			run: func(d diff.CounterDelta) bool {
				_, ok := FromLustreMDCDelta(d)
				return ok
			},
		},
		{
			name: "osc",
			run: func(d diff.CounterDelta) bool {
				_, ok := FromLustreOSCDelta(d)
				return ok
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.run(d) {
				t.Fatalf("%s converter accepted the wrong collector", tt.name)
			}
		})
	}
}

func TestLustreFeatureFunctionsHandleMissingRates(t *testing.T) {
	d := diff.CounterDelta{
		Collector:   "lustre_llite_stats",
		Object:      "lustrefs",
		IntervalSec: 3,
		Rates:       map[string]float64{},
	}

	got, ok := FromLustreLLiteDelta(d)
	if !ok {
		t.Fatal("FromLustreLLiteDelta() returned false for lustre_llite_stats")
	}

	want := LustreClientIOFeatures{
		Object:      "lustrefs",
		IntervalSec: 3,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"missing rates should produce zero-valued features:\ngot:  %+v\nwant: %+v",
			got,
			want,
		)
	}
}
