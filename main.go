package main // import "github.com/Nomon/cassandra-buddy"

import "github.com/Nomon/cassandra-buddy/buddy"

func main() {
	cfg := buddy.NewConfig()
	srv := buddy.NewServer(cfg)
	defer srv.Close()
	srv.Serve()
}
