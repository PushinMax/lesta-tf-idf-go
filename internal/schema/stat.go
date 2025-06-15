package schema


// Структура для хранения результатов
type WordStat struct {
	Word string `bson:"word" json:"word"`
	TF   float64    `bson:"tf" json:"tf"`
	IDF  float64`bson:"idf" json:"idf"`
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