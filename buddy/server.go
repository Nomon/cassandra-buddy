package buddy

import (
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Nomon/cassandra-buddy/buddy/cassandra"
	"github.com/Nomon/cassandra-buddy/buddy/nodetool"
	"github.com/Nomon/cassandra-buddy/buddy/structs"
	"github.com/Nomon/cassandra-buddy/datastore"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"gopkg.in/inconshreveable/log15.v2"
)

type Server struct {
	cfg *Config
	log log15.Logger
	// Cassandra process
	cascfg     *cassandra.Config
	cas        cassandra.Process
	casinfo    *nodetool.Info
	cascluster *nodetool.ClusterInfo
	// http
	mux *echo.Echo
	// rpc
	rpcServer *rpc.Server
	endpoints endpoints
	// data
	store datastore.Store
}

type endpoints struct {
	Snapshots *Snapshots
}

func NewServer(cfg *Config) *Server {
	srv := &Server{
		cfg:       cfg,
		rpcServer: rpc.NewServer(),
	}
	if err := srv.setupLogging(); err != nil {
		srv.log.Error("Failed to setup logging", "error", err)
		panic(err)
	}
	if err := srv.setupRPC(); err != nil {
		srv.log.Error("Failed to setup logging", "error", err)
		panic(err)
	}
	if err := srv.setupHTTP(); err != nil {
		srv.log.Error("Failed to setup logging", "error", err)
		panic(err)
	}
	if err := srv.setupSignals(); err != nil {
		srv.log.Error("Failed to setup signal handlers", "error", err)
		panic(err)
	}
	if err := srv.setupCassandra(); err != nil {
		srv.log.Error("Failed to start cassandra", "error", err)
		panic(err)
	}
	if err := srv.setupStore(); err != nil {
		srv.log.Error("Failed to setup store", "error", err)
		panic(err)
	}
	return srv
}

func (srv *Server) Serve() {
	go func() {
		time.Sleep(10 * time.Second)
		err := srv.RPC("Snapshots.Create", &structs.SnapshotsCreateRequest{}, &structs.SnapshotsCreateReply{})
		log.Println(err)
	}()
	srv.mux.Run(standard.New(":3000"))
}

func (srv *Server) Close() error {
	return nil
}

// RPC is used to make a local RPC call
func (srv *Server) RPC(method string, args interface{}, reply interface{}) error {
	logger := srv.logger(nil)
	codec := &inmemCodec{
		method: method,
		args:   args,
		reply:  reply,
	}
	if err := srv.rpcServer.ServeRequest(codec); err != nil {
		logger.Error("RPC Call error", "error", err)
		return err
	}
	if codec.err != nil {
		logger.Warn("RPC Response error", method, codec.reply, "error", codec.err)
	} else {
		logger.Debug("RPC Response", method, codec.reply)
	}

	return codec.err
}

func (srv *Server) setupRPC() error {
	srv.endpoints.Snapshots = &Snapshots{srv}
	srv.rpcServer.Register(srv.endpoints.Snapshots)
	return nil
}

func (srv *Server) setupHTTP() error {
	srv.mux = echo.New()
	srv.mux.Post("/snapshots/create", srv.CreateSnapshot)
	srv.mux.Post("/snapshots/restore", srv.RestoreSnapshot)
	return nil
}

func (srv *Server) setupCassandra() error {
	srv.cascfg = cassandra.DefaultConfig()
	srv.cas = cassandra.New(srv.cascfg)
	err := srv.cas.Start()
	if err != nil {
		return err
	}
	nt := nodetool.New()
	for {
		time.Sleep(1 * time.Second)
		info, err := nt.Info()
		if err != nil || info.ID == "" {
			continue
		}
		srv.casinfo = info
		break
	}
	for {
		time.Sleep(1 * time.Second)
		info, err := nt.ClusterInfo()
		if err != nil || info.Name == "" {
			continue
		}
		srv.cascluster = info
		break
	}
	srv.log.Info("Cassandra started", "info", srv.casinfo, "cluster", srv.cascluster)
	return nil
}

func (srv *Server) setupStore() error {
	//srv.store = datastore.NewFs("/Users/nomon/casb")
	clusterName := strings.Replace(srv.cascluster.Name, " ", "_-_", -1)
	basePath := filepath.Join(srv.cfg.S3Path, clusterName, srv.casinfo.ID)
	srv.store = datastore.NewS3(&datastore.S3Cfg{
		DataPath: srv.cascfg.DataPath,
		BasePath: basePath,
		Region:   srv.cfg.S3Region,
		Bucket:   srv.cfg.S3Bucket,
	})
	return nil
}

func (srv *Server) setupSignals() error {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)
	go func() {
		select {
		case <-ch:
			srv.cas.Stop()
			os.Exit(0)
		}
	}()
	return nil
}

func (srv *Server) setupLogging() error {
	srvLogger := log15.New(log15.Ctx{
		"service": "cassandra-buddy",
	})
	srv.log = srvLogger
	switch srv.cfg.LogLevel {
	case "debug":
		srv.log.SetHandler(log15.LvlFilterHandler(log15.LvlDebug, log15.StdoutHandler))
	case "info":
		srv.log.SetHandler(log15.LvlFilterHandler(log15.LvlInfo, log15.StdoutHandler))
	case "crit":
		srv.log.SetHandler(log15.LvlFilterHandler(log15.LvlCrit, log15.StdoutHandler))
	case "warn":
		srv.log.SetHandler(log15.LvlFilterHandler(log15.LvlWarn, log15.StdoutHandler))
	default:
		srv.log.SetHandler(log15.LvlFilterHandler(log15.LvlInfo, log15.StdoutHandler))
	}
	return nil
}

func (srv *Server) logger(obj interface{}) log15.Logger {
	return srv.log
}
