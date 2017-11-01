package dot

import (
	"github.com/awalterschulze/gographviz"
	"io/ioutil"
	"strings"
)

type Node struct {
	Name     string
	DstNodes []*Relation
	hasSrc   bool
}

type Relation struct {
	Node  *Node
	Style string
}

func (node *Node) removeGenericRelation(genericRelationMap map[string]*Node) {
	for index, relation := range node.DstNodes {
		if _, ok := genericRelationMap[relation.Node.Name]; ok {
			node.DstNodes[index] = relation.Node.DstNodes[0]
		}
	}
	for _, relation := range node.DstNodes {
		relation.Node.removeGenericRelation(genericRelationMap)
	}
}

func getMethodName(fullMethodName, split string, subDomainCallback func(string, string)) (string, string, bool) {
	if strings.Contains(fullMethodName, split) {
		tmp := strings.Split(fullMethodName, split)
		methodName := tmp[len(tmp)-1]
		methodName = strings.Replace(methodName, "\\l", "", -1)
		subDomain := strings.Replace(tmp[0], "\\l", "", -1)
		subDomainCallback(subDomain, methodName)

		return methodName, subDomain, true
	}
	return fullMethodName, "", false
}

func (node *Node) RemoveNS(subDomainCallback func(string, string)) {
	fullMethodName := node.Name
	if methodName, _, ok := getMethodName(fullMethodName, "::", subDomainCallback); ok {
		node.Name = methodName
	} else {
		node.Name, _, _ = getMethodName(fullMethodName, ".", subDomainCallback)
	}
	for _, relation := range node.DstNodes {
		relation.Node.RemoveNS(subDomainCallback)
	}
}

func (node *Node) IsIt(it string) bool {
	name := node.Name
	if methodName, _, ok := getMethodName(name, "::", func(s string, s2 string) {}); ok {
		name = methodName
	} else {
		name, _, _ = getMethodName(name, ".", func(s string, s2 string) {})
	}
	return node.isIt(it) && name != it
}

func (node *Node) isIt(it string) bool {
	result := strings.HasSuffix(node.Name, it)
	if !result {
		for _, relation := range node.DstNodes {
			if relation.Style != "\"dashed\"" {
				return relation.Node.isIt(it)
			}
		}
	}

	return result
}

func ParseDoxygenFile(file string) *Node {
	fbuf, _ := ioutil.ReadFile(file)
	g, _ := gographviz.Read(fbuf)
	nodes := nodes(g, 1)

	nodeMap := make(map[string]*Node)
	genericRelationMap := make(map[string]*Node)
	for key := range g.Edges.DstToSrcs {
		for edgesKey := range g.Edges.DstToSrcs[key] {
			for _, edge := range g.Edges.DstToSrcs[key][edgesKey] {
				//doxygen use dir is "back"
				dst := nodes[edge.Src]
				src := nodes[edge.Dst]

				if _, ok := nodeMap[src]; !ok {
					nodeMap[src] = &Node{Name: src, DstNodes: make([]*Relation, 0)}
				}
				if _, ok := nodeMap[dst]; !ok {
					nodeMap[dst] = &Node{Name: dst, DstNodes: make([]*Relation, 0), hasSrc: true}
				} else {
					nodeMap[dst].hasSrc = true
				}

				nodeMap[src].DstNodes = append(nodeMap[src].DstNodes,
					&Relation{Node: nodeMap[dst], Style: edge.Attrs["style"]})

				if strings.Contains(src, "\\<") &&
					(strings.Contains(edge.Attrs["label"], "dummy_for_doxygen") ||
						strings.Contains(edge.Attrs["label"], "elements")) {
					genericRelationMap[src] = nodeMap[src]
				}
			}
		}
	}
	var result *Node
	for key := range nodeMap {
		if !nodeMap[key].hasSrc {
			result = nodeMap[key]
		}
	}

	result.removeGenericRelation(genericRelationMap)
	return result
}

func nodes(g *gographviz.Graph, index int) map[string]string {
	nodes := make(map[string]string)
	for _, node := range g.Nodes.Nodes {
		fullMethodName := strings.Replace(node.Attrs["label"], "\"", "", 2)
		methodName := strings.Replace(fullMethodName, "\\l", "", -1)
		nodes[node.Name] = methodName
	}
	return nodes
}
