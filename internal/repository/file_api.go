package repository

import (
	"context"
	"errors"

	// "github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"fmt"
	"log"

	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FileRepo struct {
	db *mongo.Client
}

func newFileApi(db *mongo.Client) *FileRepo {
	return &FileRepo{db: db}
}

func (r *FileRepo) InsertFile(file FileDocument) error {
	result, err := r.db.Database("test").Collection("files").InsertOne(context.TODO(), file)
	if err != nil {
		return err
	}

	//insertedID := result.InsertedID
	log.Printf("Документ вставлен с ID: %v", result.InsertedID)
    log.Println(file.Content)
	return nil
}

func (r *FileRepo) GetListFiles(userID string) ([]string, error) {
	opts := options.Find().SetProjection(bson.M{
        "file_id": 1,
        "_id":     0,
    })

    cursor, err := r.db.Database("test").Collection("files").Find(
        context.TODO(),
        bson.M{"user_id": userID},
        opts,
    )
    if err != nil {
		log.Printf("ошибка поиска файлов: %v", err)
        return nil, fmt.Errorf("ошибка поиска файлов")
    }
    defer cursor.Close(context.TODO())

    var results []struct {
        FileID string `bson:"file_id"`
    }

    if err = cursor.All(context.TODO(), &results); err != nil {
		log.Printf("ошибка декодирования результатов: %v", err)
        return nil, errors.New("ошибка декодирования результатов")
    }

    fileIDs := make([]string, len(results))
    for i, r := range results {
        fileIDs[i] = r.FileID
    }

    return fileIDs, nil
	
}

func (r *FileRepo) GetFilesStats(fileID string, userID string) ([]schema.WordStat, error) {
	var result struct {
        Stats []schema.WordStat `bson:"stat"`
    }

    err := r.db.Database("test").Collection("files").FindOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID, 
        },
        options.FindOne().SetProjection(bson.M{
            "stat": 1,
            "_id":  0,
        }),
    ).Decode(&result)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("файл не найден или у вас нет прав доступа")
        }
        return nil, fmt.Errorf("ошибка при получении статистики: %v", err)
    }

    return result.Stats, nil
}

func (r *FileRepo) GetFile(fileID string, userID string) (string, error) {
    var result struct {
        File string `bson:"file"`
    }
    err := r.db.Database("test").Collection("files").FindOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID, 
        },
        options.FindOne().SetProjection(bson.M{
            "file": 1,
            "_id":  0,
        }),
    ).Decode(&result)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return "", fmt.Errorf("файл не найден или у вас нет прав доступа")
        }
        return "", fmt.Errorf("ошибка при получении файла: %v", err)
    }

    return result.File, nil
}

func (r *FileRepo) DeleteFile(fileID, userID string) error {
	result, err := r.db.Database("test").Collection("files").DeleteOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID, 
        },
    )
    if err != nil {
        return fmt.Errorf("ошибка при удалении файла: %v", err)
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("файл не найден или у вас нет прав на удаление")
    }

    return nil
}

func (r *FileRepo) DeleteUserFiles(userID string) (int64, error) {
	result, err := r.db.Database("test").Collection("files").DeleteMany(
        context.TODO(),
        bson.M{"user_id": userID},
    )
    if err != nil {
        return 0, fmt.Errorf("ошибка при удалении файлов: %v", err)
    }
    
    return result.DeletedCount, nil
}
