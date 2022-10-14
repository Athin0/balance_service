package Formatters

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// CSVformatter Хранит данные о входном и выходном файле
type CSVformatter struct {
	in  *os.File
	out *os.File
}

func CreateCSVformatter(in *os.File, out *os.File) *CSVformatter {
	return &CSVformatter{in, out}
}

// Format преобразует csv в html
//данные о входном и выходном файле хранятся в структуре CSVformatter
func (f *CSVformatter) Format() {
	var columns []string
	reader := csv.NewReader(f.in)
	reader.Comma = ',' //задаем разделитель
	n := true
	for {
		record, e := reader.Read()
		if e == io.EOF { // если конец файла выходим из цикла
			break
		}
		if e != nil {
			fmt.Println("Err in read file:", e)
			break
		}
		if n { //находим названия полей
			n = false
			columns = record
			continue
		}

		prepareArr, err := func(raw []string) (line []string, err error) {
			for i, word := range raw { // добавляем поочереди каждое поле строки (отформатированное)
				line = append(line, MakeFormat4String(columns[i], []rune(word))) //собираем выходную строку в нужном формате
			}
			return
		}(record)
		if err != nil {
			log.Println(err)
			return
		}
		ans := strings.Join(prepareArr, ",") + "\n" //объединение в одну строку для записи в файл + перенос строки

		_, err = f.out.WriteString(ans) //записываем в исходный файл
		if err != nil {
			log.Println("Err in writeString:", err)
			return
		}

	}
}

