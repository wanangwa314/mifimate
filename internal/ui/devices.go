package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) ShowDevicesDialog() {
	// Create devices list
	devicesList := widget.NewList(
		func() int {
			return len(a.cachedDevices)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Template Device"),
				widget.NewLabel("Template Details"),
				widget.NewSeparator(),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(a.cachedDevices) {
				device := a.cachedDevices[id]
				box := obj.(*fyne.Container)

				nameLabel := box.Objects[0].(*widget.Label)
				hostname := device.Hostname
				if hostname == "" {
					hostname = "Unknown Device"
				}
				nameLabel.SetText(hostname)
				nameLabel.TextStyle.Bold = true

				detailsLabel := box.Objects[1].(*widget.Label)
				detailsLabel.SetText(fmt.Sprintf("IP: %s | MAC: %s", device.IPAddress, device.MACAddress))
			}
		},
	)

	// Refresh button
	refreshBtn := widget.NewButton("Refresh", func() {
		a.refreshDevicesList(devicesList)
	})

	buttons := container.NewHBox(refreshBtn)

	// Layout
	content := container.NewBorder(
		buttons,
		nil,
		nil,
		nil,
		devicesList,
	)

	// Create dialog
	devicesDialog := dialog.NewCustom("Connected Devices", "Close", content, a.MainWindow)
	devicesDialog.Resize(fyne.NewSize(500, 400))

	// Load initial devices list
	a.refreshDevicesList(devicesList)

	devicesDialog.Show()
}

// refreshDevicesList fetches and refreshes the connected devices list
func (a *App) refreshDevicesList(list *widget.List) {
	devices, err := a.APIClient.GetConnectedDevices()
	if err != nil {
		a.Logger.Errorf("Failed to fetch connected devices: %v", err)
		dialog.ShowError(fmt.Errorf("Failed to load connected devices: %v", err), a.MainWindow)
		return
	}

	a.cachedDevices = devices
	list.Refresh()
}
