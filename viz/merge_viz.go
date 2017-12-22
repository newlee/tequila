package viz

import (
	"bufio"
	"github.com/awalterschulze/gographviz"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func dotFiles(codeDir string) []string {
	codeDotFiles := make([]string, 0)
	filepath.Walk(codeDir, func(path string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".dot") {
			codeDotFiles = append(codeDotFiles, path)
		}
		return nil
	})

	return codeDotFiles
}

func MergeDotFiles(dir string) *FullGraph {
	fullGraph = &FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*Relation),
	}
	codeDotFiles := dotFiles(dir)

	for _, codeDotfile := range codeDotFiles {
		parseDotFile(codeDotfile)
	}

	graph := gographviz.NewGraph()
	graph.SetName("G")

	nodeIndex := 1
	classIndex := 1

	nodes := make(map[string]string)
	classMap := make(map[string][]string)

	for nodeKey := range fullGraph.NodeList {
		tmp := strings.Split(nodeKey, "/")
		packageName := tmp[0]
		if packageName == nodeKey {
			packageName = "main"
		}
		if len(tmp) > 2 {
			packageName = strings.Join(tmp[0:len(tmp)-1], "/")
		}

		if _, ok := classMap[packageName]; !ok {
			classMap[packageName] = make([]string, 0)
		}
		classMap[packageName] = append(classMap[packageName], nodeKey)
	}

	for layer := range classMap {
		layerAttr := make(map[string]string)
		layerAttr["label"] = "\"" + layer + "\""
		layerName := "cluster" + strconv.Itoa(classIndex)
		graph.AddSubGraph("G", layerName, layerAttr)
		classIndex++
		for _, node := range classMap[layer] {
			attrs := make(map[string]string)
			fileName := strings.Replace(node, layer+".", "", -1)
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
			graph.AddEdge(toNode, fromNode, true, attrs)

		}
	}

	f, _ := os.Create("merged.dot")
	w := bufio.NewWriter(f)
	w.WriteString("di" + graph.String())
	w.Flush()
	return fullGraph

}

func MergePackage(dir string) *FullGraph {
	fullGraph = &FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*Relation),
	}
	codeDotFiles := dotFiles(dir)

	for _, codeDotfile := range codeDotFiles {
		parseDotFile(codeDotfile)
	}

	graph := gographviz.NewGraph()
	graph.SetName("G")

	nodes := make(map[string]string)
	//packages := make(map[string]string)
	//packgeIndex := 1
	for nodeKey := range fullGraph.NodeList {
		tmp := strings.Split(nodeKey, ".")
		packageName := tmp[0]
		if packageName == nodeKey {
			packageName = "main"
		}
		if len(tmp) > 2 {
			packageName = strings.Join(tmp[0:len(tmp)-1], ".")
		}

		if !strings.Contains(packageName, ".") {
			continue
		}
		packageName = "\"" + packageName + "\""
		nodes[nodeKey] = packageName
		attrs := make(map[string]string)

		attrs["shape"] = "box"
		graph.AddNode(packageName, packageName, attrs)
	}

	relationMap := make(map[string]string)
	for key := range fullGraph.RelationList {
		relation := fullGraph.RelationList[key]

		if nodes[relation.From] != "" && nodes[relation.To] != "" {
			fromNode := nodes[relation.From]
			toNode := nodes[relation.To]
			if _, ok := relationMap[toNode+fromNode]; !ok {
				if fromNode != toNode {
					relationMap[toNode+fromNode] = ""
					attrs := make(map[string]string)
					if strings.HasSuffix(relation.From, ".m") {
						toName := strings.Replace(relation.From, ".m", "", -1)
						if !strings.Contains(relation.To, toName) {
							relation.Style = "\"dashed\""
						}
					}
					attrs["style"] = relation.Style
					graph.AddEdge(toNode, fromNode, true, attrs)
				}

			}
		}
	}

	f, _ := os.Create("merged_package.dot")
	w := bufio.NewWriter(f)
	w.WriteString("di" + graph.String())
	w.Flush()
	return fullGraph

}
