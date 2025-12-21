package ui

import (
	"errors"
	"strconv"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) ShowSettingsDialog() {
	// Theme selection
	themeOptions := []string{"Light", "Dark", "System"}
	themeSelect := widget.NewSelect(themeOptions, nil)
	themeSelect.SetSelected(a.Config.App.Theme)

	// Auto-start checkbox
	autoStartCheck := widget.NewCheck("Start application on system login", nil)
	autoStartCheck.SetChecked(a.Config.App.AutoStart)

	// Notifications checkbox
	notificationsCheck := widget.NewCheck("Show desktop notifications", nil)
	notificationsCheck.SetChecked(a.Config.App.ShowNotifications)

	// Polling interval
	pollingEntry := widget.NewEntry()
	pollingEntry.SetText(strconv.Itoa(a.Config.Device.PollInterval))
	pollingEntry.SetPlaceHolder("Seconds between updates")

	// Log level
	logLevelOptions := []string{"debug", "info", "warn", "error"}
	logLevelSelect := widget.NewSelect(logLevelOptions, nil)
	logLevelSelect.SetSelected(a.Config.App.LogLevel)

	// Connection timeout
	timeoutEntry := widget.NewEntry()
	timeoutEntry.SetText(strconv.Itoa(a.Config.Device.ConnectionTimeout))
	timeoutEntry.SetPlaceHolder("Seconds")

	// Auto-reconnect checkbox
	autoReconnectCheck := widget.NewCheck("Automatically reconnect on network issues", nil)
	autoReconnectCheck.SetChecked(a.Config.Device.AutoReconnect)

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Theme", Widget: themeSelect},
			{Text: "Auto Start", Widget: autoStartCheck},
			{Text: "Notifications", Widget: notificationsCheck},
			{Text: "Poll Interval (s)", Widget: pollingEntry},
			{Text: "Log Level", Widget: logLevelSelect},
			{Text: "Connection Timeout (s)", Widget: timeoutEntry},
			{Text: "Auto Reconnect", Widget: autoReconnectCheck},
		},
	}

	// Create dialog
	dialog.ShowCustomConfirm(
		"Application Settings",
		"Save",
		"Cancel",
		form,
		func(save bool) {
			if !save {
				return
			}

			a.saveSettings(
				themeSelect.Selected,
				autoStartCheck.Checked,
				notificationsCheck.Checked,
				pollingEntry.Text,
				logLevelSelect.Selected,
				timeoutEntry.Text,
				autoReconnectCheck.Checked,
			)
		},
		a.MainWindow,
	)
}

// saveSettings validates and saves the application settings
func (a *App) saveSettings(theme string, autoStart, notifications bool, pollInterval, logLevel, timeout string, autoReconnect bool) {
	// Validate poll interval
	poll, err := strconv.Atoi(pollInterval)
	if err != nil || poll < 1 {
		dialog.ShowError(errors.New("invalid poll interval. Must be a number >= 1"), a.MainWindow)
		return
	}

	// Validate timeout
	t, err := strconv.Atoi(timeout)
	if err != nil || t < 1 {
		dialog.ShowError(errors.New("invalid timeout. Must be a number >= 1"), a.MainWindow)
		return
	}

	// Update config
	a.Config.App.Theme = theme
	a.Config.App.AutoStart = autoStart
	a.Config.App.ShowNotifications = notifications
	a.Config.App.LogLevel = logLevel
	a.Config.Device.PollInterval = poll
	a.Config.Device.ConnectionTimeout = t
	a.Config.Device.AutoReconnect = autoReconnect

	// Save to file
	if err := a.Config.Save(); err != nil {
		a.Logger.Errorf("Failed to save settings: %v", err)
		dialog.ShowError(err, a.MainWindow)
		return
	}

	dialog.ShowInformation("Success",
		"Settings have been saved.\\nSome changes may require restart.",
		a.MainWindow)
}
