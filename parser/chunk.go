package parser

import (
	"errors"
	"fmt"

	"github.com/johnsudaar/midi"
)

func parseChunk(p *Parser) (parseFunc, error) {
	if p.pos >= len(p.input)-1 {
		return nil, nil
	}

	t, err := p.ConsumeN(4)
	if err != nil {
		return nil, err
	}

	length, err := p.NextUint32()
	if err != nil {
		return nil, err
	}

	chunk := chunk{
		Type:   string(t),
		Length: length,
	}

	fmt.Println(chunk)
	switch string(t) {
	case "MThd":
		return chunk.parseHeader, nil
	default:
		return chunk.ignore, nil
	}
	return nil, nil
}

type chunk struct {
	Type   string
	Length uint32
}

func (c chunk) String() string {
	return fmt.Sprintf("Chunk %s:%v", c.Type, c.Length)
}

func (c chunk) ignore(p *Parser) (parseFunc, error) {
	err := p.Skip(int(c.Length))
	return parseChunk, err
}

func (c chunk) parseHeader(p *Parser) (parseFunc, error) {
	if c.Length != 6 {
		return nil, fmt.Errorf("Invalid header length: %v", c.Length)
	}

	format, err := p.NextUint16()
	if err != nil {
		return nil, err
	}

	switch format {
	case 0:
		p.file.Format = midi.SingleMultiChannelTrack
	case 1:
		p.file.Format = midi.MultipleSimultaneousTracks
	case 2:
		p.file.Format = midi.SequentialIndependantSingleTrack
	default:
		return nil, fmt.Errorf("Invalid file format: %v", format)
	}

	tracks, err := p.NextUint16()
	if err != nil {
		return nil, err
	}

	p.file.TracksCount = tracks

	timingData, err := p.NextUint16()
	if err != nil {
		return nil, err
	}

	if timingData&0x8000 == 0 { // Ticks per quarter note mode
		p.file.Timing = midi.Timing{
			TicksPerQuarterNote: timingData & 0x7fff,
		}
	} else { // SMPTE timing
		//TODO: Add SMPTE support
		return nil, errors.New("SMPTE is not supported")
	}

	return parseChunk, nil
}
