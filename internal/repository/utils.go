package repository

import (
	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID uuid.UUID `db:"id"`
	Login string `db:"login"`
	Password string `db:"password_hash"`
}


type FileDocument struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    FileID  string            `bson:"file_id"`
    Name    string            `bson:"file_name"`
    UserID  string            `bson:"user_id"`
    Content string            `bson:"file"`
    Stats   []schema.WordStat `bson:"stat"`
    Length  int               `bson:"len"`
}