# Machina Trace

Machina Trace is an analysis framework for distributed filesystems.

It transforms low-level observations into analyses of workload, client behavior, distribution, and performance, helping storage administrators understand **what a filesystem is doing** and **why it behaves the way it does**.

Rather than focusing on synthetic benchmarks or isolated metrics, Machina Trace observes real systems and produces analyses that help explain filesystem behavior.

For example:

```text
Filesystem workload

  Predominantly sequential writes
  Large average I/O size
  Moderate metadata activity

Client distribution

  96 active clients
  Workload evenly distributed

Outliers

  node100 performs 50 lookups/sec
  Most other clients are create-heavy

Performance

  Metadata service shows elevated RPC latency.
  Data path remains healthy.
```

The goal is not simply to collect statistics.

The goal is to answer questions such as:

- What workload is my filesystem serving?
- Are clients behaving similarly or performing different jobs?
- Which nodes are unusual?
- Is the workload well balanced?
- Is the filesystem responding appropriately to the workload it is serving?

## Philosophy

Machina Trace separates analysis into independent, reusable stages.

Each stage has one responsibility:

- Observe
- Measure
- Interpret
- Compare
- Evaluate
- Present

Keeping these responsibilities separate makes analyses easier to test, extend, and reuse while avoiding duplicated logic.

For more information, see [docs/architecture.md](docs/architecture.md).

## Architecture

```
Operating System
        │
        ▼
Collectors
        │
        ▼
Snapshots
        │
        ▼
Diffs
        │
        ▼
Per-node Features
        │
        ├──────────────► Aggregate Features
        │                       │
        │                       ▼
        │              Workload Classification
        │                       │
        └──────────────► Distribution Analysis
                                │
                                ▼
                     Performance Analysis
                                │
                                ▼
                            Renderers
```

Each transformation has one canonical implementation.

## Current Capabilities

- Snapshot collection
- Snapshot diffing
- Linux system metrics
- Lustre client metrics
- Workload feature extraction
- Semantic workload classification

## Roadmap

- Cross-node workload analysis
- Filesystem-wide workload summaries
- Client outlier detection
- Workload diversity analysis
- Performance evaluation relative to observed workload
- Additional filesystem backends
- Additional collectors and analyses

## Project Goals

Machina Trace aims to become a reusable analysis engine rather than a collection of independent utilities.

As the project grows, new analyses should build upon existing observations, snapshots, diffs, and features instead of duplicating collection or interpretation logic.

## Status

Machina Trace is under active development.

The architecture is intentionally evolving toward a modular analysis pipeline that preserves information at each stage while making it straightforward to add new analyses over the same collected data.

