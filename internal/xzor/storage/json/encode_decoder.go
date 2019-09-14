package json

import (
	"encoding/json"

	"github.com/xzor-dev/xzor/internal/xzor/storage"
)

var _ storage.RecordEncodeDecoder = &EncodeDecoder{}

// EncodeDecoder provides methods to encode and decode records into and from JSON strings.
type EncodeDecoder struct{}

// DecodeRecord converts a JSON byte slice into the supplied record interface.
func (ed *EncodeDecoder) DecodeRecord(data []byte, record interface{}) error {
	return json.Unmarshal(data, record)
}

// EncodeRecord converts a record interface into a JSON byte slice.
func (ed *EncodeDecoder) EncodeRecord(record interface{}) ([]byte, error) {
	return json.Marshal(record)
}
