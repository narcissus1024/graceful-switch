package data

import (
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/tools"
	"io/ioutil"
	"log"
	"os"
)

//type MetaData struct {
//	Host     string
//	HostName string
//}

//type SSHConfig struct {
//	MetaDataList map[string]MetaData
//}

type SSHConfig struct {
	SSHConfigData string
}

func (s *SSHConfig) Load() error {
	sysPath := config.GetConfig().SSHConfigPath
	if !tools.IsExist(sysPath) {
		if f, err := os.Create(sysPath); err != nil {
			return err
		} else {
			f.Close()
		}
	}
	data, err := ioutil.ReadFile(sysPath)
	if err != nil {
		return err
	}
	s.SSHConfigData = string(data)
	log.Println("系统数据：", data)
	return nil
}

func (s *SSHConfig) Persist() error {
	sysPath := config.GetConfig().SSHConfigPath
	if err := os.WriteFile(sysPath, []byte(s.SSHConfigData), 0666); err != nil {
		return err
	}
	return nil
}

//func (s *SSHConfig) Load() error {
//	sysPath := config.GetConfig().SSHConfigPath
//	log.Println("ssh config file path: ", sysPath)
//	var f *os.File
//	var err error
//	if !tools.IsExist(sysPath) {
//		if f, err = os.Create(sysPath); err != nil {
//			return err
//		}
//	} else {
//		f, err = os.Open(sysPath)
//		if err != nil {
//			return err
//		}
//	}
//	defer func() {
//		if f != nil {
//			f.Close()
//		}
//	}()
//	data, err := ioutil.ReadAll(f)
//	if err != nil {
//		return err
//	}
//	metaDataList, err := s.UnmarshalMetaData(string(data))
//	if err != nil {
//		return err
//	}
//
//	for i := range metaDataList {
//		s.MetaDataList[metaDataList[i].Host] = metaDataList[i]
//	}
//	// todo delete
//	log.Printf("系统数据：%+v\n", s.MetaDataList)
//	return nil
//}
//
//func (s *SSHConfig) Persist() error {
//	return nil
//}
//
//func (s *SSHConfig)  UnmarshalMetaData(content string) ([]MetaData, error) {
//	var result []MetaData
//	reader := bufio.NewReader(strings.NewReader(content))
//
//	m := make(map[string]string)
//	for {
//		l, _, err := reader.ReadLine()
//		if err != nil {
//			if err == io.EOF {
//				metaData, err := s.GenerateMetaData(m)
//				if err != nil {
//					return nil, err
//				}
//				result = append(result, metaData)
//				break
//			}
//			return nil, err
//		}
//
//		line := strings.TrimSpace(string(l))
//		if len(line) == 0 {
//			continue
//		}
//		key, value := s.ParseLineToKV(line)
//
//		if strings.ToLower(key) == "host" {
//			if len(m) > 0 {
//				metaData, err := s.GenerateMetaData(m)
//				if err != nil {
//					return nil, err
//				}
//				result = append(result, metaData)
//				m = make(map[string]string)
//			}
//		}
//		m[key] = value
//
//	}
//	return result, nil
//}
//
//func (s *SSHConfig) MarshalMetaData(metaDataList []MetaData) (string, error) {
//	var res string
//	for _, metaData := range metaDataList {
//		m := make(map[string]string)
//		metaDataByte, err := json.Marshal(metaData)
//		if err != nil {
//			return "", nil
//		}
//		if err := json.Unmarshal(metaDataByte, &m); err != nil {
//			return "", nil
//		}
//		for k, v := range m {
//			res += fmt.Sprintf("%s %s\n", k, v)
//		}
//		res += "\n"
//	}
//	return res, nil
//}
//
//// ParseLineToKV parse each line of config in ~/.ssh/config
//func (s *SSHConfig) ParseLineToKV(kv string) (string, string) {
//	kv = strings.TrimSpace(kv)
//	reg := regexp.MustCompile(`([a-zA-Z]+)[=|\s](.*)`)
//	res := reg.FindStringSubmatch(kv)
//	key := res[1]
//	value := res[2]
//	return key, value
//}
//
//// GenerateMetaData generate one ssh config item
//func (s *SSHConfig) GenerateMetaData(m map[string]string) (MetaData, error) {
//	metData := MetaData{}
//
//	mJson, err := json.Marshal(m)
//	if err != nil {
//		return metData, err
//	}
//
//	if err := json.Unmarshal(mJson, &metData); err != nil {
//		return metData, err
//	}
//	return metData, nil
//}
