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
	contentToolBar := container.NewHBox(saveBtn, cleanBtn)
	v.ContentToolCanvas = contentToolBar

	return container.NewBorder(nil, contentToolBar, nil, nil, text)
}

// SaveConfigBtnEvent save ssh config in the multiLineEntry
func (v *ViewManager) SaveConfigBtnEvent() func() {
	return func() {
		log.Printf("[save]id:%s,title:%s,content:%s\n", v.SelectedId, v.SelectedTitle, v.DataContentCanvas.Text)
		text := v.DataContentCanvas.Text
		if len(strings.TrimSpace(text)) == 0 || len(v.SelectedId) == 0 {
			v.NewDialog("text does not is empty").Show()
			return
		}
		if err := v.DataManager.AddData(v.SelectedId, text); err != nil {
			log.Println(err)
		} else {
			dialog.NewInformation("save ssh config", "save success", v.Window).Show()
		}

		// refresh
		v.DataContentCanvas.Refresh()
	}
}

func (v *ViewManager) CleanConfigBtnEvent() func() {
	return func() {
		v.DataContentCanvas.SetText("")
		//v.DataContentCanvas.Refresh()
	}
}
