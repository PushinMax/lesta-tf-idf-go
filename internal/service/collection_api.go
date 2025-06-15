package service

import (
	"github.com/PushinMax/lesta-tf-idf-go/internal/repository"
	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
)


type CollectionService struct {
	repos *repository.Repository
}

func newCollectionApi(repos *repository.Repository) *CollectionService {
	return &CollectionService{repos: repos}
}

func (s *CollectionService) CreateCollection(userID, name string) error {
	 return s.repos.CollectionApi.CreateCollection(userID, name)
}

func (s *CollectionService) GetListCollections(userID string) ([]string, error) {
  	return s.repos.CollectionApi.GetListCollections(userID)
}

func (s *CollectionService) GetDocumentsInCollection(userID, collectionName string) ([]string, error) {
  	return s.repos.CollectionApi.GetDocumentsInCollection(userID, collectionName)
}

func (s *CollectionService) GetCollectionStats(userID, collectionName string) ([]schema.WordStat, error) {
  	return s.repos.CollectionApi.GetCollectionStats(userID, collectionName)
}

func (s *CollectionService) AddDocumentToCollection(userID, collectionName, fileID string) error {
  	return s.repos.CollectionApi.AddDocumentToCollection(userID, collectionName, fileID)
}

func (s *CollectionService) DeleteDocumentFromCollection(userID, collectionName, fileID string) error {
  	return s.repos.CollectionApi.DeleteDocumentFromCollection(userID, collectionName, fileID)
}

func (s *CollectionService) DeleteCollection(userID, collectionName string) error {
  	return s.repos.CollectionApi.DeleteCollection(userID, collectionName)
}

func (s *CollectionService) DeleteAllCollections(userID string) error {
  	return s.repos.CollectionApi.DeleteAllCollections(userID)
}