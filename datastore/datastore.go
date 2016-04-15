package datastore

type Store interface {
	Put(m *Manifest) error
	Get(path string) error
}
