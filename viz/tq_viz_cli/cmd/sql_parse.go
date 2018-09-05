package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xwb1989/sqlparser"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/newlee/tequila/viz"
	"reflect"
	"regexp"
)

func printJoin(table sqlparser.TableExpr, currentQuery *viz.Query) {
	switch table := table.(type) {
	case *sqlparser.JoinTableExpr:
		printJoin(table.RightExpr, currentQuery)
		printJoin(table.LeftExpr, currentQuery)
		printWhen(table.Condition.On, currentQuery)
	case *sqlparser.AliasedTableExpr:
		tName := sqlparser.String(table.Expr)
		if strings.HasPrefix(tName, "ODSUSER.T_") {
			currentQuery.AddTable(tName, table.As.String())
		}

		if strings.HasPrefix(tName, "(select") && strings.Contains(tName, "ODSUSER.") {
			tName = strings.Trim(tName, " \n")
			parseQuery(tName[1 : len(tName)-1])
		}
	default:
		fmt.Println("----------------")
	}
}
func printWhen(when sqlparser.Expr, currentQuery *viz.Query) {
	switch when := when.(type) {
	case *sqlparser.ComparisonExpr:
		left := sqlparser.String(when.Left)
		if strings.Contains(left, ".") {
			currentQuery.AddColumn(left)
		}
		right := sqlparser.String(when.Right)
		if strings.Contains(right, ".") {
			currentQuery.AddColumn(right)
		}

		switch l := when.Left.(type) {
		case *sqlparser.Subquery:
			parseQuery(sqlparser.String(l.Select))
		default:
			printWhen(l, currentQuery)
		}
		switch l := when.Right.(type) {
		case *sqlparser.Subquery:
			parseQuery(sqlparser.String(l.Select))
		default:
			printWhen(l, currentQuery)
		}
	case *sqlparser.OrExpr:
		printWhen(when.Left, currentQuery)
		printWhen(when.Right, currentQuery)
	case *sqlparser.ExistsExpr:
		parseQuery(sqlparser.String(when.Subquery.Select))
	case *sqlparser.IsExpr:
		currentQuery.AddColumn(sqlparser.String(when.Expr))
	case *sqlparser.ParenExpr:
		printWhen(when.Expr, currentQuery)
	case *sqlparser.AndExpr:
		printWhen(when.Left, currentQuery)
		printWhen(when.Right, currentQuery)
	case *sqlparser.ColName:
		tname := when.Qualifier.Name.String()
		col := when.Name.String()
		if tname != "" {
			currentQuery.AddColumn(fmt.Sprintf("%s.%s", tname, col))
		}
	case *sqlparser.SubstrExpr:
		tname := when.Name.Qualifier.Name.String()
		col := when.Name.Name.String()
		if tname != "" {
			currentQuery.AddColumn(fmt.Sprintf("%s.%s", tname, col))
		}
	case *sqlparser.FuncExpr:
		tname := when.Qualifier.String()
		col := when.Name.String()
		if tname != "" {
			currentQuery.AddColumn(fmt.Sprintf("%s.%s", tname, col))
		}
	case *sqlparser.SQLVal:
	case *sqlparser.BinaryExpr:
	case sqlparser.ValTuple:
	case *sqlparser.UnaryExpr:

	default:
		fmt.Println(reflect.TypeOf(when).String() + " --- " + sqlparser.String(when))
	}
}

var queryArr = make([]*viz.Query, 0)

