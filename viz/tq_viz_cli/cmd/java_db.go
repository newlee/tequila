package cmd

import (
	"github.com/spf13/cobra"
	"strings"
	"path/filepath"
	"os"
	"bufio"
	"fmt"
	"github.com/newlee/tequila/viz"
)

var javaDbCmd *cobra.Command = &cobra.Command{
	Use:   "jd",
	Short: "icall grpah",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		filter := cmd.Flag("filter").Value.String()
		call := cmd.Flag("call").Value.String()
		tableKey := cmd.Flag("table").Value.String()

		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".java")  {
				codeFiles = append(codeFiles, path)
			}
			split := strings.Split(filter, ",")
			if call == "true" {
				for _, key := range split {
					if strings.Contains(strings.ToUpper(path), key) {
						codeFiles = append(codeFiles, path)
					}
				}
			}else{
				for _, key := range split {
					if strings.Contains(strings.ToUpper(path), key) {
						return nil
					}
				}
				//codeFiles = append(codeFiles, path)
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
				hasCallKey := false
				if call != "true" {
					split := strings.Split(call, ",")

					for _, key := range split {
						if strings.Contains(line,key) {
							hasCallKey = true
							break
						}
					}
				}else{
					hasCallKey = true
				}

				if strings.Contains(line,"PKG_")  && hasCallKey  {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r== '(' || r== ',' || r== '\''
					})

					for _, key := range tmp {
						if strings.HasPrefix(key,"PKG_") {
							sk := strings.Replace(key,"\"" ,"", -1)
							spss := strings.Split(sk, ".")
							pkg := spss[0]

							if len(spss) > 1 && strings.Contains(pkg,"PKG_") {
								split := strings.Split(codeFileName, "/")
								pkgCallerFiles[split[len(split) - 1]] = ""
								sp := spss[1]
								if strings.HasPrefix(sp, "P_"){
									allPkg.Add(pkg, sp)
								}

							}

						}
					}
				}

				if strings.Contains(line," T_") || strings.Contains(line, ",T_") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r == ',' || r == '.' || r=='"' || r == ':' || r == '(' ||  r == ')' ||  r == '）' ||  r == '%' ||  r == '!' ||  r == '\''
					})
					for _, t3 := range tmp {
						if strings.HasPrefix(t3,tableKey) {
							split := strings.Split(codeFileName, "/")
							tableCallerFiles[split[len(split) - 1]] = ""
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
	javaDbCmd.Flags().StringP("filter", "f", "life", "file filter")
	javaDbCmd.Flags().StringP("table", "t", "T_", "table filter")
	javaDbCmd.Flags().StringP("call", "c", "true", "filter  by call")
	javaDbCmd.Flags().StringP("output", "o", "java_db.dot", "output dot file name")
}
