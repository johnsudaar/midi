package midi

type Format int

const (
	SingleMultiChannelTrack          Format = 0
	MultipleSimultaneousTracks       Format = 1
	SequentialIndependantSingleTrack Format = 1
)

type File struct {
	TracksCount uint16
	Format      Format
	Timing      Timing
}

type Timing struct {
	TicksPerQuarterNote uint16
}
