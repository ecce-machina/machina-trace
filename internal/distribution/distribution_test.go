package distribution

import (
	"math"
	"testing"

	"github.com/ecce-machina/machina-trace/internal/aggregate"
	"github.com/ecce-machina/machina-trace/internal/features"
)

func TestAnalyzeWriteConcentration(t *testing.T) {
	fs := aggregate.FilesystemFeatures{
		Nodes: map[string]features.WorkloadFeatures{
			"client-01": {
				WriteBytesPerSec: 100,
			},
			"client-02": {
				WriteBytesPerSec: 300,
			},
		},
	}

	result := Analyze(fs)

	want := 0.75

	if math.Abs(result.WriteConcentration-want) > 1e-9 {
		t.Fatalf(
			"WriteConcentration = %v, want %v",
			result.WriteConcentration,
			want,
		)
	}
}
