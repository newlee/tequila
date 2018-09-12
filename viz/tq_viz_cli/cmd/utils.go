package cmd

import (
	"bufio"
	"github.com/newlee/tequila/viz"
	"os"
	"strings"
)

func doFiles(fileNames []string, fileCallback func(), callback func(string, string)) {
	for _, fileName := range fileNames {
		file, _ := os.Open(fileName)
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		fileCallback()
		for scanner.Scan() {
			line := scanner.Text()
			callback(line, fileName)
		}

		file.Close()
	}
}

var emptyFilter = func(line string) bool {
	return true
}

var pkgSplit = func(r rune) bool {
	return r == ' ' || r == '(' || r == ',' || r == '\'' || r == '"' || r == ')'
}
var tableSplit = func(r rune) bool {
	return r == ' ' || r == ',' || r == '.' || r == '"' || r == ':' || r == '(' || r == ')' || r == '）' || r == '%' || r == '!' || r == '\''
}

func isComment(first string) bool {
	return strings.HasPrefix(first, "/*") || strings.HasPrefix(first, "*") || strings.HasPrefix(first, "//")
}

func doSplit(line string, f func(rune) bool, callback func(string)) {
	tmp := strings.FieldsFunc(line, f)
	for _, s := range tmp {
		callback(s)
	}
}

func doPkg(s string, pkgFilter func(line string) bool, callback func(string, string)) {
	pkg := ""
	sp := ""
	if strings.HasPrefix(s, "PKG_") && pkgFilter(strings.Split(s, ".")[0]) {
		s = strings.Replace(s, "\"", "", -1)
		tmp := strings.Split(s, ".")
		pkg = tmp[0]
		if len(tmp) > 1 && strings.HasPrefix(tmp[1], "P_") {
			sp = tmp[1]
			callback(pkg, sp)
		}
	}
}

func doPkgLine(line string, pkgFilter func(line string) bool, callback func(string, string)) {
	doSplit(line, pkgSplit, func(s string) {
		doPkg(s, pkgFilter, callback)
	})
}

func doTableLine(line string, tableFilter func(line string) bool, callback func(string)) {
	doSplit(line, tableSplit, func(s string) {
		if strings.HasPrefix(s, "T_") && tableFilter(s) && !viz.IsChineseChar(s) {
			callback(s)
		}
	})
}

func doCreatePkg(line string, callback func(string)) {
	tmp := strings.FieldsFunc(line, pkgSplit)
	for _, key := range tmp {
		if strings.HasPrefix(key, "PKG_") {
			callback(key)
		}
	}
}
func doCreateProcedure(line string, callback func(string)) {
	tmp := strings.FieldsFunc(line, pkgSplit)
	for _, key := range tmp {
		if strings.HasPrefix(key, "P_") {
			callback(key)
		}
	}
}
