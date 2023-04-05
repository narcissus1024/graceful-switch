package data

import (
	"bufio"
	"github.com/google/uuid"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"
)

var (
	once      sync.Once
	manager   *DataManager
	lineBreak = "\n"
)

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
			InnerDataList:      &InnerDataList{contents: make(map[string]InnerMetaData)},
			InnerDataIndexList: &InnerDataIndexList{contentIndexes: make([]InnerDataIndex, 0, 10)},
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

	if err := d.MergeConfig(); err != nil {
		return err
	}

	if err := d.InnerDataIndexList.Load(); err != nil {
		return err
	}

	return nil
}

// MergeConfig merge sys ssh config and inner ssh config
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
		isHost, value := d.IsHostConfigItem(line)
		if isHost {
			findConfig, err := d.FindEquivalentHost(innerReader, value)
			if err != nil {
				return err
			}
			// todo 优化
			innerReader.Reset(strings.NewReader(innerData))

			// 在inner中找到相同host，使用inner配置。sys则跳到下一个host
			if len(findConfig) > 0 {
				skip = true
				foundHost[value] = struct{}{}
				continue
			} else {
				// 使用sys的
				skip = false
				res += line + lineBreak
			}
		} else {
			if skip {
				continue
			} else {
				res += line + lineBreak
			}
		}
	}

	skip = false
	for {
		l, _, err := innerReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line := string(l)
		isHost, value := d.IsHostConfigItem(line)
		if isHost {
			skip = false
			_, exist := foundHost[value]
			if exist {
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

func (d *DataManager) FindEquivalentHost(reader *bufio.Reader, targetHost string) (string, error) {
	// todo inner data 必须host开头
	// found为true，说明sys中的某个host，在inner中存在
	found := false
	// tmpHostValue 记录sys中不存在inner中的host配置 的host值。
	// 以该值为key，所属该host的所有配置为value存储在map中
	//tmpHostValue := ""
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
		//if isHost {
		//	tmpHostValue = hostValue
		//}
		//if len(tmpHostValue) > 0 {
		//	foundConfig[tmpHostValue] += line + lineBreak
		//}
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

// IsHostConfigItem parse each line of config in ~/.ssh/config
func (d *DataManager) IsHostConfigItem(str string) (bool, string) {
	reg := regexp.MustCompile(`\s*(?i)host[\s|\s*=\s*](.*)`)
	res := reg.FindStringSubmatch(str)
	if len(res) > 0 {
		return true, res[1]
	}
	return false, ""
}

// MergeConfig merge sys ssh config and inner ssh config
//func (d *DataManager) MergeConfig() error {
//	sshConfig := d.SSHConfig
//	innerDataList := d.InnerDataList
//
//	for _, innerData := range innerDataList.contents {
//		for i, item := range innerData.Contents {
//			// sys存在，则覆盖；不存在则添加
//			sshConfig.MetaDataList[item.Host] = innerData.Contents[i]
//			//if _, exist := sshConfig.MetaDataList[item.Host]; exist {
//			//	sshConfig.MetaDataList[item.Host] = item
//			//} else {
//			//	sshConfig.MetaDataList[item.Host] = item
//			//}
//		}
//	}
//	return nil
//}

func (d *DataManager) PersistSSHConfig() {

}

func (d *DataManager) AddData(id, text string) error {
	content := InnerMetaData{
		Id:       id,
		Contents: text,
	}
	if err := d.InnerDataList.UpdateAndPersist(content); err != nil {
		return err
	}
	if err := d.MergeConfig(); err != nil {
		return err
	}
	if err := d.SSHConfig.Persist(); err != nil {
		return err
	}
	return nil
}

func (d *DataManager) AddDataIndex(t string, title string, ids ...string) error {
	switch t {
	case SIMPLE.String():
		indexList := d.InnerDataIndexList
		uid, _ := uuid.NewUUID()
		index := InnerDataIndex{
			Id:    uid.String(),
			Title: title,
			Open:  true,
		}
		indexList.Append(index)
		if err := indexList.Persist(); err != nil {
			log.Println()
			return err
		}
	case COMBINE.String():
		log.Println(ids)
	}
	return nil
}
