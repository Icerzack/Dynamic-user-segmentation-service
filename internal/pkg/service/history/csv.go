package history

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

type CSVHistoryService struct {
}

func NewCSVHistoryService() *CSVHistoryService {
	return &CSVHistoryService{}
}

func (s *CSVHistoryService) WriteToFile(userID int, segmentTitle string, operationName string, date time.Time) {
	record := []string{
		strconv.Itoa(userID), segmentTitle, operationName, date.String(),
	}

	f, err := os.OpenFile("docs/history.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to file", err)
	}

}
