package modules

import (
	"testing"
)

func TestGetFileCreateTime(t *testing.T) {

	tm := GetFileCreateTime("env_test.go")

	t.Log(tm)
}
