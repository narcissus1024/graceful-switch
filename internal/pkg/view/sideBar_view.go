package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
)

func (v *ViewManager) InitSideBarView() *widget.List {
	// sidebar
	sideBar := widget.NewList(
		func() int {
			return len(v.DataManager.InnerDataIndexList.Get())
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), layout.NewSpacer(), widget.NewLabel(""))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			c := object.(*fyne.Container)
			title := c.Objects[0].(*widget.Label)
			switchBtn := c.Objects[2].(*widget.Label)

			title.SetText(v.DataManager.InnerDataIndexList.Get()[id].Title)
			open := "off"
			if v.DataManager.InnerDataIndexList.Get()[id].Open {
				open = "on"
			}
			switchBtn.SetText(open)
		})
	// set selected handler
	sideBar.OnSelected = func(id widget.ListItemID) {
		itemId := v.DataManager.InnerDataIndexList.Get()[id].Id
		v.SelectedId = itemId
		v.SelectedTitle = v.DataManager.InnerDataIndexList.Get()[id].Title
		contents := v.DataManager.InnerDataList.Get(v.SelectedId).Contents
		if v.SelectedId == data.SYSTEM_ID_SSH {
			contents = v.DataManager.SSHConfig.SSHConfigData
		}

		v.DataContentCanvas.SetText(contents)
	}
	sideBar.ScrollToTop()
	return sideBar
}
