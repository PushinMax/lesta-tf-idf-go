package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CollectionRepo struct {
	db *mongo.Database
}

func NewCollectionApi(db *mongo.Database) *CollectionRepo {
 return &CollectionRepo{db: db}
}

func (r *CollectionRepo) CreateCollection(userID, name string) error {
	_, err := r.db.Collection("collections").InsertOne(context.TODO(), Collection{
		Name: name,
		UserID: userID,	
		Stats: make([]schema.WordStat, 0),
		IsValid: false,
		Length: 0,
		Words: make(map[string]struct {
			   amount_w int `bson:"amount_w"`
			   amount_d int `bson:"amount_d"`
		}),
		DocumentsID: make([]string, 0),
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *CollectionRepo) GetListCollections(userID string) ([]string, error) {
	opts := options.Find().SetProjection(bson.M{
		"name": 1,
		"_id":   0,
	})

	cursor, err := r.db.Collection("collections").Find(
		context.TODO(),
		bson.M{"user_id": userID},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []struct {
		Name string `bson:"name"`
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	collectionNames := make([]string, len(results))
	for i, result := range results {
		collectionNames[i] = result.Name
	}
	return collectionNames, nil	
}

func (r *CollectionRepo) GetDocumentsInCollection(userID, collectionName string) ([]string, error) {
	opts := options.Find().SetProjection(bson.M{
		"documents": 1,
		"_id":       0,
	})

	cursor, err := r.db.Collection("collections").Find(
		context.TODO(),
		bson.M{"user_id": userID, "name": collectionName},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []struct {
		Documents []string `bson:"documents"`
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil // No documents found in the collection
	}

	return results[0].Documents, nil
}

func (r *CollectionRepo) GetCollectionStats(userID, collectionName string) ([]schema.WordStat, error) {
	var collection Collection
    err := r.db.Collection("collections").FindOne(
        context.TODO(),
        bson.M{
            "name": collectionName,
            "user_id": userID,
        },
    ).Decode(&collection)
    
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("collection not found")
        }
        return nil, err
    }

    if collection.IsValid {
        return collection.Stats, nil
    }

    words := collection.Words
    length := collection.Length
	documentAmount := len(collection.DocumentsID)



	var newStats []schema.WordStat
	for word, stat := range words {
		newStats = append(newStats, schema.WordStat{
			Word: word,
			TF:  float64(stat.amount_w) / float64(length),
			IDF: math.Log(float64(documentAmount) / float64(stat.amount_d)),
		})
	}
	sort.Slice(newStats, func(i, j int) bool {
		return newStats[i].IDF > newStats[j].IDF})


    update := bson.M{
        "$set": bson.M{
            "stat": newStats,
            "isvalid": true,
        },
    }
    
    _, err = r.db.Collection("collections").UpdateOne(
        context.TODO(),
        bson.M{
            "name": collectionName,
            "user_id": userID,
        },
        update,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to update collection stats: %v", err)
    }

    return newStats, nil
}

func (r *CollectionRepo) AddDocumentToCollection(userID, collectionName, fileID string) error {
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

	var collection Collection
    err = r.db.Collection("collections").FindOne(
        context.TODO(),
        bson.M{
            "name": collectionName,
            "user_id": userID,
        },
    ).Decode(&collection)
    if err != nil {
        return fmt.Errorf("collection not found: %v", err)
    }

    for word, amount := range doc.Words {
        if _, exists := collection.Words[word]; !exists {
            collection.Words[word] = struct {
                amount_w int `bson:"amount_w"`
                amount_d int `bson:"amount_d"`
            }{0, 0}
        }
        wordStat := collection.Words[word]
        wordStat.amount_w += amount
        wordStat.amount_d += 1
        collection.Words[word] = wordStat
    }

	_, err = r.db.Collection("files").UpdateOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID,
        },
        bson.M{
            "$addToSet": bson.M{
                "collections": collectionName,
            },
            "$set": bson.M{
                "isvalid": false,
            },
        },
    )
    if err != nil {
        return fmt.Errorf("failed to update document: %v", err)
    }

    update := bson.M{
        "$addToSet": bson.M{
            "documents": fileID,
        },
        "$set": bson.M{
            "words": collection.Words,
            "length": collection.Length + doc.Length,
            "isvalid": false,
        },
    }

    result, err := r.db.Collection("collections").UpdateOne(
        context.TODO(),
        bson.M{
            "name": collectionName,
            "user_id": userID,
        },
        update,
    )

    if err != nil {
        return fmt.Errorf("failed to update collection: %v", err)
    }

    if result.MatchedCount == 0 {
        return errors.New("collection not found")
    }

    return nil
}

func (r *CollectionRepo) DeleteDocumentFromCollection(userID, collectionName, fileID string) error {
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

    var collection Collection
    err = r.db.Collection("collections").FindOne(
        context.TODO(),
        bson.M{
            "name": collectionName,
            "user_id": userID,
        },
    ).Decode(&collection)
    if err != nil {
        return fmt.Errorf("collection not found: %v", err)
    }

    for word, amount := range doc.Words {
        if wordStat, exists := collection.Words[word]; exists {
            wordStat.amount_w -= amount
            wordStat.amount_d -= 1
            if wordStat.amount_d == 0 {
                delete(collection.Words, word)
            } else {
                collection.Words[word] = wordStat
            }
        }
    }


	_, err = r.db.Collection("files").UpdateOne(
        context.TODO(),
        bson.M{
            "file_id": fileID,
            "user_id": userID,
        },
        bson.M{
            "$pull": bson.M{
                "collections": collectionName,
            },
            "$set": bson.M{
                "isvalid": false,
            },
        },
    )
    if err != nil {
        return fmt.Errorf("failed to update document: %v", err)
    }

    update := bson.M{
        "$pull": bson.M{
            "documents": fileID,
        },
        "$set": bson.M{
            "words": collection.Words,
            "length": collection.Length - doc.Length,
            "isvalid": false,
        },
    }

    result, err := r.db.Collection("collections").UpdateOne(
        context.TODO(),
        bson.M{
            "name": collectionName,
            "user_id": userID,
        },
        update,
    )

    if err != nil {
        return fmt.Errorf("failed to update collection: %v", err)
    }

    if result.MatchedCount == 0 {
        return errors.New("collection not found")
    }

    return nil
}

func (r *CollectionRepo) DeleteCollection(userID, collectionName string) error {
	var collection Collection
    err := r.db.Collection("collections").FindOne(
        context.TODO(),
        bson.M{
            "name": collectionName,
            "user_id": userID,
        },
    ).Decode(&collection)
    if err != nil {
        return fmt.Errorf("collection not found: %v", err)
    }

    _, err = r.db.Collection("files").UpdateMany(
        context.TODO(),
        bson.M{
            "file_id": bson.M{
                "$in": collection.DocumentsID,
            },
            "user_id": userID,
        },
        bson.M{
            "$pull": bson.M{
                "collections": collectionName,
            },
            "$set": bson.M{
                "isvalid": false,
            },
        },
    )
    if err != nil {
        return fmt.Errorf("failed to update documents: %v", err)
    }

    _, err = r.db.Collection("collections").DeleteOne(
        context.TODO(),
        bson.M{
            "name": collectionName,
            "user_id": userID,
        },
    )
    if err != nil {
        return fmt.Errorf("failed to delete collection: %v", err)
    }

    return nil
}


func (r *CollectionRepo) DeleteAllCollections(userID string) error {
	_, err := r.db.Collection("collections").DeleteMany(
		context.TODO(),
		bson.M{"user_id": userID},
	)
	if err != nil {
		return fmt.Errorf("failed to delete all collections: %v", err)
	}

	_, err = r.db.Collection("files").UpdateMany(
		context.TODO(),
		bson.M{"user_id": userID},
		bson.M{
			"$set": bson.M{"isvalid": false},
			"$pull": bson.M{"collections": bson.M{"$exists": true}},
	},
	)
	if err != nil {
		return fmt.Errorf("failed to update documents after deleting collections: %v", err)
	}

	return nil
}