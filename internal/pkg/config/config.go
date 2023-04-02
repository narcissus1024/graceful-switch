package config

import (
	"encoding/json"
	"github.com/narcissus1024/graceful-switch/tools"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

var (
	homeDir    string
	configPath string
	Conf       = new(Config)
)

func init() {
	u, _ := user.Current()
	homeDir = u.HomeDir

	Conf.StoreRoot = path.Join(homeDir, ".graceful-switch")
	Conf.DataDirRoot = path.Join(Conf.StoreRoot, "data")
	configPath = path.Join(Conf.StoreRoot, "config.json")
}

type Config struct {
	StoreRoot   string `json:"storeRoot"`
	DataDirRoot string `json:"dataDirRoot"`
}

func (c *Config) Load() error {
	if tools.IsExist(configPath) {
		if configContent, err := ioutil.ReadFile(configPath); err != nil {
			return err
		} else if err := json.Unmarshal(configContent, Conf); err != nil {
			return err
		}
	} else {
		defaultConfig, err := json.Marshal(Conf)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(path.Dir(configPath), 0777); err != nil {
			return err
		}
		if err := os.WriteFile(configPath, defaultConfig, 0777); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) GetDataPath() string {
	return path.Join(c.DataDirRoot, "data")
}

func (c *Config) GetDataIndexPath() string {
	return path.Join(c.DataDirRoot, "index.json")
}
