package cmd

import (
	//"fmt"
	"fmt"
	. "github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var collCmd *cobra.Command = &cobra.Command{
	Use:   "coll",
	Short: "full collaboration grpahh",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		filter := cmd.Flag("filter").Value.String()
		ignore := cmd.Flag("ignore").Value.String()
		result := ParseColl(source, filter)

		if cmd.Flag("mergePackage").Value.String() == "true" {
			Level, _ = strconv.Atoi(cmd.Flag("mergeLevel").Value.String())

			result = result.MergeHeaderFile(MergePackageFunc)
		}
		if cmd.Flag("entryPoints").Value.String() == "true" {
			entryPoints := result.EntryPoints(MergePackageFunc)
			for _, cf := range entryPoints {
				fmt.Println(cf)
			}
			return
		}

		if cmd.Flag("fanInFanOut").Value.String() == "true" {
			fans := result.SortedByFan(MergePackageFunc)
			fmt.Println("Name\tTotal\tFan-In\tFan-Out")
			for _, fan := range fans {
				fmt.Printf("%s\t%v\t%v\t%v\n", fan.Name, fan.FanIn+fan.FanOut, fan.FanIn, fan.FanOut)
			}
			return
		}

		ignores := strings.Split(ignore, ",")
		var nodeFilter = func(key string) bool {
			for _, f := range ignores {
				if key == f {
					return true
				}
			}
			return false
		}
		if cmd.Flag("java").Value.String() != "" && cmd.Flag("common").Value.String() != "" {
			javaFilterFile := cmd.Flag("java").Value.String()
			commonFilterFile := cmd.Flag("common").Value.String()
			javaFilter := CreatePrefixFilter(javaFilterFile)

			commonFilter := CreatePrefixFilter(commonFilterFile)
			printRelation(result, javaFilter.Match, commonFilter.Match)
			return

		}
		result.ToDot(cmd.Flag("output").Value.String(), ".", nodeFilter)
	},
}

func init() {
	rootCmd.AddCommand(collCmd)

	collCmd.Flags().StringP("source", "s", "", "source code directory")
	collCmd.Flags().StringP("filter", "f", "coll__graph.dot", "dot file filter")
	collCmd.Flags().StringP("ignore", "i", "main.cpp,main", "ignore")
	collCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
	collCmd.Flags().BoolP("entryPoints", "E", false, "list entry points")
	collCmd.Flags().BoolP("fanInFanOut", "F", false, "sorted fan-in and fan-out")
	collCmd.Flags().BoolP("mergePackage", "P", false, "merge package/folder for include dependencies")
	collCmd.Flags().Int32P("mergeLevel", "L", 3, "merge package/folder level")
	collCmd.Flags().StringP("java", "j", "", "java class filter")
	collCmd.Flags().StringP("common", "c", "", "common java class")

}

func printRelation(f *FullGraph, javaFilter func(string) bool, commonFilter func(string) bool) {

	class2other, other2class, class2common, other2common, common2 := make([]*Relation, 0), make([]*Relation, 0), make([]*Relation, 0), make([]*Relation, 0), make([]*Relation, 0)

	var isOther = func(s string) bool {
		return !javaFilter(s) && !commonFilter(s)
	}
	for _, relation := range f.RelationList {
		if !strings.Contains(relation.To, ".") {
			continue
		}
		if javaFilter(relation.From) && !javaFilter(relation.To) {
			if commonFilter(relation.To) {
				class2common = append(class2common, relation)
			}
			if !commonFilter(relation.To) {
				class2other = append(class2other, relation)
			}
		}

		if isOther(relation.From) && !isOther(relation.To) {
			if javaFilter(relation.To) {
				other2class = append(other2class, relation)
			}
			if commonFilter(relation.To) {
				other2common = append(other2common, relation)
			}
		}

		if commonFilter(relation.From) && !commonFilter(relation.To) {
			common2 = append(common2, relation)
		}
	}

	fmt.Println("class2other:")
	for _, r := range class2other {
		fmt.Printf("%s -> %s\n", r.From, r.To)
	}
	fmt.Println("-----------")
	fmt.Println("")
	fmt.Println("other2class:")
	for _, r := range other2class {
		fmt.Printf("%s -> %s\n", r.From, r.To)
	}
	fmt.Println("-----------")
	fmt.Println("")
	fmt.Println("class2common:")
	for _, r := range class2common {
		fmt.Printf("%s -> %s\n", r.From, r.To)
	}
	fmt.Println("-----------")
	fmt.Println("")
	fmt.Println("other2common:")
	for _, r := range other2common {
		fmt.Printf("%s -> %s\n", r.From, r.To)
	}
	fmt.Println("-----------")
	fmt.Println("")
	fmt.Println("common2...:")
	for _, r := range common2 {
		fmt.Printf("%s -> %s\n", r.From, r.To)
	}
	fmt.Println("-----------")

}
