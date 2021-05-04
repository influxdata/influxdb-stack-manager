package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

type object struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   yaml.Node `yaml:"metadata"`
	Spec       yaml.Node `yaml:"spec"`
}

const (
	kindDashboard string = "Dashboard"
	kindLabel     string = "Label"
	kindTask      string = "Task"
)

type queryNode struct {
	Name string
	Node *yaml.Node
}

func walkDashboard(spec *yaml.Node) []queryNode {
	var queryNodes []queryNode

	charts := walkNode(spec, "charts").Content
	for _, c := range charts {
		chartName := walkNode(c, "name").Value
		chartKind := walkNode(c, "kind").Value
		queries := walkNode(c, "queries").Content

		var nodes []*yaml.Node
		for _, q := range queries {
			nodes = append(nodes, walkNode(q, "query"))
		}

		name := fmt.Sprintf("%s_%s", chartName, chartKind)
		for _, node := range nodes {
			queryNodes = append(queryNodes, queryNode{Name: name, Node: node})
		}
	}

	return queryNodes
}

func walkTask(spec *yaml.Node) []queryNode {
	return []queryNode{{
		Name: "query",
		Node: walkNode(spec, "query"),
	}}
}

// walkNodes will walk a node's children, looking for a given key,
// and returning the matching node.
// If the node cannot be found, it will exit.
func walkNode(node *yaml.Node, key string) *yaml.Node {
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}

	log.Fatalf("unable to find key %q in template", key)
	return nil
}
