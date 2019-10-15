package ojson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Unmarshal unmarshals/decodes data as JSON.
func Unmarshal(data []byte) (interface{}, error) {
	d := NewDecoder(bytes.NewReader(data))
	return d.Decode()
}

// Decoder reads and decodes JSON vlaues from an input stream
type Decoder struct {
	d   *json.Decoder
	err error
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		d: json.NewDecoder(r),
	}
}

// UseNumber causes the Decoder to unmarshal a number into as json.Number
// instead of as a flot64.
func (d *Decoder) UseNumber() {
	d.d.UseNumber()
}

// Decode decodes a JSON value.
func (d *Decoder) Decode() (interface{}, error) {
	tok, err := d.token()
	if err != nil {
		return nil, err
	}
	return d.decode(tok)
}

func (d *Decoder) token() (json.Token, error) {
	if d.err != nil {
		return nil, d.err
	}
	tok, err := d.d.Token()
	if err != nil {
		d.err = err
		return nil, err
	}
	return tok, nil
}

func (d *Decoder) decode(tok json.Token) (interface{}, error) {
	v, ok := tok.(json.Delim)
	if !ok {
		return tok, nil
	}
	switch v {
	case '[':
		a, err := d.decodeArray()
		if err != nil {
			return nil, err
		}
		return a, nil
	case '{':
		o, err := d.decodeObject()
		if err != nil {
			return nil, err
		}
		return o, nil
	default:
		return nil, fmt.Errorf("unexpected delim %+v", v)
	}
}

func (d *Decoder) decodeObject() (Object, error) {
	o := Object{}
	for {
		tok, err := d.token()
		if err != nil {
			return nil, err
		}
		if tok == json.Delim('}') {
			break
		}
		k, ok := tok.(string)
		if !ok {
			return nil, fmt.Errorf("not key where expected key: %+v", tok)
		}
		v, err := d.Decode()
		if err != nil {
			return nil, err
		}
		o.Put(k, v)
	}
	return o, nil
}

func (d *Decoder) decodeArray() (Array, error) {
	a := Array{}
	for {
		tok, err := d.token()
		if err != nil {
			return nil, err
		}
		if tok == json.Delim(']') {
			break
		}
		v, err := d.decode(tok)
		if err != nil {
			return nil, err
		}
		a.Add(v)
	}
	return a, nil
}
