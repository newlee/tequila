package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"github.com/newlee/tequila/viz"
	"sort"
	"regexp"
	"bufio"
)

func readFilterFile(fileName string) []string {
	result := make([]string, 0)
	f, _ := os.Open(fileName)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			result = append(result, strings.Trim(line," "))
		}
	}
	return result
}

func matchByRegexps(name string, regexps []string) bool  {
	for _, reg := range regexps {
		re, _ := regexp.Compile(reg)
		if re.MatchString(name) {
			return true
		}
	}
	return false
}

func unMatchByRegexps(name string, regexps []string) bool  {
	for _, reg := range regexps {
		re, _ := regexp.Compile(reg)
		if re.MatchString(name) {
			return false
		}
	}
	return true
}

var dbDepCmd *cobra.Command = &cobra.Command{
	Use:   "dd",
	Short: "database dependencies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//defer profile.Start().Stop()
		source := cmd.Flag("source").Value.String()
		filterFile := cmd.Flag("filter").Value.String()
		regexps := readFilterFile(filterFile)

		var match func(name string) bool
		if cmd.Flag("reverse").Value.String() == "true" {
			match = func(name string) bool {
				return unMatchByRegexps(name, regexps)
			}

		} else {
			match = func(name string) bool {
				return matchByRegexps(name, regexps)
			}
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
