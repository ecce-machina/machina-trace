package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ecce-machina/machina-trace/internal/providers/linux"
	lustreprovider "github.com/ecce-machina/machina-trace/internal/providers/lustre"
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
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  machina-agent snapshot\n")
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

	lustreSources, err := lustreprovider.New().Collect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lustre provider unavailable: %v\n", err)
	} else {
		sources = append(sources, lustreSources...)
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


