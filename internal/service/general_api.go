package service

import (
	"bufio"
	"errors"
	"math"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PushinMax/lesta-tf-idf-go/internal/schema"
	"github.com/PushinMax/lesta-tf-idf-go/internal/session"
)

type GenenalService struct {
	session *session.Session
}

func newGeneralApi(session *session.Session) *GenenalService {
	return &GenenalService{session: session}
}

func (s *GenenalService) UploadFile(file *multipart.FileHeader) (*schema.UploadResponse, error) {
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("Cannot open file")
	}
	defer src.Close()

	scanner := bufio.NewScanner(src)
	
	documents := make([]map[string]bool, 0)
	wordCount := make(map[string]int)
	totalWords := 0

	for scanner.Scan() {
		line := scanner.Text()
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

	sessionID := generateSessionID()
	/*
	dataMutex.Lock()
	sessionData[sessionID] = stats
	dataMutex.Unlock()
	*/
	err = s.session.SetState(sessionID, stats)
	if err != nil {
		return nil, err
	}

	return &schema.UploadResponse{
		SessionID: sessionID,
		Page: 1,
		Words: getPage(stats, 1),
		Total: len(stats),
	}, nil
}

func (s *GenenalService) GetPageData(sessionID string, page int) (*schema.PageResponse, error) {
	stats, exists := s.session.GetState(sessionID)

	if !exists {
		return nil, errors.New("Session not found")
	}

	return &schema.PageResponse{
		Words: getPage(stats, page),
		Total: len(stats),
	}, nil
}


func cleanWord(s string) string {
	s = strings.ToLower(s)
	return strings.Trim(s, ".,!?;:\"()[]{}")
}

func generateSessionID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func getPage(stats []schema.WordStat, page int) []schema.WordStat {
	itemsPerPage := 50
	start := (page - 1) * itemsPerPage
	if start >= len(stats) {
		return []schema.WordStat{}
	}
	
	end := start + itemsPerPage
	if end > len(stats) {
		end = len(stats)
	}
	return stats[start:end]
}