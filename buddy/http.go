package buddy

import (
	"github.com/Nomon/cassandra-buddy/buddy/structs"
	"github.com/labstack/echo"
)

func (srv *Server) CreateSnapshot(c echo.Context) error {
	var args structs.SnapshotsCreateRequest
	var reply structs.SnapshotsCreateReply
	if err := c.Bind(&args); err != nil {
		return err
	}
	if err := srv.RPC("Snapshots.Create", &args, &reply); err != nil {
		return err
	}
	return c.JSON(200, reply)
}

func (srv *Server) RestoreSnapshot(c echo.Context) error {
	var args structs.SnapshotsRestoreRequest
	var reply structs.SnapshotsRestoreReply
	if err := c.Bind(&args); err != nil {
		return err
	}
	if err := srv.RPC("Snapshots.Restore", &args, &reply); err != nil {
		return err
	}
	return c.JSON(200, reply)
}
