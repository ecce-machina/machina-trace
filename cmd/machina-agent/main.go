package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ecce-machina/machina-trace/internal/diff"
	"github.com/ecce-machina/machina-trace/internal/providers/linux"
	"github.com/ecce-machina/machina-trace/internal/render"
	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "snapshot":
		runSnapshot()
	case "diff":
		if len(os.Args) == 4 {
			runDiff(os.Args[2], os.Args[3], false)
		} else if len(os.Args) == 5 && os.Args[2] == "--raw" {
			runDiff(os.Args[3], os.Args[4], true)
		} else {
			usage()
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  machina-agent snapshot\n")
	fmt.Fprintf(os.Stderr, "  machina-agent diff before.json after.json\n")
	fmt.Fprintf(os.Stderr, "  machina-agent diff --raw before.json after.json\n")
}

func runSnapshot() {
	host, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to determine hostname: %v\n", err)
		os.Exit(1)
	}

	provider := linux.New()

	sources, err := provider.Collect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "linux provider failed: %v\n", err)
		os.Exit(1)
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

func runDiff(beforePath, afterPath string, raw bool) {
	before, err := snapshot.ReadFile(beforePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", beforePath, err)
		os.Exit(1)
	}

	after, err := snapshot.ReadFile(afterPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", afterPath, err)
		os.Exit(1)
	}

	deltas := diff.DiffSnapshots(before, after)

	render.WriteDiskFeaturesText(os.Stdout, deltas)
	if raw {
		render.WriteDiffText(os.Stdout, deltas)
	}

}
