package textdb

import "encoding/json"

type jsonMarshaller struct{}

var _ Marshaller = &jsonMarshaller{}

func NewJSONMarshaller() Marshaller {
	return &jsonMarshaller{}
}

func (m *jsonMarshaller) Marshal(val interface{}) ([]byte, error) {
	return json.MarshalIndent(val, "", " ")
}

func (m *jsonMarshaller) Unmarshal(buf []byte, v interface{}) error {
	return json.Unmarshal(buf, v)
}
