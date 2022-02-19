package tool

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	Ghttp = NewGHttp()
)

type ghttp struct {
}

func NewGHttp() *ghttp {
	return &ghttp{}
}

func (t *ghttp) Get(url string, header map[string]string) (map[string]string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respMap := map[string]string{}
	respMap["response"] = string((body))
	respMap["code"] = strconv.Itoa(resp.StatusCode)
	return respMap, nil
}

func (t *ghttp) Post(url string, param string, header map[string]string) (map[string]string, error) {
	respMap := map[string]string{}
	reader := bytes.NewBufferString(param)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return respMap, err
	}
	defer request.Body.Close()
	for k, v := range header {
		request.Header.Set(k, v)
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return respMap, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respMap, err
	}
	defer resp.Body.Close()
	respMap["cookie"] = resp.Header.Get("Set-Cookie")
	respMap["response"] = string((respBytes))
	respMap["code"] = strconv.Itoa(resp.StatusCode)
	return respMap, nil
}
