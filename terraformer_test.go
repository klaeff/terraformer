package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestMainNoArgs(t *testing.T) {
	s := captureStdout(main_)
	assert.Contains(t, s, "usage")
}

func TestMain1StArgWrong(t *testing.T) {
	os.Args = []string{"cmd", "aaa", "bbb"}
	s := captureStdout(main_)
	assert.Contains(t, s, "aaa")
}

func TestMain2ndArgWrong(t *testing.T) {
	os.Args = []string{"cmd", "./examples/tf.template", "bbb"}
	s := captureStdout(main_)
	assert.Contains(t, s, "bbb")
}

func TestMainCorrectArgs(t *testing.T) {
	os.Args = []string{"cmd", "./examples/tf.template", "./examples/context.yml"}
	s := captureStdout(main_)
	assert.Contains(t, s, "map[context:map")
}

func captureStdout(f func() int) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

