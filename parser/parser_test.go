package parser

import (
	"testing"

	"github.com/johnsudaar/midi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParserTwinkle(t *testing.T) {
	parser, err := New("./fixtures/twinkle.mid")
	require.NoError(t, err)
	file, err := parser.Parse()
	require.NoError(t, err)
	assert.Equal(t, uint16(1), file.TracksCount)
	assert.Equal(t, midi.SingleMultiChannelTrack, file.Format)
	assert.Equal(t, uint16(128), file.Timing.TicksPerQuarterNote)
}
