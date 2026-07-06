package snapshot

type Snapshot struct {
	SchemaVersion string   `json:"schema_version"`
	Node          string   `json:"node"`
	TimestampNS   int64    `json:"timestamp_ns"`
	Sources       []Source `json:"sources"`
}

type Source struct {
	Collector   string            `json:"collector"`
	Object      string            `json:"object,omitempty"`
	Mount       string            `json:"mount,omitempty"`
	TimestampNS int64             `json:"timestamp_ns"`
	Counters    map[string]int64  `json:"counters,omitempty"`
	Values      map[string]string `json:"values,omitempty"`
}
