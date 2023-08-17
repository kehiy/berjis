package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type Foo struct{}

func (f Foo) String() string {
	return "foo"
}

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	Trace("a", "ok", "!ok")
	Info("b", nil)
	Info("b", "a", nil)
	Info("c", "b", []byte{1, 2, 3})
	Warn("d", "x")
	Error("e", "y", Foo{})

	out := buf.String()

	fmt.Println(out)
	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "010203")
	assert.Contains(t, out, "!INVALID-KEY!")
	assert.Contains(t, out, "!MISSING-VALUE!")
	assert.Contains(t, out, "null")
	assert.Contains(t, out, "trace")
	assert.NotContains(t, out, "debug")
	assert.Contains(t, out, "info")
	assert.Contains(t, out, "warn")
	assert.Contains(t, out, "error")
}
