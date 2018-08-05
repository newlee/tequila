package viz

import (
	"sort"
	"strings"
	"fmt"
)

type Table struct {
	Name string
	Count int
}

type AllTable struct {
	Tables map[string]*Table
}

func NewAllTable() *AllTable  {
	return &AllTable{Tables:make(map[string]*Table)}
}

func (all *AllTable) Add(tableName string) {
	if _, ok := all.Tables[tableName]; !ok {
		all.Tables[tableName] = &Table{Name: tableName}
	}

	all.Tables[tableName].Count++
}

func (all *AllTable) Print() {
	tables := make([]*Table,0)
	for key := range all.Tables {
		tables= append(tables, all.Tables[key])
	}
	sort.Slice(tables, func(i, j int) bool {
		return strings.Compare(tables[i].Name, tables[j].Name ) < 0
	})

	for _, table := range tables {
		fmt.Printf("%s : %d\n", table.Name, table.Count)
	}
}

