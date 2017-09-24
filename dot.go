package main

import (
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
)

type Entity struct {
	name         string
	entities     []*Entity
	vos          []*ValueObject
	Refs         []*Entity
	callEntities []*Entity
}

type ValueObject struct {
	name string
}

type Repository struct {
	name string
	For  *Entity
}
type Provider struct {
	name string
}
type SubDomain struct {
	ARs       map[string]*Entity
	Repos     map[string]*Repository
	Providers map[string]*Provider
	es        map[string]*Entity
	vos       map[string]*ValueObject
}
type Model struct {
	SubDomains map[string]*SubDomain
}

func (subDomain *SubDomain) Validate() bool {
	for key := range subDomain.ARs {
		ar := subDomain.ARs[key]
		for _, cEntity := range ar.callEntities {
			if _, ok := subDomain.ARs[cEntity.name]; !ok {
				return false
			}
		}
	}
	return true
}

func (model *Model) Validate() bool {
	for key := range model.SubDomains {
		if !model.SubDomains[key].Validate() {
			return false
		}
	}
	return true
}

func (entity *Entity) findEntity(name string) (*Entity, bool) {
	if entity.name == name {
		return entity, true
	} else {
		for _, et := range entity.entities {
			if finded, ok := et.findEntity(name); ok {
				return finded, true
			}
		}
	}
	return nil, false
}

func (entity *Entity) ChildrenEntities() []*Entity {
	return entity.entities
}

func (entity *Entity) ChildrenValueObjects() []*ValueObject {
	return entity.vos
}
func (entity *Entity) Compare(other *Entity) bool {
	if len(entity.entities) != len(other.entities) {
		return false
	}
	if len(entity.vos) != len(other.vos) {
		return false
	}
	em := make(map[string]*Entity)
	for _, childEntity := range entity.entities {
		em[childEntity.name] = childEntity
	}
	for _, childEntity := range other.entities {
		if !em[childEntity.name].Compare(childEntity) {
			return false
		}
	}
	vom := make(map[string]*ValueObject)
	for _, vo := range entity.vos {
		vom[vo.name] = vo
	}
	for _, vo := range other.vos {
		if _, ok := vom[vo.name]; !ok {
			return false
		}
	}
	return true
}

func (repo *Repository) Compare(other *Repository) bool {
	return repo.For.Compare(other.For)
}

func (subDomain *SubDomain) Compare(other *SubDomain) bool {
	if len(subDomain.ARs) != len(other.ARs) {
		return false
	}
	if len(subDomain.Repos) != len(other.Repos) {
		return false
	}
	for key := range subDomain.ARs {
		ar := subDomain.ARs[key]
		if !ar.Compare(other.ARs[key]) {
			return false
		}
	}
	for key := range subDomain.Repos {
		repo := subDomain.Repos[key]
		if !repo.Compare(other.Repos[key]) {
			return false
		}
	}
	return true
}
func createSubDomain() *SubDomain {
	ars := make(map[string]*Entity)
	es := make(map[string]*Entity)
	vos := make(map[string]*ValueObject)
	repos := make(map[string]*Repository)
	providers := make(map[string]*Provider)
	return &SubDomain{ARs: ars, Repos: repos, Providers: providers, es: es, vos: vos}
}
func (model *Model) Compare(other *Model) bool {
	if len(model.SubDomains) != len(model.SubDomains) {
		return false
	}

	for key := range model.SubDomains {
		ar := model.SubDomains[key]
		if !ar.Compare(other.SubDomains[key]) {
			return false
		}
	}

	return true
}
func Parse(dotFile string) *Model {
	fbuf, _ := ioutil.ReadFile(dotFile)
	g, _ := gographviz.Read(fbuf)

	// fmt.Println(g.Nodes.Nodes[0].Attrs["comment"])

	c2pMap := make(map[string]string)
	p2c := g.Relations.ParentToChildren

	subDomains := make(map[string]*SubDomain)

	if _, ok := p2c["g"]; ok {
		for key := range p2c["g"] {
			c2pMap[key] = "subdomain"
		}
		subDomains["subdomain"] = createSubDomain()
	} else {
		for clusterKey := range p2c {
			subDomainName := g.SubGraphs.SubGraphs[clusterKey].Attrs["label"]
			for key := range p2c[clusterKey] {
				c2pMap[key] =subDomainName
			}
			subDomains[subDomainName] = createSubDomain()
		}
	}

	for _, node := range g.Nodes.Nodes {
		subDomain := subDomains[c2pMap[node.Name]]
		if node.Attrs["comment"] == "AR" {
			subDomain.ARs[node.Name] = &Entity{name: node.Name}
		}
		if node.Attrs["comment"] == "E" {
			subDomain.es[node.Name] = &Entity{name: node.Name}
		}
		if node.Attrs["comment"] == "VO" {
			subDomain.vos[node.Name] = &ValueObject{name: node.Name}
		}
		if node.Attrs["comment"] == "Repo" {
			subDomain.Repos[node.Name] = &Repository{name: node.Name}
		}
		if node.Attrs["comment"] == "Provider" {
			subDomain.Providers[node.Name] = &Provider{name: node.Name}
		}
	}

	for key := range g.Edges.SrcToDsts {
		subDomain := subDomains[c2pMap[key]]
		ars := subDomain.ARs
		es := subDomain.es
		vos := subDomain.vos
		repos := subDomain.Repos
		if ar, ok := ars[key]; ok {
			for ckey := range g.Edges.SrcToDsts[key] {
				if ref, ok := ars[ckey]; ok {
					ar.Refs = append(ar.Refs, ref)
				}
				if et, ok := es[ckey]; ok {
					ar.entities = append(ar.entities, et)
				}
				if vo, ok := vos[ckey]; ok {
					ar.vos = append(ar.vos, vo)
				}
			}
		}

		if entity, ok := es[key]; ok {
			for ckey := range g.Edges.SrcToDsts[key] {
				if et, ok := es[ckey]; ok {
					entity.entities = append(entity.entities, et)
				}
				if vo, ok := vos[ckey]; ok {
					entity.vos = append(entity.vos, vo)
				}
			}
		}
		if repo, ok := repos[key]; ok {
			for ckey := range g.Edges.SrcToDsts[key] {
				repo.For = ars[ckey]
			}
		}
	}

	return &Model{SubDomains: subDomains}
}
