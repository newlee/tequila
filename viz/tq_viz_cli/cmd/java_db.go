package cmd

import (
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
		commonTableFile := cmd.Flag("commonTable").Value.String()

		pf := viz.CreateRegexpFilter(pkgFilterFile)
		tf := viz.CreateRegexpFilter(tableFilterFile).AddExcludes(commonTableFile)
		jf := viz.CreateRegexpFilter(filterFile)

		if cmd.Flag("reverse").Value.String() == "false" {
			pkgFilter = pf.NotMatch
			tableFilter = tf.UnMatch
			javaFilter = jf.Match
		} else {
			pkgFilter = pf.Match
			tableFilter = tf.Match
			javaFilter = jf.NotMatch
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

		doFiles(codeFiles, func(line, codeFileName string) {
			line = strings.ToUpper(line)
			fields := strings.Fields(line)
			if len(fields) == 0 || isComment(fields[0]) {
				return
			}

			if strings.Contains(line, "PKG_") {
				doPkgLine(line, pkgFilter, func(pkg string, sp string) {
					split := strings.Split(codeFileName, "/")
					pkgCallerFiles[split[len(split)-1]] = ""
					allPkg.Add(pkg, sp)
				})
			}

			if strings.Contains(line, " T_") || strings.Contains(line, ",T_") {
				doTableLine(line, tableFilter, func(table string) {
					split := strings.Split(codeFileName, "/")
					s := split[len(split)-1]
					if _, ok := tableCallerFiles[s]; !ok {
						tableCallerFiles[s] = ""
					}
					tableCallerFiles[s] = tableCallerFiles[s] + "," + table
					//tableCallerFiles[codeFileName] = ""
					tables.Add(table)
				})
			}
		})

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
		for key, value := range tableCallerFiles {
			fmt.Println(key + " -- " + value)

		}
	},
}

func init() {
	rootCmd.AddCommand(javaDbCmd)

	javaDbCmd.Flags().StringP("source", "s", "", "source code directory")
	javaDbCmd.Flags().StringP("filter", "f", "java", "file filter")
	javaDbCmd.Flags().StringP("table", "t", "table", "table filter file")
	javaDbCmd.Flags().StringP("commonTable", "c", "table_common", "common table file")
	javaDbCmd.Flags().StringP("package", "p", "pkg", "package filter file")
	javaDbCmd.Flags().BoolP("reverse", "R", false, "reverse dep")
}
