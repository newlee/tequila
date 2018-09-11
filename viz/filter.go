package viz

import (
	"os"
	"bufio"
	"strings"
)

type RegexpFilter struct {
	writeList []string
	blackList []string
}

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

func NewRegexpFilter() *RegexpFilter  {
	return &RegexpFilter{writeList:make([]string, 0), blackList:make([]string, 0)}
}

func (r *RegexpFilter) addReg(reg string)  {
	if strings.HasPrefix(reg, "- ") {
		r.blackList = append(r.blackList, reg[2:])
	}else {
		r.writeList = append(r.writeList, reg)
	}
}
func CreateRegexpFilter(fileName string) *RegexpFilter {
	rf := NewRegexpFilter()
	regs := readFilterFile(fileName)
	for _, reg := range regs {
		rf.addReg(reg)
	}
	return rf
}