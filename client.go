package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Opts is the options to specify the chosen currency pair
type Opts struct {
	// Metal used in currency pair
	Metal Metal
	// Currency used in currency pair
	Currency Currency
}

// GoldapiClient is the client wrapper for querying goldapi.io
type GoldapiClient struct {
	baseURL    string
	apiKey     string
	opts       Opts
	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}
}

// NewGoldapiClient creates a new GoldapiClient object
func NewGoldapiClient(baseURL, apiKey string, metal Metal, currency Currency) *GoldapiClient {
	return &GoldapiClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		opts: Opts{
			Metal:    metal,
			Currency: currency,
		},
		httpClient: &http.Client{},
	}
}

func (gc GoldapiClient) url() string {
	return fmt.Sprintf("%s/%s/%s/", gc.baseURL, gc.opts.Metal, gc.opts.Currency)
}

func (gc GoldapiClient) buildReq() (*http.Request, error) {
	req, err := http.NewRequest("GET", gc.url(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-access-token", gc.apiKey)

	return req, nil
}

func (gc GoldapiClient) get() (*Response, error) {
	log.Println(fmt.Sprintf("Getting data from %q...", gc.url()))

	req, err := gc.buildReq()
	if err != nil {
		return &Response{}, err
	}

	resp, err := gc.httpClient.Do(req)
	if err != nil {
		return &Response{}, err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{}, err
	}

	var respObj Response
	json.Unmarshal(respData, &respObj)

	return &respObj, nil
}

// Response defines the data returned from a successful goldapi query
type Response struct {
	Timestamp      int64   `json:"timestamp"`
	Metal          string  `json:"metal"`
	Currency       string  `json:"currency"`
	Exchange       string  `json:"exchange"`
	Symbol         string  `json:"symbol"`
	PrevClosePrice float32 `json:"prev_close_price"`
	OpenPrice      float32 `json:"open_price"`
	LowPrice       float32 `json:"low_price"`
	HighPrice      float32 `json:"high_price"`
	OpenTime       int64   `json:"open_time"`
	Price          float32 `json:"price"`
	Ch             float32 `json:"ch"`
	Chp            float32 `json:"chp"`
	Ask            float32 `json:"ask"`
	Bid            float32 `json:"bid"`
}

// Metal is used to define the metal selected in the currency pair
type Metal int

// Valid Metal
const (
	Gold Metal = iota
	Silver
	Platinum
	Palladium
)

func (m Metal) String() string {
	switch m {
	case Gold:
		return "XAU"
	case Silver:
		return "XAG"
	case Platinum:
		return "XPT"
	case Palladium:
		return "XPD"
	default:
		return "unknown"
	}
}

// Set sets the Metal to a value represented by the string s. Set Implements the flag.Value interface.
func (m *Metal) Set(s string) error {
	switch s {
	case "XAU":
		*m = Gold
	case "XAG":
		*m = Silver
	case "XPT":
		*m = Platinum
	case "XPD":
		*m = Palladium
	default:
		return fmt.Errorf("unknown metal %q: expected either XAU, XAG, XPT, or XPD", s)
	}

	return nil
}

// Currency is used to define the fiat selected in the currency pair
type Currency int

// Valid currency
const (
	USD Currency = iota
	AUD
	GBP
	EUR
	CHF
	CAD
	JPY
	INR
	SGD
	BTC
	CZK
	RUB
	PLN
	MYR
	XAG // Silver - Returns the gold/silver ratio
)

func (c Currency) String() string {
	switch c {
	case USD:
		return "USD"
	case AUD:
		return "AUD"
	case GBP:
		return "GBP"
	case EUR:
		return "EUR"
	case CHF:
		return "CHF"
	case CAD:
		return "CAD"
	case JPY:
		return "JPY"
	case INR:
		return "INR"
	case SGD:
		return "SGD"
	case BTC:
		return "BTC"
	case CZK:
		return "CZK"
	case RUB:
		return "RUB"
	case PLN:
		return "PLN"
	case MYR:
		return "MYR"
	case XAG:
		return "XAG"
	default:
		return "unknown currency code %q"
	}
}

// Set sets the Currency value represented by the string s. Set implements the flag.Value interface.
func (c *Currency) Set(s string) error {
	switch s {
	case "USD":
		*c = USD
	case "AUD":
		*c = AUD
	case "GBP":
		*c = GBP
	case "EUR":
		*c = EUR
	case "CHF":
		*c = CHF
	case "CAD":
		*c = CAD
	case "JPY":
		*c = JPY
	case "INR":
		*c = INR
	case "SGD":
		*c = SGD
	case "BTC":
		*c = BTC
	case "CZK":
		*c = CZK
	case "RUB":
		*c = RUB
	case "PLN":
		*c = PLN
	case "MYR":
		*c = MYR
	case "XAG":
		*c = XAG
	default:
		return fmt.Errorf("unknown currency code %q", s)
	}

	return nil
}
