package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestMain1StArgWrong(t *testing.T) {
	s := captureStdout(main_, "aaa", "bbb")
	assert.Contains(t, s, "aaa")
}

func TestMain2ndArgWrong(t *testing.T) {
	s := captureStdout(main_, "./examples/unit-testing/tf.template", "bbb")
	assert.Contains(t, s, "bbb")
}

func TestMainCorrectArgs(t *testing.T) {
	s := captureStdout(main_, "./examples/unit-testing/tf.template", "./examples/unit-testing/context.yml")
	assert.Contains(t, s, "map[context:map")
}


func TestTfStringListFormater(t *testing.T) {
	// var array []interface{}{1, 2, "4", 1.4}
	var result string

	result = tfStringListFormater(nil)
	assert.Equal(t, "[]", result)

	result = tfStringListFormater([]interface{}{})
	assert.Equal(t, "[]", result)

	result = tfStringListFormater([]interface{}{1})
	assert.Equal(t, "[\"1\"]", result)

	result = tfStringListFormater([]interface{}{1.3, "1.1.1.1", 2})
	assert.Equal(t, "[\"1.3\", \"1.1.1.1\", \"2\"]", result)
}

func captureStdout(f func(string, string) int, s1 string, s2 string) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f(s1, s2)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
