package datastore

import (
	"io/ioutil"
	"log"
	"os"
	path "path/filepath"
	"strings"
)

// Manifest fully describes a cassandra node and can be used to restore a node.
type Manifest struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	ClusterName string   `json:"cluster_name"`
	Hosts       []string `json:"hosts"`
	Compression string   `json:"compression"`
	Keyspaces   []string `json:"keyspaces"`
	Directories []string `json:"-"`
	Tables      []string `json:"-"`
	Paths       []string `json:"paths"`
}

// NewManifest creates a new manifest
func NewManifest(baseDir, name, path string) (m *Manifest, err error) {
	m = &Manifest{
		Name:        name,
		Path:        path,
		Directories: make([]string, 0),
		Paths:       make([]string, 0),
	}
	if m.Keyspaces, err = readKeyspaces(baseDir); err != nil {
		return
	}
	for _, ks := range m.Keyspaces {
		dirs, err := readSnapshotDirs(baseDir, ks, name)
		if err != nil {
			return nil, err
		}
		m.Directories = append(m.Directories, dirs...)
	}

	return m, nil
}

// readSnapshotDirs will look for any subfolders under cassandra data folder that match the
// backup name.
func readSnapshotDirs(dir, keyspace, backupName string) ([]string, error) {
	dirs := make([]string, 0)

	if isSkippedKeyspace(keyspace) {
		return dirs, nil
	}

	dataDir, err := ioutil.ReadDir(path.Join(dir, keyspace))
	if err != nil {
		return nil, err
	}

	for _, file := range dataDir {
		if file.IsDir() && !isSkippedTable(keyspace, file.Name()) {
			snapPath := path.Join(dir, keyspace, file.Name(), "snapshots", backupName)
			_, err := os.Stat(snapPath)
			if os.IsNotExist(err) {
				continue
			} else if err != nil {
				return nil, err
			}
			dirs = append(dirs, snapPath)
		}
	}
	return dirs, nil
}

func readKeyspaces(dir string) ([]string, error) {
	keyspaces := make([]string, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			keyspaces = append(keyspaces, file.Name())
		}
	}
	return keyspaces, nil
}

func isBackupDir(dir string) bool {
	log.Println("isBackupDir", dir)
	return true
}

var skippedKeyspaces = []string{"OpsCenter"}
var skippedTables = map[string][]string{
	"system": {"local", "peers", "LocationInfo"},
}

func isSkippedTable(ks, table string) bool {
	if skipped, ok := skippedTables[ks]; ok {
		for _, v := range skipped {
			if strings.Contains(table, v) {
				return true
			}
		}
	}
	return false
}

func isSkippedKeyspace(ks string) bool {
	for _, keyspace := range skippedKeyspaces {
		if ks == keyspace {
			return true
		}
	}
	return false
}
