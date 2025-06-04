package service

import (
	"errors"

	"github.com/spf13/viper"
)

type StatusService struct {
}

func newStatusApi() *StatusService {
	return &StatusService{}
}

func (s *StatusService) Status() error {
	// TODO Добавить наличие подключения к бд или другие состояния
	return nil
}

func (s *StatusService) Version() (string, error) {
	version := viper.GetString("server.version")
	if version == "" {
		return "", errors.New("Fail")
	}
	return version, nil
}