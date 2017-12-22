package viz

import (
	"strings"
)

var MergeHeaderFunc = func(input string) string {
	tmp := strings.Split(input, ".")
	if len(tmp) > 1 {
		return strings.Join(tmp[0:len(tmp)-1], ".")
	}
	return input
}

var MergePackageFunc = func(input string) string {
	tmp := strings.Split(input, "/")
	packageName := tmp[0]
	if packageName == input {
		packageName = "main"
	}
	if len(tmp) > 2 {
		packageName = strings.Join(tmp[0:len(tmp)-1], "/")
	}

	return packageName
}
