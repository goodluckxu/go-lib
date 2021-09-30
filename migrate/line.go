package migrate

import (
	"io/ioutil"
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
	defer f.Close()
	l.lineList = l.readFileList(f)
}

func (l *Line) readFileList(f *os.File) (rs [][]interface{}) {
	by, _ := ioutil.ReadAll(f)
	content := string(by)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	for index, lineText := range strings.Split(content, "\n") {
		rs = append(rs, []interface{}{index + 1, lineText})
	}
	return
}
