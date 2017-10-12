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

type Relation struct {
	key string
	edgesKey string
	edges     []*gographviz.Edge
}

var codeArs = make(map[string]*Entity)
var repos = make(map[string]*Repository)
var providers = make(map[string]*Provider)
var subDomainMap = make(map[string][]string)

func isAggregateRoot(className string) bool {
	tmp := strings.Split(className, "::")
	return tmp[len(tmp)-1] == "AggregateRoot"
}
func (result *CodeDotFileParseResult) parse(edge *gographviz.Edge, nodes map[string]string) {
	tmp := make([]string, 0)
	if strings.Contains(nodes[edge.Dst], "::") {
		tmp = strings.Split(nodes[edge.Dst], "::")
	} else {
		tmp = strings.Split(nodes[edge.Dst], ".")
	}

	dst := strings.Replace(tmp[len(tmp)-1], "\\l", "", -1)

	subDomain := strings.Replace(tmp[0], "\\l", "", -1)
	if _, ok := subDomainMap[subDomain]; ok {
		subDomainMap[subDomain] = append(subDomainMap[subDomain], dst)
	}

	if strings.Contains(nodes[edge.Src], "::") {
		tmp = strings.Split(nodes[edge.Src], "::")
	} else {
		tmp = strings.Split(nodes[edge.Src], ".")
	}
	src := strings.Replace(tmp[len(tmp)-1], "\\l", "", -1)
	if edge.Attrs["style"] == "\"dashed\"" {
		if _, ok := result.edges[dst]; !ok {
			result.edges[dst] = make([]string, 0)
		}
		result.edges[dst] = append(result.edges[dst], src)
	} else {
		if !isAggregateRoot(dst) {
			if strings.HasSuffix(src, "AggregateRoot") {
				codeArs[dst] = &Entity{name: dst}
			}
			if strings.HasSuffix(src, "Entity") {
				result.es[dst] = &Entity{name: dst}
			}
			if strings.HasSuffix(src, "ValueObject") {
				result.vos[dst] = &ValueObject{name: dst}
			}
			if strings.HasSuffix(src, "Repository") {
				repos[dst] = &Repository{name: dst}
			}
			if strings.HasSuffix(src, "Provider") {
				providers[dst] = &Provider{name: dst}
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
				ar.Entities = append(ar.Entities, et)
			}
			if vo, ok := result.vos[edge]; ok {
				ar.VOs = append(ar.VOs, vo)
			}
		}
	}
}

func (result *CodeDotFileParseResult) parseEntity(key string) {
	if entity, ok := result.es[key]; ok {
		for _, edge := range result.edges[key] {
			if et, ok := result.es[edge]; ok {
				entity.Entities = append(entity.Entities, et)
			}
			if vo, ok := result.vos[edge]; ok {
				entity.VOs = append(entity.VOs, vo)
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

			codeDotFiles = append(codeDotFiles, path)
		}

		return nil
	})

	return codeDotFiles
}

func callDotFiles(codeDir string) []string {
	callDotFiles := make([]string, 0)
	filepath.Walk(codeDir, func(path string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_cgraph.dot") {
			callDotFiles = append(callDotFiles, path)
		}

		return nil
	})

	return callDotFiles
}

func nodes(g *gographviz.Graph) map[string]string {
	nodes := make(map[string]string)
	for _, node := range g.Nodes.Nodes {
		nodes[node.Name] = strings.Replace(node.Attrs["label"], "\"", "", 2)
	}
	return nodes
}

func parseDotFile(codeDotfile string) *CodeDotFileParseResult {
	fbuf, _ := ioutil.ReadFile(codeDotfile)
	g, _ := gographviz.Read(fbuf)

	nodes := nodes(g)
	result := &CodeDotFileParseResult{
		edges: make(map[string][]string),
		es:    make(map[string]*Entity),
		vos:   make(map[string]*ValueObject),
	}
	for key := range g.Edges.DstToSrcs {
		for edgesKey := range g.Edges.DstToSrcs[key] {
			for _, edge := range g.Edges.DstToSrcs[key][edgesKey] {
				result.parse(edge, nodes)
			}
		}
	}
	return result
}

func parseCode(codeDotfile string) {
	codeDotFileParseResult := parseDotFile(codeDotfile)
	for key := range codeDotFileParseResult.edges {
		codeDotFileParseResult.parseAggregateRoot(key)
		codeDotFileParseResult.parseEntity(key)
	}
}

func doCallRelation(src string, dst string) {
	for arKey := range codeArs {
		if srcEntity, ok := codeArs[arKey].findEntity(src); ok {
			for arKey2 := range codeArs {
				if arKey != arKey2 {
					if dstEntity, ok := codeArs[arKey2].findEntity(dst); ok {
						srcEntity.callEntities = append(srcEntity.callEntities, dstEntity)
					}
				}
			}
		}
	}
}

func parseCall(codeDotfile string) {

	fbuf, _ := ioutil.ReadFile(codeDotfile)
	g, _ := gographviz.Read(fbuf)

	nodes := nodes(g)
	for key := range nodes {
		fullMethodName := nodes[key]
		if strings.Contains(fullMethodName, "::") {
			tmp := strings.Split(fullMethodName, "::")
			methodName := tmp[len(tmp)-2]
			nodes[key] = strings.Replace(methodName, "\\l", "", -1) //, "\\l", "", -1)
		} else {
			tmp := strings.Split(fullMethodName, ".")
			methodName := tmp[len(tmp)-2]
			nodes[key] = strings.Replace(methodName, "\\l", "", -1) //, "\\l", "", -1)
		}

	}
	for key := range g.Edges.DstToSrcs {
		for edgesKey := range g.Edges.DstToSrcs[key] {
			for _, edge := range g.Edges.DstToSrcs[key][edgesKey] {

				dst := nodes[edge.Dst]
				src := nodes[edge.Src]

				if repo, ok := repos[src]; ok {
					repo.For = codeArs[dst]

				} else {
					doCallRelation(src, dst)
				}
			}
		}
	}

}
func ParseCodeDir(codeDir string, subs []string) *Model {
	codeDotFiles := codeDotFiles(codeDir)
	codeArs = make(map[string]*Entity)
	repos = make(map[string]*Repository)
	providers = make(map[string]*Provider)
	subDomainMap = make(map[string][]string)
	for _, sub := range subs {
		subDomainMap[sub] = make([]string, 0)
	}
	for _, codeDotfile := range codeDotFiles {
		parseCode(codeDotfile)
	}
	callDotFiles := callDotFiles(codeDir)
	for _, callDotFile := range callDotFiles {
		parseCall(callDotFile)
	}
	subDomains := make(map[string]*SubDomain)
	if len(subDomainMap) > 0 {
		for key := range subDomainMap {
			t_ars := make(map[string]*Entity)
			t_repos := make(map[string]*Repository)
			t_providers := make(map[string]*Provider)
			for _, ekey := range subDomainMap[key] {
				if ar, ok := codeArs[ekey]; ok {

					t_ars[ekey] = ar
				}
				if repo, ok := repos[ekey]; ok {
					t_repos[ekey] = repo
				}
				if provider, ok := providers[ekey]; ok {
					t_providers[ekey] = provider
				}
			}
			subDomain := &SubDomain{ARs: t_ars, Repos: t_repos, Providers: t_providers}
			subDomains[key] = subDomain
		}
	} else {

		subDomain := &SubDomain{ARs: codeArs, Repos: repos, Providers: providers}
		subDomains["subdomain"] = subDomain
	}

	return &Model{SubDomains: subDomains}

}
