package cqlsh

import (
	"log"
	"os"
	"os/exec"
)

type Cql struct {
}

func NeqCql() *Cql {
	return &Cql{}
}

func (cql *Cql) DescribeKeyspace(keyspace string) (string, error) {
	args := []string{"-e", "DESCRIBE KEYSPACE " + keyspace}
	d, err := cql.exec(args)
	if err != nil {
		return "", err
	}
	return string(d), err
}

func (cql *Cql) exec(args []string) ([]byte, error) {
	cmd := exec.Command("/usr/local/bin/cqlsh", args...)
	cmd.Env = os.Environ()
	log.Printf("cmd: %#v", cmd)
	return cmd.CombinedOutput()
}
