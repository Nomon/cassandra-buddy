package datastore

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

type s3Store struct {
	bucket      string
	region      string
	base        string
	dataPath    string
	maxParallel int
	auth        *aws.Auth
	s3bucket    *s3.Bucket
}

type S3Cfg struct {
	DataPath    string
	Bucket      string
	BasePath    string
	Region      string
	MaxParallel int
}

func NewS3(cfg *S3Cfg) Store {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}
	amz := s3.New(auth, aws.Regions[cfg.Region])
	bucket := amz.Bucket(cfg.Bucket)
	return &s3Store{
		bucket:      cfg.Bucket,
		region:      cfg.Region,
		base:        cfg.BasePath,
		dataPath:    cfg.DataPath,
		s3bucket:    bucket,
		auth:        &auth,
		maxParallel: 20,
	}
}

func (s *s3Store) Put(m *Manifest) error {
	var size int64
	sem := make(chan bool, s.maxParallel)
	var wg sync.WaitGroup

	for _, dir := range m.Directories {
		path := s.getS3Path(m.Name, dir)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		// relPath is relative to snapshot manifest location that already has the name in path
		relPath, err := filepath.Rel(filepath.Join(s.base, m.Name), path)
		if err != nil {
			return err
		}

		m.Paths = append(m.Paths, relPath)
		for _, file := range files {
			wg.Add(1)
			go func(src, dst string) {
				defer wg.Done() // complete wg
				defer func() {
					<-sem // decrease max parallel semaphore
				}()
				// aquire semaphore
				sem <- true
				upSize, err := s.putFile(dst, src)
				if err != nil {
					log.Println(err)
					return
				}
				size += upSize
				log.Println("File uploaded", src, dst)
			}(filepath.Join(dir, file.Name()), filepath.Join(path, file.Name()))
		}
	}
	md, err := json.Marshal(m)
	if err != nil {
		return err
	}
	wg.Wait()
	p := filepath.Join(s.base, m.Name, "manifest.json")
	log.Println("uploading manifest to", p)
	s.s3bucket.Put(p, md, "application/json", s3.Private)
	log.Println("Snapshot uploaded, size:", size)
	return nil
}

func (s *s3Store) putFile(dst, src string) (int64, error) {
	stat, err := os.Stat(src)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	f, err := os.Open(src)
	if err != nil {
		log.Println(err)
		return stat.Size(), err
	}
	defer f.Close()
	return stat.Size(), s.s3bucket.PutReader(dst, f, stat.Size(), "application/octet-stream", s3.Private)
}

func (s *s3Store) Get(p string) error {
	reader, err := s.getFile(p)
	if err != nil {
		return err
	}
	var m Manifest
	err = json.NewDecoder(reader).Decode(&m)
	if err != nil {
		return err
	}
	m.Path = filepath.Dir(p)
	err = s.downloadManifest(&m)
	return err
}

func (s *s3Store) getFile(path string) (io.ReadCloser, error) {
	return s.s3bucket.GetReader(path)
}

func (s *s3Store) downloadManifest(m *Manifest) error {
	log.Println("downloadManifest", m)
	for _, dir := range m.Paths {
		log.Println("Downloading directory", dir)
	}
	return nil
}

func (s *s3Store) getLocalPath(name, dir string) string {
	rel, err := filepath.Rel(s.dataPath, dir)
	if err != nil {
		panic(err)
	}
	return rel
}

func (s *s3Store) getS3Path(name, dir string) string {
	rel, err := filepath.Rel(s.dataPath, dir)
	if err != nil {
		panic(err)
	}
	rel = filepath.Dir(rel)
	rel = filepath.Dir(rel)
	log.Println(rel, dir)
	return filepath.Join(s.base, name, rel)
}
