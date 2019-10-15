package ojson

import (
	"bytes"
	"encoding/json"
)

// KV is a pair of key and value for Object
type KV struct {
	K string
	V interface{}
}

// Object represents an ordered JSON object.
type Object []*KV

func (o Object) key(k string) int {
	for i, kv := range o {
		if kv.K == k {
			return i
		}
	}
	return -1
}

// Put puts a KV.
func (o *Object) Put(k string, v interface{}) *Object {
	n := o.key(k)
	if n < 0 {
		*o = append(*o, &KV{k, v})
		return o
	}
	(*o)[n].V = v
	return o
}

// Delete deletes a key if exists.
func (o *Object) Delete(k string) *Object {
	x := o.key(k)
	if x < 0 {
		return o
	}
	n := len(*o)
	copy((*o)[x:n-1], (*o)[x+1:n])
	(*o)[n-1] = nil
	*o = (*o)[:n-1]
	return o
}

// Get gets a value for the "k" key.
func (o Object) Get(k string) (interface{}, bool) {
	n := o.key(k)
	if n < 0 {
		return nil, false
	}
	return o[n].V, true
}

// MarshalJSON implements json.Marshaler
func (o Object) MarshalJSON() ([]byte, error) {
	if o == nil {
		return []byte("null"), nil
	}
	bb := &bytes.Buffer{}
	enc := json.NewEncoder(bb)
	bb.WriteRune('{')
	for i, kv := range o {
		if i > 0 {
			bb.WriteRune(',')
		}
		enc.Encode(kv.K)
		bb.WriteRune(':')
		enc.Encode(kv.V)
	}
	bb.WriteRune('}')
	return bb.Bytes(), nil
}

// Array represents JSON array.
type Array []interface{}

// Add appends elements at last of the Array.
func (a *Array) Add(values ...interface{}) *Array {
	*a = append(*a, values...)
	return a
}
