package collectors

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

type NetdevCollector struct {
	path string
}

func NewNetdevCollector(path string) *NetdevCollector {
	return &NetdevCollector{
		path: path,
	}
}

func (c *NetdevCollector) Name() string {
	return "proc_net_dev"
}

func (c *NetdevCollector) Collect() ([]snapshot.Source, error) {
	f, err := os.Open(c.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sources []snapshot.Source

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip the two header lines.
		if line == "" || strings.HasPrefix(line, "Inter-|") ||
			strings.HasPrefix(line, "face |") {
			continue
		}

		left, right, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		device := strings.TrimSpace(left)
		fields := strings.Fields(right)

		// /proc/net/dev has 16 numeric fields per interface:
		// 8 receive fields, followed by 8 transmit fields.
		if len(fields) < 16 {
			continue
		}

		receiveBytes, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			continue
		}

		receivePackets, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}

		receiveErrors, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			continue
		}

		receiveDropped, err := strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			continue
		}

		transmitBytes, err := strconv.ParseInt(fields[8], 10, 64)
		if err != nil {
			continue
		}

		transmitPackets, err := strconv.ParseInt(fields[9], 10, 64)
		if err != nil {
			continue
		}

		transmitErrors, err := strconv.ParseInt(fields[10], 10, 64)
		if err != nil {
			continue
		}

		transmitDropped, err := strconv.ParseInt(fields[11], 10, 64)
		if err != nil {
			continue
		}

		sources = append(sources, snapshot.Source{
			Collector:   "proc_net_dev",
			Object:      device,
			TimestampNS: time.Now().UnixNano(),
			Counters: map[string]int64{
				"receive_bytes":    receiveBytes,
				"receive_packets":  receivePackets,
				"receive_errors":   receiveErrors,
				"receive_dropped":  receiveDropped,
				"transmit_bytes":   transmitBytes,
				"transmit_packets": transmitPackets,
				"transmit_errors":  transmitErrors,
				"transmit_dropped": transmitDropped,
			},
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}
