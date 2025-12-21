package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DefaultTimeout = 10 * time.Second
	DefaultBaseURL = "http://192.168.1.1"
)

// Client represents the HTTP client for MiFi API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     *logrus.Logger
	sessionID  string
}

func NewClient(baseURL string, logger *logrus.Logger) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if logger == nil {
		logger = logrus.New()
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Warnf("Failed to create cookie jar: %v", err)
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
			Jar:     jar,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		Logger: logger,
	}
}

func (c *Client) SetSessionID(sessionID string) {
	c.sessionID = sessionID
}

func (c *Client) GetSessionID() string {
	return c.sessionID
}

func (c *Client) Get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	u, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Referer", c.BaseURL+"/index.html")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")

	if c.sessionID != "" {
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: c.sessionID,
		})
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}

func (c *Client) Post(endpoint string, data map[string]string) (map[string]interface{}, error) {
	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, value)
	}

	encodedData := formData.Encode()

	req, err := http.NewRequest("POST", c.BaseURL+endpoint, strings.NewReader(encodedData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", c.BaseURL+"/index.html")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")

	if c.sessionID != "" {
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: c.sessionID,
		})
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}

func (c *Client) Ping() error {
	req, err := http.NewRequest("GET", c.BaseURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	return fmt.Errorf("device unreachable, status: %d", resp.StatusCode)
}
