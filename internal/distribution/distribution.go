package distribution

import "github.com/ecce-machina/machina-trace/internal/aggregate"

type Result struct {
	WriteConcentration float64
}

func Analyze(fs aggregate.FilesystemFeatures) Result {
	var totalWrites float64
	var busiestNodeWrites float64

	for _, node := range fs.Nodes {
		totalWrites += node.WriteBytesPerSec

		if node.WriteBytesPerSec > busiestNodeWrites {
			busiestNodeWrites = node.WriteBytesPerSec
		}
	}

	if totalWrites == 0 {
		return Result{}
	}

	return Result{
		WriteConcentration: busiestNodeWrites / totalWrites,
	}
}
