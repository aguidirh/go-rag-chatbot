package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

// GetQueryParameterAsString retrieves the query parameter value for the specified key from the request URL.
//
// If the key is not present or has an empty value, the defaultValue is returned.
// This function parses the URL's query parameters using r.URL.Query().
func GetQueryParameterAsString(r *http.Request, key string, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryParameterAsInt retrieves the query parameter value for the specified key from the request URL and converts it to an integer.
//
// If the key is not present or has an empty value, the defaultValue is returned.
// If the conversion to an integer fails, an error is returned.
func GetQueryParameterAsInt(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// RequiredParameterAsString retrieves the query parameter value for the specified key from the request URL.
//
// If the key is not present or has an empty value, it returns an error with a message indicating that the required parameter was missing.
func RequiredParameterAsString(r *http.Request, param string) (string, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return "", errors.New("missing required parameter: " + param)
	}
	return value, nil
}
