/**
e2e_test.go

Makes HTTP calls to the app endpoints to test various scenarios
Assumes that the server is running on 'http://localhost:3000'
*/

package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/igor-barinov/fetch-receipt-processor/src/controller"
	"github.com/igor-barinov/fetch-receipt-processor/src/models"
	"github.com/stretchr/testify/assert"
)

const ServerEndpoint = "http://localhost:3000"

// Will hold ID of processed receipts
var Receipt1ID string
var Receipt2ID string

func TestProcessReceiptWithoutItems(t *testing.T) {

	payload := &models.Receipt{
		Retailer:     "Target",
		Total:        "35.35",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items:        []models.Item{},
	}

	resp, err := processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "The receipt is invalid.", string(b))
}

func TestProcessReceiptWithInvalidItems(t *testing.T) {

	payload := &models.Receipt{
		Retailer:     "Target",
		Total:        "35.35",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
			{
				ShortDescription: "Emils Cheese Pizza",
				Price:            "12.25",
			},
			{
				ShortDescription: "Knorr Creamy Chicken",
				Price:            "1.26",
			},
			{
				ShortDescription: "Doritos Nacho Cheese",
				Price:            "3.35",
			},
			{
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price:            "NOT a price",
			},
		},
	}

	resp, err := processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "The receipt is invalid.", string(b))

	payload = &models.Receipt{
		Retailer:     "Target",
		Total:        "35.35",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
			{
				ShortDescription: "Emils Cheese Pizza",
				Price:            "12.25",
			},
			{
				ShortDescription: "Knorr Creamy Chicken",
				Price:            "1.26",
			},
			{
				ShortDescription: "",
				Price:            "3.35",
			},
			{
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price:            "12.00",
			},
		},
	}

	resp, err = processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "The receipt is invalid.", string(b))
}

func TestProcessReceiptWithInvalidRetailer(t *testing.T) {

	payload := &models.Receipt{
		Retailer:     "",
		Total:        "35.35",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
		},
	}

	resp, err := processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "The receipt is invalid.", string(b))
}

func TestProcessReceiptWithInvalidTotal(t *testing.T) {

	payload := &models.Receipt{
		Retailer:     "Target",
		Total:        "NOT a total",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
		},
	}

	resp, err := processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "The receipt is invalid.", string(b))
}

func TestProcessReceiptWithInvalidDate(t *testing.T) {

	payload := &models.Receipt{
		Retailer:     "Target",
		Total:        "35.35",
		PurchaseDate: "Jan 01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
		},
	}

	resp, err := processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "The receipt is invalid.", string(b))
}

func TestProcessReceiptWithInvalidTime(t *testing.T) {

	payload := &models.Receipt{
		Retailer:     "Target",
		Total:        "35.35",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "8:00 PM",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
		},
	}

	resp, err := processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "The receipt is invalid.", string(b))
}

func TestProcessValidReceipt1(t *testing.T) {

	payload := &models.Receipt{
		Retailer:     "Target",
		Total:        "35.35",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
			{
				ShortDescription: "Emils Cheese Pizza",
				Price:            "12.25",
			},
			{
				ShortDescription: "Knorr Creamy Chicken",
				Price:            "1.26",
			},
			{
				ShortDescription: "Doritos Nacho Cheese",
				Price:            "3.35",
			},
			{
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price:            "12.00",
			},
		},
	}

	resp, err := processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	var respInfo models.ProcessReceiptResponse
	err = json.Unmarshal(respBytes, &respInfo)
	if !assert.NoError(t, err) {
		t.Errorf("Recieved unexpected response: %v", string(respBytes))
	}
	assert.NotEmpty(t, respInfo.Id)

	Receipt1ID = respInfo.Id
}

func TestProcessValidReceipt2(t *testing.T) {

	payload := &models.Receipt{
		Retailer:     "M&M Corner Market",
		Total:        "9.00",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []models.Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
	}

	resp, err := processReceipt(payload)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	var respInfo models.ProcessReceiptResponse
	err = json.Unmarshal(respBytes, &respInfo)
	if !assert.NoError(t, err) {
		t.Errorf("Recieved unexpected response: %v", string(respBytes))
	}
	assert.NotEmpty(t, respInfo.Id)

	Receipt2ID = respInfo.Id
}

func TestGetPointsOfInvalidID(t *testing.T) {

	resp, err := getPoints("NOT A UUID")
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "No receipt found for that ID.", string(b))
}

func TestGetPointsOfNonexistantReceipt(t *testing.T) {

	fakeID := uuid.New().String()
	resp, err := getPoints(fakeID)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	assert.Equal(t, "No receipt found for that ID.", string(b))
}

func TestGetPointsOfValidID1(t *testing.T) {
	resp, err := getPoints(Receipt1ID)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	var respInfo models.GetPointsResponse
	err = json.Unmarshal(respBytes, &respInfo)
	if !assert.NoError(t, err) {
		t.Errorf("Recieved unexpected response: %v", string(respBytes))
	}

	assert.Equal(t, int64(28), respInfo.Points)
}

func TestGetPointsOfValidID2(t *testing.T) {
	resp, err := getPoints(Receipt2ID)
	if err != nil {
		t.Errorf("Failed to make HTTP request: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read HTTP response: %v", err)
	}

	var respInfo models.GetPointsResponse
	err = json.Unmarshal(respBytes, &respInfo)
	if !assert.NoError(t, err) {
		t.Errorf("Recieved unexpected response: %v", string(respBytes))
	}

	assert.Equal(t, int64(109), respInfo.Points)
}

// Helper function to abstract logic of making call to ProcessReceipt
func processReceipt(receipt *models.Receipt) (*http.Response, error) {
	buf, err := json.Marshal(receipt)
	if err != nil {
		return nil, err
	}

	url := ServerEndpoint + controller.ProcessReceiptPath
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	return client.Do(req)
}

// Helper function to abstract logic of making call to GetPoints
func getPoints(id string) (*http.Response, error) {
	url := ServerEndpoint + "/receipts/" + id + "/points"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	return client.Do(req)
}
