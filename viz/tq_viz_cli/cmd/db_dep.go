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
		filterFile := cmd.Flag("filter").Value.String()

		pf := viz.CreateRegexpFilter(filterFile)

		var match func(name string) bool
		if cmd.Flag("reverse").Value.String() == "true" {
			match = pf.NotMatch
		} else {
			match = pf.Match
		}

		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			codeFiles = append(codeFiles, path)
			return nil
		})

		allP := parseAllPkg(codeFiles)

		ps := make([]*viz.Procedure, 0)
		for name, p := range allP.Procedures {
			if match(strings.Split(name, ".")[0]) {
				ps = append(ps, p)

			}
		}
		sort.Slice(ps, func(i, j int) bool {
			return strings.Compare(ps[i].Name, ps[j].Name) < 0
		})

		for _, p := range ps {
			cs := make([]string, 0)
			for cname := range p.CallProcedures {
				if !match(strings.Split(cname, ".")[0]) {
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
	dbDepCmd.Flags().StringP("filter", "f", "pkg", "pkg regexp filter file")
	dbDepCmd.Flags().BoolP("reverse", "R", false, "reverse dep")
}
