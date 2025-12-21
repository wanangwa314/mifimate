package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"

	"mifi_app/internal/api"
	"mifi_app/internal/config"
	"mifi_app/internal/utils"
)

type App struct {
	FyneApp    fyne.App
	MainWindow fyne.Window
	APIClient  *api.Client
	Config     *config.Config
	Logger     *logrus.Logger

	statusLabel     *widget.Label
	connectBtn      *widget.Button
	disconnectBtn   *widget.Button
	refreshBtn      *widget.Button
	wifiSettingsBtn *widget.Button
	smsBtn          *widget.Button
	devicesBtn      *widget.Button
	settingsBtn     *widget.Button
	restartBtn      *widget.Button
	shutdownBtn     *widget.Button

	networkTypeLabel  *widget.Label
	signalLabel       *widget.Label
	batteryLabel      *widget.Label
	ipAddressLabel    *widget.Label
	connectedDevLabel *widget.Label
	txSpeedLabel      *widget.Label
	rxSpeedLabel      *widget.Label

	pollingTicker *time.Ticker
	stopPolling   chan bool

	cachedSMSMessages  []api.SMSMessage
	lastSMSCount       int
	recentSMSContainer *fyne.Container

	cachedDevices []api.ConnectedDevice
}

func NewApp(fyneApp fyne.App, client *api.Client, cfg *config.Config, logger *logrus.Logger) *App {
	return &App{
		FyneApp:   fyneApp,
		APIClient: client,
		Config:    cfg,
		Logger:    logger,
	}
}

func (a *App) CreateMainWindow() {
	a.MainWindow = a.FyneApp.NewWindow("MiFiMate")
	a.MainWindow.Resize(fyne.NewSize(800, 600))

	a.createComponents()

	content := a.createLayout()

	a.MainWindow.SetContent(content)
}

func (a *App) AutoLogin() {
	a.statusLabel.SetText("Status: Connecting...")

	// Check if device is reachable
	if err := a.APIClient.Ping(); err != nil {
		a.Logger.Errorf("Device unreachable during auto-login: %v", err)
		a.statusLabel.SetText("Status: Device Unreachable")
		return
	}

	// Attempt login
	username := a.Config.Device.Username
	password := a.Config.Device.Password

	if password == "" {
		a.statusLabel.SetText("Status: No Password Configured")
		return
	}

	if err := a.APIClient.Login(username, password); err != nil {
		a.Logger.Errorf("Auto-login failed: %v", err)
		a.statusLabel.SetText("Status: Login Failed")
		return
	}

	a.statusLabel.SetText("Status: Connected")
	a.connectBtn.Disable()
	a.disconnectBtn.Enable()

	// Fetch initial status
	a.onRefresh()

	// Start auto-refresh polling
	a.startPolling()
}

func (a *App) createComponents() {
	a.statusLabel = widget.NewLabel("Status: Disconnected")
	a.statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	a.connectBtn = widget.NewButton("Connect", a.onConnect)
	a.connectBtn.Importance = widget.HighImportance

	a.disconnectBtn = widget.NewButton("Disconnect", a.onDisconnect)
	a.disconnectBtn.Importance = widget.MediumImportance
	a.disconnectBtn.Disable()

	a.refreshBtn = widget.NewButton("Refresh", a.onRefresh)

	a.wifiSettingsBtn = widget.NewButton("WiFi Settings", a.ShowWiFiSettingsDialog)
	a.smsBtn = widget.NewButton("SMS Messages", a.ShowSMSDialog)
	a.devicesBtn = widget.NewButton("Connected Devices", a.ShowDevicesDialog)
	a.settingsBtn = widget.NewButton("Settings", a.ShowSettingsDialog)

	a.restartBtn = widget.NewButton("Restart Device", a.onRestart)
	a.restartBtn.Importance = widget.WarningImportance

	a.shutdownBtn = widget.NewButton("Shutdown Device", a.onShutdown)
	a.shutdownBtn.Importance = widget.DangerImportance

	a.networkTypeLabel = widget.NewLabel("N/A")
	a.signalLabel = widget.NewLabel("N/A")
	a.batteryLabel = widget.NewLabel("N/A")
	a.ipAddressLabel = widget.NewLabel("N/A")
	a.connectedDevLabel = widget.NewLabel("0")
	a.txSpeedLabel = widget.NewLabel("0 B/s")
	a.rxSpeedLabel = widget.NewLabel("0 B/s")
}

