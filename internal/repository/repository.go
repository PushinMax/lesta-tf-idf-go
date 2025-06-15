package repository

import (
	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthApi interface{
	Authentication(login, password string) (uuid.UUID, error)
	Register(login, password string) error
	ChangePassword(id, password string) error
	//Logout()
	//DeleteUser()
	//ChangePassword()
}

type DocumentApi interface{
	InsertDocument(document Document) error
	GetDocument(fileID string, userID string) (string, error)
	GetListDocuments(userID string) ([]string, error)
	GetDocumentStats(fileID string, userID string) ([]schema.WordStat, error)
	DeleteDocument(fileID, userID string) error
	DeleteAllDocuments(userID string) error
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

type Repository struct {
	AuthApi
	DocumentApi
	CollectionApi
}

func New(db *sqlx.DB, mongoDB *mongo.Database) *Repository {
	return &Repository{
		AuthApi: newAuthApi(db),
		DocumentApi: newDocumentApi(mongoDB),
		CollectionApi: NewCollectionApi(mongoDB),
	}
}