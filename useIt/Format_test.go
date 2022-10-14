package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"testing"
)

func TestCSVformatter_Format(t *testing.T) {
	fileOutName := "test_Out.html"
	fileInName := "test_In.csv"
	text := "Name,Address,Postcode,Mobile,Limit,Birthday\n\"Oliver, El\",\"Via Archimede, 103-91\",2343aa,000 1119381,6000000,01/01/1999\n\"Harry\",Leonardo da Vinci 1,4532 AA,010 1118986,343434,31/12/1965\n\"Jack\",\"Via Rocco Chinnici 4d\",3423 ba,0313-111475,22,05/04/1984\n\"Noah\",\"Via Giannetti, 4-32\",2340 CC,28932222,434,03/10/1964\n\"Charlie\",\"Via Aldo Moro, 7\",3209 DD,30-34563332,343.8,04/10/1954\n\"Mia\",\"Via Due Giugno, 12-1\",4220 EE,43433344329,6343.6,10/08/1980\n\"Lilly\",ArcisstraЯe 21,12343,+44 728 343434,34342.3,20/10/1997"
	outText := "\"Name\":\"Oliver, El\",\"Address\":\"Via Archimede, 103-91\",\"Postcode\":\"2343aa\",\"Mobile\":\"000 1119381\",\"Limit\":\"6000000\",\"Birthday\":\"01/01/1999\"\n\"Name\":\"Harry\",\"Address\":\"Leonardo da Vinci 1\",\"Postcode\":\"4532 AA\",\"Mobile\":\"010 1118986\",\"Limit\":\"343434\",\"Birthday\":\"31/12/1965\"\n\"Name\":\"Jack\",\"Address\":\"Via Rocco Chinnici 4d\",\"Postcode\":\"3423 ba\",\"Mobile\":\"0313-111475\",\"Limit\":\"22\",\"Birthday\":\"05/04/1984\"\n\"Name\":\"Noah\",\"Address\":\"Via Giannetti, 4-32\",\"Postcode\":\"2340 CC\",\"Mobile\":\"28932222\",\"Limit\":\"434\",\"Birthday\":\"03/10/1964\"\n\"Name\":\"Charlie\",\"Address\":\"Via Aldo Moro, 7\",\"Postcode\":\"3209 DD\",\"Mobile\":\"30-34563332\",\"Limit\":\"343.8\",\"Birthday\":\"04/10/1954\"\n\"Name\":\"Mia\",\"Address\":\"Via Due Giugno, 12-1\",\"Postcode\":\"4220 EE\",\"Mobile\":\"43433344329\",\"Limit\":\"6343.6\",\"Birthday\":\"10/08/1980\"\n\"Name\":\"Lilly\",\"Address\":\"ArcisstraЯe 21\",\"Postcode\":\"12343\",\"Mobile\":\"+44 728 343434\",\"Limit\":\"34342.3\",\"Birthday\":\"20/10/1997\"\n"
	in, err := CreateFile(fileInName) //создаем входной файл
	if err != nil {
		os.Exit(1)
	}
	defer in.Close()
	_, err = in.WriteString(text) //добавляем во входной файл данные
	if err != nil {
		log.Println(err)
		return
	}

	out, err := CreateFile(fileOutName) //создаем выходной файл (пустой)
	if err != nil {
		os.Exit(1)
	}
	defer os.Remove(fileOutName)
	defer out.Close()
	//_, err = out.Write([]byte(text))

	in, _ = os.Open(fileInName)
	MakeConvert(in, out, fileInName)
	out, _ = os.Open(fileOutName)
	raw := make([]byte, 64)
	ans := ""
	for {
		n, err := out.Read(raw)
		if err == io.EOF {
			break
		}
		ans += string(raw[:n])
	}
	if ans == outText {

	} else {
		log.Println("Dif text:", ans)
		assert.Error(t, fmt.Errorf("different text"))
	}

}

func TestPRNformatter_Format(t *testing.T) {
	fileOutName := "test_Out.html"
	fileInName := "test_In.prn"
	text := "First name      Address               Postcode Mobile               Limit Birthday\nOliver          Via Archimede, 103-91 2343aa   000 1119381        6000000 19570101\nHarry           Leonardo da Vinci 1   4532 AA  010 1118986       10433301 19751203\nJack            Via Rocco Chinnici 4d 3423 ba  0313-111475          93543 19740604\nNoah            Via Giannetti, 4-32   2340 CC  28932222                34 19940906\nCharlie         Via Aldo Moro, 7      3209 DD  30-34563332           4531 19981107\nMia             Via Due Giugno, 12-1  4220 EE  43433344329           9087 19700515\nLily            Arcisstra�e 21        12343    +44 728 343434      765599 19971003\n"
	outText := "\"First name\":\"Oliver\",\"Address\":\"Via Archimede, 103-91\",\"Postcode\":\"2343aa\",\"Mobile\":\"000 1119381\",\"Limit\":\"6000000\",\"Birthday\":\"19570101\"\n\"First name\":\"Harry\",\"Address\":\"Leonardo da Vinci 1\",\"Postcode\":\"4532 AA\",\"Mobile\":\"010 1118986\",\"Limit\":\"10433301\",\"Birthday\":\"19751203\"\n\"First name\":\"Jack\",\"Address\":\"Via Rocco Chinnici 4d\",\"Postcode\":\"3423 ba\",\"Mobile\":\"0313-111475\",\"Limit\":\"93543\",\"Birthday\":\"19740604\"\n\"First name\":\"Noah\",\"Address\":\"Via Giannetti, 4-32\",\"Postcode\":\"2340 CC\",\"Mobile\":\"28932222\",\"Limit\":\"34\",\"Birthday\":\"19940906\"\n\"First name\":\"Charlie\",\"Address\":\"Via Aldo Moro, 7\",\"Postcode\":\"3209 DD\",\"Mobile\":\"30-34563332\",\"Limit\":\"4531\",\"Birthday\":\"19981107\"\n\"First name\":\"Mia\",\"Address\":\"Via Due Giugno, 12-1\",\"Postcode\":\"4220 EE\",\"Mobile\":\"43433344329\",\"Limit\":\"9087\",\"Birthday\":\"19700515\"\n\"First name\":\"Lily\",\"Address\":\"Arcisstra�e 21\",\"Postcode\":\"12343\",\"Mobile\":\"+44 728 343434\",\"Limit\":\"765599\",\"Birthday\":\"19971003\"\n"
	in, err := CreateFile(fileInName) //создаем входной файл
	if err != nil {
		os.Exit(1)
	}
	defer os.Remove(fileInName)
	defer in.Close()
	_, err = in.WriteString(text) //добавляем во входной файл данные
	if err != nil {
		log.Println(err)
		return
	}

	out, err := CreateFile(fileOutName) //создаем выходной файл (пустой)
	if err != nil {
		os.Exit(1)
	}
	defer os.Remove(fileOutName)
	defer out.Close()

	in, _ = os.Open(fileInName)
	MakeConvert(in, out, fileInName)
	out, _ = os.Open(fileOutName)
	raw := make([]byte, 64)
	ans := ""
	for {
		n, err := out.Read(raw)
		if err == io.EOF {
			break
		}
		ans += string(raw[:n])
	}
	if ans == outText {

	} else {
		log.Println("Dif text:", ans)
		assert.Error(t, fmt.Errorf("different text"))
	}

}
