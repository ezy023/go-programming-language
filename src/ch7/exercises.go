package main

import (
	"bufio"
	"bytes"
	"io"
)

// Ex 7.1 word and line counter writer
type WordCounter int
type LineCounter int

func (w *WordCounter) Write(p []byte) (int, error) {
	r := bytes.NewReader(p)
	scan := bufio.NewScanner(r)
	scan.Split(bufio.ScanWords)
	for scan.Scan() {
		*w += 1
	}
	return len(p), nil
}

func (l *LineCounter) Write(p []byte) (int, error) {
	r := bytes.NewReader(p)
	scan := bufio.NewScanner(r)
	scan.Split(bufio.ScanLines)
	for scan.Scan() {
		*l += 1
	}
	return len(p), nil
}

// Ex 7.2
type byteCounter struct {
	w     io.Writer
	count int64
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	b := byteCounter{w, 0}
	return &b, &b.count
}

func (b *byteCounter) Write(p []byte) (int, error) {
	b.count += int64(len(p))
	return b.w.Write(p)
}

// Ex 7.5
type MyLimitReader struct {
	r   io.Reader
	rem int64
}

func NewLimitReader(r io.Reader, n int64) io.Reader {
	lr := MyLimitReader{r, n}
	return &lr
}

func (l *MyLimitReader) Read(p []byte) (n int, err error) {
	if l.rem-int64(len(p)) > 0 {
		n, err = l.r.Read(p)
		if err != nil {
			return n, err
		}
		l.rem -= int64(n)
	} else {
		n, err = l.r.Read(p[:l.rem])
		return n, io.EOF
	}
	return
}
