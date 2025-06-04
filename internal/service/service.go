package service

import (
	"mime/multipart"

	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"github.com/PushinMax/lesta-tf-idf-go/internal/session"
)

type GeneralApi interface {
	UploadFile(file *multipart.FileHeader) (*schema.UploadResponse, error)
	GetPageData(sessionID string, page int) (*schema.PageResponse, error)
}

type MetricsApi interface {
	
}

type StatusApi interface {
	Status() error
	Version() (string, error)
}

type Service struct {
	GeneralApi
	MetricsApi
	StatusApi
}

func New(session *session.Session) *Service {
	
	return &Service{
		GeneralApi: newGeneralApi(session),
		StatusApi: newStatusApi(),
	}
}