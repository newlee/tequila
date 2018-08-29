package cmd

import (
	"bufio"
	"fmt"
	"github.com/awalterschulze/gographviz"
	. "github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var tableCmd *cobra.Command = &cobra.Command{
	Use:   "table",
	Short: "table grpah",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		tableFileName := fmt.Sprintf("%s/table.sql", source)
		keyFileName := fmt.Sprintf("%s/pk_fk.sql", source)

		tableFile, _ := os.Open(tableFileName)

		defer tableFile.Close()
		scanner := bufio.NewScanner(tableFile)
		scanner.Split(bufio.ScanLines)

		tables := make(map[string]string, 0)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "CREATE TABLE ") && strings.Contains(line, "CLAIM") {
				tmp := strings.Split(line, "\"")
				tables[tmp[3]] = ""

			}
		}

		keyFile, _ := os.Open(keyFileName)

		defer keyFile.Close()
		scanner = bufio.NewScanner(keyFile)
		scanner.Split(bufio.ScanLines)

		relations := make([]*Relation, 0)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "FOREIGN KEY") {
				tmp := strings.Split(line, "\"")
				table := tmp[3]
				scanner.Scan()
				line = scanner.Text()
				tmp = strings.Split(line, "\"")
				rTable := tmp[3]
				relations = append(relations, &Relation{From: table, To: rTable})
			}
		}

		graph := gographviz.NewGraph()
		graph.SetName("G")

		for t := range tables {
			attrs := make(map[string]string)
			if strings.Contains(t, "CLAIM") {
				graph.AddNode("G", t, attrs)
			}
		}

		for _, r := range relations {
			attrs := make(map[string]string)
			if strings.Contains(r.From, "CLAIM") || strings.Contains(r.To, "CLAIM") {
				if !strings.Contains(r.From, "CLAIM") {
					fmt.Println("from: " + r.From + "  -> " + r.To)
				}
				if !strings.Contains(r.To, "CLAIM") {
					fmt.Println("to: " + r.From + "  -> " + r.To)
				}
				if _, ok := tables[r.From]; !ok {
					tables[r.From] = ""
					graph.AddNode("G", r.From, attrs)
				}
				if _, ok := tables[r.To]; !ok {
					tables[r.To] = ""
					graph.AddNode("G", r.To, attrs)
				}
				graph.AddEdge(r.From, r.To, true, attrs)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(tableCmd)

	tableCmd.Flags().StringP("source", "s", "", "source code directory")
	tableCmd.Flags().StringP("output", "o", "table.dot", "output dot file name")
}
