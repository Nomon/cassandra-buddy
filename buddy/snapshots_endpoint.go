package buddy

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nomon/cassandra-buddy/buddy/nodetool"
	"github.com/Nomon/cassandra-buddy/buddy/structs"
	"github.com/Nomon/cassandra-buddy/datastore"
)

type Snapshots struct {
	srv *Server
}

// Create is the RPC endpoint for creating a snapshot
func (s *Snapshots) Create(args *structs.SnapshotsCreateRequest, reply *structs.SnapshotsCreateReply) error {
	logger := s.srv.logger(args)
	if args.Name == "" {
		args.Name = createManifestName()
	}
	log.Println("Creating snapshot", "path", s.srv.cascfg.BackupPath+"/"+args.Name)
	nt := nodetool.New()
	snapshot, err := nt.Snapshot(args.Name, nil, nil)
	log.Println(snapshot, err)
	if err != nil {
		logger.Error("Nodetool error", "error", err)
		return err
	}

	// s3 path is /configured_path_prefix/cluster_name/host_id/backup_name
	path := filepath.Join(s.srv.cfg.S3Path, strings.Replace(s.srv.cascluster.Name, " ", "_-_", -1), s.srv.casinfo.ID, args.Name)

	manifest, err := datastore.NewManifest(s.srv.cascfg.DataPath, args.Name, path)
	if err != nil {
		return err
	}

	logger.Info("Putting manifest into store", "manifest", manifest)
	if err = s.srv.store.Put(manifest); err != nil {
		return err
	}

	return nil
}

func (s *Snapshots) Restore(args *structs.SnapshotsRestoreRequest, reply *structs.SnapshotsRestoreReply) error {
	logger := s.srv.logger(args)

	if err := args.Validate(); err != nil {
		logger.Error("Snapshots.Restore Validation failed", "error", err)
		return err
	}
	if err := s.srv.cas.Stop(); err != nil {
		logger.Error("Failed to stop cassandra", "error", err)
		return err
	}
	if err := s.srv.cas.ClearData(args.Keyspaces); err != nil {
		logger.Error("Failed to clear cassandra data")
	}
	if err := s.srv.store.Get(args.Path); err != nil {
		logger.Error("Failed to download backups", "error", err)
		return err
	}
	if err := s.srv.setupCassandra(); err != nil {
		return err
	}
	return nil
}

func createManifestName() string {
	now := time.Now()
	y, m, d := now.Date()
	h, mi, s := now.Clock()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d", y, m, d, h, mi, s)
}
