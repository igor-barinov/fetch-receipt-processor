/**
models.go

Describes the request/response struct definitions along with some helper methods
*/

package models

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Regexes for validating properties
var (
	shortDescRgx = regexp.MustCompile(`^[\w\s\-]+$`)
	dollarAmtRgx = regexp.MustCompile(`^\d+\.\d{2}$`)
	retailerRgx  = regexp.MustCompile(`^[\w\s\-&]+$`)
	DateFormat   = "2006-01-02"
	TimeFormat   = "15:04"
)

// Describes a purchased item in a receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Returns an error if any of the properties are invalid
// Returns `nil` otherwise
func (item *Item) ValidateProperties() error {
	ok := shortDescRgx.MatchString(item.ShortDescription)
	if !ok {
		return fmt.Errorf("property ShortDescription didn't follow pattern: %v", shortDescRgx.String())
	}

	ok = dollarAmtRgx.MatchString(item.Price)
	if !ok {
		return fmt.Errorf("property Price didn't follow pattern: %v", dollarAmtRgx.String())
	}

	return nil
}

// Describes a receipt of a transaction
type Receipt struct {
	Retailer     string `json:"retailer"`
	Total        string `json:"total"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
}

// Returns an error if any of the properties are invalid
// Returns `nil` otherwise
func (r *Receipt) ValidateProperties() error {

	if len(r.Items) < 1 {
		return fmt.Errorf("property Items must have at least 1 item")
	}

	ok := retailerRgx.MatchString(r.Retailer)
	if !ok {
		return fmt.Errorf("property Retailer didn't follow pattern: %v", retailerRgx.String())
	}

	ok = dollarAmtRgx.MatchString(r.Total)
	if !ok {
		return fmt.Errorf("property Total didn't follow pattern: %v", dollarAmtRgx.String())
	}

	_, err := time.Parse(DateFormat, r.PurchaseDate)
	if err != nil {
		return fmt.Errorf("property PruchaseDate is not a valid date; %v", err)
	}

	_, err = time.Parse(TimeFormat, r.PurchaseTime)
	if err != nil {
		return fmt.Errorf("property PurchaseTime is not a valid time; %v", err)
	}

	for _, item := range r.Items {
		err := item.ValidateProperties()
		if err != nil {
			return err
		}
	}

	return nil
}

// Calculates the points based on the receipt details
// Assumes properties are valid, if not then calling this will result in UB since conversion errors are unchecked
func (r *Receipt) CalculatePoints() int64 {

	totalPoints := int64(0)

	// Retailer alphanum rule
	for _, r := range r.Retailer {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			totalPoints += 1
		}
	}

	// Parse total
	s := dollarAmtRgx.FindString(r.Total)
	parts := strings.Split(s, ".")
	// dollars, _ := strconv.Atoi(parts[0])
	cents, _ := strconv.Atoi(parts[1])

	// Round dollar rule
	if cents == 0 {
		totalPoints += 50
	}

	// Multiple of 0.25 rule
	if cents%25 == 0 {
		totalPoints += 25
	}

	// Every 2 items rule
	totalPoints += 5 * int64(len(r.Items)/2)

	// Trimmed description rule
	for _, item := range r.Items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if len(trimmed)%3 == 0 && len(trimmed) > 0 {
			priceAmt, _ := strconv.ParseFloat(item.Price, 64)
			totalPoints += int64(math.Ceil(priceAmt * 0.2))
		}
	}

	// Odd day rule
	d, _ := time.Parse(DateFormat, r.PurchaseDate)
	if d.Day()%2 != 0 {
		totalPoints += 6
	}

	// 2-4 PM rule
	t, _ := time.Parse(TimeFormat, r.PurchaseTime)
	twoPM, _ := time.Parse(TimeFormat, "14:00")
	fourPM, _ := time.Parse(TimeFormat, "18:00")
	if t.After(twoPM) && t.Before(fourPM) {
		totalPoints += 10
	}

	return totalPoints

}

// Same method as above but with print statements, for debugging
func (r *Receipt) CalculatePointsVerbose() int64 {

	totalPoints := int64(0)

	// Retailer alphanum rule
	addition := int64(0)
	for _, r := range r.Retailer {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			addition += 1
		}
	}

	log.Printf("%v alphanum chars +%v", r.Retailer, addition)
	totalPoints += addition

	// Parse total
	s := dollarAmtRgx.FindString(r.Total)
	parts := strings.Split(s, ".")
	// dollars, _ := strconv.Atoi(parts[0])
	cents, _ := strconv.Atoi(parts[1])

	// Round dollar rule
	if cents == 0 {
		log.Printf("%v is a round dollar amount +50", r.Total)
		totalPoints += 50
	}

	// Multiple of 0.25 rule
	if cents%25 == 0 {
		log.Printf("%v is a multiple of 0.25 +50", r.Total)
		totalPoints += 25
	}

	// Every 2 items rule
	addition = 5 * int64(len(r.Items)/2)
	log.Printf("For every two items +%v", addition)
	totalPoints += addition

	// Trimmed description rule
	for _, item := range r.Items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if len(trimmed)%3 == 0 && len(trimmed) > 0 {
			priceAmt, _ := strconv.ParseFloat(item.Price, 64)
			addition := int64(math.Ceil(priceAmt * 0.2))

			log.Printf("%v trimmed len is multiple of 3 +%v", item.ShortDescription, addition)
			totalPoints += addition
		}
	}

	// Odd day rule
	d, _ := time.Parse(DateFormat, r.PurchaseDate)
	if d.Day()%2 != 0 {
		log.Printf("%v contains an odd day +6", r.PurchaseDate)
		totalPoints += 6
	}

	// 2-4 PM rule
	t, _ := time.Parse(TimeFormat, r.PurchaseTime)
	twoPM, _ := time.Parse(TimeFormat, "14:00")
	fourPM, _ := time.Parse(TimeFormat, "18:00")
	if t.After(twoPM) && t.Before(fourPM) {
		log.Printf("%v is between 2 and 4 +10", r.PurchaseTime)
		totalPoints += 10
	}

	return totalPoints

}

// Describes the response structure for the `ProcessReceipt` endpoint
type ProcessReceiptResponse struct {
	Id string `json:"id"`
}

// Describes the response structure for the `GetPoints` endpoint
type GetPointsResponse struct {
	Points int64 `json:"points"`
}
