package main

import (
	. "github.com/newlee/tequila/dot"
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
		if isRepo(node) {
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
func (result *CodeDotFileParseResult) parseRepo(node *Node) {
	if repo, ok := repos[node.Name]; ok {
		for _, relation := range node.DstNodes {
			dst := relation.Node.Name

			if ar, ok := codeArs[dst]; ok {
				repo.For = ar
			}
		}
	}
	for _, relation := range node.DstNodes {
		result.parseRepo(relation.Node)
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
	codeDotFileParseResult.parseRepo(node)
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
