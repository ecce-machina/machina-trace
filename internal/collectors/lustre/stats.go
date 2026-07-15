package lustre

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

const (
	mdcMDStatsCollectorName = "lustre_mdc_md_stats"
	mdcStatsCollectorName   = "lustre_mdc_stats"
	oscStatsCollectorName   = "lustre_osc_stats"
)

type statsCollector struct {
	name      string
	parameter string
	prefix    string
	suffix    string
	lctlPath  string
	runner    commandRunner
}

func NewMDCMDStatsCollector(lctlPath string) *statsCollector {
	return newStatsCollector(mdcMDStatsCollectorName, "mdc.*.md_stats", "mdc.", ".md_stats=", lctlPath)
}

func NewMDCStatsCollector(lctlPath string) *statsCollector {
	return newStatsCollector(mdcStatsCollectorName, "mdc.*.stats", "mdc.", ".stats=", lctlPath)
}

func NewOSCStatsCollector(lctlPath string) *statsCollector {
	return newStatsCollector(oscStatsCollectorName, "osc.*.stats", "osc.", ".stats=", lctlPath)
}

func newStatsCollector(name, parameter, prefix, suffix, lctlPath string) *statsCollector {
	return &statsCollector{
		name:      name,
		parameter: parameter,
		prefix:    prefix,
		suffix:    suffix,
		lctlPath:  lctlPath,
		runner:    execRunner{},
	}
}

func (c *statsCollector) Name() string { return c.name }

func (c *statsCollector) Collect() ([]snapshot.Source, error) {
	result, runErr := c.runner.Run(c.lctlPath, "get_param", c.parameter)
	if len(result.stdout) == 0 {
		if runErr != nil {
			return nil, fmt.Errorf(
				"read %s: %w: %s",
				c.parameter,
				runErr,
				strings.TrimSpace(string(result.stderr)),
			)
		}
		return nil, fmt.Errorf("read %s: empty output", c.parameter)
	}

	sources, err := parseParameterizedStats(result.stdout, c.name, c.prefix, c.suffix)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", c.parameter, err)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("read %s: no matching sources found", c.parameter)
	}

	return sources, nil
}

func parseParameterizedStats(data []byte, collector, prefix, suffix string) ([]snapshot.Source, error) {
	var sources []snapshot.Source
	var current *snapshot.Source

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if object, ok := parseParameterizedHeader(line, prefix, suffix); ok {
			sources = append(sources, snapshot.Source{
				Collector: collector,
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
				return nil, fmt.Errorf("invalid snapshot_time %q: %w", fields[1], err)
			}
			current.TimestampNS = timestampNS
		case "start_time", "elapsed_time":
			continue
		default:
			parseCounter(current.Counters, fields)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan stats: %w", err)
	}

	return sources, nil
}

func parseParameterizedHeader(line, prefix, suffix string) (string, bool) {
	if !strings.HasPrefix(line, prefix) || !strings.HasSuffix(line, suffix) {
		return "", false
	}

	object := strings.TrimSuffix(strings.TrimPrefix(line, prefix), suffix)
	if object == "" {
		return "", false
	}
	return object, true
}
