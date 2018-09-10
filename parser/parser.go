package parser

import (
	"encoding/binary"
	"errors"
	"io/ioutil"

	"github.com/johnsudaar/midi"
)

var eof = errors.New("Unexpected EOF")

type parseFunc func(*Parser) (parseFunc, error)

type Parser struct {
	input []byte
	pos   int
	file  *midi.File
}

func New(file string) (*Parser, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return &Parser{
		input: content,
		pos:   0,
		file:  &midi.File{},
	}, nil
}

func (p *Parser) Consume() (byte, error) {
	if p.pos >= len(p.input)-1 {
		return 0, eof
	}
	res := p.input[p.pos]
	p.pos++
	return res, nil
}

func (p *Parser) ConsumeN(n int) ([]byte, error) {
	buffer := make([]byte, n)
	for i := 0; i < n; i++ {
		res, err := p.Consume()
		if err != nil {
			return buffer, err
		}
		buffer[i] = res
	}
	return buffer, nil
}

func (p *Parser) NextUint32() (uint32, error) {
	if p.pos+4 >= len(p.input) {
		return 0, eof
	}
	res := binary.BigEndian.Uint32(p.input[p.pos:])
	p.pos += 4
	return res, nil
}

func (p *Parser) NextUint16() (uint16, error) {
	if p.pos+2 >= len(p.input) {
		return 0, eof
	}

	res := binary.BigEndian.Uint16(p.input[p.pos:])
	p.pos += 2
	return res, nil
}

func (p *Parser) Skip(n int) error {
	if p.pos+n > len(p.input) {
		return eof
	}
	p.pos += n
	return nil
}

func (p *Parser) Parse() (*midi.File, error) {
	next := parseChunk
	var err error
	for {
		next, err = next(p)
		if err != nil {
			return nil, err
		}
		if next == nil {
			return p.file, nil
		}
	}
}
