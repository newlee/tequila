package viz

import (
	"fmt"
	"sort"
	"strings"
)

type Table struct {
	Name  string
	Count int
}

type AllTable struct {
	Tables map[string]*Table
}

func NewAllTable() *AllTable {
	return &AllTable{Tables: make(map[string]*Table)}
}

func (all *AllTable) Add(tableName string) {
	if _, ok := all.Tables[tableName]; !ok {
		all.Tables[tableName] = &Table{Name: tableName}
	}

	all.Tables[tableName].Count++
}

func (all *AllTable) Print() {
	tables := make([]*Table, 0)
	for key := range all.Tables {
		tables = append(tables, all.Tables[key])
	}
	sort.Slice(tables, func(i, j int) bool {
		return strings.Compare(tables[i].Name, tables[j].Name) < 0
	})

	for _, table := range tables {
		fmt.Printf("%s : %d\n", table.Name, table.Count)
	}
}

type QueryTable struct {
	Name    string
	Alias   string
	Columns map[string]string
}

type Query struct {
	Sql    string
	Tables map[string]*QueryTable
}

func NewQuery(sql string) *Query {
	return &Query{Sql: sql, Tables: make(map[string]*QueryTable)}
}

func (query *Query) AddTable(name, alias string) {
	query.Tables[name] = &QueryTable{Name: name, Alias: alias, Columns: make(map[string]string)}
}

func (query *Query) AddColumn(name string) {
	if !strings.Contains(name, ".") {
		for key := range query.Tables {
			query.Tables[key].Columns[name] = name + "?"
			return
		}
	}

	for _, qt := range query.Tables {
		alias := strings.ToUpper(qt.Alias)
		if strings.HasPrefix(strings.ToUpper(name), alias+".") {
			tmp := strings.Split(name, ".")
			qt.Columns[tmp[1]] = name
			return
		}
		if strings.Contains(strings.ToUpper(name), "("+alias+".") {
			qt.Columns[name] = name
			return
		}
	}
}

func (qt *QueryTable) Merge(other *QueryTable) {
	for key := range other.Columns {
		if _, ok := qt.Columns[key]; !ok {
			qt.Columns[key] = key
		}
	}
}
func (query *Query) Merge(other *Query) {
	for key, qt := range other.Tables {
		if _, ok := query.Tables[key]; !ok {
			query.Tables[key] = qt
		} else {
			query.Tables[key].Merge(qt)
		}
	}
}

func (query *Query) ToString() {
	for key, qt := range query.Tables {
		fmt.Println(key)
		columns := make([]string, 0)
		for key := range qt.Columns {
			columns = append(columns, key)
		}

		sort.Slice(columns, func(i, j int) bool {
			return strings.Compare(columns[i], columns[j]) < 0
		})

		for _, key := range columns {
			fmt.Printf("  %s\n", key)
		}
	}
}
