package main

import (
	"github.com/awalterschulze/gographviz"
	. "github.com/newlee/tequila/dot"
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

var codeArs = make(map[string]*Entity)
var repos = make(map[string]*Repository)
var providers = make(map[string]*Provider)
var subDomainMap = make(map[string][]string)

func isAggregateRoot(className string) bool {
	return className == "AggregateRoot"
}

func codeDotFiles(codeDir string) []string {
	codeDotFiles := make([]string, 0)
	filepath.Walk(codeDir, func(path string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".dot") {
			if strings.HasSuffix(path, "class_domain_1_1_aggregate_root__coll__graph.dot") {
				return nil
			}
			if strings.Contains(path, "inherit_") {
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


func isAR(node *Node) bool {
	return node.IsIt("AggregateRoot")
}
func isEntity(node *Node) bool {
	return node.IsIt("Entity")
}

func isValueObject(node *Node) bool {
	return node.IsIt("ValueObject")
}
func isRepo(node *Node) bool {
	return node.IsIt("Repository")
}

func (result *CodeDotFileParseResult) parseRelation(node *Node) {
	src := node.Name
	if !isAggregateRoot(src) {
		if isAR(node) {
			codeArs[src] = &Entity{name: src}
			for _, relation := range node.DstNodes {
				result.parseRelation(relation.Node)
			}
			return
		}
		if isEntity(node) {
			result.es[src] = &Entity{name: src}
		}
		if isValueObject(node) {
			result.vos[src] = &ValueObject{name: src}
		}
		if isRepo(node)  {
			repos[src] = &Repository{name: src}
		}
		if strings.HasSuffix(src, "Provider") {
			providers[src] = &Provider{name: src}
		}
	}

	for _, relation := range node.DstNodes {
		result.parseRelation(relation.Node)
	}
}
func (result *CodeDotFileParseResult) parseAR(node *Node) {
	if ar, ok := codeArs[node.Name]; ok {
		for _, relation := range node.DstNodes {
			dst := relation.Node.Name

			if ref, ok := codeArs[dst]; ok {
				ar.Refs = append(ar.Refs, ref)
			}
			if et, ok := result.es[dst]; ok {
				ar.Entities = append(ar.Entities, et)
			}
			if vo, ok := result.vos[dst]; ok {

				ar.VOs = append(ar.VOs, vo)
			}
		}
	}
	
	for _, relation := range node.DstNodes {
		result.parseAR(relation.Node)
	}
}

func (result *CodeDotFileParseResult) parseEntity(node *Node) {
	if entity, ok := result.es[node.Name]; ok {
		for _, relation := range node.DstNodes {
			dst := relation.Node.Name
			if et, ok := result.es[dst]; ok {
				entity.Entities = append(entity.Entities, et)
			}
			if vo, ok := result.vos[dst]; ok {
				entity.appendVO(vo)
			}
		}
	}
	for _, relation := range node.DstNodes {
		result.parseEntity(relation.Node)
	}
}
func subDomainCallback(subDomain, methodName string) {
	if _, ok := subDomainMap[subDomain]; ok {
		subDomainMap[subDomain] = append(subDomainMap[subDomain], methodName)
	}
}
func parseCode(codeDotfile string) {
	node := ParseDoxygenFile(codeDotfile)
	node.RemoveNS(subDomainCallback)
	codeDotFileParseResult := &CodeDotFileParseResult{
		edges: make(map[string][]string),
		es:    make(map[string]*Entity),
		vos:   make(map[string]*ValueObject),
	}

	codeDotFileParseResult.parseRelation(node)

	codeDotFileParseResult.parseAR(node)
	codeDotFileParseResult.parseEntity(node)
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

func getMethodName(fullMethodName, split string, index int) (string, string, bool) {
	if strings.Contains(fullMethodName, split) {
		tmp := strings.Split(fullMethodName, split)
		methodName := tmp[len(tmp)-index]
		methodName = strings.Replace(methodName, "\\l", "", -1)
		subDomain := strings.Replace(tmp[0], "\\l", "", -1)
		if _, ok := subDomainMap[subDomain]; ok {
			subDomainMap[subDomain] = append(subDomainMap[subDomain], methodName)
		}
		return methodName, subDomain, true
	}
	return fullMethodName, "", false
}

func nodes(g *gographviz.Graph, index int) map[string]string {
	nodes := make(map[string]string)
	for _, node := range g.Nodes.Nodes {
		fullMethodName := strings.Replace(node.Attrs["label"], "\"", "", 2)

		if methodName, _, ok := getMethodName(fullMethodName, "::", index); ok {
			nodes[node.Name] = methodName
		} else {
			nodes[node.Name], _, _ = getMethodName(fullMethodName, ".", index)
		}
	}
	return nodes
}

func parseCall(codeDotfile string) {
	fbuf, _ := ioutil.ReadFile(codeDotfile)
	g, _ := gographviz.Read(fbuf)

	nodes := nodes(g, 2)

	for key := range g.Edges.DstToSrcs {
		for edgesKey := range g.Edges.DstToSrcs[key] {
			dst := nodes[key]
			src := nodes[edgesKey]

			if repo, ok := repos[src]; ok {
				repo.For = codeArs[dst]

			} else {
				doCallRelation(src, dst)
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
