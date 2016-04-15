package datastore

/*
import (
	"io"
	"io/ioutil"
)

type mockDataStore struct {
	files map[string][]byte
}

func NewMockStore() Store {
	return &mockDataStore{
		files: make(map[string][]byte),
	}
}

func (ms *mockDataStore) PutReader(path string, r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	ms.files[path] = data
	return nil
}
func (ms *mockDataStore) PutBytes(path string, b []byte) error {
	ms.files[path] = make([]byte, len(b))
	copy(b, ms.files[path])
	return nil
}
*/
