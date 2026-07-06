package collectors

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

type DiskstatsCollector struct {
	path string
}

func NewDiskstatsCollector(path string) *DiskstatsCollector {
	return &DiskstatsCollector{
		path: path,
	}
}

func (c *DiskstatsCollector) Name() string {
	return "proc_diskstats"
}

func (c *DiskstatsCollector) Collect() ([]snapshot.Source, error) {

	f, err := os.Open(c.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sources []snapshot.Source

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())

		if len(fields) < 14 {
			continue
		}

		device := fields[2]

		readIOs, _ := strconv.ParseInt(fields[3], 10, 64)
		readSectors, _ := strconv.ParseInt(fields[5], 10, 64)

		writeIOs, _ := strconv.ParseInt(fields[7], 10, 64)
		writeSectors, _ := strconv.ParseInt(fields[9], 10, 64)

		source := snapshot.Source{
			Collector:   "proc_diskstats",
			Object:      device,
			TimestampNS: time.Now().UnixNano(),
			Counters: map[string]int64{
				"reads_completed":  readIOs,
				"sectors_read":     readSectors,
				"writes_completed": writeIOs,
				"sectors_written":  writeSectors,
			},
		}

		sources = append(sources, source)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}
