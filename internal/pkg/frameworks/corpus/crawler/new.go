package crawler

import (
	"os"
	"path"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/util"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type Crawler struct {
	collector       *colly.Collector
	cachePath       string
	collectionsPath string
	http            *util.HttpAccessor
	log             *logrus.Logger
}

/*
Crawler is a service that crawls a website and extracts text content. It uses the colly library to perform web scraping.
It takes an HTML document source data (htmlDoc) and a callback function (cb) as input.
The htmlDoc contains metadata about the URL and allowed/allowed-regexs/skip/skip-regexs filters, while cb is called for each extracted text content.
*/
func New(log *logrus.Logger) *Crawler {
	cachePath := ".cache"
	collectionsPath := "collections"
	_ = os.MkdirAll(path.Join(cachePath, collectionsPath), 0777)

	collector := colly.NewCollector(colly.CacheDir(cachePath))
	collector.AllowURLRevisit = true

	return &Crawler{collector: collector, cachePath: cachePath, collectionsPath: collectionsPath, http: util.NewHttpAccessor(), log: log}
}
