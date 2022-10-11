package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	file, err := os.Create("hello.html")

	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("Done.")

	in, err := os.Open("data/data.prn")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer in.Close()

	data := make([]byte, 64)

	for {
		n, err := in.Read(data)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		text := string(data[:n])
		//text = strings.Replace(text, " ", ",", -1)
		//text = strings.Trim(text, ",")
		r := regexp.MustCompile("\\s+")
		r2 := regexp.MustCompile("\\n+")
		r3 := regexp.MustCompile(";")
		text = r2.ReplaceAllString(text, ";")
		text = r.ReplaceAllString(text, " ")
		text = r.ReplaceAllString(text, ",")
		text = r3.ReplaceAllString(text, "\n")
		file.WriteString(text)
		fmt.Print(strings.TrimSpace(text))

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
