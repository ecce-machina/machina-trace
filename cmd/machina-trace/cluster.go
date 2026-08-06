package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ecce-machina/machina-trace/internal/aggregate"
	"github.com/ecce-machina/machina-trace/internal/diff"
	"github.com/ecce-machina/machina-trace/internal/snapshot"
	"github.com/ecce-machina/machina-trace/internal/render"
	"github.com/ecce-machina/machina-trace/internal/workload"
)

func runCluster(beforeDir, afterDir string) {
	pairs, err := snapshotPairs(beforeDir, afterDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "discover snapshots: %v\n", err)
		os.Exit(1)
	}

	window, err := buildClusterWindow(pairs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build cluster window: %v\n", err)
		os.Exit(1)
	}

	fs := window.WorkloadFeatures()

	profile := workload.BuildProfile(fs.Aggregate)
	render.WriteWorkloadProfileText(os.Stdout, profile)

	fmt.Printf("loaded %d nodes\n", len(window.Nodes))
}

func buildClusterWindow(
	pairs []pair,
) (aggregate.ClusterWindow, error) {
	observations := make([]aggregate.NodeObservation, 0, len(pairs))

	var windowStart int64
	var windowEnd int64

	for _, p := range pairs {
		before, err := snapshot.ReadFile(p.Before)
		if err != nil {
			return aggregate.ClusterWindow{}, fmt.Errorf(
				"read before snapshot %q: %w",
				p.Before,
				err,
			)
		}

		after, err := snapshot.ReadFile(p.After)
		if err != nil {
			return aggregate.ClusterWindow{}, fmt.Errorf(
				"read after snapshot %q: %w",
				p.After,
				err,
			)
		}

		if before.Node != after.Node {
			return aggregate.ClusterWindow{}, fmt.Errorf(
				"snapshot node mismatch: before=%q after=%q",
				before.Node,
				after.Node,
			)
		}

		deltas := diff.DiffSnapshots(before, after)

		observation, ok := aggregate.NewNodeObservation(deltas)
		if !ok {
			return aggregate.ClusterWindow{}, fmt.Errorf(
				"build observation from %q and %q",
				p.Before,
				p.After,
			)
		}

		observations = append(observations, observation)

		if len(observations) == 1 || observation.StartNS < windowStart {
			windowStart = observation.StartNS
		}

		if len(observations) == 1 || observation.EndNS > windowEnd {
			windowEnd = observation.EndNS
		}
	}

	window := aggregate.NewClusterWindow(windowStart, windowEnd)

	for _, observation := range observations {
		if err := window.Add(observation); err != nil {
			return aggregate.ClusterWindow{}, fmt.Errorf(
				"add node %q: %w",
				observation.Node,
				err,
			)
		}
	}

	return window, nil
}

type pair struct {
	Before string
	After  string
}

func snapshotPairs(
	beforeDir,
	afterDir string,
) ([]pair, error) {
	entries, err := os.ReadDir(beforeDir)
	if err != nil {
		return nil, err
	}

	pairs := make([]pair, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".json" {
			continue
		}

		beforePath := filepath.Join(beforeDir, name)
		afterPath := filepath.Join(afterDir, name)

		if _, err := os.Stat(afterPath); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		pairs = append(pairs, pair{
			Before: beforePath,
			After:  afterPath,
		})
	}

	return pairs, nil
}
