package diff

import "github.com/ecce-machina/machina-trace/internal/snapshot"

type CounterDelta struct {
	Node        string
	Collector   string
	Object      string
	Mount       string
	StartNS     int64
	EndNS       int64
	IntervalSec float64
	Deltas      map[string]int64
	Rates       map[string]float64
}

func DiffSnapshots(a, b *snapshot.Snapshot) []CounterDelta {
	old := make(map[string]snapshot.Source)

	for _, src := range a.Sources {
		key := sourceKey(src)
		old[key] = src
	}

	var out []CounterDelta

	for _, newer := range b.Sources {
		key := sourceKey(newer)

		older, ok := old[key]
		if !ok {
			continue
		}

		interval := float64(newer.TimestampNS-older.TimestampNS) / 1_000_000_000
		if interval <= 0 {
			continue
		}

		deltas := make(map[string]int64)
		rates := make(map[string]float64)

		for name, newVal := range newer.Counters {
			oldVal, ok := older.Counters[name]
			if !ok {
				continue
			}

			delta := newVal - oldVal
			if delta < 0 {
				continue
			}

			deltas[name] = delta
			rates[name] = float64(delta) / interval
		}

		out = append(out, CounterDelta{
			Node:        b.Node,
			Collector:   newer.Collector,
			Object:      newer.Object,
			Mount:       newer.Mount,
			StartNS:     older.TimestampNS,
			EndNS:       newer.TimestampNS,
			IntervalSec: interval,
			Deltas:      deltas,
			Rates:       rates,
		})
	}

	return out
}

func sourceKey(src snapshot.Source) string {
	return src.Collector + "|" + src.Object + "|" + src.Mount
}
