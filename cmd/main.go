package main

import (
	"log"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/httpserver"
)

func main() {
	err := httpserver.Run()

	if err != nil {
		log.Fatal(err)
	}
}
