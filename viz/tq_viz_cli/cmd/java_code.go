package cmd

import (
	"bufio"
	"github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var javaCodeCmd *cobra.Command = &cobra.Command{
	Use:   "jc",
	Short: "icall grpah",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		call := cmd.Flag("call").Value.String()
		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".java") {
				codeFiles = append(codeFiles, path)
			}

			return nil
		})
		allPackage := viz.NewAllPackage()
		for _, codeFileName := range codeFiles {
			codeFile, _ := os.Open(codeFileName)
			scanner := bufio.NewScanner(codeFile)
			scanner.Split(bufio.ScanLines)
			var pkg *viz.Package
			for scanner.Scan() {
				line := scanner.Text()

				if strings.HasPrefix(line, "package") && !strings.Contains(line, call) {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r == ';'
					})
					//fmt.Println(tmp[1])
					pkg = allPackage.Add(tmp[1])
				}

				if pkg != nil && strings.HasPrefix(line, "import") && strings.Contains(line, call) {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r == ';'
					})
					importPkg := tmp[1]

					if strings.HasPrefix(importPkg, "com") {
						pkg.AddImport(importPkg)
					}

				}

				if strings.HasPrefix(line, "public") {
					break
				}
			}

			codeFile.Close()
		}
		allPackage.Print()
	},
}

func init() {
	rootCmd.AddCommand(javaCodeCmd)

	javaCodeCmd.Flags().StringP("source", "s", "", "source code directory")
	javaCodeCmd.Flags().StringP("call", "c", "", "filter  by call")
	javaCodeCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
}
