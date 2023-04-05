package controller

import (
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
)



type SSHManager struct {
	Config      *config.Config
	DataManager *data.DataManager
}

func NewSSHManager(opts ...Option) *SSHManager {
	sshManager := new(SSHManager)
	for _, opt := range opts {
		opt(sshManager)
	}
	return sshManager
}

//func (s *SSHManager) AddDataIndex(t string, title string, ids ...string) error {
//	switch t {
//	case SIMPLE.String():
//		indexList := s.DataManager.InnerDataIndexList
//		uid, _ := uuid.NewUUID()
//		index := data.InnerDataIndex{
//			Id:    uid.String(),
//			Title: title,
//			Open:  true,
//		}
//		indexList.Append(index)
//		if err := indexList.Persist(); err != nil {
//			log.Println()
//			return err
//		}
//	case COMBINE.String():
//		log.Println(ids)
//	}
//	return nil
//}

//func (s *SSHManager) AddData(id, text string) error {
//	//contents, err := s.DataManager.SSHConfig.UnmarshalMetaData(text)
//	//if err != nil {
//	//	return err
//	//}
//	content := data.InnerMetaData{
//		Id:       id,
//		Contents: text,
//	}
//	if err := s.DataManager.InnerDataList.UpdateAndPersist(content); err != nil {
//		return err
//	}
//	return nil
//}

type Option func(manager *SSHManager)

func WithConfig(conf *config.Config) Option {
	return func(manager *SSHManager) {
		manager.Config = conf
	}
}

func WithDataManager(dataManager *data.DataManager) Option {
	return func(manager *SSHManager) {
		manager.DataManager = dataManager
	}
}
