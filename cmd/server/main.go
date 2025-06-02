package main

import (
	"bufio"
	//"fmt"
	"log"
	"math"

	// "os"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Структура для хранения результатов
type WordStat struct {
	Word string
	TF   int     // Term Frequency (частота слова в документе)
	IDF  float64 // Inverse Document Frequency (обратная частота документа)
}

var (
	sessionData = make(map[string][]WordStat)
	dataMutex   sync.RWMutex
)

func ceilDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return (a + b - 1) / b
}

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"ceilDiv": ceilDiv,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.POST("/upload", handleUpload)
	r.GET("/data/:session/:page", getPageData)

	log.Println("Server started on :8080")
	r.Run(":8080")
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot open file"})
		return
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
		c.JSON(400, gin.H{"error": "No valid words found in file"})
		return
	}

	stats := make([]WordStat, 0, len(wordCount))
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
		//fmt.Println(totalDocuments)
		stats = append(stats, WordStat{
			Word: word,
			TF:   tf,
			IDF:  idf,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].IDF > stats[j].IDF
	})

	sessionID := generateSessionID()
	dataMutex.Lock()
	sessionData[sessionID] = stats
	dataMutex.Unlock()

	c.HTML(200, "results.html", gin.H{
		"SessionID": sessionID,
		"Page":      1,
		"Words":     getPage(stats, 1),
		"Total":     len(stats),
	})
}

func getPageData(c *gin.Context) {
	sessionID := c.Param("session")
	page, _ := strconv.Atoi(c.Param("page"))

	dataMutex.RLock()
	stats, exists := sessionData[sessionID]
	dataMutex.RUnlock()

	if !exists {
		c.JSON(404, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(200, gin.H{
		"words": getPage(stats, page),
		"total": len(stats),
	})
}

func cleanWord(s string) string {
	s = strings.ToLower(s)
	return strings.Trim(s, ".,!?;:\"()[]{}")
}

func generateSessionID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func getPage(stats []WordStat, page int) []WordStat {
	itemsPerPage := 50
	start := (page - 1) * itemsPerPage
	if start >= len(stats) {
		return []WordStat{}
	}
	
	end := start + itemsPerPage
	if end > len(stats) {
		end = len(stats)
	}
	return stats[start:end]
}