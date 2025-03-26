package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func WriteHttpResponseToFile(resp *http.Response, file string) error {
	// Open the specified file for writing. If the file already exists, it will be truncated.
	outFile, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("unable to create file %s. %v", file, err)
	}
	defer outFile.Close()

	// Copy the response body from the HTTP response to the output file.
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to write response body to file %s. %v", file, err)
	}

	// Return nil if no error occurred during the operation.
	return nil
}
