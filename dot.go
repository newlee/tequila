package main

import (
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
)

type Entity struct {
	name     string
	entities []*Entity
	vos      []*ValueObject
}

type ValueObject struct {
	name string
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

func Parse(dotFile string) map[string]*Entity {
	fbuf, _ := ioutil.ReadFile(dotFile)
	g, _ := gographviz.Read(fbuf)

	// fmt.Println(g.Nodes.Nodes[0].Attrs["comment"])
	ars := make(map[string]*Entity)
	es := make(map[string]*Entity)
	vos := make(map[string]*ValueObject)
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
