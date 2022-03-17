package textdb

import "gopkg.in/yaml.v2"

type yamlMarshaller struct{}

var _ Marshaller = &yamlMarshaller{}

func NewYAMLMarshaller() Marshaller {
	return &yamlMarshaller{}
}

func (y *yamlMarshaller) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y *yamlMarshaller) Unmarshal(in []byte, v interface{}) error {
	return yaml.Unmarshal(in, v)
}
