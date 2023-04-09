package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
	"github.com/narcissus1024/graceful-switch/tools/constant"
)

type ViewManager struct {
	App                fyne.App
	Window             fyne.Window
	DataIndexCanvas    *widget.List
	DataContentCanvas  *widget.Entry
	ContentToolCanvas  *fyne.Container
	SelectedListItemId widget.ListItemID
	SelectedId         string
	SelectedTitle      string
	//SSHManager        *controller.SSHManager
	DataManager *data.DataManager
}

func NewViewManager(opts ...Option) *ViewManager {
	viewManager := new(ViewManager)
	for _, opt := range opts {
		opt(viewManager)
	}
	return viewManager
}

func (v *ViewManager) Init() {
	v.App = app.New()
	v.Window = v.App.NewWindow(constant.TITLE)
	v.Window.Resize(fyne.NewSize(constant.WIDTH, constant.HEIGHT))
	v.Window.CenterOnScreen()
	v.InitView()
}

func (v *ViewManager) Start() {
	v.Window.ShowAndRun()
}

func (v *ViewManager) InitView() {
	topToolBar := v.InitTopToolBar()
	mainView := v.InitMainView()

	v.Window.SetContent(container.NewBorder(topToolBar, nil, nil, nil, mainView))
}

// InitMainView InitContentView must be before InitSideBarView,
// because sideBarBarView need to show ssh config for default select listItem
func (v *ViewManager) InitMainView() fyne.CanvasObject {
	// ssh config content
	contentView := v.InitContentView()

	// sidebar
	v.DataIndexCanvas = v.InitSideBarView()

	result := container.NewHSplit(v.DataIndexCanvas, contentView)
	result.SetOffset(0.2)

	return result
}

type Option func(manager *ViewManager)

//func WithSSHManager(sshManager *controller.SSHManager) Option {
//	return func(manager *ViewManager) {
//		manager.SSHManager = sshManager
//	}
//}

func WithDataManager(dataManager *data.DataManager) Option {
	return func(manager *ViewManager) {
		manager.DataManager = dataManager
	}
}
