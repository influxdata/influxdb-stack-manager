package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	// Filename for the template file in any directory.
	templateFile = "template.yml"

	// prefix added to the filename to indicate a query has
	// been moved to it's own file.
	queryPrefix = "file://"
)

// An object is the basic type for all templates.
// We use yaml.Nodes for the metadata and spec because the exact
// format of these varies depending on the kind.
type object struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   yaml.Node `yaml:"metadata"`
	Spec       yaml.Node `yaml:"spec"`
}

// The different kinds of object that we can receive.
// Only those objects which contain queries are included here.
const (
	kindCheck     string = "CheckThreshold"
	kindDashboard string = "Dashboard"
	kindLabel     string = "Label"
	kindTask      string = "Task"
)

// A queryNode contains the name for the query, generated from the chart/task/check
// it belongs to, and a pointer to the node so it may be updated.
type queryNode struct {
	Name string
	Node *yaml.Node
}

// walkDashboard walks a dashboard spec, and finds all of the query nodes.
// We name each query node after the chart it is in, using the chart name
// and type. If a chart has multiple queries, or two charts have the same
// name and type, the names will not be unique.
func walkDashboard(spec *yaml.Node) []queryNode {
	var queryNodes []queryNode

	// The query nodes can be found at spec.charts[].queries[].query
	// where a [] indicates there is an array of charts/queries.
	charts := walkNode(spec, "charts").Content
	for _, c := range charts {
		queries := walkNode(c, "queries").Content
		var nodes []*yaml.Node
		for _, q := range queries {
			nodes = append(nodes, walkNode(q, "query"))
		}

		chartName := walkNode(c, "name").Value
		chartKind := walkNode(c, "kind").Value
		name := fmt.Sprintf("%s_%s", chartName, chartKind)
		for _, node := range nodes {
			queryNodes = append(queryNodes, queryNode{Name: name, Node: node})
		}
	}

	return queryNodes
}

// walkTask walks a task spec, and finds the query node.
// This returns a list to match the other walk function, but will only
// ever contain one member.
func walkTask(spec *yaml.Node) []queryNode {
	return []queryNode{{
		Name: "query",
		Node: walkNode(spec, "query"),
	}}
}

// walkCheck walks a check spec and finds the query node.
// This returns a list to match the other walk function, but will only
// ever contain one member.
func walkCheck(spec *yaml.Node) []queryNode {
	return []queryNode{{
		Name: "query",
		Node: walkNode(spec, "query"),
	}}
}

// walkNode will walk a node's children, looking for the value
// node that matches the key.
func walkNode(node *yaml.Node, key string) *yaml.Node {
	// Within a node's content, nodes are grouped in pairs.
	// The first node in a pair is a scalar string node, with the key as value.
	// The second node in the pair is the value.
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}

	return &yaml.Node{}
}
