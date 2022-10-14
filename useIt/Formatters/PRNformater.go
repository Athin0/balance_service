package Formatters

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type PRNformatter struct {
	in  *os.File
	out *os.File
}

func CreatePRNformatter(in *os.File, out *os.File) *PRNformatter {
	return &PRNformatter{in, out}
}

func (f *PRNformatter) Format() {
	columns, pre, err := GetStrings(f.in) //читаем первую строку, чтобы найти заголовки
	if err != nil {
		log.Println("Err in Get String from file:", err)
		return
	}
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	splittedText := re.FindAllString(columns, -1) //получаем список заголовков

	for i, el := range splittedText {
		splittedText[i] = strings.TrimSpace(el) //избавляемся от ненужных пробелов
	}
	//fmt.Println(splittedText)
	indexesEndOfColumns := []int{16, 38, 47, 61, 74, 82}
	for {
		a1, a2, err := GetStrings(f.in)
		if err == io.EOF { // если конец файла выходим из цикла
			break
		}
		textLine := pre + a1 //объединяем считанное в предыдущей итерации начало строки и конец, считанный в текущей
		pre = a2             //обновляем оставшийся хвост(начало строки след итерации)

		prepareArr, err := func(raw string) (line []string, err error) {
			runes := []rune(raw)
			if len(runes) < 74 { //не хватает данных для данного формата файла
				err = fmt.Errorf("ReadPrnLine detected Wrong data -> %s", raw)
				return
			}

			begin := 0
			for i, title := range splittedText { //добавляем поочереди каждое поле строки(отформатированное)
				end := indexesEndOfColumns[i] //конец поля
				line = append(line, MakeFormat4String(title, runes[begin:end]))
				begin = end
			}
			return
		}(textLine)
		if err != nil {
			log.Println(err)
			return
		}
		ans := strings.Join(prepareArr, ",") + "\n" //объединение в одну строку для записи в файл + перенос строки

		_, err = f.out.WriteString(ans)
		if err != nil {
			log.Println("Err in writeString:", err)
			return
		}

	}
}

// GetStrings получает часть искомой строки и часть следующей, все переносится вследствие странного считывания файла
func GetStrings(in *os.File) (string, string, error) {
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

// MakeFormat4String преобразует строку в соответствии с выходным форматом, а именно "имя_поля":"данные_поля"
func MakeFormat4String(title string, data []rune) (ans string) {
	ans = "\"" + title + "\"" + ":\"" + strings.TrimSpace(string(data)) + "\""
	return
}
