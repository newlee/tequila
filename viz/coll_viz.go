package viz

import (
	"strings"
)

func ParseColl(codeDir string, filter string) *FullGraph {
	fullGraph = &FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*Relation),
	}
	codeDotFiles := codeDotFiles(codeDir, func(path string) bool {
		return strings.HasSuffix(path, filter)
	})

	for _, codeDotfile := range codeDotFiles {
		parseDotFile(codeDotfile)
	}

	return fullGraph
}

func ParseICallGraph(codeDir string, filter string) *FullGraph {
	fullGraph = &FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*Relation),
	}
	codeDotFiles := codeDotFiles(codeDir, func(path string) bool {
		return strings.HasSuffix(path, "_icgraph.dot") && strings.Contains(path, filter)
	})

	for _, codeDotfile := range codeDotFiles {
		parseDotFile(codeDotfile)
	}

	return fullGraph
}

func ParseICallGraphStart() {
	fullGraph = &FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*Relation),
	}
}

func ParseICallGraphByBuffer(buf []byte) {
	parseFromBuffer(buf)
}

func ParseICallGraphEnd() *FullGraph {
	return fullGraph
}
