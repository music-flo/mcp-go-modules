package modules

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type RotateFileLogger struct {
	logPath 		string
	id 				string
	fpCreateTime 	time.Time
	fpLog 			*os.File
}

func NewRotateLogger(logPath string , id string) *RotateFileLogger{
	logger := &RotateFileLogger{
		logPath:logPath,
		id:id,
		fpLog:nil,
	}
	logger.createLogFile()

	log.SetOutput(logger)

	return logger
}

func (r *RotateFileLogger) Write(p []byte) (n int, err error) {

	r.rotateLog()

	if r.fpLog != nil {
		return r.fpLog.Write(p)
	}
	return 0, io.ErrClosedPipe
}

func (r *RotateFileLogger) Close() {
	log.SetOutput(os.Stdout)

	if r.fpLog != nil {
		r.fpLog.Close()
		r.fpLog = nil
	}
}


func (r *RotateFileLogger)rotateLog() {
	now := time.Now()

	if r.fpCreateTime.Day() != now.Day() ||
		r.fpCreateTime.Month() != now.Month() ||
		r.fpCreateTime.Year() != now.Year() {

		r.Rotate()
	}
}

func (r *RotateFileLogger)Rotate() {
	logfilePath := path.Join(r.logPath, r.id +".log")

	logfileBackupPath := path.Join(r.logPath, fmt.Sprintf("%s_%s.log" ,r.id, r.fpCreateTime.Format("2006-01-02")))
	os.Rename(logfilePath , logfileBackupPath)

	r.createLogFile()
}

func (r *RotateFileLogger)createLogFile(){
	var err error
	typePath := pathTypeCheck(r.logPath)
	if typePath != PtDirectory {
		os.MkdirAll(r.logPath , os.ModePerm)
	}

	logfilePath := path.Join(r.logPath, r.id +".log")
	if r.fpLog != nil {
		r.fpLog.Close()
		r.fpLog = nil
	}

	r.fpLog, err = os.OpenFile(logfilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		r.fpCreateTime = GetFileCreateTime(logfilePath)
	}

	if err != nil {
		r.Close()
	}
}

type pathType int
const (
	PtUnknown pathType = iota
	PtFile
	PtDirectory
)

func pathTypeCheck( path  string ) pathType{
	stat , err := os.Stat(path)

	if err != nil {
		return PtUnknown
	}
	if stat.IsDir() {
		return PtDirectory
	}
	return PtFile
}

