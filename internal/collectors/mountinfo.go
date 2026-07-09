package collectors

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

type MountinfoCollector struct {
	path string
}

func NewMountinfoCollector(path string) *MountinfoCollector {
	return &MountinfoCollector{
		path: path,
	}
}

func (c *MountinfoCollector) Name() string {
	return "proc_mountinfo"
}

func (c *MountinfoCollector) Collect() ([]snapshot.Source, error) {
	f, err := os.Open(c.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sources []snapshot.Source

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		left, right, ok := strings.Cut(line, " - ")
		if !ok {
			continue
		}

		leftFields := strings.Fields(left)
		rightFields := strings.Fields(right)

		if len(leftFields) < 5 || len(rightFields) < 3 {
			continue
		}

		mountID := leftFields[0]
		parentID := leftFields[1]
		majorMinor := leftFields[2]
		root := leftFields[3]
		mountPoint := leftFields[4]

		fsType := rightFields[0]
		source := rightFields[1]

		sources = append(sources, snapshot.Source{
			Collector:   "proc_mountinfo",
			Object:      mountPoint,
			Mount:       mountPoint,
			TimestampNS: time.Now().UnixNano(),
			Values: map[string]string{
				"mount_id":    mountID,
				"parent_id":   parentID,
				"major_minor": majorMinor,
				"root":        root,
				"mount_point": mountPoint,
				"fs_type":     fsType,
				"source":      source,
			},
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}
