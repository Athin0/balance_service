package Formatters

import "strings"

type Formatter interface {
	Format()
}

// MakeFormat4String преобразует строку в соответствии с выходным форматом, а именно "имя_поля":"данные_поля"
func MakeFormat4String(title string, data []rune) (ans string) {
	ans = "\"" + title + "\"" + ":\"" + strings.TrimSpace(string(data)) + "\""
	return
}
