package crawler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func (c *Crawler) Crawl(baseUrl string, levels int, cb adapters.Crawlback) error {
	fakeChrome := req.DefaultClient().ImpersonateChrome()

	crawler := colly.NewCollector(
		colly.MaxDepth(2),
		colly.UserAgent(fakeChrome.Headers.Get("user-agent")),
	)
	crawler.SetClient(&http.Client{
		Transport: fakeChrome.Transport,
	})

	crawler.OnHTML("section.section", func(e *colly.HTMLElement) {
		parts := e.ChildTexts("*")

		for _, part := range parts {
			part = strings.TrimSpace(part)
			docs, err := documentloaders.NewText(strings.NewReader(part)).LoadAndSplit(context.TODO(), textsplitter.NewRecursiveCharacter())
			if err != nil {
				fmt.Errorf("unable to load and split the text part of the document. %v", err)
				return
			}

			cb(part, docs, e)
		}
	})

	err := crawler.Visit(baseUrl)
	if err != nil {
		return fmt.Errorf("unable to visit the provided URL %s. %v", baseUrl, err)
	}

	return nil
}
