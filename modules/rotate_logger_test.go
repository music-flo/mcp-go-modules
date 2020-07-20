package modules

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestRotateFileLogger(t *testing.T) {

	logPath := "./test"

	logger := NewRotateLogger(logPath , "1")

	log.Println("TEST")

	logger.Rotate()
	log.Println("TEST1")

	logger.Close()

	newLog , err := ioutil.ReadFile(path.Join(logPath , "1.log"))
	if err !=nil{
		t.Error(err)
	}

	if !strings.HasSuffix(string(newLog) , "TEST1\n") {
		t.Error("Log missmatch 1")
	}
	log.Println(string(newLog))

	oldLog , err2 := ioutil.ReadFile(path.Join(logPath , fmt.Sprintf("1_%s.log" ,time.Now().Format("2006-01-02"))))
	if err2 !=nil{
		t.Error(err2)
	}
	if !strings.HasSuffix(string(oldLog) , "TEST\n") {
		t.Error("Log missmatch 1")
	}

	log.Println(string(oldLog))

	err = os.RemoveAll(logPath)
	if err != nil {
		t.Error(err)
	}

}
