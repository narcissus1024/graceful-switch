package data

import (
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/tools"
	"io/ioutil"
	"log"
	"os"
)

type SSHConfig struct {
	SSHConfigData string
}

func (s *SSHConfig) Load() error {
	sysPath := config.GetConfig().SSHConfigPath
	if !tools.IsExist(sysPath) {
		if f, err := os.Create(sysPath); err != nil {
			return err
		} else {
			f.Close()
		}
	}
	data, err := ioutil.ReadFile(sysPath)
	if err != nil {
		return err
	}
	s.SSHConfigData = string(data)
	log.Println("系统数据：", data)
	return nil
}

func (s *SSHConfig) Persist() error {
	sysPath := config.GetConfig().SSHConfigPath
	if err := os.WriteFile(sysPath, []byte(s.SSHConfigData), 0666); err != nil {
		return err
	}
	return nil
}
