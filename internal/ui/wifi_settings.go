package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"mifi_app/internal/api"
)

func (a *App) ShowWiFiSettingsDialog() {
	currentConfig, err := a.APIClient.GetWiFiConfig()
	if err != nil {
		a.Logger.Errorf("Failed to get WiFi config: %v", err)
		dialog.ShowError(err, a.MainWindow)
		return
	}

	ssidEntry := widget.NewEntry()
	ssidEntry.SetText(currentConfig.SSID)
	ssidEntry.SetPlaceHolder("WiFi Network Name")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(currentConfig.Password)
	passwordEntry.SetPlaceHolder("WiFi Password")

	securityModes := []string{"OPEN", "WPA2PSK", "WPA2/WPA3PSK"}
	securitySelect := widget.NewSelect(securityModes, nil)
	securitySelect.SetSelected(currentConfig.SecurityMode)

	hideSSIDCheck := widget.NewCheck("Hide Network Name (SSID)", nil)
	hideSSIDCheck.SetChecked(currentConfig.HideSSID)

	channelEntry := widget.NewEntry()
	if currentConfig.Channel > 0 {
		channelEntry.SetText(fmt.Sprintf("%d", currentConfig.Channel))
	}
	channelEntry.SetPlaceHolder("Auto (leave empty)")

	maxClientsEntry := widget.NewEntry()
	if currentConfig.MaxClients > 0 {
		maxClientsEntry.SetText(fmt.Sprintf("%d", currentConfig.MaxClients))
	}
	maxClientsEntry.SetPlaceHolder("Maximum Connected Devices")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Network Name (SSID)", Widget: ssidEntry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "Security Mode", Widget: securitySelect},
			{Text: "", Widget: hideSSIDCheck},
			{Text: "WiFi Channel", Widget: channelEntry},
			{Text: "Max Clients", Widget: maxClientsEntry},
		},
	}

	formDialog := dialog.NewCustomConfirm(
		"WiFi Settings",
		"Save",
		"Cancel",
		form,
		func(save bool) {
			if !save {
				return
			}

			a.saveWiFiSettings(ssidEntry.Text, passwordEntry.Text, securitySelect.Selected,
				hideSSIDCheck.Checked, channelEntry.Text, maxClientsEntry.Text)
		},
		a.MainWindow,
	)

	formDialog.Resize(fyne.NewSize(500, 400))
	formDialog.Show()
}

func (a *App) saveWiFiSettings(ssid, password, security string, hideSSID bool, channel, maxClients string) {
	if ssid == "" {
		dialog.ShowError(fmt.Errorf("SSID cannot be empty"), a.MainWindow)
		return
	}

	if security != "OPEN" && len(password) < 8 {
		dialog.ShowInformation("Invalid Password",
			"Password must be at least 8 characters for WPA2/WPA3 security",
			a.MainWindow)
		return
	}

	config := &api.WiFiConfig{
		SSID:         ssid,
		Password:     password,
		SecurityMode: security,
		HideSSID:     hideSSID,
	}

	if channel != "" {
		if ch, err := strconv.Atoi(channel); err == nil && ch > 0 {
			config.Channel = ch
		}
	}
	if maxClients != "" {
		if mc, err := strconv.Atoi(maxClients); err == nil && mc > 0 {
			config.MaxClients = mc
		}
	}

	if err := a.APIClient.SetWiFiConfig(config); err != nil {
		a.Logger.Errorf("Failed to update WiFi config: %v", err)
		dialog.ShowError(err, a.MainWindow)
		return
	}

	dialog.ShowInformation("Success",
		"WiFi settings have been updated successfully.\nDevices may need to reconnect.",
		a.MainWindow)
}
