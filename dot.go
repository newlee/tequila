package main

import (
	"github.com/awalterschulze/gographviz"
	. "github.com/newlee/tequila/model"
	"io/ioutil"
)

func edgesKey(edges map[string][]*gographviz.Edge) []string {
	result := make([]string, 0)
	for key := range edges {
		result = append(result, key)
	}
	return result
}

func Parse(dotFile string) *ProblemModel {
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
