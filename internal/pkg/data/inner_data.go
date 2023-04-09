package data

import (
	"bufio"
	"encoding/json"
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/tools"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

// InnerMetaData inner ssh config item
type InnerMetaData struct {
	Id       string `json:"id"`
	Contents string `json:"contents"`
}

// InnerDataList inner ssh config
type InnerDataList struct {
	contents map[string]*InnerMetaData
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
		content := &InnerMetaData{}
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

func (d *InnerDataList) UpdateAndPersist(content *InnerMetaData) error {
	if len(content.Id) == 0 || len(strings.TrimSpace(content.Contents)) == 0 {
		return nil
	}
	// sys config
	if content.Id == SYSTEM_ID_SSH {
		return nil
	}
	d.contents[content.Id] = content
	contentPath := path.Join(config.GetConfig().GetDataPath(), content.Id+".json")
	contentByte, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(contentPath, contentByte, 0666)
}

func (d *InnerDataList) Get(id string) *InnerMetaData {
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

// todo 优化
func (d *InnerDataList) FindId(hostValue string) (string,error) {
	for _, v := range d.contents {
		reader := bufio.NewReader(strings.NewReader(v.Contents))
		for {
			li, _, err := reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				return "", err
			}
			line := string(li)
			isHost, hv := d.IsHostConfigItem(line)
			if isHost && hv == hostValue {
				return v.Id, nil
			}
		}
	}
	return "", nil
}

// TODO 优化，重复代码
func (d *InnerDataList) IsHostConfigItem(str string) (bool, string) {
	reg := regexp.MustCompile(`^\s*(?i)host\s*=?\s*`)
	res := reg.Split(str, -1)

	if len(res) > 1 {
		return true, strings.TrimSpace(res[1])
	}
	return false, ""
}
