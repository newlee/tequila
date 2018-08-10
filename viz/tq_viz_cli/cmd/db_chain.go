package cmd

import (
	"github.com/spf13/cobra"
	"strings"
	"path/filepath"
	"os"
	"bufio"
	"github.com/newlee/tequila/viz"
	"github.com/awalterschulze/gographviz"
	"strconv"
)

var DbChainCmd *cobra.Command = &cobra.Command{
	Use:   "dc",
	Short: "icall grpah",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		point := cmd.Flag("point").Value.String()
		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			codeFiles = append(codeFiles, path)
			return nil
		})

		allP := viz.NewAllProcedure()

		for _, codeFileName := range codeFiles {
			codeFile, _ := os.Open(codeFileName)
			scanner := bufio.NewScanner(codeFile)
			scanner.Split(bufio.ScanLines)

			pkg := ""
			for scanner.Scan() {
				line := scanner.Text()
				line = strings.ToUpper(line)

				if strings.Contains(line,"PKG_")  && strings.Contains(line,"CREATE") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r== '(' || r== ',' || r== '\''|| r== '"'
					})

					for _, key := range tmp {
						if strings.HasPrefix(key,"PKG_") {
							pkg = key
						}
					}
				}

				if strings.Contains(line,"PROCEDURE")  && strings.Contains(line,"P_") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r== '(' || r== ',' || r== '\''|| r== '"'
					})

					for _, key := range tmp {
						if strings.HasPrefix(key,"P_") {
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
			for scanner.Scan() {
				line := scanner.Text()
				line = strings.ToUpper(line)

				if strings.Contains(line,"PKG_")  && strings.Contains(line,"CREATE") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r== '(' || r== ',' || r== '\''|| r== '"'
					})

					for _, key := range tmp {
						if strings.HasPrefix(key,"PKG_") {
							pkg = key
						}
					}
				}

				if strings.Contains(line,"PROCEDURE")  && strings.Contains(line,"P_") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r== '(' || r== ',' || r== '\''|| r== '"'
					})

					for _, key := range tmp {
						if strings.HasPrefix(key,"P_") {
							procedure = key
						}
					}
				}

				if strings.Contains(line,"PKG_")  && strings.Contains(line,"(") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r== '(' || r== ',' || r== '\''|| r== '"' || r== ')'
					})

					for _, key := range tmp {
						if strings.HasPrefix(key,"PKG_") {
							sk := strings.Replace(key,"\"" ,"", -1)
							spss := strings.Split(sk, ".")
							if len(spss) > 1 && strings.Contains(pkg,"PKG_") {
								sp := spss[1]
								if strings.HasPrefix(sp, "P_"){
									allP.AddCall(pkg, procedure, spss[0], spss[1])
								}

							}
						}
					}
				}

				if strings.Contains(line,"P_")  && strings.Contains(line,"(") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r== '(' || r== ',' || r== '\''|| r== '"' || r== ')'
					})

					for _, key := range tmp {
						if strings.HasPrefix(key,"P_") {
							allP.AddCall(pkg, procedure, pkg, key)
						}
					}
				}
			}

			codeFile.Close()
		}

		pTree := allP.Print(point)
		graph := gographviz.NewGraph()
		graph.SetName("G")

		nodeIndex := 1
		nodes := make(map[string]string)

		for key := range pTree {
			tmp := strings.Split(key, " -> ")
			nodes[tmp[0]] = ""
			nodes[tmp[1]] = ""
		}

		for node := range nodes {
			attrs := make(map[string]string)
			attrs["label"] = "\"" + node + "\""
			attrs["shape"] = "box"
			graph.AddNode("G", "node"+strconv.Itoa(nodeIndex), attrs)
			nodes[node] = "node" + strconv.Itoa(nodeIndex)
			nodeIndex++
		}

		for key := range pTree {
			attrs := make(map[string]string)
			tmp := strings.Split(key, " -> ")
			graph.AddEdge(nodes[tmp[0]], nodes[tmp[1]], true, attrs)
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
	DbChainCmd.Flags().StringP("output", "o", "tree.dot", "output dot file name")
}
