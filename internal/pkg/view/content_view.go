package view

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
)

func (v *ViewHandler) SaveConfigBtnEvent() func() {
	return func() {
		if err := v.SSHController.AddData(v.SelectedId, v.DataContentCanvas.Text); err != nil {
			fmt.Println(err)
		} else {
			dialog.NewInformation("save ssh config", "save success", v.Window).Show()
		}
	}
}

func (v *ViewHandler) CleanConfigBtnEvent() func() {
	return func() {
		v.DataContentCanvas.Text = ""
		v.DataContentCanvas.Refresh()
	}
}
