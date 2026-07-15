package lustre

import (
	"errors"
	"reflect"
	"testing"
)

type fakeRunner struct {
	result commandResult
	err    error
	name   string
	args   []string
}

func (r *fakeRunner) Run(name string, args ...string) (commandResult, error) {
	r.name = name
	r.args = append([]string(nil), args...)
	return r.result, r.err
}

func TestMDCMDStatsCollector(t *testing.T) {
	collector := NewMDCMDStatsCollector("/usr/sbin/lctl")
	runner := &fakeRunner{result: commandResult{stdout: []byte(`mdc.lustrefs-MDT0000-mdc-abc.md_stats=
snapshot_time             1784136426.641291855 secs.nsecs
start_time                1783705564.337158373 secs.nsecs
elapsed_time              430862.304133482 secs.nsecs
create                    1 samples [reqs]
intent_lock               218 samples [reqs]
`)}}
	collector.runner = runner

	sources, err := collector.Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if runner.name != "/usr/sbin/lctl" || !reflect.DeepEqual(runner.args, []string{"get_param", "mdc.*.md_stats"}) {
		t.Fatalf("unexpected command: %q %v", runner.name, runner.args)
	}
	if len(sources) != 1 {
		t.Fatalf("len(sources) = %d, want 1", len(sources))
	}
	got := sources[0]
	if got.Collector != mdcMDStatsCollectorName || got.Object != "lustrefs-MDT0000-mdc-abc" {
		t.Fatalf("unexpected source identity: %+v", got)
	}
	if got.TimestampNS != 1784136426641291855 {
		t.Fatalf("TimestampNS = %d", got.TimestampNS)
	}
	if got.Counters["create"] != 1 || got.Counters["intent_lock"] != 218 {
		t.Fatalf("unexpected counters: %#v", got.Counters)
	}
}

func TestMDCStatsCollectorPreservesCountAndLatency(t *testing.T) {
	collector := NewMDCStatsCollector("lctl")
	collector.runner = &fakeRunner{result: commandResult{stdout: []byte(`mdc.lustrefs-MDT0000-mdc-abc.stats=
snapshot_time             1784136426.680044958 secs.nsecs
req_waittime              29095 samples [usecs] 454 1005807 30046183 2057647607287
req_active                29095 samples [reqs] 1 3 29099 29109
`)}}

	sources, err := collector.Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	counters := sources[0].Counters
	if counters["req_waittime"] != 29095 || counters["req_waittime_usecs"] != 30046183 {
		t.Fatalf("unexpected req_waittime counters: %#v", counters)
	}
	if counters["req_active"] != 29095 {
		t.Fatalf("req_active = %d", counters["req_active"])
	}
}

func TestOSCStatsCollectorParsesMultipleTargets(t *testing.T) {
	collector := NewOSCStatsCollector("lctl")
	collector.runner = &fakeRunner{result: commandResult{stdout: []byte(`osc.lustrefs-OST0000-osc-abc.stats=
snapshot_time             1784136426.715420479 secs.nsecs
req_waittime              2 samples [usecs] 1411 2698 4109 9270125
osc.lustrefs-OST0001-osc-abc.stats=
snapshot_time             1784136426.715451336 secs.nsecs
req_active                2 samples [reqs] 1 1 2 2
`)}}

	sources, err := collector.Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(sources) != 2 {
		t.Fatalf("len(sources) = %d, want 2", len(sources))
	}
	if sources[0].Object != "lustrefs-OST0000-osc-abc" || sources[1].Object != "lustrefs-OST0001-osc-abc" {
		t.Fatalf("unexpected objects: %q, %q", sources[0].Object, sources[1].Object)
	}
}

func TestStatsCollectorAcceptsParsedStdoutWithCommandError(t *testing.T) {
	collector := NewMDCMDStatsCollector("lctl")
	collector.runner = &fakeRunner{
		result: commandResult{stdout: []byte(`mdc.lustrefs-MDT0000-mdc-abc.md_stats=
snapshot_time 1784136426.1 secs.nsecs
create 1 samples [reqs]
`)},
		err: errors.New("exit status 2"),
	}

	sources, err := collector.Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(sources) != 1 || sources[0].Counters["create"] != 1 {
		t.Fatalf("unexpected sources: %#v", sources)
	}
}
