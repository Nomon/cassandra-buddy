package nodetool

import (
	"regexp"
	"strings"
)

type ClusterInfo struct {
	Name           string
	Snitch         string
	Partitioner    string
	SchemaVersions map[string][]string
}

var nameRegex = regexp.MustCompile(`^\s*Name\s*:\s*(.+)$`)
var snitchRegex = regexp.MustCompile(`^\s*Snitch\s*:\s*(.+)$`)
var partitionerRegex = regexp.MustCompile(`^\s*Partitioner\s*:\s*(.+)$`)
var schemaRegex = regexp.MustCompile(`^\s+(.+-.+-.+-.+-.+):\s+\[(.+)\]$`)

func NewClusterInfo(d []byte) (*ClusterInfo, error) {
	c := &ClusterInfo{}
	for _, line := range strings.Split(string(d), "\n") {
		if parts := nameRegex.FindAllStringSubmatch(line, 2); parts != nil {
			c.Name = parts[0][1]
		} else if parts := snitchRegex.FindAllStringSubmatch(line, 2); parts != nil {
			c.Snitch = parts[0][1]
		} else if parts := partitionerRegex.FindAllStringSubmatch(line, 2); parts != nil {
			c.Partitioner = parts[0][1]
		} else if parts := schemaRegex.FindAllStringSubmatch(line, 2); parts != nil {
			if c.SchemaVersions == nil {
				c.SchemaVersions = make(map[string][]string)
			}
			c.SchemaVersions[parts[0][1]] = []string{}
			for _, ip := range strings.Split(parts[0][2], ",") {
				c.SchemaVersions[parts[0][1]] = append(c.SchemaVersions[parts[0][1]], strings.Trim(ip, " "))
			}
		}
	}
	return c, nil
}
