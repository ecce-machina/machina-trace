package collectors

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

type MeminfoCollector struct {
	path string
}

func NewMeminfoCollector(path string) *MeminfoCollector {
	return &MeminfoCollector{
		path: path,
	}
}

func (c *MeminfoCollector) Name() string {
	return "proc_meminfo"
}

func (c *MeminfoCollector) Collect() ([]snapshot.Source, error) {

	f, err := os.Open(c.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	counters := make(map[string]int64)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())

		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")

		val, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}

		// Most values in /proc/meminfo are reported in kB.
		if len(fields) >= 3 && fields[2] == "kB" {
			val *= 1024
			key += "_bytes"
		}

		counters[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return []snapshot.Source{
		{
			Collector:   "proc_meminfo",
			TimestampNS: time.Now().UnixNano(),
			Counters:    counters,
		},
	}, nil
}
