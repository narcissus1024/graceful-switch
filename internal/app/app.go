package app

import (
	"context"
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
	"github.com/narcissus1024/graceful-switch/internal/pkg/view"
)

type SwitchApp struct {
	ViewHandler *view.ViewManager
}

func Run(ctx context.Context) error {
	s := SwitchApp{}
	if err := s.init(); err != nil {
		return err
	}
	s.Start()
	return nil
}

func (s *SwitchApp) init() error {
	// load config
	conf := config.GetConfig()
	if err := conf.Load(); err != nil {
		return err
	}

	// load data
	dataManager := data.GetDataManager()
	if err := dataManager.Init(); err != nil {
		return err
	}

	// necessary?
	//sshManager := controller.NewSSHManager(controller.WithConfig(conf), controller.WithDataManager(dataManager))

	// load view
	viewManager := view.NewViewManager(view.WithDataManager(dataManager))
	viewManager.Init()
	s.ViewHandler = viewManager

	return nil
}

func (s *SwitchApp) Start() {
	s.ViewHandler.Start()
}
