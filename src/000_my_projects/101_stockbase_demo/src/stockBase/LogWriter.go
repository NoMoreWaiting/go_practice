package stockBase

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var logWriter      LogWriter

const (
	DEBUG = iota
	INFO
	ERROR
	WARN
)

type LogInfo struct {
	current time.Time
	level   int
	info    string
}

type LogWriter struct {
	outchan  chan LogInfo
	lasttime time.Time
	logfile  *os.File
	logger   *log.Logger
}

func (l *LogWriter) Info(format string, v ...interface{}) {
	var logInfo LogInfo
	logInfo.level = INFO
	logInfo.current = time.Now()
	logInfo.info = fmt.Sprintf(format, v...)
	l.outchan <- logInfo
}

func (_self *LogWriter) createFolder(current time.Time) {
	if current.Day() != _self.lasttime.Day() || current.Month() != _self.lasttime.Month() || current.Year() != _self.lasttime.Year() {
		os.MkdirAll(SysConfig.Basic.Logdir, 0777)
                
		_self.logfile.Close()
		fileName := fmt.Sprintf("stockbase.run.%04d-%02d-%02d.log", current.Year(), current.Month(), current.Day())
		var err error
		_self.logfile, err = os.OpenFile(SysConfig.Basic.Logdir+"/"+fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
                
		if err != nil {
			fmt.Printf("%s\r\n", err.Error())
			return
		}
		_self.lasttime = current
		_self.logger = log.New(_self.logfile, "[PID:"+strconv.Itoa(os.Getpid())+"] ", log.Ldate|log.Ltime|log.Lshortfile)
		fmt.Println("******************" + strconv.Itoa(os.Getpid()) + "********************")
	}
}

func (_self *LogWriter) WriteInfo() {
	_self.outchan = make(chan LogInfo)
	// 创建协程
	go func() {
		for logInfo := range _self.outchan {
			//向通道内写入数据，如果无人读取会等待
			_self.createFolder(logInfo.current)
			_self.logger.Printf(logInfo.info)
		}
	}()
}
