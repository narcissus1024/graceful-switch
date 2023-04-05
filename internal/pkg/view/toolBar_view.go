package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
	"log"
)

func (v *ViewManager) InitTopToolBar() fyne.CanvasObject {
	leftToolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), v.AddConfigBtnEvent()),
	)
	rightToolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			log.Println("setting")
		}),
	)
	title := binding.NewString()
	title.Set(v.SelectedTitle)
	topToolContainer := container.NewBorder(nil, nil, leftToolBar, rightToolBar, widget.NewLabelWithData(title))
	return topToolContainer
}

func (v *ViewManager) DrawAddConfigView() dialog.Dialog {
	result := &fyne.Container{}
	// ssh type
	sshType := widget.NewRadioGroup([]string{data.SIMPLE.String(), data.COMBINE.String()}, func(str string) {

	})
	sshType.Horizontal = true
	sshType.SetSelected(data.SIMPLE.String())
	sshTypeContainer := container.NewVBox(widget.NewLabel("ssh type options"), sshType)

	// ssh title
	titleEntry := widget.NewEntry()
	titleContainer := container.NewVBox(widget.NewLabel("ssh title"), titleEntry)

	// result
	result = container.NewVBox(sshTypeContainer, titleContainer)
	addPanel := dialog.NewCustomConfirm("Add SSH Config", "save", "cancel", result, func(confirm bool) {
		if confirm {
			if err := v.DataManager.AddDataIndex(sshType.Selected, titleEntry.Text); err != nil {
				log.Printf("add ssh config err: %v\n", err)
			} else {
				titleEntry.Text = ""
				v.DataIndexCanvas.Refresh()
			}

		}
	}, v.Window)
	return addPanel
}

// AddConfigBtnEvent add ssh config
func (v *ViewManager) AddConfigBtnEvent() func() {
	obj := v.DrawAddConfigView()
	return func() {
		obj.Show()
	}
}
