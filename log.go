package foolgo

import (
	"log"
	"os"
)

var (
	access_f *os.File
	error_f  *os.File
	run_f    *os.File
)

type Log struct {
	AccessLogger *log.Logger
	ErrorLogger  *log.Logger
	RunLogger    *log.Logger
}

func NewLog(access_log, error_log, run_log string) (logger *Log) { /*{{{*/
	logger = &Log{}

	if access_log != "" {
		access_f, _ = os.OpenFile(access_log, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		logger.AccessLogger = log.New(access_f, "", 0)
	}
	if error_log != "" {
		error_f, _ = os.OpenFile(error_log, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		logger.ErrorLogger = log.New(error_f, "", log.Ldate|log.Ltime)
	}
	if run_log != "" {
		run_f, _ = os.OpenFile(run_log, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		logger.RunLogger = log.New(run_f, "", log.Ldate|log.Ltime)
	}

	return logger
} /*}}}*/

// Request access log like nginx
func (this *Log) AccessLog(msg interface{}) { /*{{{*/
	if this.AccessLogger != nil {
		this.AccessLogger.Println(msg)
	}
} /*}}}*/

func (this *Log) ErrorLog(msg interface{}) { /*{{{*/
	//in fact,this error may never happend
	if this.ErrorLogger != nil {
		this.ErrorLogger.Println(msg)
	}
} /*}}}*/

// System run log
func (this *Log) RunLog(msg interface{}) { /*{{{*/
	if this.RunLogger != nil {
		this.RunLogger.Println(msg)
	} else {
		//log.Println(msg)
	}
} /*}}}*/
