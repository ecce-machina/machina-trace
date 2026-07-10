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

		readsCompleted, _ := strconv.ParseInt(fields[3], 10, 64)
		readsMerged, _ := strconv.ParseInt(fields[4], 10, 64)
		sectorsRead, _ := strconv.ParseInt(fields[5], 10, 64)
		readTimeMS, _ := strconv.ParseInt(fields[6], 10, 64)

		writesCompleted, _ := strconv.ParseInt(fields[7], 10, 64)
		writesMerged, _ := strconv.ParseInt(fields[8], 10, 64)
		sectorsWritten, _ := strconv.ParseInt(fields[9], 10, 64)
		writeTimeMS, _ := strconv.ParseInt(fields[10], 10, 64)

		iosInProgress, _ := strconv.ParseInt(fields[11], 10, 64)
		ioTimeMS, _ := strconv.ParseInt(fields[12], 10, 64)
		weightedIOTimeMS, _ := strconv.ParseInt(fields[13], 10, 64)

		source := snapshot.Source{
			Collector:   "proc_diskstats",
			Object:      device,
			TimestampNS: time.Now().UnixNano(),
			Counters: map[string]int64{
				"reads_completed":     readsCompleted,
				"reads_merged":        readsMerged,
				"sectors_read":        sectorsRead,
				"read_time_ms":        readTimeMS,
				"writes_completed":    writesCompleted,
				"writes_merged":       writesMerged,
				"sectors_written":     sectorsWritten,
				"write_time_ms":       writeTimeMS,
				"ios_in_progress":     iosInProgress,
				"io_time_ms":          ioTimeMS,
				"weighted_io_time_ms": weightedIOTimeMS,
			},
		}

		sources = append(sources, source)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}
