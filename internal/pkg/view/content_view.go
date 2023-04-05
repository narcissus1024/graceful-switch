package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
)

func (v *ViewManager) InitContentView() fyne.CanvasObject {
	text := widget.NewMultiLineEntry()
	v.DataContentCanvas = text

	// button
	saveBtn := widget.NewButton("save", v.SaveConfigBtnEvent())
	cleanBtn := widget.NewButton("clean", v.CleanConfigBtnEvent())

	return container.NewBorder(nil, container.NewHBox(saveBtn, cleanBtn), nil, nil, text)
}

// SaveConfigBtnEvent save ssh config in the multiLineEntry
func (v *ViewManager) SaveConfigBtnEvent() func() {
	return func() {
		text := v.DataContentCanvas.Text
		if len(strings.TrimSpace(text)) == 0 || len(v.SelectedId) == 0 {
			return
		}
		if err := v.DataManager.AddData(v.SelectedId, v.DataContentCanvas.Text); err != nil {
			log.Println(err)
		} else {
			dialog.NewInformation("save ssh config", "save success", v.Window).Show()
		}
	}
}

func (v *ViewManager) CleanConfigBtnEvent() func() {
	return func() {
		v.DataContentCanvas.Text = ""
		v.DataContentCanvas.Refresh()
	}
}
