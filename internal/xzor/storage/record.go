package storage

// RecordDecoder decodes a record's encoded data.
type RecordDecoder interface {
	DecodeRecord([]byte, interface{}) error
}

// RecordEncoder converts a record's interface into a byte slice.
type RecordEncoder interface {
	EncodeRecord(interface{}) ([]byte, error)
}

// RecordEncodeDecoder combines RecordDecoder and RecordEncoder.
type RecordEncodeDecoder interface {
	RecordDecoder
	RecordEncoder
}

// RecordID is used to identify individual stored records.
type RecordID string

// RecordStore is used to delete, read, and write record data.
type RecordStore interface {
	Delete(RecordID) error
	Read(RecordID) ([]byte, error)
	Write(RecordID, []byte) error
}
