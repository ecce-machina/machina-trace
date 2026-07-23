# Machina Trace Architecture

## Philosophy

Machina Trace is a pipeline of analysis stages.

Each stage has one responsibility:

- consume information from the previous stage
- produce a richer representation
- never duplicate work already performed
- never destroy information needed by later stages

Every reusable transformation should have one canonical implementation.

---

# Pipeline

```
Operating System
        ↓
Collectors
        ↓
Snapshots
        ↓
Diffs
        ↓
Per-node Features
        ↓
Filesystem Aggregate
        ↓
Workload Classification
        ↓
Distribution Analysis
        ↓
Performance Analysis
        ↓
Renderers
```

Each arrow represents a reusable domain transformation.

---

# Collectors

Collectors preserve observations.

Examples:

- /proc
- sysfs
- Lustre statistics

Collectors answer:

> What did the system report?

Collectors never interpret.

---

# Snapshots

Snapshots capture a point in time.

Snapshots answer:

> What existed at this instant?

They preserve counters, timestamps and identifiers.

---

# Diffs

Diffs compare snapshots.

They answer:

> What changed during the observation interval?

Diffs compute:

- deltas
- elapsed time
- rates directly implied by counters

Diffs do not classify workloads.

---

# Features

Features describe measurable behaviour.

Examples:

- read operations/sec
- write operations/sec
- average write size
- metadata operations/sec
- lookup rate
- bytes/sec

A feature answers:

> What measurable behaviour occurred?

Features remain factual.

---

# Per-node Features

Node identity is preserved.

Every client receives its own feature set.

```
node001
    WorkloadFeatures

node002
    WorkloadFeatures

node003
    WorkloadFeatures
```

No aggregation should destroy node identity.

---

# Aggregate Features

Filesystem-wide features are derived from all nodes.

This answers:

> What workload is the filesystem serving overall?

Aggregate features are **not** averages of classifications.

They are computed from the underlying measurements.

Example:

```
Overall write bandwidth

=
Σ node write bandwidth
```

---

# Workload Classification

Classification converts measurable behaviour into semantic descriptions.

Examples:

- read dominated
- write dominated
- balanced
- small I/O
- large I/O
- metadata intensive

Classification produces assessments.

```
Assessment

- value
- confidence
- evidence
```

Classification consumes Features.

Classification never reads raw counters.

---

# Distribution Analysis

Distribution compares nodes.

It answers questions such as:

- Is the workload evenly distributed?
- Are there outliers?
- Are clients doing different jobs?
- What is the typical client?

Examples:

```
Node100 performs
50 lookups/sec

Most nodes perform
40 creates/sec
```

or

```
One client generates
85% of write bandwidth.
```

Distribution analysis is concerned with variation.

It should preserve:

- active node count
- percentiles
- dominant behaviours
- diversity
- concentration

---

# Performance Analysis

Performance is interpreted in the context of the workload.

The workload itself does not determine whether the filesystem performs well.

Performance analysis combines:

- workload
- observed response

Examples:

- latency
- RPC queues
- lock contention
- throughput
- retransmits
- server saturation

Performance analysis answers:

> Given this workload,
> is the filesystem behaving as expected?

---

# Renderers

Renderers present completed analyses.

Renderers may:

- choose formatting
- choose units
- group information
- sort information

Renderers never:

- derive features
- classify workloads
- compare nodes
- evaluate performance

Renderers answer:

> How should this information be presented?

---

# Information Preservation

Each stage enriches information.

No stage should destroy information needed later.

For example:

Per-node features must exist before aggregate features.

Aggregate features must not replace per-node features.

Node identity is preserved until rendering.

---

# Canonical Ownership

| Stage | Package |
|---------|----------|
| Collect | collectors |
| Snapshot | snapshot |
| Diff | diff |
| Features | features |
| Classification | workload |
| Distribution | distribution |
| Performance | performance |
| Rendering | render |

Every transformation has exactly one owner.

---

# Design Rule

If a new feature requires copying logic from another package,
the architecture is drifting.

The correct solution is to expose the existing transformation rather than duplicate it.

