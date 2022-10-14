package main

import (
	"log"
	"os"
	"useIt/Formatters"
)

func main() {
	fileOutName := "result/hello.html"
	fileInName := "data/data.csv"

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

func CreateFile(nameFile string) (*os.File, error) {
	file, err := os.Create(nameFile)
	if err != nil {
		log.Println("Unable to create file:", err)
		return nil, err
	}
	return file, nil
}
func OpenFile(nameFile string) (*os.File, error) {
	in, err := os.Open(nameFile)
	if err != nil {
		log.Println("Unable to open file:", err)
		return nil, err
	}
	return in, nil
}

func MakeConvert(in *os.File, out *os.File, NameInpFile string) {
	var formatter Formatters.Formatter
	lenN := len(NameInpFile)
	if lenN <= 3 {
		return
	}
	format := NameInpFile[lenN-3:]
	switch format {
	case "csv":
		formatter = Formatters.CreateCSVformatter(in, out)
	case "prn":
		formatter = Formatters.CreatePRNformatter(in, out)
	}
	formatter.Format()
}
