package view

import (
	"fyne.io/fyne/v2/dialog"
)

func (v *ViewManager) NewDialog(message string) dialog.Dialog {
	return dialog.NewInformation("information", message, v.Window)
}
