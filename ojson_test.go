package ojson

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testMarshal(t *testing.T, v interface{}, exp string) {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal failed: %s", err)
	}
	act := string(b)
	if act != exp {
		t.Fatalf("unexpected:\nexpect=%s\nactual=%s", exp, act)
	}
}

func TestObject_Marshal(t *testing.T) {
	testMarshal(t, nil, "null")
	testMarshal(t, Object{}, "{}")
	testMarshal(t, Object{
		{"foo", 123},
		{"bar", "zzz"},
	}, `{"foo":123,"bar":"zzz"}`)
	testMarshal(t, Object{
		{"bar", "zzz"},
		{"foo", 123},
	}, `{"bar":"zzz","foo":123}`)
	testMarshal(t, Object{
		{"foo", 123},
		{"bar", Object{
			{"baz", 999},
			{"qux", "zzz"},
		}},
	}, `{"foo":123,"bar":{"baz":999,"qux":"zzz"}}`)
}

func TestObject_Put(t *testing.T) {
	var o Object
	testMarshal(t, o, "null")
	o.Put("foo", 123).Put("bar", "zzz")
	testMarshal(t, o, `{"foo":123,"bar":"zzz"}`)
	o.Put("baz", 999)
	testMarshal(t, o, `{"foo":123,"bar":"zzz","baz":999}`)
	o.Put("bar", "xyz")
	testMarshal(t, o, `{"foo":123,"bar":"xyz","baz":999}`)
}

func TestObject_Delete(t *testing.T) {
	o := Object{
		{"foo", 123},
		{"bar", "zzz"},
		{"baz", 999},
	}
	testMarshal(t, o, `{"foo":123,"bar":"zzz","baz":999}`)
	o.Delete("bar")
	testMarshal(t, o, `{"foo":123,"baz":999}`)
	o.Delete("baz")
	testMarshal(t, o, `{"foo":123}`)
	o.Delete("foo")
	testMarshal(t, o, `{}`)
}

func objectUnmarshal(t *testing.T, data string, exp Object) {
	t.Helper()
	var act Object
	err := json.Unmarshal([]byte(data), &act)
	if err != nil {
		t.Fatalf("unmarshal failed: %s", err)
	}
	if diff := cmp.Diff(exp, act); diff != "" {
		t.Fatalf("unexpected unmarshal -:expected +:actual\n%s", diff)
	}
}

func TestObject_Unmarshal(t *testing.T) {
	objectUnmarshal(t, `null`, nil)
	objectUnmarshal(t, `{}`, Object{})
	objectUnmarshal(t, `{"foo":123,"bar":"zzz","baz":999}`, Object{
		{"foo", 123.0},
		{"bar", "zzz"},
		{"baz", 999.0},
	})
	objectUnmarshal(t, `{"foo":123,"bar":{"baz":999,"qux":"zzz"}}`, Object{
		{"foo", 123.0},
		{"bar", Object{
			{"baz", 999.0},
			{"qux", "zzz"},
		}},
	})
}

func testGet(t *testing.T, o Object, k string, expV interface{}, expOK bool) {
	t.Helper()
	actV, actOK := o.Get(k)
	if actOK != expOK {
		t.Fatalf("unexpected OK: expect=%t actual=%t", expOK, actOK)
	}
	if diff := cmp.Diff(expV, actV); diff != "" {
		t.Fatalf("unexpected get -:expected +:actual\n%s", diff)
	}
}

func TestObject_Get(t *testing.T) {
	o := Object{
		{"foo", 123},
		{"bar", "zzz"},
		{"baz", 999},
	}
	testGet(t, o, "foo", 123, true)
	testGet(t, o, "bar", "zzz", true)
	testGet(t, o, "baz", 999, true)
	testGet(t, o, "qux", nil, false)
}

func TestArray_Marshal(t *testing.T) {
	testMarshal(t, nil, "null")
	testMarshal(t, Array{}, "[]")
	testMarshal(t, Array{1, 2, 3}, "[1,2,3]")
	testMarshal(t, Array{
		Object{{"id", 1}, {"name", "foo"}},
		Object{{"id", 2}, {"name", "bar"}},
		Object{{"id", 3}, {"name", "baz"}},
	}, `[{"id":1,"name":"foo"},{"id":2,"name":"bar"},{"id":3,"name":"baz"}]`)
	testMarshal(t, Array{
		Object{{"id", 1}, {"name", "foo"}},
		Object{{"name", "bar"}, {"id", 2}},
		Object{{"id", 3}, {"name", "baz"}},
	}, `[{"id":1,"name":"foo"},{"name":"bar","id":2},{"id":3,"name":"baz"}]`)
}

func TestArray_Add(t *testing.T) {
	var a Array
	testMarshal(t, a, "null")
	a.Add(1)
	testMarshal(t, a, `[1]`)
	a.Add(2).Add(3)
	testMarshal(t, a, `[1,2,3]`)
	a.Add("xxx", "yyy", "zzz")
	testMarshal(t, a, `[1,2,3,"xxx","yyy","zzz"]`)
}

func arrayUnmarshal(t *testing.T, data string, exp Array) {
	t.Helper()
	var act Array
	err := json.Unmarshal([]byte(data), &act)
	if err != nil {
		t.Fatalf("unmarshal failed: %s", err)
	}
	if diff := cmp.Diff(exp, act); diff != "" {
		t.Fatalf("unexpected unmarshal -:expected +:actual\n%s", diff)
	}
}

func TestArray_Unmarshal(t *testing.T) {
	arrayUnmarshal(t, `null`, nil)
	arrayUnmarshal(t, `[]`, Array{})
	arrayUnmarshal(t, `[1,2,3]`, Array{1.0, 2.0, 3.0})
	arrayUnmarshal(t,
		`[{"id":1,"name":"foo"},{"id":2,"name":"bar"},{"id":3,"name":"baz"}]`,
		Array{
			Object{{"id", 1.0}, {"name", "foo"}},
			Object{{"id", 2.0}, {"name", "bar"}},
			Object{{"id", 3.0}, {"name", "baz"}},
		})
}
