package main

import (
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
)

type Entity struct {
	name     string
	entities []*Entity
	vos      []*ValueObject
	Refs     []*Entity
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
type Model struct {
	ARs       map[string]*Entity
	Repos     map[string]*Repository
	Providers map[string]*Provider
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

func (model *Model) Compare(other *Model) bool {
	if len(model.ARs) != len(other.ARs) {
		return false
	}
	if len(model.Repos) != len(other.Repos) {
		return false
	}
	for key := range model.ARs {
		ar := model.ARs[key]
		if !ar.Compare(other.ARs[key]) {
			return false
		}
	}
	for key := range model.Repos {
		repo := model.Repos[key]
		if !repo.Compare(other.Repos[key]) {
			return false
		}
	}
	return true
}
func Parse(dotFile string) *Model {
	fbuf, _ := ioutil.ReadFile(dotFile)
	g, _ := gographviz.Read(fbuf)

	// fmt.Println(g.Nodes.Nodes[0].Attrs["comment"])
	ars := make(map[string]*Entity)
	es := make(map[string]*Entity)
	vos := make(map[string]*ValueObject)
	repos := make(map[string]*Repository)
	providers := make(map[string]*Provider)

	for _, node := range g.Nodes.Nodes {
		if node.Attrs["comment"] == "AR" {
			ars[node.Name] = &Entity{name: node.Name}
		}
		if node.Attrs["comment"] == "E" {
			es[node.Name] = &Entity{name: node.Name}
		}
		if node.Attrs["comment"] == "VO" {
			vos[node.Name] = &ValueObject{name: node.Name}
		}
		if node.Attrs["comment"] == "Repo" {
			repos[node.Name] = &Repository{name: node.Name}
		}
		if node.Attrs["comment"] == "Provider" {
			providers[node.Name] = &Provider{name: node.Name}
		}
	}

	for key := range g.Edges.SrcToDsts {
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

	return &Model{ARs: ars, Repos: repos, Providers: providers}
}
