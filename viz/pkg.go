package viz

import (
	"sort"
	"strings"
	"fmt"
)

type Procedure struct {
	Name string
	FullName string
	Count int
	CallProcedures map[string]*Procedure
	BePrint bool
}

type Pkg struct {
	Name string
	Procedures map[string]*Procedure
}

type AllPkg struct {
	Pkgs map[string]*Pkg
}

type AllProcedure struct {
	Procedures map[string]*Procedure
}

func NewAllPkg() *AllPkg  {
	return &AllPkg{Pkgs:make(map[string]*Pkg)}
}

func NewAllProcedure() *AllProcedure  {
	return &AllProcedure{Procedures:make(map[string]*Procedure)}
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

func (all *AllProcedure) Add(pkgName, procedure string) {
	fullName := procedure
	if pkgName != "" {
		fullName = pkgName + "." + procedure
	}

	if _, ok := all.Procedures[fullName]; !ok {
		all.Procedures[fullName] = &Procedure {Name:procedure, FullName: fullName, CallProcedures:make(map[string]*Procedure)}
	}
}

func (all *AllProcedure) AddCall(pkgName, procedure,callPkgName,callProcedure string) {
	fullName := procedure
	if pkgName != "" {
		fullName = pkgName + "." + procedure
	}
	cFullName := callProcedure
	if callPkgName != "" {
		cFullName = callPkgName + "." + callProcedure
	}

	if _, ok := all.Procedures[fullName]; ok {
		p:= all.Procedures[fullName]
		if _, ok := all.Procedures[cFullName]; ok {
			p.CallProcedures[cFullName] =all.Procedures[cFullName]
		}
	}
}

func (all *AllPkg) Exist(name string) bool {
	if _, ok := all.Pkgs[name]; ok {
		return true
	}

	return false
}

func (all *AllPkg) ExistSp(name string, procedure string) bool {
	if _, ok := all.Pkgs[name]; ok {
		if _, ok := all.Pkgs[name].Procedures[procedure]; ok {
			return true
		}
	}

	return false
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

var pTree map[string]string

func (all *AllProcedure) Print(fullName string) map[string]string{
	pTree = make(map[string]string)
	if _, ok := all.Procedures[fullName]; ok {
		all.Procedures[fullName].Print(fullName)
	}
	fmt.Println(len(pTree))
	return pTree
}

func (p *Procedure) Print(fullName string)  {
	if p.BePrint {
		return
	}
	p.BePrint = true
	for key, procedure := range p.CallProcedures {
		if key != fullName {
			pTree[fmt.Sprintf("%s -> %s", fullName, key)] = ""
			procedure.Print(key)
		}
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