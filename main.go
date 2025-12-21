package main

import (
	"mifi_app/internal/api"
	"mifi_app/internal/config"
	"mifi_app/internal/ui"
	"mifi_app/internal/utils"

	"fyne.io/fyne/v2/app"
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

	mifiApp := ui.NewApp(fyneApp, apiClient, cfg, logger)
	mifiApp.CreateMainWindow()

	mifiApp.AutoLogin()

	logger.Info("Showing main window")
	mifiApp.MainWindow.ShowAndRun()

	logger.Info("Application closed")
}
