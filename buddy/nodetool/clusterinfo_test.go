package nodetool

import "testing"

var testInfo = []byte(`Cluster Information:
	Name: challenge-cassandra
	Snitch: org.apache.cassandra.locator.DynamicEndpointSnitch
	Partitioner: org.apache.cassandra.dht.Murmur3Partitioner
	Schema versions:
		e094133b-e2f4-32f0-a359-5557f36a7259: [172.18.35.82, 172.18.36.63, 172.18.37.220]
`)

func TestNewClusterInfo(t *testing.T) {
	info, err := NewClusterInfo(testInfo)
	if err != nil {
		t.Fatal(err)
	}
	if info.Name != "challenge-cassandra" {
		t.Fatal("failed to parse name")
	}
	if info.Partitioner != "org.apache.cassandra.dht.Murmur3Partitioner" {
		t.Fatal("Failed to parse partitioner")
	}
	if info.Snitch != "org.apache.cassandra.locator.DynamicEndpointSnitch" {
		t.Fatal("Fauiled to parse snitch")
	}
	if info.SchemaVersions == nil {
		t.Fatal("Failed to parse schema versions")
	}
	if len(info.SchemaVersions["e094133b-e2f4-32f0-a359-5557f36a7259"]) != 3 {
		t.Fatal("Failed to parse schema versions")
	}
}
