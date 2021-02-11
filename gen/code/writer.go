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
	if w.d > 0 {
		fmt.Fprint(w.buf, strings.Repeat(w.indent, w.d))
	}
}

func (w *Writer) Print(a ...interface{}) {
	w.printIndent()
	fmt.Fprint(w.buf, a...)
}

func (w *Writer) Println(a ...interface{}) {
	w.printIndent()
	fmt.Fprintln(w.buf, a...)
}

func (w *Writer) Printf(format string, a ...interface{}) {
	w.printIndent()
	fmt.Fprintf(w.buf, format, a...)
}
