package textdb

import "gopkg.in/yaml.v2"

type yamlMarshaller struct{}

var _ Marshaller = &yamlMarshaller{}

var y *yamlMarshaller

func YAMLMarshaller() Marshaller {
	if y == nil {
		return &yamlMarshaller{}
	}
	return y
}

func (y *yamlMarshaller) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y *yamlMarshaller) Unmarshal(in []byte, v interface{}) error {
	return yaml.Unmarshal(in, v)
}
