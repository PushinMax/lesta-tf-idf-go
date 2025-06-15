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

type Document struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    FileID  string            `bson:"file_id"`
    Name    string            `bson:"file_name"`
    UserID  string            `bson:"user_id"`
    Content string            `bson:"file"`
    Stats   []schema.WordStat `bson:"stat"`
    IsValid bool              `bson:"isvalid"`
    Length  int               `bson:"len"`
    Words map[string]int      `bson:"words"`
    Collections []string      `bson:"collections"`
}

type Collection struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    Name    string            `bson:"name"`
    UserID  string            `bson:"user_id"`
    Stats   []schema.WordStat `bson:"stat"`
    IsValid bool              `bson:"isvalid"`
    Length  int               `bson:"len"`
    Words map[string]struct {
        amount_w int `bson:"amount_w"`
        amount_d int `bson:"amount_d"`
    }      `bson:"words"`
    DocumentsID []string      `bson:"documents"`
}

