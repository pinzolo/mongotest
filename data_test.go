package mongotest_test

import (
	"testing"

	"github.com/pinzolo/mongotest"
)

func TestDocData_StringValue(t *testing.T) {
	data := mongotest.DocData{
		"foo": "foo",
		"bar": 1,
	}
	testdata := []struct {
		key  string
		val  string
		ok   bool
		memo string
	}{
		{key: "foo", val: "foo", ok: true, memo: "string value"},
		{key: "bar", val: "", ok: false, memo: "not string"},
		{key: "baz", val: "", ok: false, memo: "not exists"},
	}

	for _, d := range testdata {
		t.Run(d.memo, func(t *testing.T) {
			v, ok := data.StringValue(d.key)
			if d.ok != ok {
				t.Errorf("invalid result (want: %t, got: %t)", d.ok, ok)
			}
			if ok && v != d.val {
				t.Errorf("invalid value (want: %s, got: %s)", d.val, v)
			}
		})
	}
}
