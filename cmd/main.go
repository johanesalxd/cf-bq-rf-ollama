package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	bqrfollama "github.com/johanesalxd/cf-bq-rf-ollama"
)

// main starts the function framework server on the specified port
func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	funcframework.RegisterHTTPFunction("/", bqrfollama.BQRFOllama)
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
