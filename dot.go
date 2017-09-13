package main

import (
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
)

type AggregateRoot struct {
	name     string
	entities []*Entity
	vos      []*ValueObject
}
type Entity struct {
	name     string
	entities []*Entity
	vos      []*ValueObject
}

type ValueObject struct {
	name string
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
func (ar *AggregateRoot) Compare(other *AggregateRoot) bool {
	if len(ar.entities) != len(other.entities) {
		return false
	}
	if len(ar.vos) != len(other.vos) {
		return false
	}
	em := make(map[string]*Entity)
	for _, entity := range ar.entities {
		em[entity.name] = entity
	}
	for _, entity := range other.entities {
		if !em[entity.name].Compare(entity) {
			return false
		}
	}
	vom := make(map[string]*ValueObject)
	for _, vo := range ar.vos {
		vom[vo.name] = vo
	}
	for _, vo := range other.vos {
		if _, ok := vom[vo.name]; !ok {
			return false
		}
	}
	return true
}
func Parse(dotFile string) map[string]*AggregateRoot {
	fbuf, _ := ioutil.ReadFile(dotFile)
	g, _ := gographviz.Read(fbuf)

	// fmt.Println(g.Nodes.Nodes[0].Attrs["comment"])
	ars := make(map[string]*AggregateRoot)
	es := make(map[string]*Entity)
	vos := make(map[string]*ValueObject)
	for _, node := range g.Nodes.Nodes {
		if node.Attrs["comment"] == "AR" {
			ars[node.Name] = &AggregateRoot{name: node.Name}
		}
		if node.Attrs["comment"] == "E" {
			es[node.Name] = &Entity{name: node.Name}
		}
		if node.Attrs["comment"] == "VO" {
			vos[node.Name] = &ValueObject{name: node.Name}
		}
	}

	for key, _ := range g.Edges.SrcToDsts {
		if ar, ok := ars[key]; ok {
			for ckey, _ := range g.Edges.SrcToDsts[key] {
				if et, ok := es[ckey]; ok {
					ar.entities = append(ar.entities, et)
				}
				if vo, ok := vos[ckey]; ok {
					ar.vos = append(ar.vos, vo)
				}
			}
		}
		if entity, ok := es[key]; ok {
			for ckey, _ := range g.Edges.SrcToDsts[key] {
				if et, ok := es[ckey]; ok {
					entity.entities = append(entity.entities, et)
				}
				if vo, ok := vos[ckey]; ok {
					entity.vos = append(entity.vos, vo)
				}
			}
		}
	}
	return ars
}
