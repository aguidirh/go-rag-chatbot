package crawler

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/util"
	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func (c *Crawler) Crawl(baseUrl string, levels int, cb adapters.Crawlback, allowedDomains []string) error {
	buffer := util.NewCircularBuffer(3)
	fakeChrome := req.DefaultClient().ImpersonateChrome()

	crawler := colly.NewCollector(
		colly.MaxDepth(levels),
		colly.AllowedDomains(allowedDomains...),
		colly.UserAgent(fakeChrome.Headers.Get("user-agent")),
	)
	crawler.SetClient(&http.Client{
		Transport: fakeChrome.Transport,
	})

	// On every a element which has href attribute call callback
	crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Visit link found on page
		e.Request.Visit(link)
	})
	var lastURL string
	crawler.OnHTML("section.section", func(e *colly.HTMLElement) {
		parts := e.ChildTexts("*")
		space := regexp.MustCompile(`\s+`)
		for _, part := range parts {
			if e.Request.URL.Path != lastURL {
				// try to keep relationships within a page
				c.log.Infof("checking path: %s", e.Request.URL)
				lastURL = e.Request.URL.Path
				buffer = util.NewCircularBuffer(3)
			}

			part = strings.TrimSpace(part)
			part = space.ReplaceAllString(part, " ")
			buffer.Add(part)
			doc, err := buffer.Join() // Join all parts in the buffer into a single string
			if err != nil {
				c.log.Errorf("unable to join parts in the buffer. %v", err)
				continue
			}
			docs, err := documentloaders.NewText(strings.NewReader(doc)).LoadAndSplit(context.TODO(), textsplitter.NewRecursiveCharacter())
			if err != nil {
				c.log.Errorf("unable to load and split the text part of the document. %v", err)
				return
			}

			cb(part, docs, e)
		}
	})

	err := crawler.Visit(baseUrl)
	if err != nil {
		return fmt.Errorf("unable to visit the provided URL %s. %v", baseUrl, err)
	}
	crawler.Wait()
	return nil
}
