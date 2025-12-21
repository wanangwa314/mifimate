package api

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (c *Client) GetDeviceStatus() (*DeviceStatus, error) {
	params := map[string]string{
		"cmd": "modem_main_state,pin_status,network_type,signalbar,battery_value," +
			"battery_charging,wifi_status,ssid1,station_mac,network_provider," +
			"wan_ipaddr,wan_apn,ppp_status,realtime_tx_bytes,realtime_rx_bytes," +
			"realtime_time,realtime_tx_thrpt,realtime_rx_thrpt,sta_count",
		"multi_data": "1",
		"isTest":     "false",
	}

	resp, err := c.Get(StatusEndpoint, params)
	if err != nil {
		return nil, err
	}

	status := &DeviceStatus{}

	if val, ok := resp["network_type"].(string); ok && val != "" {
		status.NetworkType = val
	}
	if val, ok := resp["signalbar"].(string); ok && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			status.SignalStrength = i
		}
	}
	if val, ok := resp["battery_value"].(string); ok && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			status.BatteryLevel = i
		}
	}
	if val, ok := resp["wan_ipaddr"].(string); ok && val != "" {
		status.WanIPAddress = val
	}
	if val, ok := resp["sta_count"].(string); ok && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			status.ConnectedDevs = i
		}
	}
	if val, ok := resp["realtime_tx_thrpt"].(string); ok && val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			status.TxSpeed = f
		}
	}
	if val, ok := resp["realtime_rx_thrpt"].(string); ok && val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			status.RxSpeed = f
		}
	}
	if val, ok := resp["realtime_tx_bytes"].(string); ok && val != "" {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			status.TxBytes = u
		}
	}
	if val, ok := resp["realtime_rx_bytes"].(string); ok {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			status.RxBytes = u
		}
	}
	if val, ok := resp["imei"].(string); ok {
		status.IMEI = val
	}
	if val, ok := resp["iccid"].(string); ok {
		status.ICCID = val
	}
	if val, ok := resp["model_name"].(string); ok {
		status.ModelName = val
	}
	if val, ok := resp["hardware_version"].(string); ok {
		status.HardwareVersion = val
	}
	if val, ok := resp["software_version"].(string); ok {
		status.SoftwareVersion = val
	}

	return status, nil
}

func (c *Client) GetWiFiConfig() (*WiFiConfig, error) {
	params := map[string]string{
		"cmd": "wifi_ssid,wifi_password,security_mode,hide_ssid,wifi_channel,max_client_num",
	}

	resp, err := c.Get(StatusEndpoint, params)
	if err != nil {
		return nil, err
	}

	config := &WiFiConfig{}

	if val, ok := resp["wifi_ssid"].(string); ok {
		config.SSID = val
	}
	if val, ok := resp["wifi_password"].(string); ok {
		config.Password = val
	}
	if val, ok := resp["security_mode"].(string); ok {
		config.SecurityMode = val
	}
	if val, ok := resp["hide_ssid"].(string); ok {
		config.HideSSID = val == "1"
	}
	if val, ok := resp["wifi_channel"].(string); ok {
		if i, err := strconv.Atoi(val); err == nil {
			config.Channel = i
		}
	}
	if val, ok := resp["max_client_num"].(string); ok {
		if i, err := strconv.Atoi(val); err == nil {
			config.MaxClients = i
		}
	}

	return config, nil
}

func (c *Client) SetWiFiConfig(config *WiFiConfig) error {
	data := map[string]string{
		"goformId":      "SET_WIFI_SSID_PASSWORD",
		"wifi_ssid":     config.SSID,
		"wifi_password": config.Password,
		"security_mode": config.SecurityMode,
	}

	if config.HideSSID {
		data["hide_ssid"] = "1"
	} else {
		data["hide_ssid"] = "0"
	}

	if config.Channel > 0 {
		data["wifi_channel"] = strconv.Itoa(config.Channel)
	}

	if config.MaxClients > 0 {
		data["max_client_num"] = strconv.Itoa(config.MaxClients)
	}

	resp, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return err
	}

	if result, ok := resp["result"].(string); ok {
		if result == "0" || strings.ToLower(result) == "success" {
			return nil
		}
		return fmt.Errorf("failed to update WiFi config: %s", result)
	}

	return fmt.Errorf("unexpected response format")
}

// ConnectNetwork attempts to connect to the network
func (c *Client) ConnectNetwork() error {
	data := map[string]string{
		"goformId": "CONNECT_NETWORK",
	}

	resp, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return err
	}

	if result, ok := resp["result"].(string); ok {
		if result == "0" || strings.ToLower(result) == "success" {
			return nil
		}
		return fmt.Errorf("failed to connect: %s", result)
	}

	return fmt.Errorf("unexpected response format")
}

// DisconnectNetwork disconnects from the network
func (c *Client) DisconnectNetwork() error {
	data := map[string]string{
		"goformId": "DISCONNECT_NETWORK",
	}

	resp, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return err
	}

	if result, ok := resp["result"].(string); ok {
		if result == "0" || strings.ToLower(result) == "success" {
			return nil
		}
		return fmt.Errorf("failed to disconnect: %s", result)
	}

	return fmt.Errorf("unexpected response format")
}

