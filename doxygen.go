package main

import (
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CodeDotFileParseResult struct {
	edges map[string][]string
	es    map[string]*Entity
	vos   map[string]*ValueObject
}

func (result *CodeDotFileParseResult) parse(edge *gographviz.Edge, nodes map[string]string) {
	if edge.Attrs["style"] == "\"dashed\"" {
		if _, ok := result.edges[nodes[edge.Dst]]; !ok {
			result.edges[nodes[edge.Dst]] = make([]string, 0)
		}
		result.edges[nodes[edge.Dst]] = append(result.edges[nodes[edge.Dst]], nodes[edge.Src])
	} else {
		if nodes[edge.Dst] != "AggregateRoot" {
			if nodes[edge.Src] == "AggregateRoot" {
				codeArs[nodes[edge.Dst]] = &Entity{name: nodes[edge.Dst]}
			}
			if nodes[edge.Src] == "Entity" {
				result.es[nodes[edge.Dst]] = &Entity{name: nodes[edge.Dst]}
			}
			if nodes[edge.Src] == "ValueObject" {
				result.vos[nodes[edge.Dst]] = &ValueObject{name: nodes[edge.Dst]}
			}
		}

	}
}

func (result *CodeDotFileParseResult) parseAggregateRoot(key string) {
	if ar, ok := codeArs[key]; ok {
		for _, edge := range result.edges[key] {
			if ref, ok := codeArs[edge]; ok {
				ar.Refs = append(ar.Refs, ref)
			}
			if et, ok := result.es[edge]; ok {
				ar.entities = append(ar.entities, et)
			}
			if vo, ok := result.vos[edge]; ok {
				ar.vos = append(ar.vos, vo)
			}
		}
	}
}

func (result *CodeDotFileParseResult) parseEntity(key string) {
	if entity, ok := result.es[key]; ok {
		for _, edge := range result.edges[key] {
			if et, ok := result.es[edge]; ok {
				entity.entities = append(entity.entities, et)
			}
			if vo, ok := result.vos[edge]; ok {
				entity.vos = append(entity.vos, vo)
			}
		}
	}
}

var codeArs = make(map[string]*Entity)

func codeDotFiles(codeDir string) []string {
	codeDotFiles := make([]string, 0)
	filepath.Walk(codeDir, func(path string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".dot") {
			if strings.HasSuffix(path, "class_aggregate_root__coll__graph.dot") {
				return nil
			}
			if strings.Contains(path, "inherit") {
				return nil
			}

			codeDotFiles = append(codeDotFiles, path)
		}

		return nil
	})
	return codeDotFiles
}

func nodes(g *gographviz.Graph) map[string]string {
	nodes := make(map[string]string)
	for _, node := range g.Nodes.Nodes {
		nodes[node.Name] = strings.Replace(node.Attrs["label"], "\"", "", 2)
	}
	return nodes
}

func parseDotFile(g *gographviz.Graph) *CodeDotFileParseResult {
	nodes := nodes(g)
	result := &CodeDotFileParseResult{
		edges: make(map[string][]string),
		es:    make(map[string]*Entity),
		vos:   make(map[string]*ValueObject),
	}

	for key, _ := range g.Edges.DstToSrcs {
		for edgesKey, _ := range g.Edges.DstToSrcs[key] {
			for _, edge := range g.Edges.DstToSrcs[key][edgesKey] {
				result.parse(edge, nodes)
			}
		}
	}
	return result
}
func parseCode(codeDotfile string) {
	fbuf, _ := ioutil.ReadFile(codeDotfile)

	g, _ := gographviz.Read(fbuf)

	codeDotFileParseResult := parseDotFile(g)

	if len(codeArs) > 0 {
		for key, _ := range codeDotFileParseResult.edges {
			codeDotFileParseResult.parseAggregateRoot(key)
			codeDotFileParseResult.parseEntity(key)
		}
	}
}
func ParseCodeDir(codeDir string) map[string]*Entity {
	codeDotFiles := codeDotFiles(codeDir)

	for _, codeDotfile := range codeDotFiles {
		parseCode(codeDotfile)
	}

	return codeArs
}
