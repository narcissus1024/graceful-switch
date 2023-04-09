package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
	"log"
)

func (v *ViewManager) InitSideBarView() *widget.List {
	// sidebar
	sideBar := widget.NewList(
		func() int {
			return len(v.DataManager.InnerDataIndexList.GetAll())
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""),
				layout.NewSpacer(),
				widget.NewCheck("", func(checked bool) {}),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
					dialog.NewConfirm("", "", func(confirm bool) {

					}, v.Window).Show()
				}))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			c := object.(*fyne.Container)
			title := c.Objects[0].(*widget.Label)
			switchBtn := c.Objects[2].(*widget.Check)

			index := v.DataManager.InnerDataIndexList.GetAll()[id]
			c.Objects[3] = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				dialog.NewConfirm("Delete", fmt.Sprintf("Are you confirm to delete the config: %s?", index.Title), func(confirm bool) {
					// todo delete config
				}, v.Window).Show()
			})
			title.SetText(index.Title)
			if index.Id == data.SYSTEM_ID_SSH {
				c.Objects[3].Hide()
				switchBtn.Hide()
			} else {
				if index.Open {
					switchBtn.SetChecked(true)
				}
				switchBtn.OnChanged = v.IndexSwitchChanged(index.Id)
			}
		})
	// set selected handler
	sideBar.OnSelected = v.SideBarSelected()
	sideBar.Select(0)
	return sideBar
}

func (v *ViewManager) IndexSwitchChanged(id string) func(checked bool) {
	return func(checked bool) {
		if err := v.DataManager.InnerDataIndexList.SetOpen(id, checked); err != nil {
			log.Println("[IndexSwitchChanged]", err)
			return
		}
		if err := v.DataManager.PersistSSHConfig(); err != nil {
			log.Println("[IndexSwitchChanged]", err)
			return
		}
		// refresh current selected listItem text
		index := v.DataManager.InnerDataIndexList.GetAll()[v.SelectedListItemId]
		v.DataContentCanvas.SetText(v.DataManager.GetSSHConfigById(index.Id))
	}
}

func (v *ViewManager) SideBarSelected() func(id widget.ListItemID) {
	return func(id widget.ListItemID) {
		index := v.DataManager.InnerDataIndexList.GetAll()[id]
		v.SelectedId = index.Id
		v.SelectedListItemId = id
		v.SelectedTitle = index.Title

		contents := v.DataManager.GetSSHConfigById(v.SelectedId)
		v.DataContentCanvas.SetText(contents)

		if v.SelectedId == data.SYSTEM_ID_SSH {
			v.ContentToolCanvas.Hide()
		} else {
			v.ContentToolCanvas.Show()
		}
		// todo debug
		//log.Printf("select: %s, content: %s\n", v.SelectedTitle, v.DataContentCanvas.Text)
	}
}
