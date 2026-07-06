package collectors

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

type VMStatCollector struct {
	path string
}

func NewVMStatCollector(path string) *VMStatCollector {
	return &VMStatCollector{
		path: path,
	}
}

func (c *VMStatCollector) Name() string {
	return "proc_vmstat"
}

func (c *VMStatCollector) Collect() ([]snapshot.Source, error) {

	f, err := os.Open(c.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	counters := make(map[string]int64)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())

		if len(fields) != 2 {
			continue
		}

		val, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}

		counters[fields[0]] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return []snapshot.Source{
		{
			Collector:   "proc_vmstat",
			TimestampNS: time.Now().UnixNano(),
			Counters:    counters,
		},
	}, nil
}
