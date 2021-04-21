package blogger

import (
	"runtime"
	"strings"
	"time"
)

type logModel struct {
	timestamp string
	message string
	funFilePath string
	line int
	level int
	logid int
	rootPath string
}

func (l logModel) AddLog(message string, level int) logModel {
	_, funFilePath, line, _ := runtime.Caller(2)
	l.message = message
	l.funFilePath = strings.Replace(funFilePath,l.rootPath,"",-1)
	l.line = line
	l.timestamp = time.Now().Format("2006-01-02 15:04:05")
	l.level = level
	if l.logid == 0 {
		l.logid = GetLogid()
	}
	return l
}