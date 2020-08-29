package main

import "testing"

var mockResponse = &Response{
	Timestamp:      1598525693,
	Metal:          "XAU",
	Currency:       "USD",
	Exchange:       "FOREXCOM",
	Symbol:         "FOREXCOM:XAUUSD",
	PrevClosePrice: 1954.27,
	OpenPrice:      1954.27,
	LowPrice:       1937.04,
	HighPrice:      1955.8,
	OpenTime:       1598475600,
	Price:          1939.11,
	Ch:             -15.16,
	Chp:            -0.78,
	Ask:            1939.56,
	Bid:            1938.83,
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
