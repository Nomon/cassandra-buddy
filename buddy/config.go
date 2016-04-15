package buddy

type Config struct {
	LogLevel     string
	CqlshAddr    string
	NodetoolAddr string
	Hostname     string

	// S3 settings
	// bucket to place backups in
	// region where the bucket lives
	// S3Path to store the backups under
	S3Bucket string
	S3Region string
	S3Path   string
}

func NewConfig() *Config {
	return &Config{
		LogLevel:     "debug",
		S3Bucket:     "us-west-staging-media",
		S3Region:     "us-west-1",
		S3Path:       "/cassandra-backups",
		CqlshAddr:    "192.168.33.100",
		NodetoolAddr: "",
	}
}
