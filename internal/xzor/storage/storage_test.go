package storage_test

import (
	"errors"
	"os"
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/storage"
	"github.com/xzor-dev/xzor/internal/xzor/storage/file"
	"github.com/xzor-dev/xzor/internal/xzor/storage/json"
)

func TestStorageService(t *testing.T) {
	s := &storage.Service{
		EncodeDecoder: &testEncodeDecoder{},
		Store:         &testRecordStore{},
	}
	id := storage.RecordID("test-id")
	valueA := "test value"
	err := s.Write(id, valueA)
	if err != nil {
		t.Fatalf("%v", err)
	}

	var valueB string
	err = s.Read(id, &valueB)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if valueA != valueB {
		t.Fatalf("unexpected data read from storage: wanted '%s', got '%s'", valueA, valueB)
	}
	err = s.Delete(id)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestJSONFileStorage(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}
	s := &storage.Service{
		EncodeDecoder: &json.EncodeDecoder{},
		Store: &file.RecordStore{
			RootDir: dir + "/testdata",
		},
	}

	type testRecord struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}
	recordID := storage.RecordID("foo-bar")
	recordA := &testRecord{
		Foo: "hello",
		Bar: 42,
	}

	err = s.Write(recordID, recordA)
	if err != nil {
		t.Fatalf("%v", err)
	}

	recordB := &testRecord{}
	err = s.Read(recordID, recordB)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if recordA.Foo != recordB.Foo || recordA.Bar != recordB.Bar {
		t.Fatalf("records do not match: wanted %v, got %v", recordA, recordB)
	}

	err = s.Delete(recordID)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

var _ storage.RecordEncodeDecoder = &testEncodeDecoder{}

type testEncodeDecoder struct{}

func (ed *testEncodeDecoder) DecodeRecord(data []byte, record interface{}) error {
	p, ok := record.(*string)
	if !ok {
		return errors.New("could not convert record to a string pointer")
	}
	*p = string(data)
	return nil
}

func (ed *testEncodeDecoder) EncodeRecord(record interface{}) ([]byte, error) {
	str, ok := record.(string)
	if !ok {
		return nil, errors.New("could not convert record to a string")
	}
	return []byte(str), nil
}

var _ storage.RecordStore = &testRecordStore{}

type testRecordStore struct {
	records map[storage.RecordID][]byte
}

func (s *testRecordStore) Delete(id storage.RecordID) error {
	if s.records != nil && s.records[id] != nil {
		delete(s.records, id)
	}
	return nil
}

func (s *testRecordStore) Read(id storage.RecordID) ([]byte, error) {
	if s.records == nil || s.records[id] == nil {
		return nil, errors.New("unknown record ID")
	}
	return s.records[id], nil
}

func (s *testRecordStore) Write(id storage.RecordID, data []byte) error {
	if s.records == nil {
		s.records = make(map[storage.RecordID][]byte)
	}
	s.records[id] = data
	return nil
}
