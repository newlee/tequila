package main

import (
	. "github.com/newlee/tequila/dot"
	. "github.com/newlee/tequila/model"
	"os"
	"path/filepath"
	"strings"
)

type parseResult struct {
	es  map[string]*Entity
	vos map[string]*ValueObject
}

var codeArs = make(map[string]*Entity)
var repos = make(map[string]*Repository)
var services = make(map[string]*Service)
var providers = make(map[string]*Provider)
var subDomainMap = make(map[string][]string)
var layerMap = make(map[string][]string)

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

func isProvider(node *Node) bool {
	return node.IsIt("Provider") && !strings.HasPrefix(node.Name, "Stub") && !strings.HasPrefix(node.Name, "Fake")
}

func isSerice(node *Node) bool {
	return strings.HasSuffix(node.Name, "Service")
}

func (result *parseResult) parseNode(node *Node, todo func(*Node, *parseResult)) *parseResult {
	todo(node, result)
	for _, relation := range node.DstNodes {
		result.parseNode(relation.Node, todo)
	}
	return result
}

func parseRelation(node *Node, result *parseResult) {
	src := node.Name
	if !isAggregateRoot(src) {
		isAr := isAR(node)
		if isAr {
			codeArs[src] = NewEntity(src)
		}
		if isEntity(node) && !isAr {
			result.es[src] = NewEntity(src)
		}
		if isValueObject(node) {
			result.vos[src] = NewValueObject(src)
		}
		if isSerice(node) {
			services[src] = NewService(src)
		}
		if isRepo(node) {
			repos[src] = NewRepository(src)
		}
		if isProvider(node) {
			providers[src] = NewProvider(src)
		}
	}
}

func parseAR(node *Node, result *parseResult) {
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
}

func parseEntity(node *Node, result *parseResult) {
	if entity, ok := result.es[node.Name]; ok {
		for _, relation := range node.DstNodes {
			dst := relation.Node.Name
			if et, ok := result.es[dst]; ok {
				entity.Entities = append(entity.Entities, et)
			}
			if vo, ok := result.vos[dst]; ok {
				entity.AppendVO(vo)
			}
		}
	}
}
func parseRepo(node *Node, result *parseResult) {
	if repo, ok := repos[node.Name]; ok {
		for _, relation := range node.DstNodes {
			dst := relation.Node.Name
			if dst != "Repository" {
				repo.For = dst
			}
		}
	}
}
func parseService(node *Node, result *parseResult) {
	if service, ok := services[node.Name]; ok {
		for _, relation := range node.DstNodes {
			dst := relation.Node.Name
			service.Refs = append(service.Refs, dst)
		}
	}
}
func fullNameCallback(fullName, methodName string) {
	tmp := strings.Split(fullName, "::")
	subDomain := tmp[0]
	if _, ok := subDomainMap[subDomain]; ok {
		subDomainMap[subDomain] = append(subDomainMap[subDomain], methodName)
	}
}
func parseCode(codeDotfile string) {
	node := ParseDoxygenFile(codeDotfile)
	node.RemoveNS(fullNameCallback)

	codeDotFileParseResult := &parseResult{
		es:  make(map[string]*Entity),
		vos: make(map[string]*ValueObject),
	}

	codeDotFileParseResult.parseNode(node, parseRelation).
		parseNode(node, parseAR).
		parseNode(node, parseEntity).
		parseNode(node, parseRepo)
}

func fullNameCallbackForLayer(fullName, methodName string) {
	tmp := strings.Split(fullName, "::")
	layer := tmp[0]
	if _, ok := layerMap[layer]; ok {
		layerMap[layer] = append(layerMap[layer], methodName)
	}
}

func parseCodeForSolution(codeDotfile string) {
	node := ParseDoxygenFile(codeDotfile)
	node.RemoveNS(fullNameCallbackForLayer)

	codeDotFileParseResult := &parseResult{
		es:  make(map[string]*Entity),
		vos: make(map[string]*ValueObject),
	}

	codeDotFileParseResult.parseNode(node, parseRelation).
		parseNode(node, parseAR).
		parseNode(node, parseEntity).
		parseNode(node, parseRepo).
		parseNode(node, parseService)
}

func ParseCodeProblemModel(codeDir string, subs []string) *ProblemModel {
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

	return &ProblemModel{SubDomains: subDomains}

}

func ParseCodeSolutionModel(codeDir string, layers []string) *BCModel {
	codeDotFiles := codeDotFiles(codeDir)
	codeArs = make(map[string]*Entity)
	repos = make(map[string]*Repository)
	services = make(map[string]*Service)
	providers = make(map[string]*Provider)
	layerMap = make(map[string][]string)
	for _, layer := range layers {
		layerMap[layer] = make([]string, 0)
	}
	for _, codeDotfile := range codeDotFiles {
		parseCodeForSolution(codeDotfile)
	}
	model := NewBCModel()
	for key := range layerMap {
		model.AppendLayer(key)

		for _, o := range layerMap[key] {
			if ar, ok := codeArs[o]; ok {
				model.AddARToLayer(key, ar)
			}
			if repo, ok := repos[o]; ok {
				model.AddRepoToLayer(key, repo)
			}
			if service, ok := services[o]; ok {
				model.AddServiceToLayer(key, service)
			}
			if provider, ok := providers[o]; ok && key == "services" {
				model.AddProviderToLayer(key, provider)
			}
		}
	}
	return model
}
