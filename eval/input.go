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

func newInput(input io.Reader) (*Input, error) {
	p := Input{}
	sc := bufio.NewScanner(input)
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = scanTokens(data, atEOF)
		r, _ := utf8.DecodeRune(token)
		if r == '\n' {
			p.Cursor.Ln++
			p.Cursor.Col = 0
		} else {
			p.Cursor.Col += advance
		}
		return
	})
	p.Scanner = sc
	p.Cursor = Cursor{1, 0}
	return &p, nil
}

func (p *Input) readBool() (bool, error) {
	b, err := p.next()
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(string(b))
	if err != nil {
		return false, ErrBadParse{Want: "bool", Got: b, Cursor: p.Cursor}
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
		return 0, ErrBadParse{Want: "int", Got: b, Cursor: p.Cursor}
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
		return 0, ErrBadParse{Want: "int64", Got: b, Cursor: p.Cursor}
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
		return 0, ErrBadParse{Want: "float32", Got: b, Cursor: p.Cursor}
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
		return 0, ErrBadParse{Want: "float64", Got: b, Cursor: p.Cursor}
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
	b := []byte{}
	for {
		err := p.scan()
		if err != nil {
			return "", err
		}
		t := p.Scanner.Bytes()
		r, _ := utf8.DecodeRune(t)
		if r == '\n' {
			break
		}
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = append(b, t...)
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
	err := p.scan()
	if err != nil {
		return nil, err
	}
	b := p.Scanner.Bytes()
	return b, nil
}

func (p *Input) scan() error {
	if !p.Scanner.Scan() {
		err := p.Scanner.Err()
		if err == nil {
			err = io.EOF
		}
		return err
	}
	return nil
}

func (p *Input) skipSpace() (bool, error) {
	col := p.Cursor.Col
	err := p.scan()
	if err != nil {
		return false, err
	}
	b := p.Scanner.Bytes()
	r, _ := utf8.DecodeRune(b)
	return col > 0 && r == ' ', nil
}

// scanTokens is a split function for a Scanner that returns spaces and words
// separately. It consumes a blank space (' ') if it appears immediately
// after a word. The definition of space is set by isSpace.
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
			if r == ' ' {
				return i + 1, data[:i], nil
			}
			return i, data[:i], nil
		}
	}
	// If we're not at EOF, request more data.
	if atEOF && len(data) > 0 {
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
