package viz

import (
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	//"fmt"
	"bufio"
	"fmt"
	"strconv"
)

type Relation struct {
	From  string
	To    string
	Style string
}

type FullGraph struct {
	NodeList     map[string]string
	RelationList map[string]*Relation
}

var fullGraph *FullGraph

func parseRelation(edge *gographviz.Edge, nodes map[string]string) *Relation {
	dst := nodes[edge.Dst]
	src := nodes[edge.Src]
	return &Relation{
		From:  dst,
		To:    src,
		Style: edge.Attrs["style"],
	}
}

func parseDotFile(codeDotfile string) {
	fbuf, _ := ioutil.ReadFile(codeDotfile)
	g, _ := gographviz.Read(fbuf)
	nodes := make(map[string]string)
	for _, node := range g.Nodes.Nodes {
		fullMethodName := strings.Replace(node.Attrs["label"], "\"", "", 2)
		if strings.Contains(fullMethodName, "::") {
			methodName := strings.Replace(fullMethodName, "\\l", "", -1)
			fullGraph.NodeList[methodName] = methodName
			nodes[node.Name] = methodName
		}
	}
	for key := range g.Edges.DstToSrcs {
		for edgesKey := range g.Edges.DstToSrcs[key] {
			for _, edge := range g.Edges.DstToSrcs[key][edgesKey] {
				relation := parseRelation(edge, nodes)
				fullGraph.RelationList[relation.From+relation.To] = relation
			}
		}
	}
}

func codeDotFiles(codeDir string) []string {
	codeDotFiles := make([]string, 0)
	filepath.Walk(codeDir, func(path string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".dot") {
			if strings.HasSuffix(path, "class_domain_1_1_aggregate_root__coll__graph.dot") {
				return nil
			}
			if strings.Contains(path, "inherit") {
				return nil
			}
			if strings.HasSuffix(path, "_cgraph.dot") {
				return nil
			}

			if strings.HasSuffix(path, "_incl.dot") {
				return nil
			}

			codeDotFiles = append(codeDotFiles, path)
		}

		return nil
	})

	return codeDotFiles
}

func ParseCodeDir(codeDir string) *FullGraph {
	fullGraph = &FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*Relation),
	}
	codeDotFiles := codeDotFiles(codeDir)

	for _, codeDotfile := range codeDotFiles {
		parseDotFile(codeDotfile)
	}

	graph := gographviz.NewGraph()
	graph.SetName("G")

	nodeIndex := 1
	layerIndex := 1
	nodes := make(map[string]string)

	layerMap := make(map[string][]string)

	for nodeKey := range fullGraph.NodeList {
		tmp := strings.Split(nodeKey, "::")

		if _, ok := layerMap[tmp[0]]; !ok {
			layerMap[tmp[0]] = make([]string, 0)
		}
		layerMap[tmp[0]] = append(layerMap[tmp[0]], nodeKey)
	}

	for layer := range layerMap {
		layerAttr := make(map[string]string)
		layerAttr["label"] = layer
		layerName := "cluster" + strconv.Itoa(layerIndex)
		graph.AddSubGraph("G", layerName, layerAttr)
		layerIndex++
		for _, node := range layerMap[layer] {
			attrs := make(map[string]string)
			attrs["label"] = "\"" + node + "\""
			attrs["shape"] = "box"
			graph.AddNode(layerName, "node"+strconv.Itoa(nodeIndex), attrs)
			nodes[node] = "node" + strconv.Itoa(nodeIndex)
			nodeIndex++
		}
	}

	for key := range fullGraph.RelationList {
		relation := fullGraph.RelationList[key]
		fmt.Println(relation)

		if nodes[relation.From] != "" {
			fromNode := nodes[relation.From]
			toNode := nodes[relation.To]

			attrs := make(map[string]string)

			attrs["style"] = relation.Style
			graph.AddEdge(fromNode, toNode, true, attrs)
		}
	}

	f, _ := os.Create("../examples/bc-code/dep.dot")
	w := bufio.NewWriter(f)
	w.WriteString(graph.String())
	w.Flush()
	return fullGraph

}
