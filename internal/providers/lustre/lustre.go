package lustre

import (
	"fmt"

	"github.com/ecce-machina/machina-trace/internal/collectors"
	lustrecollectors "github.com/ecce-machina/machina-trace/internal/collectors/lustre"
	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

type Provider struct {
	collectors []collectors.Collector
}

func New() *Provider {
	return NewWithLctl("lctl")
}

func NewWithLctl(lctlPath string) *Provider {
	return &Provider{
		collectors: []collectors.Collector{
			lustrecollectors.NewLLiteCollector(lctlPath),
			lustrecollectors.NewMDCMDStatsCollector(lctlPath),
			lustrecollectors.NewMDCStatsCollector(lctlPath),
			lustrecollectors.NewOSCStatsCollector(lctlPath),
		},
	}
}

func (p *Provider) Collect() ([]snapshot.Source, error) {
	var sources []snapshot.Source

	for _, collector := range p.collectors {
		collected, err := collector.Collect()
		if err != nil {
			return nil, fmt.Errorf(
				"collector %s failed: %w",
				collector.Name(),
				err,
			)
		}

		sources = append(sources, collected...)
	}

	return sources, nil
}
