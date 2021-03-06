package stockService

import (
	"testing"
)

// TestStockQuote calls GetStockQuote with a code, checking
// for a valid return value.
func TestStockCode(t *testing.T) {
	code := "aapl.us"
	msg, err := GetStockQuote(code)
	if msg == "Could not get stock quote." || err != nil {
		t.Fatalf(`GetStockQuote("aapl.us") = %q, %v, want "AAPL.US quote is 148.69 per share.", nil`, msg, err)
	}
}

// TestStockQuote calls GetStockQuote with a code, checking
// for a valid return value.
func TestStockEmpty(t *testing.T) {
	code := ""
	msg, err := GetStockQuote(code)
	if msg != "" || err == nil {
		t.Fatalf(`GetStockQuote("") = %q, %v, want "Could not get stock quote.", invalid code`, msg, err)
	}
}

// TestStockQuote calls GetStockQuote with a code, checking
// for a invalid return value.
func TestWrongStockCode(t *testing.T) {
	code := "HELLO WORLD"
	msg, err := GetStockQuote(code)
	if msg == "" || err != nil {
		t.Fatalf(`GetStockQuote("HELLO WORLD") = %q, %v, want "Could not get stock quote.", nil`, msg, err)
	}
}
