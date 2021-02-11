// Copyright 2020 Furqan Software Ltd. All rights reserved.

package code

import (
	"bytes"
	"fmt"
	"strings"
)

type Writer struct {
	buf *bytes.Buffer

	indent string
	d      int

	r, c int
}

func NewWriter(indent string) *Writer {
	return &Writer{
		buf:    &bytes.Buffer{},
		indent: indent,
	}
}

func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Writer) Indent(d int) {
	w.d += d
}

func (w *Writer) printIndent() {
	if w.d > 0 && w.c == 0 {
		n, _ := fmt.Fprint(w.buf, strings.Repeat(w.indent, w.d))
		w.c += n
	}
}

func (w *Writer) Print(a ...interface{}) {
	w.printIndent()
	n, _ := fmt.Fprint(w.buf, a...)
	w.c += n
}

func (w *Writer) Printf(format string, a ...interface{}) {
	w.printIndent()
	n, _ := fmt.Fprintf(w.buf, format, a...)
	w.c += n
}

func (w *Writer) Println(a ...interface{}) {
	w.printIndent()
	fmt.Fprintln(w.buf, a...)
	w.r++
	w.c = 0
}
