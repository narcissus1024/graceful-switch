package data

import (
	"encoding/json"
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/tools"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// InnerMetaData inner ssh config item
type InnerMetaData struct {
	Id       string `json:"id"`
	Contents string `json:"contents"`
}

//type InnerMetaData struct {
//	Id       string     `json:"id"`
//	Contents []MetaData `json:"contents"`
//}

// InnerDataList inner ssh config
type InnerDataList struct {
	contents map[string]InnerMetaData
}

func (d *InnerDataList) Load() error {
	dataPath := config.GetConfig().GetDataPath()
	if !tools.IsExist(dataPath) {
		if err := os.MkdirAll(dataPath, 0777); err != nil {
			return err
		}
	}
	fslist, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return err
	}
	for _, fs := range fslist {
		content := InnerMetaData{}
		if !fs.IsDir() {
			dataItem, err := ioutil.ReadFile(path.Join(dataPath, fs.Name()))
			if err != nil {
				return err
			}
			if err := json.Unmarshal(dataItem, &content); err != nil {
				return err
			}
			d.contents[content.Id] = content
		}
	}

	// todo delete
	log.Printf("内部数据：%+v\n", d.contents)
	return nil
}

func (d *InnerDataList) UpdateAndPersist(content InnerMetaData) error {
	d.contents[content.Id] = content
	contentPath := path.Join(config.GetConfig().GetDataPath(), content.Id+".json")
	contentByte, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(contentPath, contentByte, 0666)
}

func (d *InnerDataList) Get(id string) InnerMetaData {
	return d.contents[id]
}

func (d *InnerDataList) GetAll() string {
	res := ""
	for _, v := range d.contents {
		res += v.Contents
		res += "\n"
	}
	return res
}

//func (d *InnerDataList) Update(content InnerMetaData) {
//	d.contents[content.Id] = content
//}
//
//func (d *InnerDataList) Persist(id string) error {
//	contentPath := path.Join(config.Conf.GetDataPath(), id+".json")
//	contentByte, err := json.Marshal(d.Get(id))
//	if err != nil {
//		return err
//	}
//	return ioutil.WriteFile(contentPath, contentByte, 0666)
//}
