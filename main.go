package main

import (
	"mifi_app/internal/api"
	"mifi_app/internal/config"
	"mifi_app/internal/ui"
	"mifi_app/internal/utils"

	"fyne.io/fyne/v2/app"
	"fyne.io/systray"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	logger := utils.InitLogger(cfg.App.LogLevel)
	logger.Info("Starting MiFiMate")

	apiClient := api.NewClient(
		"http://"+cfg.Device.DefaultIP,
		logger,
	)

	fyneApp := app.New()
	fyneApp.SetIcon(ui.GetAppIcon())

	mifiApp := ui.NewApp(fyneApp, apiClient, cfg, logger)
	mifiApp.CreateMainWindow()

	ready := make(chan struct{})

	// Start system tray in a goroutine so it doesn't block the main thread
	go func() {
		mifiApp.StartSystemTray(func() {
			close(ready)
		})
	}()

	// Wait for tray to be ready before showing window
	<-ready

	mifiApp.AutoLogin()

	logger.Info("Showing main window")
	mifiApp.MainWindow.ShowAndRun()

	logger.Info("Application closed")
	systray.Quit()
}
