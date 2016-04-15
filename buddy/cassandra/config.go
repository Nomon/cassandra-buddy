package cassandra

import "fmt"

type Config struct {
	Executable   string
	JoinRing     bool
	NewHeapSize  string
	HeapSize     string
	DataPath     string
	CommitPath   string
	BackupPath   string
	CachePath    string
	JmxPort      int
	MaxDirectMem string
}

func (c *Config) Env() []string {
	e := make([]string, 0)
	e = append(e, fmt.Sprintf("HEAP_NEWSIZE=%s", c.NewHeapSize))
	e = append(e, fmt.Sprintf("MAX_HEAP_SIZE=%s", c.HeapSize))
	e = append(e, fmt.Sprintf("DATA_DIR=%s", c.DataPath))
	e = append(e, fmt.Sprintf("COMMIT_LOG_DIR=%s", c.CommitPath))
	e = append(e, fmt.Sprintf("LOCAL_BACKUP_DIR=%s", c.BackupPath))
	e = append(e, fmt.Sprintf("CACHE_DIR=%s", c.CachePath))
	e = append(e, fmt.Sprintf("JMX_PORT=%d", c.JmxPort))
	e = append(e, fmt.Sprintf("MAX_DIRECT_MEMORY=%s", c.MaxDirectMem))
	jr := "false"
	if c.JoinRing {
		jr = "true"
	}
	e = append(e, fmt.Sprintf("cassandra.join_ring=%s", jr))
	return e
}

func DefaultConfig() *Config {
	return &Config{
		Executable:   "/usr/local/bin/cassandra",
		JoinRing:     false,
		HeapSize:     "2G",
		NewHeapSize:  "200M",
		DataPath:     "/usr/local/var/lib/cassandra/data",
		CommitPath:   "/usr/local/var/lib/cassandra/commitlog",
		BackupPath:   "/usr/local/var/lib/cassandra/backups",
		CachePath:    "/usr/local/var/lib/cassandra/cache",
		JmxPort:      7199,
		MaxDirectMem: "1G",
	}
}
