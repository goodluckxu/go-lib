package migrate

import (
	"strings"
)

type Line struct {
}

func (l *Line) SetLine(search string) {
	setNum := 0
	numList := l.searchFileLine(search)
	for _, lineNum := range numList {
		if setNum == 0 {
			if lineNum >= LastLineNum {
				setNum = lineNum
			}
		}
	}
	if setNum != 0 {
		LastLineNum = setNum
	}
}

func (l *Line) searchFileLine(search string) []int {
	rs := []int{}
	l.getLineList()
	for sKey, sLine := range strings.Split(search, "\n") {
		for _, line := range LineList {
			if strings.Index(line[1].(string), sLine) != -1 {
				if sKey == 0 {
					rs = append(rs, line[0].(int))
				}
			}
		}
	}
	return rs
}

func (l *Line) getLineList() {
	if len(LineList) == 0 {
		LineList = l.readFileList(Content)
	}
}

func (l *Line) readFileList(content string) (rs [][]interface{}) {
	for index, lineText := range strings.Split(content, "\n") {
		rs = append(rs, []interface{}{index + 1, lineText})
	}
	return
}
