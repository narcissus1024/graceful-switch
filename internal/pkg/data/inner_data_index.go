package data

import (
	"encoding/json"
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"io/ioutil"
	"log"
	"os"
)

// InnerDataIndex inner ssh config item index
type InnerDataIndex struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Open  bool   `json:"open"`
}

// InnerDataIndexList inner ssh config index list
type InnerDataIndexList struct {
	contentIndexes []*InnerDataIndex
}

func (d *InnerDataIndexList) Load() error {
	contentIndexes := []*InnerDataIndex{}
	contentIndexes = append(contentIndexes, &InnerDataIndex{
		Id:    SYSTEM_ID_SSH,
		Title: "System SSH Config",
		Open:  true,
	})

	innerIndex := []*InnerDataIndex{}
	file, err := os.OpenFile(config.GetConfig().GetDataIndexPath(), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	indexData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(indexData) > 0 {
		if err := json.Unmarshal(indexData, &innerIndex); err != nil {
			return err
		}
	}

	contentIndexes = append(contentIndexes, innerIndex...)
	d.contentIndexes = contentIndexes
	// todo delete
	log.Printf("索引数据：%+v\n", d.contentIndexes)
	return nil
}

func (d *InnerDataIndexList) Persist() error {
	indexPath := config.GetConfig().GetDataIndexPath()
	// 去除系统配置索引
	innerIndex := d.contentIndexes[1:]
	indexBytes, err := json.Marshal(innerIndex)
	if err != nil {
		return err
	}
	if err := os.WriteFile(indexPath, indexBytes, 0666); err != nil {
		return err
	}
	return nil
}

func (d *InnerDataIndexList) Append(index *InnerDataIndex) error {
	d.contentIndexes = append(d.contentIndexes, index)
	if err := d.Persist(); err != nil {
		return err
	}
	return nil
}

func (d *InnerDataIndexList) GetAll() []*InnerDataIndex {
	return d.contentIndexes
}

func (d *InnerDataIndexList) Get(id string) *InnerDataIndex {
	for _, index := range d.contentIndexes {
		if index.Id == id {
			return index
		}
	}
	return nil
}

func (d *InnerDataIndexList) SetOpen(id string, open bool) error {
	index := d.Get(id)
	index.Open = open
	log.Printf("[setOpen]id: %s, title: %s, set open:%v\n",index.Id, index.Title,index.Open)
	return d.Persist()
}

func (d *InnerDataIndexList) IsOpen(id string) bool {
	for _, index := range d.contentIndexes {
		if index.Id == id && index.Open {
			return true
		}
	}
	return false
}
