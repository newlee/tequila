package viz

import (
	"fmt"
	"sort"
	"strings"
)

var Level = 4

type Package struct {
	Name    string
	Imports map[string]string
}

type AllPackage struct {
	Packages map[string]*Package
}

func NewAllPackage() *AllPackage {
	return &AllPackage{Packages: make(map[string]*Package)}
}

func (all *AllPackage) Add(name string) *Package {
	tmp := strings.Split(name, ".")
	levelName := name
	if len(tmp) > Level {
		levelName = strings.Join(tmp[:(Level)], ".")
	}

	if _, ok := all.Packages[levelName]; !ok {
		all.Packages[levelName] = &Package{Name: levelName, Imports: make(map[string]string)}
	}
	return all.Packages[levelName]
}

func (p *Package) AddImport(importPkg string) {
	tmp := strings.Split(importPkg, ".")
	levelName := importPkg
	if len(tmp) > Level {
		levelName = strings.Join(tmp[:(Level)], ".")
	}

	if _, ok := p.Imports[levelName]; !ok {
		p.Imports[levelName] = ""
	}
}

func (all *AllPackage) Print() {
	pkgs := make([]*Package, 0)
	for key := range all.Packages {
		pkgs = append(pkgs, all.Packages[key])
	}
	sort.Slice(pkgs, func(i, j int) bool {
		return strings.Compare(pkgs[i].Name, pkgs[j].Name) < 0
	})

	for _, pkg := range pkgs {
		pkg.Print()
	}
}

func (p *Package) Print() {
	if len(p.Imports) > 0 {
		fmt.Printf("%s\n", p.Name)
		return
	}
	//pkgs := make([]string,0)
	//for key := range p.Imports {
	//	pkgs= append(pkgs, key)
	//}
	//sort.Slice(pkgs, func(i, j int) bool {
	//	return strings.Compare(pkgs[i], pkgs[j] ) < 0
	//})
	//
	//fmt.Printf("%s :\n", p.Name)
	//for _, importPkg := range pkgs {
	//	fmt.Printf("  %s\n",importPkg)
	//}
}
