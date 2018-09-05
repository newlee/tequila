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

	for _, codeFileName := range codeFiles {
		codeFile, _ := os.Open(codeFileName)
		scanner := bufio.NewScanner(codeFile)
		scanner.Split(bufio.ScanLines)

		pkg := ""
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.ToUpper(line)

			if strings.Contains(line, "PKG_") && strings.Contains(line, "CREATE") {
				tmp := strings.FieldsFunc(line, func(r rune) bool {
					return r == ' ' || r == '(' || r == ',' || r == '\'' || r == '"'
				})

				for _, key := range tmp {
					if strings.HasPrefix(key, "PKG_") {
						pkg = key
					}
				}
			}

			if strings.Contains(line, "PROCEDURE") && strings.Contains(line, "P_") {
				tmp := strings.FieldsFunc(line, func(r rune) bool {
					return r == ' ' || r == '(' || r == ',' || r == '\'' || r == '"'
				})

				for _, key := range tmp {
					if strings.HasPrefix(key, "P_") {
						allP.Add(pkg, key)
					}
				}
			}
		}

		codeFile.Close()
	}

	for _, codeFileName := range codeFiles {
		codeFile, _ := os.Open(codeFileName)
		scanner := bufio.NewScanner(codeFile)
		scanner.Split(bufio.ScanLines)

		pkg := ""
		procedure := ""
		isComments := false

		for scanner.Scan() {
			line := scanner.Text()
			line = strings.ToUpper(line)
			line = strings.Trim(line, " ")
			if strings.HasPrefix(line, "/*") {
				isComments = true
			}

			if strings.HasSuffix(line, "*/") || strings.HasSuffix(line, "*/;") {
				isComments = false
			}
			if isComments {
				continue
			}
			if strings.Contains(line, "PKG_") && strings.Contains(line, "CREATE") {
				tmp := strings.FieldsFunc(line, func(r rune) bool {
					return r == ' ' || r == '(' || r == ',' || r == '\'' || r == '"'
				})

				for _, key := range tmp {
					if strings.HasPrefix(key, "PKG_") {
						pkg = key
					}
				}
			}

			if strings.Contains(line, "PROCEDURE") && strings.Contains(line, "P_") {
				tmp := strings.FieldsFunc(line, func(r rune) bool {
					return r == ' ' || r == '(' || r == ',' || r == '\'' || r == '"'
				})

				for _, key := range tmp {
					if strings.HasPrefix(key, "P_") {
						procedure = key
					}
				}
			}

			if strings.HasPrefix(line, "--") {
				continue
			}

			if strings.Contains(line, "PKG_") && strings.Contains(line, "(") {
				tmp := strings.FieldsFunc(line, func(r rune) bool {
					return r == ' ' || r == '(' || r == ',' || r == '\'' || r == '"' || r == ')'
				})

				for _, key := range tmp {
					if strings.HasPrefix(key, "PKG_") {
						sk := strings.Replace(key, "\"", "", -1)
						spss := strings.Split(sk, ".")
						if len(spss) > 1 && strings.Contains(pkg, "PKG_") {
							sp := spss[1]
							if strings.HasPrefix(sp, "P_") {
								allP.AddCall(pkg, procedure, spss[0], spss[1])
							}

						}
					}
				}
			}

			if strings.Contains(line, "P_") && strings.Contains(line, "(") {
				tmp := strings.FieldsFunc(line, func(r rune) bool {
					return r == ' ' || r == '(' || r == ',' || r == '\'' || r == '"' || r == ')'
				})

				for _, key := range tmp {
					if strings.HasPrefix(key, "P_") {
						allP.AddCall(pkg, procedure, pkg, key)
					}
				}
			}

			if strings.Contains(line, " T_") || strings.Contains(line, ",T_") {
				tmp := strings.FieldsFunc(line, func(r rune) bool {
					return r == ' ' || r == ',' || r == '.' || r == '"' || r == ':' || r == '(' || r == ')' || r == '）' || r == '%' || r == '!' || r == '\''
				})
				for _, t3 := range tmp {
					if strings.HasPrefix(t3, "T_") && !viz.IsChineseChar(t3) && !strings.Contains(t3, ";") && !strings.Contains(t3, "、") {
						isWrite := strings.Contains(line, "INSERT ") || strings.Contains(line, "UPDATE ") || strings.Contains(line, "DELETE ")
						allP.AddTable(pkg, procedure, t3, isWrite)
					}

				}
			}
		}

		codeFile.Close()
	}
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
