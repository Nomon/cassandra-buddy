package datastore

/*
import (
	"io/ioutil"
	"log"
)

type fsStore struct {
	Base string
}

func NewFs(basePath string) Store {
	return &fsStore{
		Base: basePath,
	}
}

func (fs *fsStore) Put(m *Manifest) error {
	for _, dir := range m.Directories {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, file := range files {
			log.Println("UPLOAD", file)
		}
	}
	return nil
}

func (fs *fsStore) Get(m *Manifest) error {
	return nil
}*/

/*
func (fs *fsStore) Put(path string) error {
	return getter.Get("./backups", "file:"+path)
}

func (fs *fsStore) Get(path, dest string) error {
	return getter.Get(dest, "file:"+path)
}*/
