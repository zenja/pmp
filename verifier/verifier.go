package verifier

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Verifier interface {
	VerifyProxy(ip string, port int) (time.Duration, error)
}

type SimpleVerifier struct {
	testURL     string
	connTimeout time.Duration
	mustContain string
}

func New() Verifier {
	return &SimpleVerifier{
		testURL:     "http://www.baidu.com/",
		connTimeout: 5 * time.Second,
		mustContain: "hao123",
	}
}

func NewSimpleVerifier(testURL string, connTimeout time.Duration, mustContain string) Verifier {
	return &SimpleVerifier{
		testURL:     testURL,
		connTimeout: connTimeout,
		mustContain: mustContain,
	}
}

func (sv *SimpleVerifier) VerifyProxy(ip string, port int) (time.Duration, error) {
	proxyHost := fmt.Sprintf("%s:%d", ip, port)
	proxyURL := &url.URL{Host: proxyHost}
	startTS := time.Now()
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		Timeout:   sv.connTimeout}
	resp, err := client.Get(sv.testURL)
	if err != nil {
		return time.Duration(0), fmt.Errorf("failed to connect to proxy %s: %s", proxyHost, err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	endTS := time.Now()
	if !bytes.Contains(body, []byte(sv.mustContain)) {
		return time.Duration(0), fmt.Errorf("failed to fetch correct content with proxy %s: %s", proxyHost, err)
	}
	return endTS.Sub(startTS), nil
}
