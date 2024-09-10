package bqrfollama

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// init initializes the HTTP function handler for BQRFOllama
func init() {
	initOnce.Do(initAll)
	functions.HTTP("BQRFOllama", BQRFOllama)
}
