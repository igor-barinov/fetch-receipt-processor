/**
main.go

Entrypoint for the application. Constructs the HTTP server
*/

package main

import (
	"log"
	"net/http"

	"github.com/igor-barinov/fetch-receipt-processor/src/controller"
)

const (
	// Listen on port 3000
	ServerEndpoint = ":3000"
)

func main() {

	// Register endpoints for the server with a mux
	mux := http.NewServeMux()
	mux.Handle(controller.ProcessReceiptPath, http.HandlerFunc(controller.ProcessReceipt))
	mux.Handle(controller.GetPointsPath, http.HandlerFunc(controller.GetPoints))

	// Start the server
	err := http.ListenAndServe(ServerEndpoint, mux)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
