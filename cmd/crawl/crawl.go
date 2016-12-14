package main

import (
	"fmt"
	"log"

	"github.com/zenja/pmp/crawler"
)

func main() {
	var crawlers []crawler.Crawler
	crawlers = append(crawlers, crawler.NewGatherProxyComCrawler())
	crawlers = append(crawlers, crawler.NewProxyDBNetCrawler())
	crawlers = append(crawlers, crawler.NewProxyListOrgCrawler())
	crawlers = append(crawlers, crawler.NewUSProxyOrgCrawler())
	for _, c := range crawlers {
		ch, err := c.CrawlProxies()
		if err != nil {
			log.Printf("Failed to crawl using proxy %v", c)
		}
		for p := range ch {
			fmt.Printf("%s:%d\n", p.IP, p.Port)
		}
	}
}
