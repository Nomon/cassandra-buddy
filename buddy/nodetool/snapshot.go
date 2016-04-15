package nodetool

type Snapshot struct {
	Name string
	Path string
}

func NewSnapshot(d []byte) *Snapshot {
	return &Snapshot{}
}
