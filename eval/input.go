package eval

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode/utf8"
)

type Input struct {
	sc *bufio.Scanner

	token []byte
	ahead [][]byte
	err   error

	cur, curnext Cursor
}

func newInput(input io.Reader) (*Input, error) {
	p := Input{
		curnext: Cursor{1, 0},
	}
	sc := bufio.NewScanner(input)
	sc.Split(scanTokens)
	p.sc = sc
	return &p, nil
}

func (p *Input) readBool() (bool, error) {
	b, err := p.next()
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(string(b))
	if err != nil {
		return false, ErrBadParse{Want: "bool", Got: b, Cursor: p.cur}
	}
	return v, nil
}

func (p *Input) readInt() (int, error) {
	b, err := p.next()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return 0, ErrBadParse{Want: "int", Got: b, Cursor: p.cur}
	}
	return int(n), nil
}

func (p *Input) readInt64() (int64, error) {
	b, err := p.next()
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0, ErrBadParse{Want: "int64", Got: b, Cursor: p.cur}
	}
	return n, nil
}

func (p *Input) readFloat32() (float32, error) {
	b, err := p.next()
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(string(b), 32)
	if err != nil {
		return 0, ErrBadParse{Want: "float32", Got: b, Cursor: p.cur}
	}
	return float32(f), nil
}

func (p *Input) readFloat64() (float64, error) {
	b, err := p.next()
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return 0, ErrBadParse{Want: "float64", Got: b, Cursor: p.cur}
	}
	return float64(f), nil
}

func (p *Input) readString() (string, error) {
	b, err := p.next()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (p *Input) readStringLn() (string, error) {
	b, err := p.nextLn()
	if err != nil {
		return "", nil
	}
	return string(b), nil
}

func (p *Input) isAtEOL() (bool, error) {
	b, err := p.next()
	if err != nil {
		return false, err
	}
	return bytes.Equal(b, []byte("\n")), nil
}

func (p *Input) isAtEOF() (bool, error) {
	_, err := p.next()
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func (p *Input) next() ([]byte, error) {
	p.scan()

	if len(p.ahead) > 0 {
		b := p.ahead[0]
		p.token = b

		copy(p.ahead, p.ahead[1:])
		p.ahead = p.ahead[:len(p.ahead)-1]

		r, _ := utf8.DecodeRune(b)
		p.cur = p.curnext
		if r == '\n' {
			p.curnext.Ln++
			p.curnext.Col = 0
		} else {
			p.curnext.Col += len(b)
		}

		if !isSpace(r) && len(p.ahead) == 2 {
			r0, _ := utf8.DecodeRune(p.ahead[0])
			r1, _ := utf8.DecodeRune(p.ahead[1])
			if r0 == ' ' && !isSpace(r1) {
				// Drop sandwiched space.
				copy(p.ahead, p.ahead[1:])
				p.ahead = p.ahead[:len(p.ahead)-1]

				p.curnext.Col++
			}
		}

		return b, nil
	}

	return nil, p.err
}

func (p *Input) nextLn() ([]byte, error) {
	if p.err != nil {
		return nil, p.err
	}

	p.cur = p.curnext

	b := []byte{}
	for {
		p.scan()

		t := p.ahead[0]
		copy(p.ahead, p.ahead[1:])
		p.ahead = p.ahead[:len(p.ahead)-1]

		r, _ := utf8.DecodeRune(t)
		if r == '\n' {
			p.curnext.Ln++
			p.curnext.Col = 0
			break
		} else {
			p.curnext.Col += len(t)
		}

		b = append(b, t...)
	}
	p.token = b

	return b, nil
}

func (p *Input) scan() {
	for len(p.ahead) < 3 && p.err == nil {
		if !p.sc.Scan() {
			p.err = p.sc.Err()
			if p.err == nil {
				p.err = io.EOF
			}
			return
		}
		b := p.sc.Bytes()
		t := make([]byte, len(b))
		copy(t, b)
		p.ahead = append(p.ahead, t)
	}
}

// scanTokens is a split function for a Scanner that returns spaces and words
// separately. The definition of space is set by isSpace.
func scanTokens(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Return any leading space.
	r, width := utf8.DecodeRune(data)
	if isSpace(r) {
		return width, data[:width], nil
	}
	// Scan until space, marking end of word.
	for i := width; i < len(data); i += width {
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) {
			return i, data[:i], nil
		}
	}
	if atEOF && len(data) > 0 {
		// If we're at EOF, return data.
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// isSpace reports whether the character is a Unicode white space character.
// We avoid dependency on the unicode package, but check validity of the implementation
// in the tests.
func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}
