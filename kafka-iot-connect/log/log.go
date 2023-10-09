package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

const (
	LogNone = iota
	LogInfo
	LogWarning
	LogError
	LogVerbose
	LogDebug
)

var LogMap = map[int]string{
	LogInfo:    "LogInfo   ",
	LogWarning: "LogWarning",
	LogError:   "LogError  ",
	LogVerbose: "LogVerbose",
	LogDebug:   "LogDebug  ",
}

type LoggerSystem struct {
	Logs []*MyFileLogger
	Name string
	Id   int
	mu   *sync.RWMutex
}

func RemoveElementFromSliceNew[V any](slice []V, s int) ([]V, error) {
	if s < 0 || s >= len(slice) {
		return slice, fmt.Errorf("index: %d is invalid", s)
	} else {
		return append(slice[:s], slice[s+1:]...), nil
	}
}

func InitializeLoggerSystemWithLock(name string, id int) *LoggerSystem {
	return &LoggerSystem{
		Name: name,
		Id:   id,
		Logs: make([]*MyFileLogger, 0),
		mu:   &sync.RWMutex{},
	}
}

func (ls *LoggerSystem) FindLog(name string) int {
	for ind, ele := range ls.Logs {
		if ele.logFile.Name() == name {
			return ind
		}
	}
	return -1
}

func (ls *LoggerSystem) GetLogByNameMust(name string) *MyFileLogger {
	for _, ele := range ls.Logs {
		if ele.Name == name {
			return ele
		}
	}
	panic("failed to get log by name: " + name)
}

func (ls *LoggerSystem) GetLogByName(name string) *MyFileLogger {
	for _, ele := range ls.Logs {
		if ele.Name == name {
			return ele
		}
	}
	return nil
}

func (ls *LoggerSystem) CreateLog(l *MyFileLogger) (int, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	name := l.logFile.Name()
	if ls.FindLog(name) < 1 {
		ls.Logs = append(ls.Logs, l)
		return len(ls.Logs) - 1, nil
	} else {
		return -1, fmt.Errorf(" another file with fileName: %s exists", name)
	}
}

func (ls *LoggerSystem) FindLogByName(name string) int {
	for ind, ele := range ls.Logs {
		if ele.Name == name {
			return ind
		}
	}
	return -1
}

func (ls *LoggerSystem) CreateLogsByName(lss ...*MyFileLogger) ([]string, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	var names = []string{}
	for _, l := range lss {
		if ls.FindLogByName(l.Name) < 1 {
			ls.Logs = append(ls.Logs, l)
			names = append(names, l.Name)
		} else {
			return names, fmt.Errorf(" another file with fileName: %s exists", l.Name)
		}
	}
	return names, nil
}

func (ls *LoggerSystem) CreateLogs(lss ...*MyFileLogger) ([]string, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	var names = []string{}
	for _, l := range lss {
		name := l.logFile.Name()
		if ls.FindLog(name) < 1 {
			ls.Logs = append(ls.Logs, l)
			names = append(names, l.logFile.Name())
		} else {
			return names, fmt.Errorf(" another file with fileName: %s exists", name)
		}
	}
	return names, nil
}

func (ls *LoggerSystem) ClearAllLogs() error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	for _, ele := range ls.Logs {
		if err := ele.ClearLog(); err != nil {
			return err
		}
	}
	return nil
}

func (ls *LoggerSystem) RemoveLog(l *MyFileLogger) (err error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	name := l.logFile.Name()
	if ind := ls.FindLog(name); ind < 0 {
		err = fmt.Errorf(" failed to find log by fileName: %s", name)
	} else {
		ls.Logs, err = RemoveElementFromSliceNew(ls.Logs, ind)
	}
	return
}

func (ls *LoggerSystem) RemoveLogByName(name string) (err error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if ind := ls.FindLogByName(name); ind < 0 {
		err = fmt.Errorf(" failed to find log by fileName: %s", name)
	} else {
		ls.Logs, err = RemoveElementFromSliceNew(ls.Logs, ind)
	}
	return
}

type MyFileLogger struct {
	Name           string
	logger         *log.Logger
	logFile        *os.File
	logLevel       int
	logFileMaxSize int64
}

func NewFileLogger() *MyFileLogger {
	return &MyFileLogger{
		logger:   nil,
		logFile:  nil,
		logLevel: LogNone,
	}
}

func StartLogger(file string, maxBytes int64) *MyFileLogger {
	res := &MyFileLogger{
		logger:         nil,
		logFile:        nil,
		logLevel:       LogNone,
		logFileMaxSize: maxBytes,
	}
	if err := res.StartLog(LogNone, file); err != nil {
		fmt.Println(err.Error())
		// panic(err.Error())
	}
	return res
}

func StartLoggerByName(name, file string, maxBytes int64, logLevel int) *MyFileLogger {
	res := &MyFileLogger{
		Name:           name,
		logger:         nil,
		logFile:        nil,
		logLevel:       LogNone,
		logFileMaxSize: maxBytes,
	}
	if err := res.StartLog(logLevel, file); err != nil {
		log.Println(err.Error())
	}
	return res
}

func (myLogger *MyFileLogger) StartLog(level int, file string) error {
	if len(file) == 0 {
		return nil
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	myLogger.logger = log.New(f, "", 0)
	myLogger.logLevel = level
	myLogger.logFile = f
	return nil
}
func (myLogger *MyFileLogger) StopLog() error {
	if myLogger.logFile != nil {
		return myLogger.logFile.Close()
	}
	return nil
}

func (myLogger *MyFileLogger) ClearLog() error {
	if myLogger.logFile != nil {
		if err := myLogger.logFile.Truncate(0); err != nil {
			return err
		}
	}
	return nil
}

// You can add a log of auxiliary functions here to make the log more easier
func (myLogger *MyFileLogger) Log(level int, msg string, a ...interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	msg = fmt.Sprintf("%s | %s  | %s %s\n", time.Now().Format("2006-01-02 15:04:05.000 "), LogMap[level], msg, fmt.Sprint(a))
	if myLogger.logFile == nil {
		if level >= myLogger.logLevel {
			fmt.Print(msg)
		}
		return nil
	}
	if myLogger.logger == nil {
		return errors.New("MyFileLogger is not initialized correctly")
	}
	if level >= myLogger.logLevel {
		file, err := myLogger.logFile.Stat()
		if err != nil {
			fmt.Println(err.Error())
			myLogger.logger.Print(err.Error())
		}
		if file == nil {
			return nil
		} else if file.Size() < myLogger.logFileMaxSize {
			myLogger.logger.Print(msg)
		} else {
			myLogger.logger.Print(msg)
			myLogger.logger.Printf("Log File Exceeds max size: %v bytes, close log file", myLogger.logFileMaxSize)
			myLogger.logFile.Close()
			myLogger.logFile = nil
		}
	}
	return nil
}

func (myLogger *MyFileLogger) LogInfo(pkg any, a ...interface{}) error {
	return myLogger.Log(LogInfo, reflect.TypeOf(pkg).PkgPath(), a...)
}
