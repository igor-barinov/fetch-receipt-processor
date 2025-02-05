/**
receipts.go

Contains 'business' logic for processing/querying receipts
*/

package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/igor-barinov/fetch-receipt-processor/src/models"
)

// Define the paths for the HTTP server
const (
	ProcessReceiptPath = "/receipts/process"
	GetPointsPath      = "/receipts/{id}/points"
)

var idRgx = regexp.MustCompile(`^\S+$`)

// We use a sync.Map just in case mutliple clients start making requests
// this could probably be just a map[string]int64
var pointsStore sync.Map

var bonusMap sync.Map

// Validate a request to process a receipt, then calculate and store the points for the given receipt
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {

	// Unmarshal the request bytes
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read HTTP request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var receiptData models.Receipt
	err = json.Unmarshal(bytes, &receiptData)
	if err != nil {
		log.Printf("Failed to unmarshal HTTP request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The receipt is invalid."))
		return
	}

	// Validate the request
	err = receiptData.ValidateProperties()
	if err != nil {
		log.Printf("Reciept data was invalid: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The receipt is invalid."))
		return
	}

	// Calculate and store the points
	id := uuid.New().String()

	bonusPoints := int64(0)
	n := int64(0)
	timesProcessed, ok := bonusMap.Load(receiptData.UserID)
	if ok {
		n = timesProcessed.(int64)
	}

	if n < 3 {
		bonusPoints = 250
	}
	if n < 2 {
		bonusPoints = 500
	}
	if n < 1 {
		bonusPoints = 1000
	}

	nPoints := receiptData.CalculatePoints(bonusPoints)
	pointsStore.Store(id, nPoints)
	bonusMap.Store(receiptData.UserID, n+1)
	resp := &models.ProcessReceiptResponse{
		Id: id,
	}

	log.Printf("Created entry for receipt: (%v, %v)", id, nPoints)

	// Provide the ID as a response
	buf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal HTTP response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(buf)
}

// Validate a request to query the points for a given receipt ID, then return the result of the query
func GetPoints(w http.ResponseWriter, r *http.Request) {

	// Retrieve the 'id' path parameter
	receiptID := r.PathValue("id")
	if receiptID == "" {
		log.Printf("No ID was supplied")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No receipt found for that ID."))
		return
	}

	// Validate the ID parameter
	ok := idRgx.MatchString(receiptID)
	if !ok {
		log.Printf("ID didn't match pattern: %v", idRgx.String())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No receipt found for that ID."))
		return
	}

	// Attempt to retrieve the points for the given ID
	receiptPoints, ok := pointsStore.Load(receiptID)
	if !ok {
		log.Printf("ID does not exist")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No receipt found for that ID."))
		return
	}

	// Convert to int64
	n, ok := receiptPoints.(int64)
	if !ok {
		log.Printf("Points could not be converted to int64")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the points as the response
	resp := &models.GetPointsResponse{
		Points: n,
	}

	log.Printf("Retrieved ID '%v': %v points", receiptID, n)

	buf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal HTTP response body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(buf)
}
