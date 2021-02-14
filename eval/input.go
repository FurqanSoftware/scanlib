package eval

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode/utf8"
)

type Input struct {
	Scanner *bufio.Scanner
	Cursor  Cursor
}

func NewInput(input io.Reader) (*Input, error) {
	sc := bufio.NewScanner(input)
	sc.Split(ScanTokens)
	p := &Input{
		Scanner: sc,
	}
	return p, nil
}

func (p *Input) Bool() (bool, error) {
	b, err := p.next()
	if err != nil {
		return false, p.Scanner.Err()
	}
	return strconv.ParseBool(string(b))
}

func (p *Input) Int() (int, error) {
	b, err := p.next()
	if err != nil {
		return 0, p.Scanner.Err()
	}
	n, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func (p *Input) Int64() (int64, error) {
	b, err := p.next()
	if err != nil {
		return 0, p.Scanner.Err()
	}
	return strconv.ParseInt(string(b), 10, 64)
}

func (p *Input) Float32() (float32, error) {
	b, err := p.next()
	if err != nil {
		return 0, p.Scanner.Err()
	}
	f, err := strconv.ParseFloat(string(b), 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func (p *Input) Float64() (float64, error) {
	b, err := p.next()
	if err != nil {
		return 0, p.Scanner.Err()
	}
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return 0, err
	}
	return float64(f), nil
}

func (p *Input) String() (string, error) {
	b, err := p.next()
	if err != nil {
		return "", p.Scanner.Err()
	}
	return string(b), nil
}

func (p *Input) EOL() (bool, error) {
	b, err := p.next()
	if err != nil {
		return false, p.Scanner.Err()
	}
	return bytes.Equal(b, []byte("\n")), nil
}

func (p *Input) EOF() (bool, error) {
	b, err := p.next()
	if len(b) == 0 && err == nil {
		return true, nil
	}
	return false, err
}

func (p *Input) next() ([]byte, error) {
	err := p.skipSpace()
	if err != nil {
		return nil, err
	}
	b := p.Scanner.Bytes()
	p.pushCursor(b)
	return b, nil
}

func (p *Input) skipSpace() error {
	if !p.Scanner.Scan() {
		return p.Scanner.Err()
	}
	b := p.Scanner.Bytes()
	p.pushCursor(b)
	r, _ := utf8.DecodeRune(b)
	if r == ' ' {
		if !p.Scanner.Scan() {
			return p.Scanner.Err()
		}
	}
	return nil
}

func (p *Input) pushCursor(b []byte) {
	r, _ := utf8.DecodeRune(b)
	if r == '\n' {
		p.Cursor.Ln++
		p.Cursor.Col = 0
	} else {
		p.Cursor.Col = len(b)
	}
}

// ScanTokens is a split function for a Scanner that returns each
// space-separated word of text, including the space that follows it. It will
// never return an empty string. The definition of space is set by
// isSpace.
func ScanTokens(data []byte, atEOF bool) (advance int, token []byte, err error) {
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
	// If we're not at EOF, request more data.
	if !atEOF {
		return 0, nil, nil
	}
	// Return the final token.
	return len(data), data, bufio.ErrFinalToken
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
