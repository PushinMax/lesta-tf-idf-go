package repository

import (
	"context"
	"errors"
	"math"

	// "github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"fmt"
	"log"

	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DocumentRepo struct {
	db *mongo.Database
}

func newDocumentApi(db *mongo.Database) *DocumentRepo {
	return &DocumentRepo{db: db}
}


func (r *DocumentRepo) InsertDocument(document Document) error {
    result, err := r.db.Collection("documents").InsertOne(context.TODO(), document)
    if err != nil {
        return err
    }

    log.Printf("Документ вставлен с ID: %v", result.InsertedID)
    log.Println(document.Content)
    return nil
}

func (r *DocumentRepo) GetDocument(fileID string, userID string) (string, error) {
    var result struct {
        Content string `bson:"file"`
    }
    err := r.db.Collection("documents").FindOne(
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
            return "", fmt.Errorf("документ не найден или у вас нет прав доступа")
        }
        return "", fmt.Errorf("ошибка при получении документа: %v", err)
    }

    return result.Content, nil
}

func (r *DocumentRepo) GetListDocuments(userID string) ([]string, error) {
    opts := options.Find().SetProjection(bson.M{
        "file_id": 1,
        "_id":     0,
    })

    cursor, err := r.db.Collection("documents").Find(
        context.TODO(),
        bson.M{"user_id": userID},
        opts,
    )
    if err != nil {
        log.Printf("ошибка поиска документов: %v", err)
        return nil, fmt.Errorf("ошибка поиска документов")
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

func (r *DocumentRepo) GetDocumentStats(fileID string, userID string) ([]schema.WordStat, error) {
    var doc Document
    err := r.db.Collection("files").FindOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID,
        },
    ).Decode(&doc)
    
    if err != nil {
        return nil, fmt.Errorf("failed to find document: %v", err)
    }

    if doc.IsValid {
        return doc.Stats, nil
    }

    var collections []Collection
    cursor, err := r.db.Collection("collections").Find(
        context.TODO(),
        bson.M{
            "documents": fileID,
            "user_id": userID,
        },
    )
    if err != nil {
        return nil, fmt.Errorf("failed to find collections: %v", err)
    }
    defer cursor.Close(context.TODO())

    if err = cursor.All(context.TODO(), &collections); err != nil {
        return nil, fmt.Errorf("failed to decode collections: %v", err)
    }

    collectionWords := make(map[string]struct{
        amount_w int
        amount_d int
    })
    
    for _, collection := range collections {
        for word, stats := range collection.Words {
            if _, exists := collectionWords[word]; !exists {
                collectionWords[word] = struct{
                    amount_w int
                    amount_d int
                }{0, 0}
            }
            existing := collectionWords[word]
            existing.amount_w += stats.amount_w
            existing.amount_d += stats.amount_d
            collectionWords[word] = existing
        }
    }

    newStats := make([]schema.WordStat, 0, len(collectionWords))
    for word, stats := range collectionWords {
        newStats = append(newStats, schema.WordStat{    
            Word: word,
            TF:   float64(stats.amount_w) / float64(doc.Length),
            IDF:  math.Log(float64(len(collections)) / float64(stats.amount_d)),
        })
    }   


    update := bson.M{
        "$set": bson.M{
            "stat": newStats,
            "isvalid": true,
        },
    }

    _, err = r.db.Collection("files").UpdateOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID,
        },
        update,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to update stats: %v", err)
    }

    return newStats, nil
}

func (r *DocumentRepo) DeleteDocument(fileID string, userID string) error {
    var doc Document
    err := r.db.Collection("files").FindOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID,
        },
    ).Decode(&doc)
    if err != nil {
        return fmt.Errorf("document not found: %v", err)
    }

    for _, collectionName := range doc.Collections {
        _, err := r.db.Collection("collections").UpdateOne(
            context.TODO(),
            bson.M{
                "name": collectionName,
                "user_id": userID,
            },
            bson.M{
                "$pull": bson.M{
                    "documents": fileID,
                },
                "$set": bson.M{
                    "isvalid": false,
                },
            },
        )
        if err != nil {
            return fmt.Errorf("failed to remove document from collection %s: %v", collectionName, err)
        }
    }

    _, err = r.db.Collection("files").DeleteOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID,
        },
    )
    if err != nil {
        return fmt.Errorf("failed to delete document: %v", err)
    }

    return nil
}

func (r *DocumentRepo) DeleteAllDocuments(userID string) error {
    _, err := r.db.Collection("files").DeleteMany(
        context.TODO(),
        bson.M{"user_id": userID},
    )
    if err != nil {
        return fmt.Errorf("failed to delete all documents: %v", err)
    }

    _, err = r.db.Collection("collections").UpdateMany(
        context.TODO(),
        bson.M{"user_id": userID},
        bson.M{
            "$set": bson.M{"isvalid": false},
            "$unset": bson.M{"documents": ""},
        },
    )
    if err != nil {
        return fmt.Errorf("failed to update collections after deleting documents: %v", err)
    }

    return nil
}