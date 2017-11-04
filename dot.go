package main

import (
	"github.com/awalterschulze/gographviz"
	. "github.com/newlee/tequila/model"
	"io/ioutil"
	//"fmt"
)

func edgesKey(edges map[string][]*gographviz.Edge) []string {
	result := make([]string, 0)
	for key := range edges {
		result = append(result, key)
	}
	return result
}

func ParseProblemModel(dotFile string) *ProblemModel {
	fbuf, _ := ioutil.ReadFile(dotFile)
	g, _ := gographviz.Read(fbuf)

	c2pMap := make(map[string]string)
	p2c := g.Relations.ParentToChildren

	subDomains := make(map[string]*SubDomain)

	if _, ok := p2c["g"]; ok {
		for key := range p2c["g"] {
			c2pMap[key] = "subdomain"
		}
		subDomains["subdomain"] = NewSubDomain()
	} else {
		for clusterKey := range p2c {
			subDomainName := g.SubGraphs.SubGraphs[clusterKey].Attrs["label"]
			for key := range p2c[clusterKey] {
				c2pMap[key] = subDomainName
			}
			subDomains[subDomainName] = NewSubDomain()
		}
	}
	cms := InitCommentMapping()
	for _, node := range g.Nodes.Nodes {
		subDomain := subDomains[c2pMap[node.Name]]
		subDomain.AddNode(cms, node.Name, node.Attrs["comment"])
	}

	for key := range g.Edges.SrcToDsts {
		edgeKeys := edgesKey(g.Edges.SrcToDsts[key])
		subDomain := subDomains[c2pMap[key]]
		subDomain.AddRelations(key, edgeKeys)
	}

	return &ProblemModel{SubDomains: subDomains}
}

func ParseSolutionModel(dotFile string) *BCModel {
	fbuf, _ := ioutil.ReadFile(dotFile)
	g, _ := gographviz.Read(fbuf)

	p2c := g.Relations.ParentToChildren

	model := NewBCModel()
	for clusterKey := range p2c {
		if clusterKey != "g" {
			layerName := g.SubGraphs.SubGraphs[clusterKey].Attrs["label"]
			model.AppendLayer(layerName)
			for key := range p2c[clusterKey] {
				model.AppendNode(layerName, key)
			}
		}
	}
	cms := InitCommentMapping()
	for _, node := range g.Nodes.Nodes {
		model.AddNode(cms, node.Name, node.Attrs["comment"])
	}

	for key := range g.Edges.SrcToDsts {
		edgeKeys := edgesKey(g.Edges.SrcToDsts[key])

		model.AddRelations(key, edgeKeys)
	}

	return model
}
