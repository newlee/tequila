package main

import (
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ParseCode(codeDir string) map[string]*Entity {
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
	// gographviz.Edge.
	codeArs := make(map[string]*Entity)

	// fmt.Println(code)

	for _, codeDotfile := range codeDotFiles {

		fbuf, _ := ioutil.ReadFile(codeDotfile)
		g, _ := gographviz.Read(fbuf)
		nodes := make(map[string]string)
		edges := make(map[string][]string)
		es := make(map[string]*Entity)
		vos := make(map[string]*ValueObject)

		for _, node := range g.Nodes.Nodes {
			nodes[node.Name] = strings.Replace(node.Attrs["label"], "\"", "", 2)
		}
		for key, _ := range g.Edges.DstToSrcs {
			for k, _ := range g.Edges.DstToSrcs[key] {
				for _, edge := range g.Edges.DstToSrcs[key][k] {
					if edge.Attrs["style"] == "\"dashed\"" {
						if _, ok := edges[nodes[edge.Dst]]; !ok {
							edges[nodes[edge.Dst]] = make([]string, 0)
						}
						edges[nodes[edge.Dst]] = append(edges[nodes[edge.Dst]], nodes[edge.Src])
					} else {
						if nodes[edge.Dst] != "AggregateRoot" {
							if nodes[edge.Src] == "AggregateRoot" {
								codeArs[nodes[edge.Dst]] = &Entity{name: nodes[edge.Dst]}
							}
							if nodes[edge.Src] == "Entity" {
								es[nodes[edge.Dst]] = &Entity{name: nodes[edge.Dst]}
							}
							if nodes[edge.Src] == "ValueObject" {
								vos[nodes[edge.Dst]] = &ValueObject{name: nodes[edge.Dst]}
							}
							// fmt.Println(nodes[edge.Src])
							// fmt.Println(nodes[edge.Dst])
						}

					}
				}
			}

		}
		if len(codeArs) > 0 {
			for key, _ := range edges {
				if ar, ok := codeArs[key]; ok {
					for _, edge := range edges[key] {
						if et, ok := es[edge]; ok {
							ar.entities = append(ar.entities, et)
						}
						if vo, ok := vos[edge]; ok {
							ar.vos = append(ar.vos, vo)
						}
					}
				}
				if entity, ok := es[key]; ok {
					for _, edge := range edges[key] {
						if et, ok := es[edge]; ok {
							entity.entities = append(entity.entities, et)
							// fmt.Println("eee")
						}
						if vo, ok := vos[edge]; ok {
							entity.vos = append(entity.vos, vo)
						}
					}
				}
			}
		}

	}

	return codeArs
}
