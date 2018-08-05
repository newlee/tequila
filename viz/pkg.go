package viz

import (
	"sort"
	"strings"
	"fmt"
)

type Procedure struct {
	Name string
	Count int
}

type Pkg struct {
	Name string
	Procedures map[string]*Procedure
}

type AllPkg struct {
	Pkgs map[string]*Pkg
}

func NewAllPkg() *AllPkg  {
	return &AllPkg{Pkgs:make(map[string]*Pkg)}
}

func (all *AllPkg) Add(pkgName, procedure string) {
	if _, ok := all.Pkgs[pkgName]; !ok {
		all.Pkgs[pkgName] = &Pkg{Name: pkgName, Procedures:make(map[string]*Procedure)}
	}

	pkg := all.Pkgs[pkgName]
	if _, ok := pkg.Procedures[procedure]; !ok {
		pkg.Procedures[procedure] = &Procedure{Name: procedure, Count:0}
	}

	pkg.Procedures[procedure].Count++
}

func (all *AllPkg) Print() {
	pkgs := make([]*Pkg,0)
	for key := range all.Pkgs {
		pkgs= append(pkgs, all.Pkgs[key])
	}
	sort.Slice(pkgs, func(i, j int) bool {
		return strings.Compare(pkgs[i].Name, pkgs[j].Name ) < 0
	})

	for _, pkg := range pkgs {
		pkg.Print()
	}
}
func (pkg *Pkg) Print() {
	procedures := make([]*Procedure,0)
	count := 0
	for key := range pkg.Procedures {
		procedure := pkg.Procedures[key]
		procedures= append(procedures, procedure)
		count += procedure.Count
	}
	sort.Slice(procedures, func(i, j int) bool {
		return strings.Compare(procedures[i].Name, procedures[j].Name ) < 0
	})

	fmt.Printf("%s : %d\n", pkg.Name, count)
	for _, procedure := range procedures {
		fmt.Printf("  %s : %d\n",procedure.Name, procedure.Count)
	}
}