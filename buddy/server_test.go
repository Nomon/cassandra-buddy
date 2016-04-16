package buddy

import "testing"

func TestServerStart(t *testing.T) {
	cfg := NewConfig()
	srv := NewServer(cfg)
	srv.Close()
}
