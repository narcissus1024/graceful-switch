package data

import (
	"encoding/json"
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"io/ioutil"
	"os"
)

type ContentIndex struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Open  bool   `json:"open"`
}

type ContentIndexList struct {
	contentIndexes []ContentIndex
}

func (d *ContentIndexList) load(config *config.Config) error {
	indexPath := config.GetDataIndexPath()
	contentIndexes := []ContentIndex{}
	file, err := os.OpenFile(indexPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	indexData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(indexData) > 0 {
		if err := json.Unmarshal(indexData, &contentIndexes); err != nil {
			return err
		}
	}
	d.contentIndexes = contentIndexes
	return nil
}

func (d *ContentIndexList) Persist() error {
	indexPath := config.Conf.GetDataIndexPath()
	indexBytes, err := json.Marshal(d.contentIndexes)
	if err != nil {
		return err
	}
	if err := os.WriteFile(indexPath, indexBytes, 0666); err != nil {
		return err
	}
	return nil
}

func (d *ContentIndexList) Append(index ContentIndex) {
	d.contentIndexes = append(d.contentIndexes, index)
}

func (d *ContentIndexList) Get() []ContentIndex {
	return d.contentIndexes
}

