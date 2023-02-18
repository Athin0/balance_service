package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func MakeReport(arr *[]Report) (string, error) {
	fileOutName := "report.csv"

	r, err := CreateFile(fileOutName)
	if err != nil {
		return "", fmt.Errorf("err in make report: %w", err)
	}
	defer r.Close()

	elem := []string{"название услуги ", "общая сумма выручки за отчетный период"}

	w := csv.NewWriter(r)
	defer w.Flush()
	if err := w.Write(elem); err != nil {
		return "", fmt.Errorf("error writing record to file: %w", err)
	}
	for _, record := range *arr {
		elem = []string{strconv.FormatInt(record.ServiceId, 10), strconv.FormatFloat(record.Sum, 'f', -1, 64)}
		if err := w.Write(elem); err != nil {
			return "", fmt.Errorf("error writing record to file: %w", err)
		}
	}

	return fileOutName, nil
}
func CreateFile(nameFile string) (*os.File, error) {
	file, err := os.Create(nameFile)
	if err != nil {
		log.Println("Unable to create file:", err)
		return nil, err
	}
	return file, nil
}
