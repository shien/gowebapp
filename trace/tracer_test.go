package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("New からの戻り値が nil です")
	} else {
		tracer.Trace("こんにちは trace package")
		if buf.String() != "こんにちは trace pakage" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}

	}

}
