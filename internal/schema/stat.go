package schema


// Структура для хранения результатов
type WordStat struct {
	Word string
	TF   int     // Term Frequency (частота слова в документе)
	IDF  float64 // Inverse Document Frequency (обратная частота документа)
}


type UploadResponse struct {
	SessionID string
	Page      int
	Words     []WordStat
	Total int
}

type PageResponse struct {
	Words     []WordStat
	Total int
}