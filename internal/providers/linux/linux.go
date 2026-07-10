package linux

import (
	"fmt"

	"github.com/ecce-machina/machina-trace/internal/collectors"
	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

type Provider struct {
	collectors []collectors.Collector
}

func New() *Provider {
	return &Provider{
		collectors: []collectors.Collector{
			collectors.NewMeminfoCollector("/proc/meminfo"),
			collectors.NewVMStatCollector("/proc/vmstat"),
			collectors.NewDiskstatsCollector("/proc/diskstats"),
			collectors.NewMountinfoCollector("/proc/self/mountinfo"),
			collectors.NewMountinfoCollector("/proc/net/dev"),
		},
	}
}

func (p *Provider) Collect() ([]snapshot.Source, error) {
	var sources []snapshot.Source

	for _, collector := range p.collectors {
		result, err := collector.Collect()
		if err != nil {
			return nil, fmt.Errorf(
				"collector %s failed: %w",
				collector.Name(),
				err,
			)
		}

		sources = append(sources, result...)
	}

	return sources, nil
}
