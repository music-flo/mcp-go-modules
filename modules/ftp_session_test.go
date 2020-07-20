package modules

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestFtpSession(t *testing.T)  {
	ftp := NewFtpSession("211.110.226.57:9921" , "flotmp1" , "flotmp!@#")

	err := ftp.Connect()
	if err != nil {
		t.Error(err)
	}
	defer ftp.Close()

	testPath := fmt.Sprintf("/test/%x" , time.Now().Unix())
	testFile := "test.txt"

	testRetrFile := "test_2.txt"

	var f *os.File
	f , err = os.Create(testFile)
	if err != nil {
		t.Error(err)
	}

	_, err = f.WriteString("TEST")
	if err != nil {
		t.Error(err)
	}
	err = f.Close()

	if err != nil {
		t.Error(err)
	}

	err = ftp.MakeDirs2(testPath)
	if err != nil {
		t.Error(err)
	}
	upPath := path.Join(testPath , testFile)
	err = ftp.Store(upPath ,testFile)
	if err != nil {
		t.Error(err)
	}

	list , err := ftp.List(testPath)

	find := false
	for _, f := range list {
		fmt.Println(f)
		if strings.EqualFold(f , upPath) {
			find = true
		}
	}


	if !find {
		t.Error("cannot find store file")
	}

	var fi *FtpFileInfo
	fi , err = ftp.GetFileInfo(upPath)

	if err != nil {
		t.Error(err)
	}
	fmt.Println(fi.Name)

	err = ftp.Retr(upPath , testRetrFile)
	if err != nil {
		t.Error()
	}

	var data []byte
	data , err = ioutil.ReadFile(testRetrFile)
	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold("TEST" , string(data)){
		t.Error()
	}

	err = ftp.DeleteFile(upPath)
	if err != nil {
		t.Error(err)
	}
	err = os.Remove(testFile)
	if err != nil {
		t.Error(err)
	}
	err = os.Remove(testRetrFile)
	if err != nil {
		t.Error(err)
	}

	err = ftp.DeleteDirectoryEmpty(testPath)
	if err != nil {
		t.Error(err)
	}
}
