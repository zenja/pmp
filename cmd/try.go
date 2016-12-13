package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Proxy struct {
	IP   string
	Port int
}

func main() {
	proxyDBNet()
}

func crawlProxies(url string, reg *regexp.Regexp, matchHandler func([]string) (string, string, error)) (chan Proxy, error) {
	proxies := make(chan Proxy)
	res, err := http.Get(url)
	if err != nil {
		close(proxies)
		return proxies, fmt.Errorf("failed to send Get request to %s: %s", url, err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		close(proxies)
		return proxies, fmt.Errorf("failed to read data from %s: %s", url, err)
	}
	html := string(content)
	go func() {
		for _, match := range reg.FindAllStringSubmatch(html, -1) {
			ip, portStr, err := matchHandler(match)
			if err != nil {
				log.Printf("Match handler returned an error: %s", err)
				continue
			}
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Printf("Failed to parse %d to a int port", portStr)
				continue
			}
			proxies <- Proxy{IP: ip, Port: port}
		}
		close(proxies)
	}()
	return proxies, nil
}

func multiCrawlProxies(urls []string, reg *regexp.Regexp, matchHandler func([]string) (string, string, error)) (chan Proxy, error) {
	proxies := make(chan Proxy)
	var wg sync.WaitGroup
	wg.Add(len(urls))
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(proxies)
	}(&wg)
	for _, url := range urls {
		ch, err := crawlProxies(url, reg, matchHandler)
		if err != nil {
			wg.Done()
			continue
		}
		go func(wg *sync.WaitGroup) {
			for p := range ch {
				proxies <- p
			}
			wg.Done()
		}(&wg)
	}
	return proxies, nil
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
	url := "http://www.gatherproxy.com"
	reg := regexp.MustCompile(`"PROXY_IP":"([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)","PROXY_LAST_UPDATE":"[^"]+","PROXY_PORT":"([^"]+)"`)
	matchHandler := func(match []string) (string, string, error) {
		port, err := strconv.ParseInt(match[2], 16, 32)
		if err != nil {
			return "", "", err
		}
		return match[1], strconv.Itoa(int(port)), nil
	}
	proxies, err := crawlProxies(url, reg, matchHandler)
	if err != nil {
		log.Fatal(err)
	}
	for proxy := range proxies {
		fmt.Printf("%s:%d\n", proxy.IP, proxy.Port)
	}
}

func proxyListOrg() {
	var urls []string
	for i := 1; i <= 10; i++ {
		urls = append(urls, fmt.Sprintf("http://proxy-list.org/english/index.php?p=%d", i))
	}
	reg := regexp.MustCompile(`Proxy\('([^']+)'\)`)
	matchHandler := func(match []string) (string, string, error) {
		decoded, err := base64.StdEncoding.DecodeString(match[1])
		if err != nil {
			return "", "", err
		}
		decodedStr := string(decoded)
		splits := strings.Split(decodedStr, ":")
		if len(splits) != 2 {
			log.Printf("format not correct: %s", decodedStr)
		}
		return splits[0], splits[1], nil
	}
	proxies, err := multiCrawlProxies(urls, reg, matchHandler)
	if err != nil {
		log.Fatal(err)
	}
	for proxy := range proxies {
		fmt.Printf("%s:%d\n", proxy.IP, proxy.Port)
	}
}

func proxyDBNet() {
	var urls []string
	for i := 0; i < 10; i++ {
		urls = append(urls, fmt.Sprintf("http://proxydb.net/?protocol=http&protocol=https&offset=%d", i*50))
	}
	reg := regexp.MustCompile(`([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+):([0-9]+)`)
	matchHandler := func(match []string) (string, string, error) {
		fmt.Printf("%s:%s\n", match[1], match[2])
		return match[1], match[2], nil
	}
	proxies, err := multiCrawlProxies(urls, reg, matchHandler)
	if err != nil {
		log.Fatal(err)
	}
	for proxy := range proxies {
		fmt.Printf("%s:%d\n", proxy.IP, proxy.Port)
	}
}
