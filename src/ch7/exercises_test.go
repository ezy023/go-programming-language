package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestWordCounter(t *testing.T) {
	b := []byte("one two three four")
	var w WordCounter
	w.Write(b)
	if w != 4 {
		t.Errorf("Expected %d words, got %d\n", 4, w)
	}
}

func TestLineCounter(t *testing.T) {
	b := []byte("one\ntwo\nthree")
	var l LineCounter
	l.Write(b)
	if l != 3 {
		t.Errorf("Expected %d lines, got %d\n", 3, l)
	}
}

func TestCountingWriter(t *testing.T) {
	b := new(bytes.Buffer) // new returns pointer type
	cw, count := CountingWriter(b)
	cw.Write([]byte("abcdef"))
	if b.Len() != 6 {
		t.Errorf("Expected b.Len() to be %d, actual %d\n", 6, b.Len())
	}
	if *count != 6 {
		t.Errorf("Expected %d bytes to have been written, actual %d\n", 6, *count)
	}
}

func TestMyLimitReader(t *testing.T) {
	br := bytes.NewReader([]byte("123456789"))
	lr := NewLimitReader(br, 5)
	b, _ := ioutil.ReadAll(lr)
	if len(b) > 5 {
		t.Errorf("Read more than %d bytes, actual %d bytes", 5, len(b))
	}
}
