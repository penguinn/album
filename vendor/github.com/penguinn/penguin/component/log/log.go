package log

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/penguinn/penguin/utils"
)

var logFileString = `
<seelog>
    <outputs formatid="main">
        <filter levels="critical,error,warn,info,debug">
            <buffered size="10000" flushperiod="1000">
                <rollingfile type="size" filename="log/penguin.log" maxsize="10000000" maxrolls="30"/><!---向文件输出-->
                <!---
                <rollingfile type="date" filename="log/penguin.log" datepattern="2006.01.02" maxrolls="30"/>
                -->
            </buffered>
        </filter>
        <console />  <!---向屏幕输出-->
    </outputs>
    <formats>
        <format id="main" format="%Date/%Time [%LEV] %Func:%Line %Msg%n"/>
    </formats>
</seelog>
`

type LogConfig struct {
	File string
}

type LogComponent struct{}

func (LogComponent) Init(ops ...interface{}) (err error) {
	var (
		logger log.LoggerInterface
		isExit bool
	)
	if len(ops) == 0 {
		logger, err = log.LoggerFromConfigAsString(logFileString)
	} else {
		if utils.FileExist(ops[0].(*LogConfig).File) {
			isExit = true
		} else {
			isExit = false
		}
		if isExit {
			logger, err = log.LoggerFromConfigAsFile(ops[0].(*LogConfig).File)
		} else {
			logger, err = log.LoggerFromConfigAsString(logFileString)
		}
		if err != nil {
			fmt.Println("err parsing config log file", err)
			panic(err)
		}
	}
	log.ReplaceLogger(logger)
	return nil
}

func Trace(args ...interface{}) {
	log.Trace(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Critical(args ...interface{}) {
	log.Critical(args...)
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Criticalf(format string, args ...interface{}) {
	log.Criticalf(format, args...)
}
