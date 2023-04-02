package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/narcissus1024/graceful-switch/internal/pkg/controller"
)

func (v *ViewHandler) AddConfigBtnEvent() func() {
	obj := v.DrawAddConfigView()
	return func() {
		obj.Show()
	}
}

func (v *ViewHandler) DrawAddConfigView() dialog.Dialog {
	result := &fyne.Container{}
	// ssh type
	sshType := widget.NewRadioGroup([]string{controller.SIMPLE.String(), controller.COMBINE.String()}, func(str string) {

	})
	sshType.Horizontal = true
	sshType.SetSelected(controller.SIMPLE.String())
	sshTypeContainer := container.NewVBox(widget.NewLabel("ssh type options"), sshType)

	// ssh title
	titleEntry := widget.NewEntry()
	titleContainer := container.NewVBox(widget.NewLabel("ssh title"), titleEntry)

	// result
	result = container.NewVBox(sshTypeContainer, titleContainer)
	addPanel := dialog.NewCustomConfirm("Add SSH Config", "save", "cancel", result, func(confirm bool) {
		if confirm {
			if err := v.SSHController.AddDataIndex(sshType.Selected, titleEntry.Text); err != nil {
				fmt.Printf("add ssh config err: %v\n", err)
			} else {
				titleEntry.Text = ""
				v.DataIndexCanvas.Refresh()
			}

		}
	}, v.Window)
	return addPanel
}
