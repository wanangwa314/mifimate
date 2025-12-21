package api

import "time"

// DeviceStatus represents the current status of the MiFi device
type DeviceStatus struct {
	NetworkType     string  `json:"network_type"`
	SignalStrength  int     `json:"signalbar"`
	BatteryLevel    int     `json:"battery_value"`
	WanIPAddress    string  `json:"wan_ipaddr"`
	ConnectedDevs   int     `json:"sta_count"`
	TxSpeed         float64 `json:"realtime_tx_thrpt"`
	RxSpeed         float64 `json:"realtime_rx_thrpt"`
	TxBytes         uint64  `json:"realtime_tx_bytes"`
	RxBytes         uint64  `json:"realtime_rx_bytes"`
	IMEI            string  `json:"imei"`
	ICCID           string  `json:"iccid"`
	ModelName       string  `json:"model_name"`
	HardwareVersion string  `json:"hardware_version"`
	SoftwareVersion string  `json:"software_version"`
}

// WiFiConfig represents WiFi configuration settings
type WiFiConfig struct {
	SSID         string `json:"ssid"`
	Password     string `json:"password"`
	SecurityMode string `json:"security_mode"` // WPA2PSK, WPA3PSK, etc.
	HideSSID     bool   `json:"hide_ssid"`
	Channel      int    `json:"channel"`
	MaxClients   int    `json:"max_client_num"`
}

// SMSMessage represents an SMS message
type SMSMessage struct {
	ID        string    `json:"id"`
	Number    string    `json:"number"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"date"`
	Status    string    `json:"status"` // read/unread
	Type      string    `json:"type"`   // inbox/sent
}

// ConnectedDevice represents a device connected to the MiFi
type ConnectedDevice struct {
	Hostname      string    `json:"hostname"`
	IPAddress     string    `json:"ipaddress"`
	MACAddress    string    `json:"macaddress"`
	ConnectedTime time.Time `json:"connected_time"`
	IsBlocked     bool      `json:"is_blocked"`
}

// LoginRequest represents a login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Result  string `json:"result"`
	Session string `json:"session"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Result  string                 `json:"result"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
