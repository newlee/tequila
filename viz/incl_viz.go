package viz

import (
	"bufio"
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func (f *FullGraph) FindCrossRef() []string {
	mergedRelationMap := make(map[string]string)
	result := make([]string,0)
	for key := range f.RelationList {
		relation := f.RelationList[key]
		mergedFrom := strings.Replace(strings.Replace(relation.From, ".h", "", -1), ".cpp", "", -1)
		mergedTo := strings.Replace(strings.Replace(relation.To, ".h", "", -1), ".cpp", "", -1)

		if _, ok := mergedRelationMap[mergedTo+mergedFrom]; ok {
			result = append(result, mergedFrom + " <-> " + mergedTo)
		}
		mergedRelationMap[mergedFrom+mergedTo] = ""
	}
	return result
}

var fullGraph *FullGraph

func parseRelation(edge *gographviz.Edge, nodes map[string]string) {
	if _, ok := nodes[edge.Src]; ok {
		dst := nodes[edge.Dst]
		src := nodes[edge.Src]

		relation := &Relation{
			From:  dst,
			To:    src,
			Style: "\"solid\"",
		}
		fullGraph.RelationList[relation.From+relation.To] = relation
	}
}

func parseDotFile(codeDotfile string) {
	fbuf, _ := ioutil.ReadFile(codeDotfile)
	g, _ := gographviz.Read(fbuf)
	nodes := make(map[string]string)
	for _, node := range g.Nodes.Nodes {

		fullMethodName := strings.Replace(node.Attrs["label"], "\"", "", 2)
		if strings.Contains(fullMethodName, "_test") {
			continue
		}

		if strings.Contains(fullMethodName, "Test") {
			continue
		}

		if strings.Contains(fullMethodName, "/Library/") {
			continue
		}

		methodName := strings.Replace(fullMethodName, "\\l", "", -1)
		methodName = strings.Replace(methodName, "src/", "", -1)
		methodName = strings.Replace(methodName, "include/", "", -1)
		fullGraph.NodeList[methodName] = methodName
		nodes[node.Name] = methodName
	}
	for key := range g.Edges.DstToSrcs {
		for edgesKey := range g.Edges.DstToSrcs[key] {
			for _, edge := range g.Edges.DstToSrcs[key][edgesKey] {
				parseRelation(edge, nodes)
			}
		}
	}
}

func codeDotFiles(codeDir string) []string {
	codeDotFiles := make([]string, 0)
	filepath.Walk(codeDir, func(path string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".dot") {
			if strings.HasSuffix(path, "_dep__incl.dot") {
				//return nil
				if strings.Contains(path, "_test_") {
					return nil
				}

				codeDotFiles = append(codeDotFiles, path)
			}
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
		tmp := strings.Split(nodeKey, "/")
		packageName := tmp[0]
		if packageName == nodeKey {
			packageName = "main"
		}
		if len(tmp) > 2 {
			packageName = strings.Join(tmp[0:len(tmp)-1], "/")
		}

		if _, ok := layerMap[packageName]; !ok {
			layerMap[packageName] = make([]string, 0)
		}
		layerMap[packageName] = append(layerMap[packageName], nodeKey)
	}

	for layer := range layerMap {
		layerAttr := make(map[string]string)
		layerAttr["label"] = "\"" + layer + "\""
		layerName := "cluster" + strconv.Itoa(layerIndex)
		graph.AddSubGraph("G", layerName, layerAttr)
		layerIndex++
		for _, node := range layerMap[layer] {
			attrs := make(map[string]string)
			fileName := strings.Replace(node, layer+"/", "", -1)
			attrs["label"] = "\"" + fileName + "\""
			attrs["shape"] = "box"
			graph.AddNode(layerName, "node"+strconv.Itoa(nodeIndex), attrs)
			nodes[node] = "node" + strconv.Itoa(nodeIndex)
			nodeIndex++
		}
	}

	for key := range fullGraph.RelationList {
		relation := fullGraph.RelationList[key]
		if nodes[relation.From] != "" {
			fromNode := nodes[relation.From]
			toNode := nodes[relation.To]
			attrs := make(map[string]string)
			if strings.HasSuffix(relation.From, ".m") {
				toName := strings.Replace(relation.From, ".m", "", -1)
				if !strings.Contains(relation.To, toName) {
					relation.Style = "\"dashed\""
				}
			}
			attrs["style"] = relation.Style

			graph.AddEdge(fromNode, toNode, true, attrs)

		}
	}

	f, _ := os.Create("dep.dot")
	w := bufio.NewWriter(f)
	w.WriteString("di" + graph.String())
	w.Flush()
	return fullGraph

}
