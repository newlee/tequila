package cmd

import (
	"bufio"
	"fmt"
	"github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var pkgFilter func(line string) bool
var tableFilter func(line string) bool
var javaFilter func(line string) bool

var javaDbCmd *cobra.Command = &cobra.Command{
	Use:   "jd",
	Short: "java code to database dependencies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		source := cmd.Flag("source").Value.String()
		filterFile := cmd.Flag("filter").Value.String()
		pkgFilterFile := cmd.Flag("package").Value.String()
		tableFilterFile := cmd.Flag("table").Value.String()
		pkgRegs := readFilterFile(pkgFilterFile)
		tableRegs := readFilterFile(tableFilterFile)
		javaRegs := readFilterFile(filterFile)
		javaMatchRegs := make([]string, 0)
		javaBlackRegs := make([]string, 0)

		tableMatchRegs := make([]string, 0)
		tableBlackRegs := make([]string, 0)

		for _, reg := range javaRegs {
			if strings.HasPrefix(reg, "- ") {
				javaBlackRegs = append(javaBlackRegs, reg[2:])
			}else {
				javaMatchRegs = append(javaMatchRegs, reg)
			}
		}
		for _, reg := range tableRegs {
			if strings.HasPrefix(reg, "- ") {
				tableBlackRegs = append(tableBlackRegs, reg[2:])
			}else {
				tableMatchRegs = append(tableMatchRegs, reg)
			}
		}
		if cmd.Flag("reverse").Value.String() == "false" {
			pkgFilter = func(line string) bool {
				return unMatchByRegexps(line, pkgRegs)
			}
			tableFilter = func(line string) bool {
				return unMatchByRegexps(line, tableRegs) && !matchByRegexps(line, tableBlackRegs)
			}
			javaFilter = func(line string) bool {
				return matchByRegexps(line, javaMatchRegs) && !matchByRegexps(line, javaBlackRegs)
			}
		} else {
			pkgFilter = func(line string) bool {
				return matchByRegexps(line, pkgRegs)
			}
			tableFilter = func(line string) bool {
				return matchByRegexps(line, tableRegs) && !matchByRegexps(line, tableBlackRegs)
			}

			javaFilter = func(line string) bool {
				return !(matchByRegexps(line, javaMatchRegs) && !matchByRegexps(line, javaBlackRegs))
			}
		}

		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			if !strings.HasSuffix(path, ".java") {
				return nil
			}
			if javaFilter(path) {
				codeFiles = append(codeFiles, path)
			}
			return nil
		})

		allPkg := viz.NewAllPkg()
		tables := viz.NewAllTable()
		pkgCallerFiles := make(map[string]string)
		tableCallerFiles := make(map[string]string)

		for _, codeFileName := range codeFiles {
			codeFile, _ := os.Open(codeFileName)
			scanner := bufio.NewScanner(codeFile)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				line := scanner.Text()
				line = strings.ToUpper(line)
				fields := strings.Fields(line)
				if len(fields) == 0 {
					continue
				}
				first := fields[0]
				if strings.HasPrefix(first, "/*") || strings.HasPrefix(first, "*") || strings.HasPrefix(first, "//") {
					continue
				}
				if strings.Contains(line, "PKG_") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r == '(' || r == ',' || r == '\''
					})

					for _, key := range tmp {
						if strings.HasPrefix(key, "PKG_") && pkgFilter(strings.Split(key, ".")[0]) {
							sk := strings.Replace(key, "\"", "", -1)
							spss := strings.Split(sk, ".")
							pkg := spss[0]

							if len(spss) > 1 && strings.Contains(pkg, "PKG_") {
								split := strings.Split(codeFileName, "/")
								pkgCallerFiles[split[len(split)-1]] = ""
								sp := spss[1]
								if strings.HasPrefix(sp, "P_") {
									allPkg.Add(pkg, sp)
								}
							}
						}
					}
				}

				if strings.Contains(line, " T_") || strings.Contains(line, ",T_") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r == ',' || r == '.' || r == '"' || r == ':' || r == '(' || r == ')' || r == '）' || r == '%' || r == '!' || r == '\''
					})
					for _, t3 := range tmp {
						if strings.HasPrefix(t3, "T_") && tableFilter(t3) && !viz.IsChineseChar(t3) {
							split := strings.Split(codeFileName, "/")
							tableCallerFiles[split[len(split)-1]] = ""
							//tableCallerFiles[codeFileName] = ""
							tables.Add(t3)
						}
					}
				}
			}

			codeFile.Close()
		}

		allPkg.Print()
		fmt.Println("")
		fmt.Println("-----------")
		for key := range pkgCallerFiles {
			fmt.Println(key)

		}
		fmt.Println("")
		fmt.Println("-----------")

		tables.Print()

		fmt.Println("")
		fmt.Println("-----------")
		for key := range tableCallerFiles {
			fmt.Println(key)

		}
	},
}

func init() {
	rootCmd.AddCommand(javaDbCmd)

	javaDbCmd.Flags().StringP("source", "s", "", "source code directory")
	javaDbCmd.Flags().StringP("filter", "f", "java", "file filter")
	javaDbCmd.Flags().StringP("table", "t", "table", "table filter file")
	javaDbCmd.Flags().StringP("package", "p", "pkg", "package filter file")
	javaDbCmd.Flags().BoolP("reverse", "R", false, "reverse dep")
}
