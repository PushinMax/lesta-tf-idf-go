package service

import (
	"sort"
	"errors"
	"math"

	"github.com/PushinMax/lesta-tf-idf-go/internal/repository"
	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"

	"mime/multipart"
	"bufio"
	"strings"
	"strconv"
	"time"
)


type DocumentService struct {
	repos *repository.Repository
}

func newDocumentApi(repos *repository.Repository) *DocumentService {
	return &DocumentService{
		repos: repos,
	}
}

func (s *DocumentService) UploadDocument(file *multipart.FileHeader, userID string) (*schema.UploadResponse, error) {
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("Cannot open file")
	}
	defer src.Close()

	scanner := bufio.NewScanner(src)
	
	documents := make([]map[string]bool, 0)
	wordCount := make(map[string]int)
	totalWords := 0
	var content strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		content.WriteString(line)
		words := strings.Fields(line)
		docWords := make(map[string]bool)

		for _, word := range words {
			clean := cleanWord(word)
			if clean != "" {
				wordCount[clean]++
				totalWords++
				docWords[clean] = true
			}
		}
		
		if len(docWords) > 0 {
			documents = append(documents, docWords)
		}
	}

	if len(documents) == 0 {
		return nil, errors.New("No valid words found in file")
	}

	stats := make([]schema.WordStat, 0, len(wordCount))
	totalDocuments := len(documents)
	
	for word, tf := range wordCount {
		df := 0
		for _, doc := range documents {
			if doc[word] {
				df++
			}
		}
		var idf float64
		if df == 0 {
			idf = 0
		} else {
			idf = math.Log(float64(totalDocuments) / float64(df))
		}
		
		stats = append(stats, schema.WordStat{
			Word: word,
			TF:   float64(tf),
			IDF:  idf,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].IDF > stats[j].IDF
	})

	fileID := generateFileID()
	err = s.repos.InsertDocument(repository.Document{
		FileID: fileID,
		Name: file.Filename,
		UserID: userID,
		Content: content.String(),
		Stats: stats,
		Length: totalWords,
		Collections: make([]string, 0),
		IsValid: false,
		Words: make(map[string]int),
	})

	return &schema.UploadResponse{
		SessionID: fileID,
		Page: 1,
		Words: getPage(stats, 1),
		Total: len(stats),
	}, nil
}

func (s *DocumentService) GetDocument(fileID string, userID string) (string, error) {
	return s.repos.GetDocument(fileID, userID)
}
func (s *DocumentService) GetListDocuments(userID string) ([]string, error) {
	return s.repos.GetListDocuments(userID)
}
func (s *DocumentService) GetDocumentStats(fileID string, userID string) ([]schema.WordStat, error) {
	return s.repos.GetDocumentStats(fileID, userID)
}
func (s *DocumentService) DeleteDocument(fileID, userID string) error {
	return s.repos.DeleteDocument(fileID, userID)
}
func (s *DocumentService) DeleteUserDocuments(userID string) error {
	return s.repos.DeleteAllDocuments(userID)
}

func generateFileID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}