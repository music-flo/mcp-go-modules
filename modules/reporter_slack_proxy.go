package modules

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type ReporterSlackProxy struct {
	proxyUrl string
	slackChannel string
}

func NewReporterSlackProxy(proxyUrl string , slackChannel string) Reporter {
	return &ReporterSlackProxy{
		proxyUrl: proxyUrl,
		slackChannel:slackChannel,
	}
}

func (r *ReporterSlackProxy)Report(msg string) error {
	resp , err := http.PostForm(r.proxyUrl, url.Values{"message" : {msg} , "channel" : {r.slackChannel}})
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	log.Println(string(respBody))
	return err
}
