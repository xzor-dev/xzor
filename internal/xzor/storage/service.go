package storage

import "errors"

// Service handles the IO of stored records.
type Service struct {
	EncodeDecoder RecordEncodeDecoder
	Store         RecordStore
}

// Delete removes a record by its ID from the record store.
func (s *Service) Delete(id RecordID) error {
	return s.Store.Delete(id)
}

// Read gets a record's encoded data from the store and decodes it into the provided record.
func (s *Service) Read(id RecordID, record interface{}) error {
	if s.EncodeDecoder == nil {
		return errors.New("no EncodeDecoder provided to the service")
	}

	data, err := s.Store.Read(id)
	if err != nil {
		return err
	}
	return s.EncodeDecoder.DecodeRecord(data, record)
}

// Write encodes a record and writes it to the record store.
func (s *Service) Write(id RecordID, record interface{}) error {
	if s.EncodeDecoder == nil {
		return errors.New("no EncodeDecoder provided to t he service")
	}
	data, err := s.EncodeDecoder.EncodeRecord(record)
	if err != nil {
		return err
	}
	return s.Store.Write(id, data)
}

// NewService creates a new storage service.
func NewService(encodeDecoder RecordEncodeDecoder, store RecordStore) *Service {
	return &Service{
		EncodeDecoder: encodeDecoder,
		Store:         store,
	}
}
