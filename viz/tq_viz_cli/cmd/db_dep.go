package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"github.com/newlee/tequila/viz"
	"sort"
)

var dbDepCmd *cobra.Command = &cobra.Command{
	Use:   "dd",
	Short: "database dependencies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//defer profile.Start().Stop()
		source := cmd.Flag("source").Value.String()
		filter := cmd.Flag("filter").Value.String()
		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			codeFiles = append(codeFiles, path)
			return nil
		})

		allP := parseAllPkg(codeFiles)
		var match func(name string) bool

		if cmd.Flag("reverse").Value.String() == "true" {
			match = func(name string) bool {
				return !strings.HasPrefix(name, filter)
			}

		} else {
			match = func(name string) bool {
				return strings.HasPrefix(name, filter)
			}
		}
		ps := make([]*viz.Procedure, 0)
		for name, p := range allP.Procedures {
			if match(name) {
				ps = append(ps, p)

			}
		}
		sort.Slice(ps, func(i, j int) bool {
			return strings.Compare(ps[i].Name, ps[j].Name) < 0
		})

		for _, p := range ps {
			cs := make([]string, 0)
			for cname := range p.CallProcedures {
				if !match(cname) {
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
	dbDepCmd.Flags().StringP("filter", "f", "coll__graph.dot", "dot file filter")
	dbDepCmd.Flags().BoolP("reverse", "R", false, "reverse dep")
	dbDepCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
}
