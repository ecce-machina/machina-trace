package collectors

import "github.com/ecce-machina/machina-trace/internal/snapshot"

type Collector interface {
	Name() string
	Collect() ([]snapshot.Source, error)
}
