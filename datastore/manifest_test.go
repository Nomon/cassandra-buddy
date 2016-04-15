package datastore

import (
	"log"
	"testing"
)

func TestManifest(t *testing.T) {
	m, err := NewManifest("/usr/local/var/lib/cassandra/data/", "1460628403086")
	if err != nil {
		t.Fatal(err)

	}
	log.Printf("%#v", m)
}
