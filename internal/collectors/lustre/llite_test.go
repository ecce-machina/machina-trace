package lustre

import "testing"

func TestParseLLiteStats(t *testing.T) {
	input := []byte(`
llite.lustrefs-ffff9bc043537800.stats=
snapshot_time             1783960164.279418700 secs.nsecs
start_time                1783705564.331211471 secs.nsecs
elapsed_time              254599.948207229 secs.nsecs
getattr                   1 samples [usecs] 7 7 7 49
statfs                    16977 samples [usecs] 2 3904 16809097 16988546393
getxattr                  2 samples [usecs] 594 796 1390 986452
inode_permission          5 samples [usecs] 2 3137 3657 10096913
read_bytes                10 samples [bytes] 4096 1048576 5242880 123456
write_bytes               4 samples [bytes] 4096 1048576 2097152 654321
`)

	sources, err := parseLLiteStats(input)
	if err != nil {
		t.Fatalf("parse llite stats: %v", err)
	}

	if len(sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(sources))
	}

	got := sources[0]

	if got.Collector != collectorName {
		t.Errorf(
			"expected collector %q, got %q",
			collectorName,
			got.Collector,
		)
	}

	if got.Object != "lustrefs-ffff9bc043537800" {
		t.Errorf(
			"expected object %q, got %q",
			"lustrefs-ffff9bc043537800",
			got.Object,
		)
	}

	const expectedTimestamp int64 = 1783960164279418700

	if got.TimestampNS != expectedTimestamp {
		t.Errorf(
			"expected timestamp %d, got %d",
			expectedTimestamp,
			got.TimestampNS,
		)
	}

	expectedCounters := map[string]int64{
		"getattr":                1,
		"getattr_usecs":          7,
		"statfs":                 16977,
		"statfs_usecs":           16809097,
		"getxattr":               2,
		"getxattr_usecs":         1390,
		"inode_permission":       5,
		"inode_permission_usecs": 3657,
		"read_bytes_samples":     10,
		"read_bytes":             5242880,
		"write_bytes_samples":    4,
		"write_bytes":            2097152,
	}

	for name, expected := range expectedCounters {
		actual, ok := got.Counters[name]
		if !ok {
			t.Errorf("expected counter %q", name)
			continue
		}

		if actual != expected {
			t.Errorf(
				"counter %q: expected %d, got %d",
				name,
				expected,
				actual,
			)
		}
	}
}

func TestParseLLiteStatsMultipleInstances(t *testing.T) {
	input := []byte(`
llite.fs1-aaaa.stats=
snapshot_time 10.000000001 secs.nsecs
getattr 2 samples [usecs] 1 2 3 4
llite.fs2-bbbb.stats=
snapshot_time 11.000000002 secs.nsecs
statfs 3 samples [usecs] 1 2 4 8
`)

	sources, err := parseLLiteStats(input)
	if err != nil {
		t.Fatalf("parse llite stats: %v", err)
	}

	if len(sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(sources))
	}

	if sources[0].Object != "fs1-aaaa" {
		t.Errorf("unexpected first object %q", sources[0].Object)
	}

	if sources[1].Object != "fs2-bbbb" {
		t.Errorf("unexpected second object %q", sources[1].Object)
	}
}

func TestParseTimestampNS(t *testing.T) {
	got, err := parseTimestampNS("1783960164.279418700")
	if err != nil {
		t.Fatalf("parse timestamp: %v", err)
	}

	const expected int64 = 1783960164279418700

	if got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}
