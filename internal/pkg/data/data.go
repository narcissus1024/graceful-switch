package data

import (
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
)

var (
	SSHData = &Data{
		ContentList:      &ContentList{contents: make(map[string]Content)},
		ContentIndexList: &ContentIndexList{contentIndexes: []ContentIndex{}},
	}
)

type Data struct {
	ContentList      *ContentList
	ContentIndexList *ContentIndexList
}

func (d *Data) Load(config *config.Config) error {
	if err := d.ContentList.load(config); err != nil {
		return err
	}
	if err := d.ContentIndexList.load(config); err != nil {
		return err
	}

	return nil
}

func (d *Data) PersistSSHConfig() {

}