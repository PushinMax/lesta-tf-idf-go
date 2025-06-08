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

type FileApi interface{
	InsertFile(file FileDocument) error
	GetFile(fileID string, userID string) (string, error)
	GetListFiles(userID string) ([]string, error)
	GetFilesStats(fileID string, userID string) ([]schema.WordStat, error)
	DeleteFile(fileID, userID string) error
	DeleteUserFiles(userID string) (int64, error)
}

type Repository struct{
	AuthApi
	FileApi

}



func New(db *sqlx.DB, mongoDB *mongo.Client) *Repository {
	return &Repository{
		AuthApi: newAuthApi(db),
		FileApi: newFileApi(mongoDB),
	}
}