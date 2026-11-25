package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: "https://api.weatherapi.com",
		APIKey:  apiKey,
		HTTP:    &http.Client{Timeout: 5 * time.Second},
	}
}

type response struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *Client) CurrentTempC(ctx context.Context, city string) (float64, error) {
	q := url.Values{}
	q.Set("key", c.APIKey)
	q.Set("q", city+",Brazil")
	q.Set("aqi", "no")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.BaseURL+"/v1/current.json?"+q.Encode(), nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi status %d", resp.StatusCode)
	}
	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Current.TempC, nil
}