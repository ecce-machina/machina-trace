package lustre

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

const collectorName = "lustre_llite_stats"

type commandRunner interface {
	Output(name string, args ...string) ([]byte, error)
}

type execRunner struct{}

func (execRunner) Output(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

type LLiteCollector struct {
	lctlPath string
	runner   commandRunner
}

func NewLLiteCollector(lctlPath string) *LLiteCollector {
	return &LLiteCollector{
		lctlPath: lctlPath,
		runner:   execRunner{},
	}
}

func (c *LLiteCollector) Name() string {
	return collectorName
}

func (c *LLiteCollector) Collect() ([]snapshot.Source, error) {
	output, err := c.runner.Output(
		c.lctlPath,
		"get_param",
		"llite.*.stats",
	)
	if err != nil {
		return nil, fmt.Errorf("read llite stats: %w", err)
	}

	sources, err := parseLLiteStats(output)
	if err != nil {
		return nil, fmt.Errorf("parse llite stats: %w", err)
	}

	return sources, nil
}

func parseLLiteStats(data []byte) ([]snapshot.Source, error) {
	var sources []snapshot.Source
	var current *snapshot.Source

	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if object, ok := parseHeader(line); ok {
			sources = append(sources, snapshot.Source{
				Collector: collectorName,
				Object:    object,
				Counters:  make(map[string]int64),
			})

			current = &sources[len(sources)-1]
			continue
		}

		if current == nil {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "snapshot_time":
			timestampNS, err := parseTimestampNS(fields[1])
			if err != nil {
				return nil, fmt.Errorf(
					"invalid snapshot_time %q: %w",
					fields[1],
					err,
				)
			}

			current.TimestampNS = timestampNS

		case "start_time", "elapsed_time":
			continue

		default:
			parseCounter(current.Counters, fields)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan llite stats: %w", err)
	}

	return sources, nil
}

func parseHeader(line string) (string, bool) {
	const (
		prefix = "llite."
		suffix = ".stats="
	)

	if !strings.HasPrefix(line, prefix) {
		return "", false
	}

	if !strings.HasSuffix(line, suffix) {
		return "", false
	}

	object := strings.TrimSuffix(
		strings.TrimPrefix(line, prefix),
		suffix,
	)

	if object == "" {
		return "", false
	}

	return object, true
}

func parseCounter(counters map[string]int64, fields []string) {
	if len(fields) < 4 {
		return
	}

	if fields[2] != "samples" {
		return
	}

	name := fields[0]

	samples, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return
	}

	unit := strings.Trim(fields[3], "[]")

	switch unit {
	case "bytes":
		counters[name+"_samples"] = samples

		// Lustre byte statistics are normally:
		//
		// name samples samples [bytes] min max sum sumsquare
		if len(fields) >= 7 {
			sum, err := strconv.ParseInt(fields[6], 10, 64)
			if err == nil {
				counters[name] = sum
			}
		}

	case "usecs":
		// For operation statistics, the sample count is the
		// cumulative number of operations.
		counters[name] = samples

		// Preserve cumulative latency when present.
		if len(fields) >= 7 {
			sum, err := strconv.ParseInt(fields[6], 10, 64)
			if err == nil {
				counters[name+"_usecs"] = sum
			}
		}

	default:
		counters[name] = samples
	}
}

func parseTimestampNS(value string) (int64, error) {
	parts := strings.SplitN(value, ".", 2)

	seconds, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}

	var nanoseconds int64

	if len(parts) == 2 {
		fraction := parts[1]

		if len(fraction) > 9 {
			fraction = fraction[:9]
		}

		fraction += strings.Repeat("0", 9-len(fraction))

		nanoseconds, err = strconv.ParseInt(fraction, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	return seconds*1_000_000_000 + nanoseconds, nil
}
