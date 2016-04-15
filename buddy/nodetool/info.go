package nodetool

import (
	"regexp"
	"strconv"
	"strings"
)

type Info struct {
	ID                    string
	GossipActive          bool
	ThriftActive          bool
	NativeTransportActive bool
	Load                  string
	GenerationNo          int64
	Uptime                string
	HeapUsed              float64
	HeapMax               float64
	HeapUsage             float64
	DataCenter            string
	Rack                  string
	Exceptions            int64
}

var idRegex = regexp.MustCompile(`^\s*ID\s*: ([a-z0-9\-]+)$`)
var gossipRegex = regexp.MustCompile(`^\s*Gossip active\s*: (true|false)$`)
var thriftRegex = regexp.MustCompile(`^\s*Thrift active\s*: (true|false)$`)
var nativeTransportRegex = regexp.MustCompile(`^\s*Native Transport active\s*: (true|false)$`)
var loadRegex = regexp.MustCompile(`^\s*Load\s*: ([0-9\.]+ [KMGP]?B)$`)
var generationRegex = regexp.MustCompile(`^\s*Generation No\s*: ([0-9]+)$`)
var uptimeRegex = regexp.MustCompile(`^\s*Uptime \(seconds\)\s*: ([0-9]+)$`)
var heapRegex = regexp.MustCompile(`^\s*Heap Memory \(MB\)\s*: ([0-9\.]+) / ([0-9\.]+)$`)
var datacenterRegex = regexp.MustCompile(`^\s*Data Center\s*: (.+)$`)
var rackRegex = regexp.MustCompile(`^\s*Rack\s*: (.+)$`)
var exceptionsRegex = regexp.MustCompile(`^\s*Exceptions\s*: (.+)$`)

func NewInfo(d []byte) (*Info, error) {
	info := &Info{}
	for _, line := range strings.Split(string(d), "\n") {
		if parts := idRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.ID = parts[0][1]
		} else if parts := gossipRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.GossipActive, _ = strconv.ParseBool(parts[0][1])
		} else if parts := thriftRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.ThriftActive, _ = strconv.ParseBool(parts[0][1])
		} else if parts := nativeTransportRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.NativeTransportActive, _ = strconv.ParseBool(parts[0][1])
		} else if parts := loadRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.Load = parts[0][1]
		} else if parts := generationRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.GenerationNo, _ = strconv.ParseInt(parts[0][1], 10, 64)
		} else if parts := heapRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.HeapUsed, _ = strconv.ParseFloat(parts[0][1], 64)
			info.HeapMax, _ = strconv.ParseFloat(parts[0][2], 64)
			info.HeapUsage = (info.HeapUsed / info.HeapMax) * 100
		} else if parts := datacenterRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.DataCenter = parts[0][1]
		} else if parts := rackRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.Rack = parts[0][1]
		} else if parts := exceptionsRegex.FindAllStringSubmatch(line, 2); parts != nil {
			info.Exceptions, _ = strconv.ParseInt(parts[0][1], 10, 64)
		}
	}
	return info, nil
}
