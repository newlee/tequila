package viz

func ParseColl(codeDir string) *FullGraph {
	fullGraph = &FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*Relation),
	}
	codeDotFiles := codeDotFiles(codeDir, "coll__graph.dot")

	for _, codeDotfile := range codeDotFiles {
		parseDotFile(codeDotfile)
	}

	return fullGraph
}