func (a *App) createLayout() fyne.CanvasObject {
	statusContent := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Status:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			a.statusLabel,
		),
		widget.NewSeparator(),
		container.NewHBox(
			a.connectBtn,
			a.disconnectBtn,
			a.refreshBtn,
		),
	)
	statusCard := a.createCard("Connection Status", statusContent, theme.ComputerIcon())

	deviceInfoGrid := container.New(layout.NewFormLayout(),
		widget.NewLabelWithStyle("Network:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.networkTypeLabel,
		widget.NewLabelWithStyle("Signal:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.signalLabel,
		widget.NewLabelWithStyle("Battery:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.batteryLabel,
		widget.NewLabelWithStyle("IP Address:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.ipAddressLabel,
		widget.NewLabelWithStyle("Devices:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		a.connectedDevLabel,
	)
	deviceInfoCard := a.createCard("Device Information", deviceInfoGrid, theme.InfoIcon())

	// Network Statistics Card with icon indicator
	uploadRow := container.NewHBox(
		widget.NewLabelWithStyle("Upload:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		a.txSpeedLabel,
	)
	downloadRow := container.NewHBox(
		widget.NewLabelWithStyle("Download:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		a.rxSpeedLabel,
	)
	networkStatsContent := container.NewVBox(uploadRow, downloadRow)
	networkStatsCard := a.createCard("Network Statistics", networkStatsContent, theme.NavigateNextIcon())

	quickActionsGrid := container.NewGridWithColumns(2,
		a.wifiSettingsBtn,
		a.smsBtn,
		a.devicesBtn,
		a.settingsBtn,
	)
	quickActionsCard := a.createCard("Quick Actions", quickActionsGrid, nil)

	powerGrid := container.NewGridWithColumns(2,
		a.restartBtn,
		a.shutdownBtn,
	)
	powerCard := a.createCard("Power Management", powerGrid, nil)

	a.recentSMSContainer = container.NewVBox()
	a.updateRecentSMSContent()
	recentSMSCard := a.createCard("Recent Messages", a.recentSMSContainer, theme.MailComposeIcon())

	// Layout cards in two columns with equal sizing
	leftColumn := container.NewVBox(
		statusCard,
		deviceInfoCard,
		recentSMSCard,
	)

	rightColumn := container.NewVBox(
		networkStatsCard,
		quickActionsCard,
		powerCard,
	)

	// Use GridWrap for better space distribution
	mainContent := container.NewGridWithColumns(2,
		leftColumn,
		rightColumn,
	)

	// Return with padding (no scroll needed)
	return container.NewPadded(mainContent)
}

func (a *App) createCard(title string, content fyne.CanvasObject, icon fyne.Resource) fyne.CanvasObject {
	var titleWidget fyne.CanvasObject
	if icon != nil {
		iconImg := canvas.NewImageFromResource(icon)
		iconImg.SetMinSize(fyne.NewSize(20, 20))
		titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		titleWidget = container.NewHBox(iconImg, titleLabel)
	} else {
		titleWidget = widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	}

	cardBg := canvas.NewRectangle(color.NRGBA{R: 45, G: 45, B: 55, A: 255})
	cardBg.CornerRadius = 8

	cardContent := container.NewVBox(
		titleWidget,
		widget.NewSeparator(),
		content,
	)

	paddedContent := container.NewPadded(cardContent)

	return container.NewStack(cardBg, paddedContent)
}

func (a *App) updateRecentSMSContent() {
	if a.recentSMSContainer == nil {
		return
	}

	a.recentSMSContainer.Objects = nil

	if len(a.cachedSMSMessages) == 0 {
		a.recentSMSContainer.Add(widget.NewLabel("No messages"))
		a.recentSMSContainer.Refresh()
		return
	}

	count := 3
	if len(a.cachedSMSMessages) < count {
		count = len(a.cachedSMSMessages)
	}

	for i := 0; i < count; i++ {
		msg := a.cachedSMSMessages[i]

		content := msg.Content
		if len(content) > 50 {
			content = content[:100] + "..."
		}

		senderLabel := widget.NewLabelWithStyle(msg.Number, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		contentLabel := widget.NewLabel(content)
		dateLabel := widget.NewLabelWithStyle(msg.Timestamp.Format("Jan 02 15:04"), fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

		messageBox := container.NewVBox(
			container.NewHBox(senderLabel, layout.NewSpacer(), dateLabel),
			contentLabel,
		)

		a.recentSMSContainer.Add(messageBox)

		if i < count-1 {
			a.recentSMSContainer.Add(widget.NewSeparator())
		}
	}

	a.recentSMSContainer.Refresh()
}

func (a *App) onConnect() {
	a.connectBtn.Disable()
	defer func() {
		if !a.isConnected() {
			a.connectBtn.Enable()
		}
	}()

	if err := a.APIClient.Ping(); err != nil {
		a.Logger.Errorf("Device unreachable: %v", err)
		a.statusLabel.SetText("Status: Device Unreachable")
		dialog.ShowError(fmt.Errorf("Cannot reach device at %s", a.Config.Device.DefaultIP), a.MainWindow)
		return
	}

	username := a.Config.Device.Username
	password := a.Config.Device.Password

	if password == "" {
		a.statusLabel.SetText("Status: Password Required")
		return
	}

	if err := a.APIClient.Login(username, password); err != nil {
		a.Logger.Errorf("Login failed: %v", err)
		a.statusLabel.SetText("Status: Login Failed")
		dialog.ShowError(fmt.Errorf("Login failed: %v", err), a.MainWindow)
		return
	}

	a.statusLabel.SetText("Status: Connected")
	a.connectBtn.Disable()
	a.disconnectBtn.Enable()

	a.onRefresh()
	a.startPolling()
}

func (a *App) onDisconnect() {
	a.stopPollingNow()

	if err := a.APIClient.Logout(); err != nil {
		a.Logger.Errorf("Logout failed: %v", err)
	}

	a.statusLabel.SetText("Status: Disconnected")
	a.connectBtn.Enable()
	a.disconnectBtn.Disable()

	a.resetLabels()
}

func (a *App) onRefresh() {
	if !a.isConnected() {
		return
	}

	status, err := a.APIClient.GetDeviceStatus()
	if err != nil {
		a.Logger.Errorf("Failed to get device status: %v", err)
		a.statusLabel.SetText("Status: Error fetching data")
		dialog.ShowError(fmt.Errorf("Failed to fetch device status: %v", err), a.MainWindow)
		return
	}

	a.updateStatus(status)
}

func (a *App) updateStatus(status *api.DeviceStatus) {
	a.networkTypeLabel.SetText(status.NetworkType)

	signalQuality := utils.GetSignalQuality(status.SignalStrength)
	a.signalLabel.SetText(fmt.Sprintf("%d bars (%s)", status.SignalStrength, signalQuality))

	batteryStatus := utils.GetBatteryStatus(status.BatteryLevel)
	a.batteryLabel.SetText(fmt.Sprintf("%d%% (%s)", status.BatteryLevel, batteryStatus))

	a.ipAddressLabel.SetText(status.WanIPAddress)
	a.connectedDevLabel.SetText(strconv.Itoa(status.ConnectedDevs))
	a.txSpeedLabel.SetText(utils.FormatSpeed(status.TxSpeed))
	a.rxSpeedLabel.SetText(utils.FormatSpeed(status.RxSpeed))
}

func (a *App) updateStatusSafe(status *api.DeviceStatus) {
	// Queue all UI updates on the main thread
	go func() {
		fyne.CurrentApp().Driver().CanvasForObject(a.statusLabel).Refresh(a.statusLabel)
	}()

	a.statusLabel.SetText("Status: Connected")
	a.networkTypeLabel.SetText(status.NetworkType)

	signalQuality := utils.GetSignalQuality(status.SignalStrength)
	a.signalLabel.SetText(fmt.Sprintf("%d bars (%s)", status.SignalStrength, signalQuality))

	batteryStatus := utils.GetBatteryStatus(status.BatteryLevel)
	a.batteryLabel.SetText(fmt.Sprintf("%d%% (%s)", status.BatteryLevel, batteryStatus))

	a.ipAddressLabel.SetText(status.WanIPAddress)
	a.connectedDevLabel.SetText(strconv.Itoa(status.ConnectedDevs))
	a.txSpeedLabel.SetText(utils.FormatSpeed(status.TxSpeed))
	a.rxSpeedLabel.SetText(utils.FormatSpeed(status.RxSpeed))
}

func (a *App) resetLabels() {
	a.networkTypeLabel.SetText("N/A")
	a.signalLabel.SetText("N/A")
	a.batteryLabel.SetText("N/A")
	a.ipAddressLabel.SetText("N/A")
	a.connectedDevLabel.SetText("0")
	a.txSpeedLabel.SetText("0 B/s")
	a.rxSpeedLabel.SetText("0 B/s")
}

func (a *App) isConnected() bool {
	return a.APIClient.IsAuthenticated()
}

// startPolling starts automatic status refresh every 3 seconds
func (a *App) startPolling() {
	if a.pollingTicker != nil {
		return
	}

	a.stopPolling = make(chan bool)
	a.pollingTicker = time.NewTicker(3 * time.Second)

	go func() {
		for {
			select {
			case <-a.pollingTicker.C:
				if a.isConnected() {
					status, err := a.APIClient.GetDeviceStatus()
					if err != nil {
						a.Logger.Errorf("Failed to get device status: %v", err)
						continue
					}
					a.updateStatusSafe(status)
					a.checkForNewSMS()
				} else {
					a.stopPolling <- true
				}
			case <-a.stopPolling:
				a.pollingTicker.Stop()
				a.pollingTicker = nil
				return
			}
		}
	}()
}

// stopPollingNow stops the auto-refresh polling
func (a *App) stopPollingNow() {
	if a.pollingTicker != nil {
		a.stopPolling <- true
	}
}

func (a *App) onRestart() {
	if !a.isConnected() {
		dialog.ShowError(fmt.Errorf("not connected to device"), a.MainWindow)
		return
	}

	// Show confirmation dialog
	confirmDialog := dialog.NewConfirm(
		"Restart Device",
		"Are you sure you want to restart the MiFi device? This will disconnect all connected clients temporarily.",
		func(confirmed bool) {
			if confirmed {
				a.stopPollingNow()

				err := a.APIClient.RebootDevice()
				if err != nil {
					a.Logger.Errorf("Failed to restart device: %v", err)
					dialog.ShowError(fmt.Errorf("failed to restart device: %v", err), a.MainWindow)
					return
				}

				a.statusLabel.SetText("Device Restarting...")

				dialog.ShowInformation(
					"Restart Initiated",
					"The device is restarting. Please wait about 30 seconds before reconnecting.",
					a.MainWindow,
				)

				a.disconnectBtn.Disable()
				a.connectBtn.Enable()
			}
		},
		a.MainWindow,
	)

	confirmDialog.Show()
}

func (a *App) onShutdown() {
	if !a.isConnected() {
		dialog.ShowError(fmt.Errorf("not connected to device"), a.MainWindow)
		return
	}

	// Show confirmation dialog
	confirmDialog := dialog.NewConfirm(
		"Shutdown Device",
		"Are you sure you want to shutdown the MiFi device? You will need to manually power it on again.",
		func(confirmed bool) {
			if confirmed {
				a.stopPollingNow()

				err := a.APIClient.ShutdownDevice()
				if err != nil {
					a.Logger.Errorf("Failed to shutdown device: %v", err)
					dialog.ShowError(fmt.Errorf("failed to shutdown device: %v", err), a.MainWindow)
					return
				}

				a.statusLabel.SetText("Status: Device Shutting Down...")

				dialog.ShowInformation(
					"Shutdown Initiated",
					"The device is shutting down. You will need to manually power it on again.",
					a.MainWindow,
				)

				a.disconnectBtn.Disable()
				a.connectBtn.Enable()
			}
		},
		a.MainWindow,
	)

	confirmDialog.Show()
}
