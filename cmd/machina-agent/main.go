package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ecce-machina/machina-trace/internal/collectors"
	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

func main() {
	host, _ := os.Hostname()

	allCollectors := []collectors.Collector{
		collectors.NewMeminfoCollector("/proc/meminfo"),
		collectors.NewVMStatCollector("/proc/vmstat"),
		collectors.NewDiskstatsCollector("/proc/diskstats"),
	}

	var sources []snapshot.Source

	for _, c := range allCollectors {
		result, err := c.Collect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "collector %s failed: %v\n", c.Name(), err)
			continue
		}
		sources = append(sources, result...)
	}

	s := snapshot.Snapshot{
		SchemaVersion: "0.1",
		Node:          host,
		TimestampNS:   time.Now().UnixNano(),
		Sources:       sources,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if err := enc.Encode(s); err != nil {
		fmt.Fprintf(os.Stderr, "json encode failed: %v\n", err)
		os.Exit(1)
	}
}
