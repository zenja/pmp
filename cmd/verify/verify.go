package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/zenja/pmp"
	"github.com/zenja/pmp/crawler"
	"github.com/zenja/pmp/verifier"
)

func main() {
	var crawlers []crawler.Crawler
	crawlers = append(crawlers, crawler.NewGatherProxyComCrawler())
	crawlers = append(crawlers, crawler.NewProxyDBNetCrawler())
	crawlers = append(crawlers, crawler.NewProxyListOrgCrawler())
	crawlers = append(crawlers, crawler.NewUSProxyOrgCrawler())
	var proxies []pmp.Proxy
	fmt.Println("Crawling proxies...")
	for _, c := range crawlers {
		ch, err := c.CrawlProxies()
		if err != nil {
			log.Printf("Failed to crawl using proxy %v", c)
		}
		for p := range ch {
			proxies = append(proxies, p)
		}
	}
	fmt.Printf("%d proxies crawled.\n", len(proxies))

	fmt.Println("Testing proxies...")
	vyfr := verifier.NewSimpleVerifier("http://www.bing.com/", 10*time.Second, "bing")
	sem := make(chan int, 50)
	var wg sync.WaitGroup
	wg.Add(len(proxies))
	for _, p := range proxies {
		go func(ip string, port int, wg *sync.WaitGroup) {
			sem <- 1
			defer func() {
				wg.Done()
				<-sem
			}()
			timeSpent, err := vyfr.VerifyProxy(ip, port)
			if err != nil {
				//fmt.Printf("%s:%d test failed: %s\n", ip, port, err)
			} else {
				fmt.Printf("%s:%d test succeeded, time spent: %s\n", ip, port, timeSpent)
			}
		}(p.IP, p.Port, &wg)
	}
	wg.Wait()
}
