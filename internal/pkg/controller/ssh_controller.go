package controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/narcissus1024/graceful-switch/internal/pkg/config"
	"github.com/narcissus1024/graceful-switch/internal/pkg/data"
)

type Type string

const (
	SIMPLE  Type = "simple"
	COMBINE Type = "combine"
)

func (t Type) String() string {
	return string(t)
}

type SSHController struct {
	Config *config.Config
	Data   *data.Data
}

func (s *SSHController) AddDataIndex(t string, title string, ids ...string) error {
	switch t {
	case SIMPLE.String():
		indexList := s.Data.ContentIndexList
		uid, _ := uuid.NewUUID()
		index := data.ContentIndex{
			Id:    uid.String(),
			Title: title,
			Open:  true,
		}
		indexList.Append(index)
		if err := indexList.Persist(); err != nil {
			fmt.Println()
			return err
		}
	case COMBINE.String():
		fmt.Println(ids)
	}
	return nil
}

func (s *SSHController) AddData(id, text string) error {
	content := data.Content{
		Id:      id,
		Content: text,
	}
	if err := s.Data.ContentList.UpdateAndPersist(content); err != nil {
		return err
	}
	return nil
}

