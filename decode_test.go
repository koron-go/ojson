package ojson

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testUnmarshal(t *testing.T, data string, exp interface{}) {
	t.Helper()
	act, err := Unmarshal([]byte(data))
	if err != nil {
		t.Fatalf("unmarshal failed: %s", err)
	}
	if diff := cmp.Diff(exp, act); diff != "" {
		t.Fatalf("unexpected unmarshal -:expected +:actual\n%s", diff)
	}
}

func TestUnmarshal_Object(t *testing.T) {
	testUnmarshal(t, `null`, nil)
	testUnmarshal(t, `{}`, Object{})
	testUnmarshal(t, `{"foo":123,"bar":"zzz","baz":999}`, Object{
		{"foo", 123.0},
		{"bar", "zzz"},
		{"baz", 999.0},
	})
	testUnmarshal(t, `{"foo":123,"bar":{"baz":999,"qux":"zzz"}}`, Object{
		{"foo", 123.0},
		{"bar", Object{
			{"baz", 999.0},
			{"qux", "zzz"},
		}},
	})
}

func TestUnmarshal_Array(t *testing.T) {
	testUnmarshal(t, `null`, nil)
	testUnmarshal(t, `[]`, Array{})
	testUnmarshal(t, `[1,2,3]`, Array{1.0, 2.0, 3.0})
	testUnmarshal(t,
		`[{"id":1,"name":"foo"},{"id":2,"name":"bar"},{"id":3,"name":"baz"}]`,
		Array{
			Object{{"id", 1.0}, {"name", "foo"}},
			Object{{"id", 2.0}, {"name", "bar"}},
			Object{{"id", 3.0}, {"name", "baz"}},
		})
}

func testDecode(t *testing.T, data string, exp ...interface{}) {
	t.Helper()

	dec := NewDecoder(strings.NewReader(data))
	dec.UseNumber()
	var act []interface{}
	for {
		v, err := dec.Decode()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("decode failure: %s", err)
		}
		act = append(act, v)
	}

	if diff := cmp.Diff(exp, act); diff != "" {
		t.Fatalf("unexpected decode -:expected +:actual\n:%s", diff)
	}
}

func TestDecoder_Object(t *testing.T) {
	testDecode(t, `null`, nil)
	testDecode(t, `{}`, Object{})
	testDecode(t, `{"foo":123,"bar":"zzz","baz":999}`, Object{
		{"foo", json.Number("123")},
		{"bar", "zzz"},
		{"baz", json.Number("999")},
	})
	testDecode(t, `{"foo":123,"bar":{"baz":999,"qux":"zzz"}}`, Object{
		{"foo", json.Number("123")},
		{"bar", Object{
			{"baz", json.Number("999")},
			{"qux", "zzz"},
		}},
	})
}

func TestDecoder_Array(t *testing.T) {
	testDecode(t, `null`, nil)
	testDecode(t, `[]`, Array{})
	testDecode(t, `[1,2,3]`, Array{
		json.Number("1"),
		json.Number("2"),
		json.Number("3"),
	})
	testDecode(t,
		`[{"id":1,"name":"foo"},{"id":2,"name":"bar"},{"id":3,"name":"baz"}]`,
		Array{
			Object{{"id", json.Number("1")}, {"name", "foo"}},
			Object{{"id", json.Number("2")}, {"name", "bar"}},
			Object{{"id", json.Number("3")}, {"name", "baz"}},
		})
}

func TestDecoder_continuous(t *testing.T) {
	testDecode(t, `{}{}`, Object{}, Object{})
	testDecode(t, `{}{}{}`, Object{}, Object{}, Object{})
	testDecode(t, `{}[]{}`, Object{}, Array{}, Object{})
}
