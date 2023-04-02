package data

import (
	"encoding/json"
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/tools"
	"io/ioutil"
	"os"
	"path"
)

type Content struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

type ContentList struct {
	contents map[string]Content
}

func (d *ContentList) load(config *config.Config) error {
	dataPath := config.GetDataPath()
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
		content := Content{}
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
	return nil
}

func (d *ContentList) Get(id string) Content {
	return d.contents[id]
}

func (d *ContentList) UpdateAndPersist(content Content) error {
	d.contents[content.Id] = content
	contentPath := path.Join(config.Conf.GetDataPath(), content.Id+".json")
	contentByte, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(contentPath, contentByte, 0666)
}

//func (d *ContentList) Update(content Content) {
//	d.contents[content.Id] = content
//}
//
//func (d *ContentList) Persist(id string) error {
//	contentPath := path.Join(config.Conf.GetDataPath(), id+".json")
//	contentByte, err := json.Marshal(d.Get(id))
//	if err != nil {
//		return err
//	}
//	return ioutil.WriteFile(contentPath, contentByte, 0666)
//}