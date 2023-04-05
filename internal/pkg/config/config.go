package config

import (
	"encoding/json"
	"github.com/narcissus1024/graceful-switch/tools"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"sync"
)

var (
	homeDir string
	config  *Config
	once    sync.Once
)

func init() {
	u, _ := user.Current()
	homeDir = u.HomeDir
}

func GetConfig() *Config {
	once.Do(func() {
		config = new(Config)
		config.StoreRoot = path.Join(homeDir, ".graceful-switch")
		config.DataDirRoot = path.Join(config.StoreRoot, "data")
		config.SSHConfigPath = path.Join(homeDir, ".ssh", "config")
		config.ConfigPath = path.Join(config.StoreRoot, "config.json")
	})
	return config
}

type Config struct {
	// store root path, contain config and data file
	StoreRoot string `json:"storeRoot"`
	// data store root path
	DataDirRoot   string `json:"dataDirRoot"`
	SSHConfigPath string `json:"sshConfigPath"`
	ConfigPath    string `json:"-"`
}

func (c *Config) Load() error {
	if tools.IsExist(c.ConfigPath) {
		if configContent, err := ioutil.ReadFile(c.ConfigPath); err != nil {
			return err
		} else if err := json.Unmarshal(configContent, c); err != nil {
			return err
		}
	} else {
		defaultConfig, err := json.Marshal(c)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(path.Dir(c.ConfigPath), 0777); err != nil {
			return err
		}
		if err := os.WriteFile(c.ConfigPath, defaultConfig, 0777); err != nil {
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
