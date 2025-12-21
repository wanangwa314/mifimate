package api

import (
	"encoding/base64"
	"fmt"
)

const (
	LoginEndpoint  = "/goform/goform_set_cmd_process"
	StatusEndpoint = "/goform/goform_get_cmd_process"
)

func (c *Client) Login(username, password string) error {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

	data := map[string]string{
		"isTest":   "false",
		"goformId": "LOGIN",
		"password": encodedPassword,
	}

	resp, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	if result, ok := resp["result"].(string); ok {
		if result == "0" || result == "4" {
			return nil
		}
		switch result {
		case "1":
			return fmt.Errorf("login failed: general failure")
		case "2":
			return fmt.Errorf("login failed: duplicate user (already logged in)")
		case "3":
			return fmt.Errorf("login failed: bad password")
		default:
			return fmt.Errorf("login failed: %s", result)
		}
	}

	return fmt.Errorf("unexpected response format")
}

// Helper function to get map keys for debugging
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Logout ends the current session
func (c *Client) Logout() error {
	data := map[string]string{
		"goformId": "LOGOUT",
	}

	_, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return fmt.Errorf("logout request failed: %w", err)
	}

	c.sessionID = ""
	return nil
}

// IsAuthenticated checks if the client has a valid session
func (c *Client) IsAuthenticated() bool {
	_, err := c.GetDeviceStatus()
	return err == nil
}
