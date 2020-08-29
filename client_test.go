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

func TestGoldapiClient_ValidRequest(t *testing.T) {
	t.Run("valid request returns successfull response", func(t *testing.T) {
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
}

func TestGoldapiClient_InvalidRequest(t *testing.T) {
	cases := []struct {
		name, errorText string
		code            int
	}{
		{name: "Invalid curency pair", code: http.StatusOK, errorText: "No data available for this request"},
		{name: "Invalid api key", code: http.StatusForbidden, errorText: "Invalid API Key"},
		{name: "No api key", code: http.StatusForbidden, errorText: "No API Key provided"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				errResponse := map[string]string{
					"error": test.errorText,
				}
				data, err := json.Marshal(errResponse)
				if err != nil {
					t.Fatal(err)
				}

				w.WriteHeader(test.code)

				fmt.Fprintln(w, string(data))
			}))
			defer ts.Close()

			client := NewGoldapiClient(
				ts.URL,
				"apiKey",
				Gold,
				GBP,
			)

			_, err := client.get()
			if err == nil {
				t.Fatalf("Expected an error to be returned")
			}

			assertEqual(t, err.Error(), test.errorText)
		})
	}
}

func TestMetal(t *testing.T) {
	t.Run("valid metal codes", func(t *testing.T) {
		cases := []struct {
			name, metalCode string
			metal           Metal
		}{
			{name: "'XAU' translates to Gold constant", metalCode: "XAU", metal: Gold},
			{name: "'XAG' translates to Silver constant", metalCode: "XAG", metal: Silver},
			{name: "'XPT' translates to Platinum constant", metalCode: "XPT", metal: Platinum},
			{name: "'XPD' translates to Palladium constant", metalCode: "XPD", metal: Palladium},
		}

		for _, test := range cases {
			t.Run(test.name, func(t *testing.T) {
				var m Metal
				m.Set(test.metalCode)
				if m != test.metal {
					t.Errorf("got %d want %d", m, test.metal)
				}
				if m.String() != test.metalCode {
					t.Errorf("got %q want %q", m.String(), test.metalCode)
				}
			})
		}
	})

	t.Run("invalid metal codes", func(t *testing.T) {
		invalidMetalCode := "XXX"
		var m Metal
		err := m.Set(invalidMetalCode)
		if err == nil {
			t.Errorf("Expected an error for invalid metal code %q", invalidMetalCode)
		}
	})

	t.Run("invalid metal value", func(t *testing.T) {
		var m Metal = 99
		got := m.String()
		want := fmt.Sprintf("unknown metal code: %d", m)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestCurrency(t *testing.T) {
	t.Run("valid currency codes", func(t *testing.T) {
		cases := []struct {
			name, currencyCode string
			currency           Currency
		}{
			{name: "'USD' translates to USD constant", currencyCode: "USD", currency: USD},
			{name: "'AUD' translates to AUD constant", currencyCode: "AUD", currency: AUD},
			{name: "'GBP' translates to GBP constant", currencyCode: "GBP", currency: GBP},
			{name: "'EUR' translates to EUR constant", currencyCode: "EUR", currency: EUR},
			{name: "'CHF' translates to CHF constant", currencyCode: "CHF", currency: CHF},
			{name: "'CAD' translates to CHF constant", currencyCode: "CAD", currency: CAD},
			{name: "'JPY' translates to JPY constant", currencyCode: "JPY", currency: JPY},
			{name: "'INR' translates to INR constant", currencyCode: "INR", currency: INR},
			{name: "'SGD' translates to SGD constant", currencyCode: "SGD", currency: SGD},
			{name: "'BTC' translates to BTC constant", currencyCode: "BTC", currency: BTC},
			{name: "'CZK' translates to CZK constant", currencyCode: "CZK", currency: CZK},
			{name: "'RUB' translates to RUB constant", currencyCode: "RUB", currency: RUB},
			{name: "'PLN' translates to PLN constant", currencyCode: "PLN", currency: PLN},
			{name: "'MYR' translates to MYR constant", currencyCode: "MYR", currency: MYR},
			{name: "'XAG' translates to XAG constant", currencyCode: "XAG", currency: XAG},
		}

		for _, test := range cases {
			t.Run(test.name, func(t *testing.T) {
				var c Currency
				err := c.Set(test.currencyCode)
				if err != nil {
					t.Fatal(err)
				}
				if c != test.currency {
					t.Errorf("got %d want %d", c, test.currency)
				}
				if c.String() != test.currencyCode {
					t.Errorf("got %q want %q", c.String(), test.currencyCode)
				}
			})
		}
	})

	t.Run("invalid currency codes", func(t *testing.T) {
		invalidCurrencyCode := "XXX"
		var c Currency
		err := c.Set(invalidCurrencyCode)
		if err == nil {
			t.Errorf("Expected an error for invalid currency code %q", invalidCurrencyCode)
		}
	})

	t.Run("invalid currency value", func(t *testing.T) {
		var c Currency = 99
		got := c.String()
		want := fmt.Sprintf("unknown currency code: %d", c)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
