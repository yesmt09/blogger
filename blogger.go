package blogger

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const L_TRACE = 0
const L_DEBUG = 1
const L_INFO = 2
const L_WARNING = 3
const L_FATAL = 4

var levels = map[int]string{
	0: "TRACE",
	1: "DEBUG",
	2: "INFO",
	3: "WARNING",
	4: "FATAL",
}

var LevelMap = map[string]int{
	"trace":   0,
	"debug":   1,
	"info":    2,
	"warning": 3,
	"fatal":   4}

type BLogger struct {
	mu          sync.Mutex
	bFile 		BFile
	baseList    map[string]string
	baseLog		[]map[string]string
	logMessages []logModel
	logModel    logModel
}

func NewBlogger(bFile BFile) BLogger {
	return BLogger{
		bFile: bFile,
		baseList:    map[string]string{},
		baseLog: 	[]map[string]string{},
		logMessages: []logModel{},
		logModel:    logModel{
			rootPath: bFile.dir + "/",
		},
	}
}

func createLogFile(filepath string) {
	if _, err := os.Lstat(filepath); os.IsNotExist(err) {
		f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			panic("err")
		}
		defer f.Close()
	}
}

func (l *BLogger) RequestLogid() {
	l.AddBase("logid", strconv.Itoa(GetLogid()))
}

func (l *BLogger) AddBase(key string, value string) {
	if _,ok := l.baseList[key];ok {
		return
	}
	l.baseList[key] = value
	t := map[string]string{}
	t["key"] = key
	t["value"] = value
	l.baseLog = append(l.baseLog,t)
}

func (l *BLogger) Trace(message interface{}) {
	if l.bFile.level > L_TRACE {
		return
	}
	m := l.logModel.AddLog(fmt.Sprintf("%v", message), L_TRACE)
	l.logMessages = append(l.logMessages, m)

}

func (l *BLogger) Debug(message interface{}) {
	if l.bFile.level > L_DEBUG {
		return
	}
	m := l.logModel.AddLog(fmt.Sprintf("%v", message), L_DEBUG)
	l.logMessages = append(l.logMessages, m)
}

func (l *BLogger) Info(message interface{}) {
	if l.bFile.level > L_INFO {
		return
	}
	m := l.logModel.AddLog(fmt.Sprintf("%v", message), L_INFO)
	l.logMessages = append(l.logMessages, m)
}

func (l *BLogger) Warning(message interface{}) {
	if l.bFile.level > L_WARNING {
		return
	}
	m := l.logModel.AddLog(fmt.Sprintf("%v", message), L_WARNING)
	l.logMessages = append(l.logMessages, m)
}

func (l *BLogger) Fatal(message interface{}) {
	if l.bFile.level > L_FATAL {
		return
	}
	m := l.logModel.AddLog(fmt.Sprintf("%v", message), L_FATAL)
	l.logMessages = append(l.logMessages, m)
}

func (l BLogger) writeLog(content string, filepath ...string) {
	l.mu.Lock()
	var logpath = l.bFile.filePath
	if len(filepath) != 0 {
		logpath = filepath[0]
	}
	createLogFile(logpath)
	f, err := os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	defer l.mu.Unlock()
	if err != nil {
		panic(fmt.Sprintf("can't open ths log file %v", logpath))
	}
	f.WriteString(content)
}

func GetLogid() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(999999999999)
}

func (l *BLogger) Reset() {
	l.baseList = map[string]string{}
	l.baseLog = []map[string]string{}
	l.logMessages = []logModel{}
}

func (l *BLogger) Flush() {
	var content = ""
	var wfContent = ""
	for i := 0; i <= len(l.logMessages)-1; i++ {
		c := fmt.Sprintf("[%v][%v]", l.logMessages[i].timestamp, strings.ToUpper(levels[l.logMessages[i].level]))
		if len(l.baseLog) != 0 {
			for _, v := range l.baseLog {
				c = fmt.Sprintf("%v[%v:%v]", c, v["key"], v["value"])
			}
		}
		c = fmt.Sprintf("%v[%v:%v] %v\n", c, l.logMessages[i].funFilePath, l.logMessages[i].line, l.logMessages[i].message)
		content = content + c
		if l.logMessages[i].level > L_INFO {
			wfContent = wfContent + c
		}
	}
	l.writeLog(content)
	l.writeLog(wfContent, l.bFile.filePath+".wf")
	l.logMessages = []logModel{}
}
