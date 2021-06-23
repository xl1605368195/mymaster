package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	PROD            = "prod"
	TEST            = "test"
	MASTER_VERSION  = "1.0.0"
	TimestampFormat = "2006-01-02 15:04:05"
)

type Writer struct {
	entry *logrus.Entry
}

func (w *Writer) FormatLog(logId int, message string, a ...interface{}) string {
	var codeInfo string
	var s string
	pc, filepath, line, ok := runtime.Caller(2)
	if ok {
		fPaths := strings.Split(filepath, string(os.PathSeparator))
		if len(fPaths) > 0 {
			filename := fPaths[len(fPaths)-1]
			fun := runtime.FuncForPC(pc)
			funNames := strings.Split(fun.Name(), ".")
			if len(funNames) > 0 {
				funName := funNames[len(funNames)-1]
				codeInfo = fmt.Sprintf("%s:%d:%s", filename, line, funName)
			}
		}
	}
	s = fmt.Sprintf("%s %s [logId=%d] jrasp-master 版本:%s", message, codeInfo, logId, MASTER_VERSION)
	return s
}

func NewLog(masterHome string, hostName string, moduleName string, isOnline bool, logOutputStd bool) *Writer {
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: TimestampFormat})
	outPut := os.Stdout
	if !logOutputStd {
		file, err := os.OpenFile(filepath.Join(masterHome, moduleName+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			outPut = file
		} else {
			logrus.Warning("Failed to log to file, using default stderr")
		}
	}

	logrus.SetOutput(outPut)
	var logLevel logrus.Level
	var env string
	if isOnline {
		logLevel = logrus.WarnLevel
		env = PROD
	} else {
		logLevel = logrus.DebugLevel
		env = TEST
	}
	if PROD == env {
		logLevel = logrus.WarnLevel
	}
	logrus.SetLevel(logLevel)
	// 日志代码打印行号 有bug
	// logrus.SetReportCaller(true)

	// TODO 增加es, redis 输出钩子

	// Entry实例
	entry := logrus.WithFields(logrus.Fields{
		// 全局字段
		"host_name":   hostName,
		"module_name": moduleName,
		"env":         env,
		"self_pid":    os.Getpid(),
	})

	w := &Writer{
		entry: entry,
	}
	return w
}

func (w *Writer) Info(log_id int, m string, a ...interface{}) {
	if a != nil {
		m = fmt.Sprintf(m, a...)
	}
	message := w.FormatLog(log_id, m, a)
	w.entry.Info(message)
}

func (w *Writer) Debug(log_id int, m string, a ...interface{}) {
	if a != nil {
		m = fmt.Sprintf(m, a...)
	}
	message := w.FormatLog(log_id, m, a)
	w.entry.Debug(message)
}

func (w *Writer) Notice(log_id int, m string, a ...interface{}) {
	if a != nil {
		m = fmt.Sprintf(m, a...)
	}
	message := w.FormatLog(log_id, m, a)
	w.entry.Debug(message)
}

func (w *Writer) Warning(log_id int, m string, a ...interface{}) {
	if a != nil {
		m = fmt.Sprintf(m, a...)
	}
	message := w.FormatLog(log_id, m, a)
	w.entry.Warn(message)
}

func (w *Writer) Err(log_id int, m string, a ...interface{}) {
	if a != nil {
		m = fmt.Sprintf(m, a...)
	}
	message := w.FormatLog(log_id, m, a)
	w.entry.Error(message)
}

func (w *Writer) Emerg(log_id int, m string, a ...interface{}) {
	if a != nil {
		m = fmt.Sprintf(m, a...)
	}
	message := w.FormatLog(log_id, m, a)
	w.entry.Info(message)
}

func (w *Writer) Alert(log_id int, m string, a ...interface{}) {
	if a != nil {
		m = fmt.Sprintf(m, a...)
	}
	message := w.FormatLog(log_id, m, a)
	w.entry.Warn(message)
}

func (w *Writer) Crit(log_id int, m string, a ...interface{}) {
	if a != nil {
		m = fmt.Sprintf(m, a...)
	}
	message := w.FormatLog(log_id, m, a)
	w.entry.Warn(message)
}
