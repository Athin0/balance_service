package main

import (
	"fmt"
	"log"
	"os"
	"useIt/Formatters"
)

func main() {
	fileOutName := "result/hello.html"
	fileInName := "data/data.prn"

	out, err := CreateFile(fileOutName)
	if err != nil {
		os.Exit(1)
	}
	defer out.Close()
	in, err := OpenFile(fileInName)
	if err != nil {
		os.Exit(1)
	}
	defer in.Close()
	MakeConvert(in, out, fileInName)

}

// CreateFile создает файл по заданному имени
func CreateFile(nameFile string) (*os.File, error) {
	file, err := os.Create(nameFile)
	if err != nil {
		log.Println("Unable to create file:", err)
		return nil, err
	}
	return file, nil
}

// OpenFile открывает файл по заданному имени
func OpenFile(nameFile string) (*os.File, error) {
	in, err := os.Open(nameFile)
	if err != nil {
		log.Println("Unable to open file:", err)
		return nil, err
	}
	return in, nil
}

// MakeConvert выбирает какой конвертер использовать
func MakeConvert(in *os.File, out *os.File, NameInpFile string) {
	var formatter Formatters.Formatter
	lenN := len(NameInpFile)
	if lenN <= 3 {
		fmt.Println("file name <= 3 : input format not supported")
		return
	}
	format := NameInpFile[lenN-3:] //проверяем формат по окончанию, можно было отделить конец через split

	switch format {
	case "csv":
		formatter = Formatters.CreateCSVformatter(in, out)
	case "prn":
		formatter = Formatters.CreatePRNformatter(in, out)
	default:
		fmt.Println("input format not supported ")
		return
	}
	formatter.Format()
}
