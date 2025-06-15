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

type DocumentApi interface {
	UploadDocument(file *multipart.FileHeader, userID string) (*schema.UploadResponse, error)
	GetListDocuments(userID string) ([]string, error)
	GetDocument(documentID string, userID string) (string, error)
	GetDocumentStats(fileID string, userID string) ([]schema.WordStat, error)
	DeleteDocument(fileID, userID string) error
	DeleteUserDocuments(userID string) error
}

type CollectionApi interface {
	CreateCollection(userID, name string) error
	GetListCollections(userID string) ([]string, error)
	GetDocumentsInCollection(userID, collectionName string) ([]string, error)
	GetCollectionStats(userID, collectionName string) ([]schema.WordStat, error)
	AddDocumentToCollection(userID, collectionName, fileID string) error
	DeleteDocumentFromCollection(userID, collectionName, fileID string) error
	DeleteCollection(userID, collectionName string) error
	DeleteAllCollections(userID string) error
}

type Service struct {
	GeneralApi
	// MetricsApi
	StatusApi
	AuthApi
	UserApi
	DocumentApi
	CollectionApi
}

func New(session *session.Session, repos *repository.Repository) *Service {
	return &Service{
		GeneralApi: newGeneralApi(session),
		StatusApi: newStatusApi(),
		AuthApi: newAuthApi(repos),
		UserApi: newUserApi(repos),
		DocumentApi: newDocumentApi(repos),
		CollectionApi: newCollectionApi(repos),
	}
}