package features

import "github.com/ecce-machina/machina-trace/internal/diff"

type LustreClientIOFeatures struct {
	Object           string
	IntervalSec      float64
	ReadOpsPerSec    float64
	WriteOpsPerSec   float64
	ReadBytesPerSec  float64
	WriteBytesPerSec float64
	FsyncsPerSec     float64
	OpenOpsPerSec    float64
	CloseOpsPerSec   float64
	ReadUsecsPerSec  float64
	WriteUsecsPerSec float64
	FsyncUsecsPerSec float64
}

type LustreMetadataFeatures struct {
	Object                string
	IntervalSec           float64
	CreatesPerSec         float64
	ClosesPerSec          float64
	GetattrsPerSec        float64
	GetxattrsPerSec       float64
	SetattrsPerSec        float64
	RenamesPerSec         float64
	UnlinksPerSec         float64
	IntentLocksPerSec     float64
	RevalidateLocksPerSec float64
}

type LustreMDCFeatures struct {
	Object                 string
	IntervalSec            float64
	RequestsPerSec         float64
	RequestWaitUsecsPerSec float64
	LDLMEnqueuesPerSec     float64
	LDLMCancelsPerSec      float64
	MetadataClosesPerSec   float64
	MetadataSyncsPerSec    float64
}

type LustreOSTFeatures struct {
	Object                 string
	IntervalSec            float64
	ReadRPCsPerSec         float64
	WriteRPCsPerSec        float64
	ReadBytesPerSec        float64
	WriteBytesPerSec       float64
	RequestsPerSec         float64
	RequestWaitUsecsPerSec float64
	ExtentEnqueuesPerSec   float64
	SyncsPerSec            float64
}

func FromLustreLLiteDelta(d diff.CounterDelta) (LustreClientIOFeatures, bool) {
	if d.Collector != "lustre_llite_stats" {
		return LustreClientIOFeatures{}, false
	}

	return LustreClientIOFeatures{
		Object:           d.Object,
		IntervalSec:      d.IntervalSec,
		ReadOpsPerSec:    d.Rates["read"],
		WriteOpsPerSec:   d.Rates["write"],
		ReadBytesPerSec:  d.Rates["read_bytes"],
		WriteBytesPerSec: d.Rates["write_bytes"],
		FsyncsPerSec:     d.Rates["fsync"],
		OpenOpsPerSec:    d.Rates["open"],
		CloseOpsPerSec:   d.Rates["close"],
		ReadUsecsPerSec:  d.Rates["read_usecs"],
		WriteUsecsPerSec: d.Rates["write_usecs"],
		FsyncUsecsPerSec: d.Rates["fsync_usecs"],
	}, true
}

func FromLustreMDCMetadataDelta(d diff.CounterDelta) (LustreMetadataFeatures, bool) {
	if d.Collector != "lustre_mdc_md_stats" {
		return LustreMetadataFeatures{}, false
	}

	return LustreMetadataFeatures{
		Object:                d.Object,
		IntervalSec:           d.IntervalSec,
		CreatesPerSec:         d.Rates["create"],
		ClosesPerSec:          d.Rates["close"],
		GetattrsPerSec:        d.Rates["getattr"],
		GetxattrsPerSec:       d.Rates["getxattr"],
		SetattrsPerSec:        d.Rates["setattr"],
		RenamesPerSec:         d.Rates["rename"],
		UnlinksPerSec:         d.Rates["unlink"],
		IntentLocksPerSec:     d.Rates["intent_lock"],
		RevalidateLocksPerSec: d.Rates["revalidate_lock"],
	}, true
}

func FromLustreMDCDelta(d diff.CounterDelta) (LustreMDCFeatures, bool) {
	if d.Collector != "lustre_mdc_stats" {
		return LustreMDCFeatures{}, false
	}

	return LustreMDCFeatures{
		Object:                 d.Object,
		IntervalSec:            d.IntervalSec,
		RequestsPerSec:         d.Rates["req_waittime"],
		RequestWaitUsecsPerSec: d.Rates["req_waittime_usecs"],
		LDLMEnqueuesPerSec:     d.Rates["ldlm_ibits_enqueue"],
		LDLMCancelsPerSec:      d.Rates["ldlm_cancel"],
		MetadataClosesPerSec:   d.Rates["mds_close"],
		MetadataSyncsPerSec:    d.Rates["mds_sync"],
	}, true
}

func FromLustreOSCDelta(d diff.CounterDelta) (LustreOSTFeatures, bool) {
	if d.Collector != "lustre_osc_stats" {
		return LustreOSTFeatures{}, false
	}

	return LustreOSTFeatures{
		Object:                 d.Object,
		IntervalSec:            d.IntervalSec,
		ReadRPCsPerSec:         d.Rates["ost_read"],
		WriteRPCsPerSec:        d.Rates["ost_write"],
		ReadBytesPerSec:        d.Rates["read_bytes"],
		WriteBytesPerSec:       d.Rates["write_bytes"],
		RequestsPerSec:         d.Rates["req_waittime"],
		RequestWaitUsecsPerSec: d.Rates["req_waittime_usecs"],
		ExtentEnqueuesPerSec:   d.Rates["ldlm_extent_enqueue"],
		SyncsPerSec:            d.Rates["ost_sync"],
	}, true
}
