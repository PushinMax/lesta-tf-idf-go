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


type FileService struct {
	repos *repository.Repository
}

func newFileApi(repos *repository.Repository) *FileService {
	return &FileService{
		repos: repos,
	}
}

func (s *FileService) UploadFile(file *multipart.FileHeader, userID string) (*schema.UploadResponse, error) {
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
			TF:   tf,
			IDF:  idf,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].IDF > stats[j].IDF
	})

	fileID := generateFileID()
	err = s.repos.InsertFile(repository.FileDocument{
		FileID: fileID,
		Name: file.Filename,
		UserID: userID,
		Content: content.String(),
		Stats: stats,
		Length: totalWords,
	})

	return &schema.UploadResponse{
		SessionID: fileID,
		Page: 1,
		Words: getPage(stats, 1),
		Total: len(stats),
	}, nil
}

func (s *FileService) GetFile(fileID string, userID string) (string, error) {
	return s.repos.GetFile(fileID, userID)
}
func (s *FileService) GetListFiles(userID string) ([]string, error) {
	return nil, nil
}
func (s *FileService) GetFilesStats(fileID string, userID string) ([]schema.WordStat, error) {
	return nil, nil
}
func (s *FileService) DeleteFile(fileID, userID string) error {
	return nil
}
func (s *FileService) DeleteUserFiles(userID string) (int64, error) {
	return 0, nil
}

func generateFileID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}