package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	os.Args = []string{"terraformer",
		"generate",
		"examples/unit-testing/tf.template",
		"examples/unit-testing/context.yml"}

	s := captureStdout(main)

	assert.Contains(t, s, "map[context")
	assert.Contains(t, s, "[\"1.1.1.1\", \"2.2.2.2\", \"3.3.3.3\"]")
}

func TestGenerateContext(t *testing.T) {
	os.Args = []string{"terraformer",
		"generate-context",
		"--state=examples/unit-testing/tf.state.json",
		"--callback=examples/unit-testing/callback-yaml.sh",
		"--template=examples/unit-testing/context.yml.template",
		"examples/unit-testing/config1.yml",
		"examples/unit-testing/config2.yml"}

	s := captureStdout(main)

	assert.Contains(t, s, "context:")

	assert.Contains(t, s, "env:")
	assert.Contains(t, s, "HOME:")

	assert.Contains(t, s, "state:")
	assert.Contains(t, s, "terraform_version:")

	assert.Contains(t, s, "callback:")
	assert.Contains(t, s, "callback-yaml-value")

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

func TestTfCallback(t *testing.T) {
	// var array []interface{}{1, 2, "4", 1.4}
	var result string

	result = tfCallback("callback-value.sh")
	assert.Equal(t, "4711", result)
}

func captureStdout(f func()) string {
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
