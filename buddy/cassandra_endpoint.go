package buddy

import "github.com/Nomon/cassandra-buddy/buddy/structs"

type Cassandra struct {
	srv *Server
}

func (c *Cassandra) Start(args *structs.CassandraStartRequest, reply *structs.CassandraStartReply) error {
	return nil
}

func (c *Cassandra) Stop(args *structs.CassandraStartRequest, reply *structs.CassandraStartReply) error {
	return nil
}
