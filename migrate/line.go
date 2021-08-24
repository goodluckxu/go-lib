package migrate

import (
	"bufio"
	"os"
	"strings"
)

type Line struct {
	filePath string
	lineList [][]interface{}
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
		for _, line := range l.lineList {
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
	if l.filePath != "" {
		FilePath = l.filePath
	}
	f, _ := os.Open(FilePath)
	l.lineList = l.readFileList(f)
}

func (l *Line) readFileList(f *os.File) (rs [][]interface{}) {
	input := bufio.NewScanner(f)
	num := 1
	for input.Scan() {
		rs = append(rs, []interface{}{num, input.Text()})
		num++
	}
	return
}