var sqlParseCmd *cobra.Command = &cobra.Command{
	Use:   "sp",
	Short: "query sql parse",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".bdy") || strings.HasSuffix(path, ".sql") {
				codeFiles = append(codeFiles, path)
			}

			return nil
		})

		querys := make([]string, 0)
		for _, codeFileName := range codeFiles {
			codeFile, _ := os.Open(codeFileName)
			scanner := bufio.NewScanner(codeFile)
			scanner.Split(bufio.ScanLines)
			query := ""
			begin := false
			for scanner.Scan() {

				line := scanner.Text()

				tmp := strings.Fields(strings.ToUpper(line))
				for _, word := range tmp {

					if word == "UPDATE" {
						begin = true
					}
					if !begin && strings.HasSuffix(word, "SELECT") {
						begin = true
					}

					if word != "" {
						break
					}
				}

				if begin {
					if strings.Contains(line, "--") && !strings.Contains(line, "/") {
						tmp := strings.Split(line, "--")
						line = tmp[0]
					}
					if strings.Contains(line, "F_QHH(") && strings.Contains(line, "SUBSTR") {
						re, _ := regexp.Compile("F_QHH\\(([\\S\\s]+?)\\)")
						submatch := re.FindStringSubmatch(line)

						line = strings.Replace(line, submatch[0], submatch[1], -1)
					}

					if strings.Contains(line, "TRIM(") && strings.Contains(line, "SUBSTR") {
						re, _ := regexp.Compile("TRIM\\(([\\S\\s]+?)\\)")
						submatch := re.FindStringSubmatch(line)

						line = strings.Replace(line, submatch[0], submatch[1], -1)
					}

					if strings.Contains(line, "ROW_NUMBER()") {
						line = "0"
					}

					if strings.Contains(line, "LENGTHB(") {
						line = strings.Replace(line, "LENGTHB(", "LENGTH(", -1)
					}

					if strings.Contains(strings.ToUpper(line), " DATE") && !strings.Contains(line, "SELECT") {
						line = strings.Replace(strings.ToUpper(line), " DATE", "", -1)
					}

					if strings.HasPrefix(line, "/*") {
						continue
					}

					query = query + line + "\n"
				}
				if strings.Contains(line, ";") && begin {
					begin = false
					querys = append(querys, query)
					query = ""
				}
			}
		}

		for _, query := range querys {
			query = strings.Trim(query, " ")
			if strings.HasPrefix(query, "(") {
				query = query[1 : len(query)-3]
			}
			parseQuery(query)
		}
		fq := &viz.Query{Sql: "", Tables: make(map[string]*viz.QueryTable)}
		for _, q := range queryArr {
			fq.Merge(q)
		}
		fq.ToString()

	},
}
var isRetry = false

func parseQuery(query string) {
	if !strings.Contains(strings.ToUpper(query), "ODSUSER.") {
		return
	}
	if strings.HasSuffix(query, "*/\n") {
		return
	}
	sql := strings.ToUpper(query)
	sql = strings.Replace(query, "(+)", "", -1)
	sql = strings.Replace(sql, ";", "", -1)
	sql = strings.Replace(sql, "INTO", "", -1)

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		if !isRetry {
			isRetry = true
			parseQuery(query + " AS T")
			isRetry = false
		}
		//if isRetry {
		//	fmt.Println(sql)
		//	fmt.Println("parse error: " + err.Error())
		//}
		return
	}
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		var currentQuery = viz.NewQuery(query)

		for _, node := range stmt.From {
			printJoin(node, currentQuery)
		}
		for _, node := range stmt.SelectExprs {
			switch node := node.(type) {
			case *sqlparser.AliasedExpr:
				switch expr := node.Expr.(type) {
				case *sqlparser.CaseExpr:
					for _, when := range expr.Whens {
						printWhen(when.Cond, currentQuery)
					}
				default:
					column := sqlparser.String(expr)
					column = strings.ToUpper(column)
					if column != "NULL" && column != "SYSDATE" && !strings.HasPrefix(column, "'") && !strings.HasPrefix(column, "V_") {
						_, err := strconv.Atoi(column)
						if err != nil {
							currentQuery.AddColumn(column)
						}
					}
				}

			default:
				currentQuery.AddColumn(sqlparser.String(node))
			}

		}
		if stmt.Where != nil {
			printWhen(stmt.Where.Expr, currentQuery)
		}
		queryArr = append(queryArr, currentQuery)
	case *sqlparser.Update:

		var currentQuery = viz.NewQuery(query)
		for _, node := range stmt.Exprs {
			switch node := node.Expr.(type) {
			case *sqlparser.Subquery:
				parseQuery(sqlparser.String(node.Select))
			default:
			}
		}
		printWhen(stmt.Where.Expr, currentQuery)
	case *sqlparser.Delete:
		var currentQuery = viz.NewQuery(query)
		printWhen(stmt.Where.Expr, currentQuery)
	case *sqlparser.Insert:
	}
}

func init() {
	rootCmd.AddCommand(sqlParseCmd)

	sqlParseCmd.Flags().StringP("source", "s", "", "source code directory")
	sqlParseCmd.Flags().StringP("filter", "f", "coll__graph.dot", "dot file filter")
	sqlParseCmd.Flags().StringP("output", "o", "dep.dot", "output dot file name")
}
