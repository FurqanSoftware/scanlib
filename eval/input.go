package eval

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Input struct {
	Scanner *bufio.Scanner

	Line *bytes.Reader
	EOF  bool

	Cursor Cursor
}

func NewInput(input io.Reader) (*Input, error) {
	sc := bufio.NewScanner(input)
	sc.Split(bufio.ScanLines)
	p := &Input{
		Scanner: sc,
	}
	err := p.Next()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Input) Next() error {
	if p.Scanner.Scan() {
		p.Line = bytes.NewReader(p.Scanner.Bytes())
		p.Cursor.Ln++
		p.Cursor.Col = 0
		return nil
	}
	err := p.Scanner.Err()
	if err == nil {
		p.EOF = true
	} else {
		return err
	}
	return nil
}

func (p *Input) Scanf(format string, a ...interface{}) (int, error) {
	return fmt.Fscanf(p.Line, format, a...)
}

func (p *Input) EOL() (bool, error) {
	if p.Line.Len() != 0 {
		return false, nil
	}
	err := p.Next()
	if err != nil {
		return false, err
	}
	return true, nil
}
