package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/narcissus1024/graceful-switch/internal/pkg/controller"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
	"github.com/narcissus1024/graceful-switch/tools/constant"
)

type ViewHandler struct {
	App               fyne.App
	Window            fyne.Window
	DataIndexCanvas   *widget.List
	SelectedId string
	DataContentCanvas *widget.Entry
	SSHController     controller.SSHController
	Data              *data.Data
}

func (v *ViewHandler) Init() {
	v.Data = data.SSHData
	v.App = app.New()
	v.Window = v.App.NewWindow(constant.TITLE)
	v.Window.Resize(fyne.NewSize(constant.WIDTH, constant.HEIGHT))
	v.Window.CenterOnScreen()
	v.InitView()
}

func (v *ViewHandler) Start() {
	v.Window.ShowAndRun()
}

func (v *ViewHandler) InitView() {
	topToolBar := v.InitTopToolBar()
	mainView := v.InitMainView()

	v.Window.SetContent(container.NewBorder(topToolBar, nil, nil, nil, mainView))
}

func (v *ViewHandler) InitTopToolBar() fyne.CanvasObject {
	leftToolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), v.AddConfigBtnEvent()),
	)
	rightToolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fmt.Println("setting")
		}),
	)
	topToolContainer := container.NewBorder(nil, nil, leftToolBar, rightToolBar, widget.NewLabel("system ssh"))
	return topToolContainer
}

func (v *ViewHandler) InitMainView() fyne.CanvasObject {
	// sidebar
	sideBar := widget.NewList(
		func() int {
			return len(v.Data.ContentIndexList.Get())
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), layout.NewSpacer(), widget.NewLabel(""))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			c := object.(*fyne.Container)
			title := c.Objects[0].(*widget.Label)
			switchBtn := c.Objects[2].(*widget.Label)

			title.SetText(v.Data.ContentIndexList.Get()[id].Title)
			open := "off"
			if v.Data.ContentIndexList.Get()[id].Open {
				open = "on"
			}
			switchBtn.SetText(open)
		})
	// set selected handler
	sideBar.OnSelected = func(id widget.ListItemID) {
		itemId := v.Data.ContentIndexList.Get()[id].Id
		v.SelectedId = itemId
		content := v.Data.ContentList.Get(itemId).Content
		v.DataContentCanvas.SetText(content)
	}
	sideBar.ScrollToTop()
	v.DataIndexCanvas = sideBar

	// ssh config content
	text := widget.NewMultiLineEntry()
	v.DataContentCanvas = text

	// button
	saveBtn := widget.NewButton("save", v.SaveConfigBtnEvent())
	cleanBtn := widget.NewButton("clean", v.CleanConfigBtnEvent())

	textContainer := container.NewBorder(nil, container.NewHBox(saveBtn, cleanBtn), nil, nil, text)

	result := container.NewHSplit(v.DataIndexCanvas, textContainer)
	result.Offset = 0.2

	return result
}
