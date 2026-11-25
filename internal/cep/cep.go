package cep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

var ErrNotFound = fmt.Errorf("zipcode not found")

var cepRegex = regexp.MustCompile(`^\d{8}$`)

func IsValid(cep string) bool { return cepRegex.MatchString(cep) }

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://viacep.com.br",
		HTTP:    &http.Client{Timeout: 5 * time.Second},
	}
}

type viaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       any    `json:"erro"`
}

func (c *Client) Lookup(ctx context.Context, cep string) (string, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", c.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep status %d", resp.StatusCode)
	}
	var data viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if data.Erro != nil || data.Localidade == "" {
		return "", ErrNotFound
	}
	return data.Localidade, nil
}