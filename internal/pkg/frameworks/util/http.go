package util

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

type HttpAccessor struct {
}

// GetAsDocuments retrieves a list of documents from the specified URL and returns them as schema.Document slices.
func (h *HttpAccessor) GetAsDocuments(ctx context.Context, url string) ([]schema.Document, error) {
	// Attempt to get the response from the provided URL.
	resp, err := h.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve document(s) from %s. %v", url, err)
	}

	defer resp.Body.Close() // Ensure the response body is closed after use.

	// Load and split the text content of the response into documents using Markdown splitting.
	newDocs, err := documentloaders.NewText(resp.Body).LoadAndSplit(ctx, textsplitter.NewRecursiveCharacter())
	if err != nil {
		return nil, fmt.Errorf("unable to load document(s) from %s. %v", url, err)
	}

	return newDocs, nil
}

// Get retrieves an HTTP response for the specified resource.
func (h *HttpAccessor) Get(resource string) (*http.Response, error) {
	// Attempt to get the response from the provided URL.
	resp, err := http.Get(resource)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve document(s) from %s. %v", resource, err)
	}

	return resp, nil
}

// NewHttpAccessor creates a new instance of HttpAccessor.
func NewHttpAccessor() *HttpAccessor {
	return &HttpAccessor{}
}
