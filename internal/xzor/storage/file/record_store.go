package file

import (
	"io/ioutil"
	"os"

	"github.com/xzor-dev/xzor/internal/xzor/storage"
)

var _ storage.RecordStore = &RecordStore{}

// RecordStore stores and retrieves record data from the file system.
type RecordStore struct {
	RootDir string
}

func (s *RecordStore) filename(id storage.RecordID) string {
	return s.RootDir + "/" + string(id)
}

// Delete removes a record's data file.
func (s *RecordStore) Delete(id storage.RecordID) error {
	return os.Remove(s.filename(id))
}

// Read attempts to get a record's data from a file.
func (s *RecordStore) Read(id storage.RecordID) ([]byte, error) {
	return ioutil.ReadFile(s.filename(id))
}

// Write creates or replaces a file with the supplied data.
func (s *RecordStore) Write(id storage.RecordID, data []byte) error {
	err := os.MkdirAll(s.RootDir, 0666)
	if err != nil {
		return err
	}
	f, err := os.Create(s.filename(id))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
