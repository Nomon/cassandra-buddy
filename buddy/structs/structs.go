package structs

import (
	"errors"

	"golang.org/x/net/context"
)

type RequestContext context.Context

type SnapshotsCreateRequest struct {
	RequestContext `json:"-"`
	Name           string
}

type SnapshotsRestoreRequest struct {
	RequestContext `json:"-"`
	Name           string
	Path           string
	Keyspaces      []string
	Tables         []string
}

type CassandraStartRequest struct {
	RequestContext `json:"-"`
}

type CassandraStopRequest struct {
}

type CassandraStartReply struct {
}

type SnapshotsCreateReply struct {
	Name string
	Path string
	Size int
}

type SnapshotsRestoreReply struct {
	RequestContext `json:"-"`
	ManifestPath   string `json:"manifest_path"`
}

func (s *SnapshotsRestoreRequest) Validate() error {
	if s.Path == "" && s.Name == "" {
		return errors.New("Snapshot requires path or name to be set")
	}
	return nil
}
