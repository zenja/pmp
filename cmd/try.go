package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	//"github.com/PuerkitoBio/goquery"
	"encoding/base64"
	"strconv"
)

func main() {
	proxyDBNet()
}

func usProxyOrg() {
	res, err := http.Get("https://www.us-proxy.org")
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	html := string(content)
	r := regexp.MustCompile(`<td>([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)</td><td>([0-9]+)</td>`)
	for _, match := range r.FindAllStringSubmatch(html, -1) {
		fmt.Printf("%s:%s\n", match[1], match[2])
	}
}

func gatherProxyCom() {
	res, err := http.Get("http://www.gatherproxy.com")
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	html := string(content)
	r := regexp.MustCompile(`"PROXY_IP":"([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)","PROXY_LAST_UPDATE":"[^"]+","PROXY_PORT":"([^"]+)"`)
	for _, match := range r.FindAllStringSubmatch(html, -1) {
		port, err := strconv.ParseInt(match[2], 16, 32)
		if err != nil {
			continue
		}
		fmt.Printf("%s:%d\n", match[1], port)
	}
}

func proxyListOrg() {
	res, err := http.Get("http://proxy-list.org/english/index.php?p=1")
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	html := string(content)
	r := regexp.MustCompile(`Proxy('[^']+')`)
	for _, match := range r.FindAllStringSubmatch(html, -1) {
		if err != nil {
			continue
		}
		decoded, err := base64.StdEncoding.DecodeString(match[1])
		if err != nil {
			continue
		}
		fmt.Printf("%s\n", string(decoded))
	}
}

func proxyDBNet() {
	res, err := http.Get("http://proxydb.net/")
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	html := string(content)
	r := regexp.MustCompile(`([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+):([0-9]+)`)
	for _, match := range r.FindAllStringSubmatch(html, -1) {
		if err != nil {
			continue
		}
		if err != nil {
			continue
		}
		fmt.Printf("%s:%d\n", match[1], match[2])
	}
}
