package cmd

import (
	"fmt"
	"github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var dbDepCmd *cobra.Command = &cobra.Command{
	Use:   "dd",
	Short: "database dependencies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//defer profile.Start().Stop()
		source := cmd.Flag("source").Value.String()
		pkgFilterFile := cmd.Flag("package").Value.String()
		tableFilterFile := cmd.Flag("table").Value.String()
		commonTableFile := cmd.Flag("commonTable").Value.String()

		pf := viz.CreateRegexpFilter(pkgFilterFile)
		tf := viz.CreateRegexpFilter(tableFilterFile).AddExcludes(commonTableFile)

		if cmd.Flag("reverse").Value.String() == "true" {
			pkgFilter = pf.NotMatch
			tableFilter = tf.Match
		} else {
			pkgFilter = pf.Match
			tableFilter = tf.UnMatch
		}

		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			codeFiles = append(codeFiles, path)
			return nil
		})

		allP := parseAllPkg(codeFiles)

		ps := make([]*viz.Procedure, 0)
		for name, p := range allP.Procedures {
			if pkgFilter(strings.Split(name, ".")[0]) {
				ps = append(ps, p)

			}
		}
		sort.Slice(ps, func(i, j int) bool {
			return strings.Compare(ps[i].Name, ps[j].Name) < 0
		})

		for _, p := range ps {
			cs := make([]string, 0)
			for cname := range p.CallProcedures {
				if !pkgFilter(strings.Split(cname, ".")[0]) {
					cs = append(cs, cname)
				}
			}
			if len(cs) > 0 {
				fmt.Println(p.FullName)
				for _, cname := range cs {
					fmt.Printf("  %s\n", cname)
				}
			}
		}
		fmt.Println("-----------")
		for _, p := range ps {
			cs := make([]string, 0)
			for cname := range p.Tables {
				if tableFilter(strings.Split(cname, ".")[0]) {
					cs = append(cs, cname)
				}
			}
			if len(cs) > 0 {
				fmt.Println(p.FullName)
				for _, cname := range cs {
					fmt.Printf("  %s\n", cname)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(dbDepCmd)

	dbDepCmd.Flags().StringP("source", "s", "", "source code directory")
	dbDepCmd.Flags().StringP("table", "t", "table", "table filter file")
	dbDepCmd.Flags().StringP("commonTable", "c", "table_common", "common table file")
	dbDepCmd.Flags().StringP("package", "p", "pkg", "package filter file")
	dbDepCmd.Flags().BoolP("reverse", "R", false, "reverse dep")
}
