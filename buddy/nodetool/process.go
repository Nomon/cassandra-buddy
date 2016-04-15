package nodetool

import (
	"errors"
	"os"
	"os/exec"
)

type Nodetool interface {
	Status() (*Status, error)
	Info() (*Info, error)
	ClusterInfo() (*ClusterInfo, error)
	Snapshot(name string, keyspaces, tables []string) (*Snapshot, error)
	ClearSnapshot(name string, keyspaces, tables []string) error
	Refresh(keyspace, table string) error
}

type nodetool struct {
	Nodetool string
	Host     string
	Port     string
}

// New returns nodetool instance.
func New() Nodetool {
	return &nodetool{}
}

func (n *nodetool) Status() (*Status, error) {
	data, err := n.exec([]string{"status"})
	if err != nil {
		return nil, err
	}
	return NewStatus(data), nil
}

func (n *nodetool) Snapshot(name string, keyspaces, tables []string) (*Snapshot, error) {
	args := []string{"snapshot"}
	if len(tables) > 0 {
		for _, table := range tables {
			args = append(args, "-cf", table)
		}
	}
	if name != "" {
		args = append(args, "-t", name)
	}
	args = append(args, keyspaces...)
	data, err := n.exec(args)
	if err != nil {
		return nil, err
	}
	return NewSnapshot(data), nil
}

func (n *nodetool) ClearSnapshot(name string, keyspaces, tables []string) error {
	args := []string{"clearsnapshot"}
	if len(tables) > 0 {
		for _, table := range tables {
			args = append(args, "-cf", table)
		}
	}
	if name != "" {
		args = append(args, "-t", name)
	}
	args = append(args, keyspaces...)
	_, err := n.exec(args)
	return err
}

func (n *nodetool) Refresh(keyspace, table string) error {
	if keyspace == "" || table == "" {
		return errors.New("Refresh requires keyspace and column family")
	}
	args := []string{"refresh", keyspace, table}
	_, err := n.exec(args)
	return err
}

func (n *nodetool) Info() (*Info, error) {
	data, err := n.exec([]string{"info"})
	if err != nil {
		return nil, err
	}
	return NewInfo(data)
}

func (n *nodetool) ClusterInfo() (*ClusterInfo, error) {
	data, err := n.exec([]string{"describecluster"})
	if err != nil {
		return nil, err
	}
	return NewClusterInfo(data)
}

func (n *nodetool) exec(args []string) ([]byte, error) {
	cmd := exec.Command("/usr/local/bin/nodetool", args...)
	cmd.Env = os.Environ()
	//log.Printf("cmd: %#v", cmd)
	return cmd.CombinedOutput()
}
