package app

import (
	"context"
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/internal/pkg/controller"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
	"github.com/narcissus1024/graceful-switch/internal/pkg/view"
)

type SwitchApp struct {
	ViewHandler view.ViewHandler
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
	if err := config.Conf.Load(); err != nil {
		return err
	}
	if err := data.SSHData.Load(config.Conf); err != nil {
		return err
	}

	sshController := controller.SSHController{
		Config: config.Conf,
		Data:   data.SSHData,
	}

	s.ViewHandler = view.ViewHandler{SSHController: sshController}
	s.ViewHandler.Init()
	return nil
}

func (s *SwitchApp) Start() {
	s.ViewHandler.Start()
}
