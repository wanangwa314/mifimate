package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (a *App) ShowSMSDialog() {
	messageDetail := widget.NewRichTextFromMarkdown("")
	messageDetail.Wrapping = fyne.TextWrapWord
	detailScroll := container.NewScroll(messageDetail)

	smsList := widget.NewList(
		func() int {
			return len(a.cachedSMSMessages)
		},
		func() fyne.CanvasObject {
			header := widget.NewLabel("Sender Name")
			header.TextStyle.Bold = true

			date := widget.NewLabel("Date")
			date.TextStyle.Italic = true

			preview := widget.NewLabel("Message preview...")
			preview.Wrapping = fyne.TextTruncate

			return container.NewVBox(
				header,
				date,
				preview,
				widget.NewSeparator(),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(a.cachedSMSMessages) {
				msg := a.cachedSMSMessages[id]
				box := obj.(*fyne.Container)

				headerLabel := box.Objects[0].(*widget.Label)
				headerLabel.SetText(fmt.Sprintf("From: %s", msg.Number))

				dateLabel := box.Objects[1].(*widget.Label)
				dateLabel.SetText(msg.Timestamp.Format("Mon, Jan 2, 2006 at 3:04 PM"))

				previewLabel := box.Objects[2].(*widget.Label)
				preview := msg.Content
				if len(preview) > 80 {
					preview = preview[:80] + "..."
				}
				previewLabel.SetText(preview)
			}
		},
	)

	smsList.OnSelected = func(id widget.ListItemID) {
		if id < len(a.cachedSMSMessages) {
			msg := a.cachedSMSMessages[id]

			markdown := fmt.Sprintf("**From:** %s\n\n**Date:** %s\n\n**Message:**\n\n%s",
				msg.Number,
				msg.Timestamp.Format("Monday, January 2, 2006 at 3:04 PM"),
				msg.Content,
			)

			messageDetail.ParseMarkdown(markdown)
			detailScroll.ScrollToTop()
		}
	}

	refreshBtn := widget.NewButton("Refresh", func() {
		a.refreshSMSList(smsList)
		messageDetail.ParseMarkdown("*Select a message to view its content*")
	})

	buttons := container.NewHBox(
		refreshBtn,
		layout.NewSpacer(),
		widget.NewLabel(fmt.Sprintf("Total: %d messages", len(a.cachedSMSMessages))),
	)

	split := container.NewHSplit(
		container.NewBorder(nil, nil, nil, nil, smsList),
		container.NewBorder(nil, nil, nil, nil, detailScroll),
	)
	split.Offset = 0.4 // 40% for list, 60% for detail

	content := container.NewBorder(
		buttons,
		nil,
		nil,
		nil,
		split,
	)

	smsDialog := dialog.NewCustom("SMS Messages", "Close", content, a.MainWindow)
	smsDialog.Resize(fyne.NewSize(900, 600))

	a.refreshSMSList(smsList)

	messageDetail.ParseMarkdown("*Select a message from the list to view its content*")

	smsDialog.Show()
}

func (a *App) refreshSMSList(list *widget.List) {
	messages, err := a.APIClient.GetSMSList(0, 50)
	if err != nil {
		a.Logger.Errorf("Failed to fetch SMS list: %v", err)
		dialog.ShowError(fmt.Errorf("Failed to load SMS messages: %v", err), a.MainWindow)
		return
	}

	a.cachedSMSMessages = messages
	list.Refresh()
}

func (a *App) checkForNewSMS() {
	count, err := a.APIClient.GetSMSCount()
	if err != nil {
		a.Logger.Errorf("Failed to check SMS count: %v", err)
		return
	}

	if a.lastSMSCount == 0 || count != a.lastSMSCount {
		messages, err := a.APIClient.GetSMSList(0, 50)
		if err != nil {
			a.Logger.Errorf("Failed to fetch SMS list: %v", err)
		} else {
			a.cachedSMSMessages = messages
			a.updateRecentSMSContent()
		}
	}

	if a.lastSMSCount > 0 && count > a.lastSMSCount {
		newCount := count - a.lastSMSCount
		a.FyneApp.SendNotification(&fyne.Notification{
			Title:   "New SMS",
			Content: fmt.Sprintf("You have %d new message(s)", newCount),
		})
	}

	a.lastSMSCount = count
}
