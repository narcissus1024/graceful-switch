package data

import (
	"bufio"
	"github.com/google/uuid"
	"io"
	"log"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var (
	once      sync.Once
	manager   *DataManager
	lineBreak = "\n"
)

func init() {
	switch runtime.GOOS {
	case "windows":
		lineBreak = "\r\n"
	case "linux":
		lineBreak = "\n"
	case "darwin":
		lineBreak = "\r"
	}
}

type Type string

const (
	SIMPLE  Type = "simple"
	COMBINE Type = "combine"
)

func (t Type) String() string {
	return string(t)
}

func GetDataManager() *DataManager {
	once.Do(func() {
		manager = &DataManager{
			InnerDataList:      &InnerDataList{contents: make(map[string]*InnerMetaData)},
			InnerDataIndexList: &InnerDataIndexList{contentIndexes: make([]*InnerDataIndex, 0, 10)},
			SSHConfig:          &SSHConfig{},
		}
	})
	return manager
}

type DataManager struct {
	InnerDataList      *InnerDataList
	InnerDataIndexList *InnerDataIndexList
	SSHConfig          *SSHConfig
}

func (d *DataManager) Init() error {
	if err := d.InnerDataList.Load(); err != nil {
		return err
	}

	if err := d.SSHConfig.Load(); err != nil {
		return err
	}

	if err := d.InnerDataIndexList.Load(); err != nil {
		return err
	}

	if err := d.MergeConfig(); err != nil {
		return err
	}

	return nil
}

// MergeConfig merge sys ssh config with inner ssh config to sys ssh config
// todo 优化，太复杂了
func (d *DataManager) MergeConfig() error {
	sysData := d.SSHConfig.SSHConfigData
	innerData := d.InnerDataList.GetAll()
	var res string

	sysReader := bufio.NewReader(strings.NewReader(sysData))
	innerReader := bufio.NewReader(strings.NewReader(innerData))

	// 核心思想：遍历sys每一行， 以host的值做匹配。
	//         以sys为主，host在inner中不存在的，使用sys配置；host在inner存在，使用inner的配置。
	//         还有inner中的host不存在sys，则将inner中host不存在sys的添加到最终合并结果
	//
	// 情况1【done】：inner为空；sys不为空
	// 情况2【done】：inner不为空，sys为空
	// 情况3【done】：inner，sys均为空
	// 情况4【done】：inner，sys均不为空

	// skip为true说明sys的某段配置以inner为主
	skip := false
	foundHost := make(map[string]struct{})
	for {
		l, _, err := sysReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line := string(l)
		isHost, hostValue := d.IsHostConfigItem(line)
		if isHost {
			findConfig := ""
			// todo 优化，先找id在看是否open，太麻烦了
			exist, isOpen, err := d.ExistAndOpenInInnerData(hostValue)
			if err != nil {
				return err
			}

			if exist && isOpen {
				if findConfig, err = d.FindEquivalentHost(innerReader, hostValue); err != nil {
					return err
				}
				innerReader.Reset(strings.NewReader(innerData))
			}

			// 在inner中找到相同host，并且open状态，使用inner配置
			if exist {
				foundHost[hostValue] = struct{}{}
			}
			if exist && isOpen {
				skip = true
				res += findConfig + lineBreak
			} else if exist && !isOpen {
				// inner存在host配置，并且是close状态，则不进行合并
				skip = true
			} else {
				// inner不存在，使用sys的
				skip = false
			}
		}
		if skip {
			continue
		} else {
			res += line + lineBreak
		}
	}

	// 处理inner配置不在sys配置中的情况
	skip = false
	innerReader.Reset(strings.NewReader(innerData))
	for {
		l, _, err := innerReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line := string(l)
		isHost, hostValue := d.IsHostConfigItem(line)
		if isHost {
			skip = false
			exist, isOpen, err := d.ExistAndOpenInInnerData(hostValue)
			if err != nil {
				return err
			}
			_, found := foundHost[hostValue]
			if exist && !isOpen || found {
				skip = true
			}
		}
		if skip {
			continue
		}
		res += line + lineBreak
	}
	d.SSHConfig.SSHConfigData = res
	return nil
}

func (d *DataManager) ExistAndOpenInInnerData(hostValue string) (bool, bool, error) {
	// todo 优化，先找id在看是否open，太麻烦了
	id, err := d.InnerDataList.FindId(hostValue)
	if err != nil {
		return false, false, err
	}
	if len(id) > 0 {
		if d.InnerDataIndexList.IsOpen(id) {
			return true, true, nil
		}
		return true, false, nil
	}
	return false, false, nil
}

// FindEquivalentHost find whole ssh config item which host value equal targetHost from reader
func (d *DataManager) FindEquivalentHost(reader *bufio.Reader, targetHost string) (string, error) {
	// todo inner data 必须host开头
	// found为true，说明sys中的某个host，在inner中存在
	found := false
	res := ""
	for {
		li, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return res, err
		}
		line := string(li)
		isHost, hostValue := d.IsHostConfigItem(line)
		if isHost && hostValue == targetHost {
			// sys的host能够在inner中找到
			found = true
		} else if isHost && hostValue != targetHost {
			// 遇到新的host，并且host与需要匹配的host不想等。并且存在以前找到的情况，提前结束遍历
			if found {
				break
			}
			found = false
		}
		if found {
			res += line + lineBreak
		}
	}
	return res, nil
}

//func (d *DataManager) IsComment(str string) bool {
//	if len(strings.TrimSpace(str)) == 0 || strings.HasPrefix("#", str) {
//		return true
//	}
//	return false
//}

// IsHostConfigItem check whether the row config is host config item.
// if is host config item, return true and host value. otherwise, return false and ""
func (d *DataManager) IsHostConfigItem(str string) (bool, string) {
	reg := regexp.MustCompile(`^\s*(?i)host\s*=?\s*`)
	res := reg.Split(str, -1)

	if len(res) > 1 {
		return true, strings.TrimSpace(res[1])
	}
	return false, ""
}

// PersistSSHConfig persist sys ssh config which contain inner ssh and sys ssh
func (d *DataManager) PersistSSHConfig() error {
	if err := d.MergeConfig(); err != nil {
		return err
	}
	if err := d.SSHConfig.Persist(); err != nil {
		return err
	}
	return nil
}

// AddData add inner ssh config and persist
func (d *DataManager) AddData(id, text string) error {
	if len(id) == 0 || id == SYSTEM_ID_SSH || len(strings.TrimSpace(text)) == 0 {
		return nil
	}
	content := &InnerMetaData{
		Id:       id,
		Contents: text,
	}
	if err := d.InnerDataList.UpdateAndPersist(content); err != nil {
		return err
	}
	if err := d.PersistSSHConfig(); err != nil {
		return err
	}
	return nil
}

// AddDataIndex add inner ssh config index for inner ssh config data
func (d *DataManager) AddDataIndex(t string, title string, ids ...string) error {
	switch t {
	case SIMPLE.String():
		indexList := d.InnerDataIndexList
		uid, _ := uuid.NewUUID()
		index := &InnerDataIndex{
			Id:    uid.String(),
			Title: title,
			Open:  true,
		}
		if err := indexList.Append(index); err != nil {
			return err
		}

	case COMBINE.String():
		log.Println(ids)
	}
	return nil
}

func (d *DataManager) GetSSHConfigById(id string) string {
	if id == SYSTEM_ID_SSH {
		return d.SSHConfig.SSHConfigData
	}
	innerData := d.InnerDataList.Get(id)
	if innerData != nil {
		return innerData.Contents
	}
	return ""
}