// GetConnectedDevices retrieves the list of connected devices
func (c *Client) GetConnectedDevices() ([]ConnectedDevice, error) {
	params := map[string]string{
		"cmd": "station_list",
	}

	resp, err := c.Get(StatusEndpoint, params)
	if err != nil {
		return nil, err
	}

	// The actual format may vary - this is a placeholder
	devices := []ConnectedDevice{}

	// Parse station list if available
	if stationList, ok := resp["station_list"].([]interface{}); ok {
		for _, station := range stationList {
			if s, ok := station.(map[string]interface{}); ok {
				device := ConnectedDevice{}
				if val, ok := s["hostname"].(string); ok {
					device.Hostname = val
				}
				if val, ok := s["ipaddress"].(string); ok {
					device.IPAddress = val
				}
				if val, ok := s["macaddress"].(string); ok {
					device.MACAddress = val
				}
				devices = append(devices, device)
			}
		}
	}

	return devices, nil
}

func (c *Client) GetSMSCount() (int, error) {
	params := map[string]string{
		"cmd":        "sms_data_total",
		"multi_data": "1",
		"isTest":     "false",
	}

	resp, err := c.Get(StatusEndpoint, params)
	if err != nil {
		return 0, err
	}

	if val, ok := resp["sms_data_total"].(string); ok {
		if count, err := strconv.Atoi(val); err == nil {
			return count, nil
		}
	}

	return 0, nil
}

func (c *Client) GetSMSList(page, pageSize int) ([]SMSMessage, error) {
	params := map[string]string{
		"cmd":           "sms_data_total",
		"page":          strconv.Itoa(page),
		"data_per_page": strconv.Itoa(pageSize),
		"mem_store":     "1",
		"tags":          "10",
		"order_by":      "order+by+id+desc",
	}

	resp, err := c.Get(StatusEndpoint, params)
	if err != nil {
		return nil, err
	}

	var messages []SMSMessage

	if msgList, ok := resp["messages"].([]interface{}); ok {
		for _, msg := range msgList {
			if m, ok := msg.(map[string]interface{}); ok {
				sms := SMSMessage{}
				if val, ok := m["id"].(string); ok {
					sms.ID = val
				}
				if val, ok := m["number"].(string); ok {
					sms.Number = val
				}
				if val, ok := m["content"].(string); ok {
					decoded, err := decodeHexSMS(val)
					if err != nil {
						sms.Content = val
					} else {
						sms.Content = decoded
					}
				}
				if val, ok := m["tag"].(string); ok {
					sms.Status = val
				}
				if val, ok := m["date"].(string); ok {
					// Parse date string - device returns format: YY,MM,DD,HH,MM,SS,+TZ
					// Example: "25,12,20,18,38,02,+8" means 2025-12-20 18:38:02 +8
					parts := strings.Split(val, ",")
					if len(parts) >= 6 {
						// Convert to standard format
						year := "20" + parts[0]
						month := parts[1]
						day := parts[2]
						hour := parts[3]
						minute := parts[4]
						second := parts[5]
						dateStr := fmt.Sprintf("%s-%s-%s %s:%s:%s", year, month, day, hour, minute, second)

						if t, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
							sms.Timestamp = t
						}
					}
				}
				messages = append(messages, sms)
			}
		}
	}

	return messages, nil
}

func (c *Client) SendSMS(phoneNumber, content string) error {
	data := map[string]string{
		"goformId":    "SEND_SMS",
		"notCallback": "true",
		"Number":      phoneNumber,
		"sms_time":    time.Now().Format("2006-01-02 15:04:05"),
		"MessageBody": content,
		"ID":          "-1",
		"encode_type": "GSM7_default",
		"isTest":      "false",
	}

	resp, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return err
	}

	if result, ok := resp["result"].(string); ok {
		if result == "success" {
			return nil
		}
		return fmt.Errorf("failed to send SMS: %s", result)
	}

	return fmt.Errorf("unexpected response format")
}

func (c *Client) DeleteSMS(messageIDs []string) error {
	data := map[string]string{
		"goformId":    "DELETE_SMS",
		"msg_id":      strings.Join(messageIDs, ";"),
		"notCallback": "true",
		"isTest":      "false",
	}

	resp, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return err
	}

	if result, ok := resp["result"].(string); ok {
		if result == "success" {
			return nil
		}
		return fmt.Errorf("failed to delete SMS: %s", result)
	}

	return fmt.Errorf("unexpected response format")
}

func (c *Client) RebootDevice() error {
	data := map[string]string{
		"goformId": "REBOOT_DEVICE",
	}

	resp, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return fmt.Errorf("failed to send reboot command: %w", err)
	}

	if result, ok := resp["result"].(string); ok {
		if result == "success" {
			return nil
		}
		return fmt.Errorf("failed to reboot device: %s", result)
	}

	return fmt.Errorf("unexpected response format")
}

func (c *Client) ShutdownDevice() error {
	data := map[string]string{
		"goformId": "POWEROFF_DEVICE",
	}

	resp, err := c.Post(LoginEndpoint, data)
	if err != nil {
		return fmt.Errorf("failed to send shutdown command: %w", err)
	}

	if result, ok := resp["result"].(string); ok {
		if result == "success" {
			return nil
		}
		return fmt.Errorf("failed to shutdown device: %s", result)
	}

	return fmt.Errorf("unexpected response format")
}

func decodeHexSMS(hexContent string) (string, error) {
	hexContent = strings.ReplaceAll(hexContent, " ", "")

	decoded, err := hex.DecodeString(hexContent)
	if err != nil {
		return "", fmt.Errorf("invalid hex string: %w", err)
	}

	return string(decoded), nil
}
