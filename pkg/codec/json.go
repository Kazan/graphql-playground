package codec

import (
	"encoding/json"
	"io"
)

func NewJSONCodec() *jsonCodec {
	return &jsonCodec{}
}

type jsonCodec struct{}

func (*jsonCodec) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(&v)
}

func (*jsonCodec) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(&v)
}
