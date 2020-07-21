package modules

import "testing"

func TestReporterSlackProxy(t *testing.T) {

	r := NewReporterSlackProxy( "http://pri-heartqueen.music-flo.com/message/slack/text" , "team-mcp-tc-alarm")

	err := r.Report("AAAA")
	if err != nil{
		t.Error(err)
	}
}