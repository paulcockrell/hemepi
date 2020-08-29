package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGoldapiClient(t *testing.T) {
	url := "http://api.com"
	apiKey := "apikey"
	metal := Gold
	currency := GBP
	client := NewGoldapiClient(
		url,
		apiKey,
		metal,
		currency,
	)

	t.Run("builds valid api URL", func(t *testing.T) {
		got := client.url()
		want := fmt.Sprintf("%s/%s/%s/", url, metal, currency)

		assertEqual(t, got, want)
	})

	t.Run("request has correct headers set", func(t *testing.T) {
		req, err := client.buildReq()
		if err != nil {
			t.Fatalf(err.Error())
		}

		got := req.Header.Get("x-access-token")
		want := apiKey

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		assertEqual(t, got, want)
	})
}

func TestGoldapiClient_Integration(t *testing.T) {
	t.Run("valid api key returns successfull response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := json.Marshal(mockResponse)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Fprintln(w, string(data))
		}))
		defer ts.Close()

		client := NewGoldapiClient(
			ts.URL,
			"apiKey",
			Gold,
			GBP,
		)

		want := mockResponse
		got, err := client.get()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})

	t.Run("invalid api key returns error", func(t *testing.T) {
		errorText := "Invalid API Key"

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errResponse := map[string]string{
				"error": errorText,
			}
			data, err := json.Marshal(errResponse)
			if err != nil {
				t.Fatal(err)
			}

			w.WriteHeader(http.StatusForbidden)

			fmt.Fprintln(w, string(data))
		}))
		defer ts.Close()

		client := NewGoldapiClient(
			ts.URL,
			"badapiKey",
			Gold,
			GBP,
		)

		_, err := client.get()
		if err == nil {
			t.Fatalf("Expected an error to be returned")
		}

		assertEqual(t, err.Error(), errorText)
	})
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

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
