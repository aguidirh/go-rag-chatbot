package crawler

import "github.com/gocolly/colly"

type Crawler struct {
	collector *colly.Collector
}

func New() *Crawler {
	collector := colly.NewCollector()
	return &Crawler{collector: collector}
}
