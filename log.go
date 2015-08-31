package foolgo

import (
	//"fmt"
	"log"
	"os"
)

var (
	access_f *os.File
	error_f  *os.File
)

type Log struct {
	AccessLogger *log.Logger
	ErrorLogger  *log.Logger
	StandLogger  *log.Logger
}

func NewLog(access_log, error_log, standout_log string) (logger *Log) {
	logger = &Log{}

	if standout_log != "" {
		logFile, _ := os.OpenFile(standout_log, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		logger.StandLogger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	if access_log != "" {
		access_f, _ = os.OpenFile(access_log, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		logger.AccessLogger = log.New(access_f, "", 0)
	}

	if error_log != "" {
		error_f, _ = os.OpenFile(error_log, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		logger.ErrorLogger = log.New(error_f, "", log.Ldate|log.Ltime)
	}
	return logger
}

func (this *Log) AccessLog(msg string) {
	this.AccessLogger.Println(msg)
}

func (this *Log) ErrorLog(msg string) {
	//in fact,this error may never happend
	this.ErrorLogger.Println(msg)
}
