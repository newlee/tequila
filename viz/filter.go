package viz

import (
	"bufio"
	"github.com/dlclark/regexp2"
	"os"
	"strings"
)

type RegexpFilter struct {
	whiteList []string
	blackList []string
	excludes  []string
}

func readFilterFile(fileName string) []string {
	result := make([]string, 0)
	f, _ := os.Open(fileName)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			result = append(result, strings.Trim(line, " "))
		}
	}
	return result
}

func NewRegexpFilter() *RegexpFilter {
	return &RegexpFilter{
		whiteList: make([]string, 0),
		blackList: make([]string, 0),
		excludes:  make([]string, 0)}
}

func (r *RegexpFilter) AddExclude(exclude string) {
	r.excludes = append(r.excludes, exclude)
}

func (r *RegexpFilter) AddReg(reg string) {
	if strings.HasPrefix(reg, "- ") {
		r.blackList = append(r.blackList, reg[2:])
	} else {
		r.whiteList = append(r.whiteList, reg)
	}
}

func (r *RegexpFilter) notExclude(s string) bool {
	for _, ct := range r.excludes {
		if ct == s {
			return false
		}
	}
	return true
}

func (r *RegexpFilter) Match(s string) bool {
	return r.notExclude(s) && r.matchWhiteList(s) && !r.matchBlackList(s)
}

func (r *RegexpFilter) NotMatch(s string) bool {
	return r.notExclude(s) && !r.Match(s)
}

func (r *RegexpFilter) UnMatch(s string) bool {
	return r.notExclude(s) && !r.matchWhiteList(s) && !r.matchBlackList(s)
}

func (r *RegexpFilter) matchWhiteList(s string) bool {
	for _, reg := range r.whiteList {
		re, _ := regexp2.Compile(reg, 0)
		if isMatch, _ := re.MatchString(s); isMatch {
			return true
		}
	}
	return false
}

func (r *RegexpFilter) matchBlackList(s string) bool {
	for _, reg := range r.blackList {
		re, _ := regexp2.Compile(reg, 0)
		if isMatch, _ := re.MatchString(s); isMatch {
			return true
		}
	}
	return false
}

func CreateRegexpFilter(fileName string) *RegexpFilter {
	rf := NewRegexpFilter()
	regs := readFilterFile(fileName)
	for _, reg := range regs {
		rf.AddReg(reg)
	}
	return rf
}

func (r *RegexpFilter) AddExcludes(fileName string) *RegexpFilter {
	r.excludes = readFilterFile(fileName)
	return r
}

type PrefixFilter struct {
	whiteList []string
}


func CreatePrefixFilter(fileName string) *PrefixFilter {
	pf := &PrefixFilter { whiteList: make([]string, 0)}
	lines := readFilterFile(fileName)
	for _, line := range lines {
		pf.whiteList = append(pf.whiteList, line)
	}
	return pf
}

func (p *PrefixFilter) Match(s string) bool {
	for _, prefix := range p.whiteList {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
