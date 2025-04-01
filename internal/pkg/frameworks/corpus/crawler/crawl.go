package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/util"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func (c *Crawler) Crawl(htmlDoc data.DocSourceHttp, cb adapters.Crawlback) error {
	buffer := util.NewCircularBuffer(3)
	fakeChrome := req.DefaultClient().ImpersonateChrome()
	space := regexp.MustCompile(`\s+`)

	if len(htmlDoc.UrlFilter.Allowed) > 0 && len(htmlDoc.UrlFilter.AllowedRegexs) == 0 {
		for _, pattern := range htmlDoc.UrlFilter.Allowed {
			htmlDoc.UrlFilter.AllowedRegexs = append(htmlDoc.UrlFilter.AllowedRegexs, regexp.MustCompile(pattern))
		}
	}

	if len(htmlDoc.UrlFilter.Skip) > 0 && len(htmlDoc.UrlFilter.SkipRegexs) == 0 {
		for _, pattern := range htmlDoc.UrlFilter.Skip {
			htmlDoc.UrlFilter.SkipRegexs = append(htmlDoc.UrlFilter.SkipRegexs, regexp.MustCompile(pattern))
		}
	}

	crawler := colly.NewCollector(
		colly.MaxDepth(htmlDoc.RecursionLevels),
		colly.AllowedDomains(htmlDoc.AllowedDomains...),
		colly.UserAgent(fakeChrome.Headers.Get("user-agent")),
	)
	crawler.SetClient(&http.Client{
		Transport: fakeChrome.Transport,
	})

	// On every a element which has href attribute call callback
	crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		linkUrl, err := url.Parse(link)
		if err != nil {
			c.log.Errorf("failed to parse URL: %v", err)
			return
		}
		if len(htmlDoc.AllowedDomains) > 0 {
			allowed := false
			for _, allowedDomain := range htmlDoc.AllowedDomains {
				if linkUrl.Host == allowedDomain {
					allowed = true
					break
				}
			}

			if !allowed {
				c.log.Infof("skipping link: %s, not in allowed domains", link)
				return
			}
		}

		if strings.Index(link, "/") != 0 && strings.Index(link, "http") != 0 {
			c.log.Infof("skipping link: %s, not a valid URL", link)
			return
		}

		for _, re := range htmlDoc.UrlFilter.SkipRegexs {
			if re.MatchString(link) {
				c.log.Infof("skipping link: %s, skip URL pattern matched", link)
				return
			}
		}

		var allowed bool
		for _, re := range htmlDoc.UrlFilter.AllowedRegexs {
			if re.MatchString(link) {
				c.log.Infof("allowing link: %s", link)
				allowed = true
				break
			}
			if !allowed {
				c.log.Infof("skipping link: %s, no allowed URL patterns matched", link)
				return
			}
		}
		c.log.Infof("visiting link: %s", link)
		// Visit link found on page
		e.Request.Visit(link)
	})
	var lastURL string
	crawler.OnHTML("section.section", func(e *colly.HTMLElement) {
		parts := e.ChildTexts("*")
		for _, part := range parts {
			if e.Request.URL.Path != lastURL {
				// try to keep relationships within a page
				c.log.Infof("checking path: %s", e.Request.URL)
				lastURL = e.Request.URL.Path
				buffer = util.NewCircularBuffer(3)
			}
			part = strings.TrimSpace(part)
			part = space.ReplaceAllString(part, " ")
			if len(part) < 20 {
				continue
			}
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

	err := crawler.Visit(htmlDoc.URL)
	if err != nil {
		return fmt.Errorf("unable to visit the provided URL %s. %v", htmlDoc.URL, err)
	}
	crawler.Wait()
	return nil
}
