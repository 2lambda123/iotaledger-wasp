package textdb

import "encoding/json"

type jsonMarshaller struct{}

var _ Marshaller = &jsonMarshaller{}

var j *jsonMarshaller

func JSONMarshaller() Marshaller {
	if j == nil {
		return &jsonMarshaller{}
	}
	return j
}

func (m *jsonMarshaller) Marshal(val interface{}) ([]byte, error) {
	return json.MarshalIndent(val, "", " ")
}

func (m *jsonMarshaller) Unmarshal(buf []byte, v interface{}) error {
	return json.Unmarshal(buf, v)
}
