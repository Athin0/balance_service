package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	file, err := CreateFile("hello.html")
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	in, err := OpenFile("data/data.prn")
	if err != nil {
		os.Exit(1)
	}
	defer in.Close()

	columns, pre, err := GetString(in)
	if err != nil {
		log.Println("Err in Get String from file:", err)
	}
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	splitedText := re.FindAllString(columns, -1)
	ColumnsLen := make([]int, 0, 0)
	for _, elem := range splitedText {
		ColumnsLen = append(ColumnsLen, len(elem))
	}

	fmt.Println(ColumnsLen)
	for {
		a1, a2, err := GetString(in)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		//textLine := pre + a1 + "\n"
		//pre = a2
		//text = strings.Replace(text, " ", ",", -1)
		//text = strings.Trim(text, ",")
		//r := regexp.MustCompile("\\s+")
		//r2 := regexp.MustCompile("\\n+")
		//r3 := regexp.MustCompile(";")
		//text = r2.ReplaceAllString(text, ";")
		//text = r.ReplaceAllString(text, " ")
		//text = r.ReplaceAllString(text, ",")
		//text = r3.ReplaceAllString(text, "\n")
		file.WriteString(textLine)

	}
	/*
		in, err = os.Open("data.csv")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for {
			n, err := in.Read(data)
			if err == io.EOF { // если конец файла
				break // выходим из цикла
			}
			file.Write(data[:n])
			fmt.Print(string(data[:n]))
		}

	*/
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
func GetString(in *os.File) (string, string, error) {
	text1 := ""
	data := make([]byte, 64)
	for !strings.ContainsAny(text1, "\n") {
		n, err := in.Read(data)
		if err != nil {
			return "", "", err
		}
		text1 += string(data[:n])
	}
	ans := strings.Split(text1, "\n")
	var first, second string
	if len(ans) > 0 {
		first = ans[0]
	}
	if len(ans) > 1 {
		second = ans[1]
	}
	return first, second, nil
}

func ProcessingText(text string, columns []int) string {
	arr := make([]string, 0)
	n := 0
	for _, elem := range columns {
		t := text[n : n+elem]
		println(t)
		arr = append(arr, t)
		n += elem
	}
	return text
}
