package ui

import (
	_ "embed"

	"fyne.io/fyne/v2/theme"
	"fyne.io/systray"
)

//go:embed icons/tray.png
var trayIconPNG []byte

// GetAppIcon returns the embedded tray icon as a Fyne resource
func GetAppIcon() *AppIconResource {
	return &AppIconResource{data: trayIconPNG}
}

// AppIconResource wraps the raw PNG bytes for use with Fyne
type AppIconResource struct {
	data []byte
}

func (a *AppIconResource) Name() string {
	return "appicon"
}

func (a *AppIconResource) Content() []byte {
	return a.data
}

// StartSystemTray sets up a simple system tray with Show and Quit menu items.
// It blocks until systray.Quit is called. onReady is invoked once the tray is visible.
func (a *App) StartSystemTray(onReady func()) {
	systray.Run(func() {
		iconBytes := trayIconPNG
		if len(iconBytes) == 0 {
			iconBytes = theme.FyneLogo().Content()
		}
		systray.SetIcon(iconBytes)
		systray.SetTooltip("MiFiMate")

		mShow := systray.AddMenuItem("Show MiFiMate", "Show the main window")
		mQuit := systray.AddMenuItem("Quit", "Quit MiFiMate")

		if onReady != nil {
			onReady()
		}

		for {
			select {
			case <-mShow.ClickedCh:
				select {
				case a.trayActions <- "show":
				default:
				}
			case <-mQuit.ClickedCh:
				select {
				case a.trayActions <- "quit":
				default:
				}
				systray.Quit()
				return
			}
		}
	}, func() {
		// onExit
	})
}
