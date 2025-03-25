package crawler

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

func (c *Crawler) GetDirectDescendants(baseUrl string, levels int) (map[string]any, error) {
	descendants := make(map[string]any)

	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse the provided URL %s. %v", baseUrl, err)
	}

	basePath := url.Path

	c.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		targetUrl := e.Attr("href")
		if !strings.Contains(targetUrl, basePath) {
			return
		}

		// to-do: we might want to provide additional context in the map. for now, we'll just
		// use this to ensure we have no duplicate URLs
		descendants[targetUrl] = targetUrl
	})

	defer c.collector.OnHTMLDetach("a[href]")

	err = c.collector.Visit(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to visit the provided URL %s. %v", baseUrl, err)
	}

	if levels > 1 {
		for url := range descendants {
			innerDescendants, err := c.GetDirectDescendants(url, levels-1)
			if err != nil {
				return nil, fmt.Errorf("unable to get direct descendants %s. %v", baseUrl, err)
			}
			for innerUrl, _ := range innerDescendants {
				descendants[innerUrl] = innerUrl
			}
		}
	}
	return descendants, nil
}
