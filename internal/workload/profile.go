package workload

type Assessment[T any] struct {
	Value      T
	Confidence float64
	Evidence   []string
}

type Level int

const (
	LevelNone Level = iota
	LevelVeryLow
	LevelLow
	LevelMedium
	LevelHigh
	LevelVeryHigh
)

type Balance int

const (
	BalanceUnknown Balance = iota
	MostlyRead
	ReadHeavy
	Balanced
	WriteHeavy
	MostlyWrite
)

type IOSize int

const (
	SizeUnknown IOSize = iota
	SizeTiny
	SizeSmall
	SizeMedium
	SizeLarge
	SizeHuge
)

type AccessPattern int

const (
	AccessUnknown AccessPattern = iota
	Sequential
	Random
	Mixed
)

type NamespaceProfile int

const (
	NamespaceUnknown NamespaceProfile = iota
	NamespaceIdle
	LookupHeavy
	CreateHeavy
	DeleteHeavy
	RenameHeavy
	MixedNamespace
)

type CacheProfile int

const (
	CacheUnknown CacheProfile = iota
	MostlyCached
	MixedCache
	MostlyMisses
)

type Profile struct {
	ReadWriteBalance  Assessment[Balance]
	IOSize            Assessment[IOSize]
	AccessPattern     Assessment[AccessPattern]
	DataIntensity     Assessment[Level]
	MetadataIntensity Assessment[Level]
	NamespaceBehavior Assessment[NamespaceProfile]
	CacheBehavior     Assessment[CacheProfile]
}

func (v Level) String() string {
	switch v {
	case LevelNone:
		return "none"
	case LevelVeryLow:
		return "very_low"
	case LevelLow:
		return "low"
	case LevelMedium:
		return "medium"
	case LevelHigh:
		return "high"
	case LevelVeryHigh:
		return "very_high"
	default:
		return "unknown"
	}
}
