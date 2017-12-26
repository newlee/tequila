package viz

func ParseColl(codeDir string, filter string) *FullGraph {
	fullGraph = &FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*Relation),
	}
	codeDotFiles := codeDotFiles(codeDir, filter)

	for _, codeDotfile := range codeDotFiles {
		parseDotFile(codeDotfile)
	}

	return fullGraph
}
