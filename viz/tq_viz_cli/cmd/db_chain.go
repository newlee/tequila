package cmd

import (
	"bufio"
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func parseAllPkg(codeFiles []string) *viz.AllProcedure {
	allP := viz.NewAllProcedure()
	pkg := ""
	doFiles(codeFiles, func() {
		pkg = ""
	}, func(line, codeFileName string) {
		line = strings.ToUpper(line)

		if strings.Contains(line, "PKG_") && strings.Contains(line, "CREATE") {
			doCreatePkg(line, func(s string) {
				pkg = s
			})
		}

		if strings.Contains(line, "PROCEDURE") && strings.Contains(line, "P_") {
			doCreateProcedure(line, func(s string) {
				allP.Add(pkg, s)
			})
		}
	})
	procedure := ""
	isComments := false

	doFiles(codeFiles, func() {
		pkg = ""
		procedure = ""
		isComments = false
	}, func(line, codeFileName string) {
		line = strings.ToUpper(line)
		line = strings.Trim(line, " ")
		if strings.HasPrefix(line, "/*") {
			isComments = true
		}

		if strings.HasSuffix(line, "*/") || strings.HasSuffix(line, "*/;") {
			isComments = false
		}

		if isComments {
			return
		}

		if strings.HasPrefix(line, "--") {
			return
		}

		if strings.Contains(line, "PKG_") && strings.Contains(line, "CREATE") {
			doCreatePkg(line, func(s string) {
				pkg = s
			})
		}

		if strings.Contains(line, "PROCEDURE") && strings.Contains(line, "P_") {
			doCreateProcedure(line, func(s string) {
				procedure = s
			})
		}

		if strings.Contains(line, "PKG_") && strings.Contains(line, "(") {
			doPkgLine(line, emptyFilter, func(p string, sp string) {
				allP.AddCall(pkg, procedure, p, sp)
			})
		}

		if strings.Contains(line, "P_") && strings.Contains(line, "(") {
			doCreateProcedure(line, func(s string) {
				allP.AddCall(pkg, procedure, pkg, s)
			})
		}

		if strings.Contains(line, " T_") || strings.Contains(line, ",T_") {
			doSplit(line, tableSplit, func(table string) {
				if strings.HasPrefix(table, "T_") && !viz.IsChineseChar(table) && !strings.Contains(table, ";") && !strings.Contains(table, "、") {
					isWrite := strings.Contains(line, "INSERT ") || strings.Contains(line, "UPDATE ") || strings.Contains(line, "DELETE ")
					allP.AddTable(pkg, procedure, table, isWrite)
				}
			})
		}
	})

	return allP
}

var DbChainCmd *cobra.Command = &cobra.Command{
	Use:   "dc",
	Short: "database call chain grpah",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		point := cmd.Flag("point").Value.String()
		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			codeFiles = append(codeFiles, path)
			return nil
		})

		allP := parseAllPkg(codeFiles)

		pTree, pTables := allP.Print(point)

		tables := make([]string, 0)
		for key, rw := range pTables {
			tables = append(tables, fmt.Sprintf("%s [%s]", key, rw.ToString()))
		}
		sort.Slice(tables, func(i, j int) bool {
			return strings.Compare(tables[i], tables[j]) < 0
		})

		for _, t := range tables {
			fmt.Printf("%s\n", t)
		}
		graph := gographviz.NewGraph()
		graph.SetName("G")

		nodeIndex := 1
		nodes := make(map[string]string)

		for key := range pTree {
			tmp := strings.Split(key, " -> ")
			if cmd.Flag("mergePackage").Value.String() == "true" {
				nodes[strings.Split(tmp[0], ".")[0]] = ""
				nodes[strings.Split(tmp[1], ".")[0]] = ""
			} else {
				nodes[tmp[0]] = ""
				nodes[tmp[1]] = ""
			}
		}

		owner := cmd.Flag("key").Value.String()
		for node := range nodes {
			attrs := make(map[string]string)
			attrs["label"] = "\"" + node + "\""
			attrs["shape"] = "box"
			if owner != "" {
				if strings.Contains(node, owner) {
					attrs["color"] = "greenyellow"
					attrs["style"] = "filled"
				} else {
					attrs["color"] = "orange"
					attrs["style"] = "filled"
				}
			}

			graph.AddNode("G", "node"+strconv.Itoa(nodeIndex), attrs)
			nodes[node] = "node" + strconv.Itoa(nodeIndex)
			nodeIndex++
		}
		relations := make(map[string]string)
		for key := range pTree {
			attrs := make(map[string]string)
			tmp := strings.Split(key, " -> ")
			if cmd.Flag("mergePackage").Value.String() == "true" {
				pName0 := strings.Split(tmp[0], ".")[0]
				pName1 := strings.Split(tmp[1], ".")[0]

				if pName0 != pName1 {
					if _, ok := relations[pName0+pName1]; !ok {
						relations[pName0+pName1] = ""
						graph.AddEdge(nodes[pName0], nodes[pName1], true, attrs)
					}
				}
			} else {
				graph.AddEdge(nodes[tmp[0]], nodes[tmp[1]], true, attrs)
			}
		}

		f, _ := os.Create(cmd.Flag("output").Value.String())
		w := bufio.NewWriter(f)
		w.WriteString("di" + graph.String())
		w.Flush()

	},
}

func init() {
	rootCmd.AddCommand(DbChainCmd)

	DbChainCmd.Flags().StringP("source", "s", "", "source code directory")
	DbChainCmd.Flags().StringP("point", "p", "", "input point")
	DbChainCmd.Flags().StringP("key", "k", "", "owner key")
	DbChainCmd.Flags().StringP("output", "o", "tree.dot", "output dot file name")
	DbChainCmd.Flags().BoolP("mergePackage", "P", false, "merge package/folder for include dependencies")
}
