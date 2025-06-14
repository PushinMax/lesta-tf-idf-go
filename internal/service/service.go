package service

import (
	"mime/multipart"

	"github.com/PushinMax/lesta-tf-idf-go/internal/repository"
	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"github.com/PushinMax/lesta-tf-idf-go/internal/session"
)

type GeneralApi interface {
	//UploadFile(file *multipart.FileHeader) (*schema.UploadResponse, error)
	GetPageData(sessionID string, page int) (*schema.PageResponse, error)
}

type MetricsApi interface {
	
}

type StatusApi interface {
	Status() error
	Version() (string, error)
}

type AuthApi interface{
	Login(login, password, ip string) (*TokenPairResponse, error)
	Register(login, password string) error
	ValidateToken(token string) (*CustomClaims, error)
}

type UserApi interface {
	ChangePassword(id, password string) error
}

type FileApi interface {
	UploadFile(file *multipart.FileHeader, userID string) (*schema.UploadResponse, error)
	GetListFiles(userID string) ([]string, error)
	GetFile(documentID string, userID string) (string, error)
	GetFilesStats(fileID string, userID string) ([]schema.WordStat, error)
	DeleteFile(fileID, userID string) error
	DeleteUserFiles(userID string) (int64, error)
}

type CollectionApi interface {

}

type Service struct {
	GeneralApi
	// MetricsApi
	StatusApi
	AuthApi
	UserApi
	FileApi
	CollectionApi
}

func New(session *session.Session, repos *repository.Repository) *Service {
	return &Service{
		GeneralApi: newGeneralApi(session),
		StatusApi: newStatusApi(),
		AuthApi: newAuthApi(repos),
		UserApi: newUserApi(repos),
		FileApi: newFileApi(repos),
	}
}