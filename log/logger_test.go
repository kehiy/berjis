package log_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kehiy/berjis/log"
	zlog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type Foo struct{}

func (f Foo) String() string {
	return "foo"
}

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	zlog.Logger = zlog.Output(&buf)

	log.Trace("a", "ok", "!ok")
	log.Info("b", nil)
	log.Info("b", "a", nil)
	log.Info("c", "b", []byte{1, 2, 3})
	log.Warn("d", "x")
	log.Error("e", "y", Foo{})

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
