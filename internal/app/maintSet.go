package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RB-PRO/moysklad/pkg/directelectric"
)

// Спасить ВСЕ [товары] из ДиректЭлектрика и загрузить из в МойСклад в [список товаров].
//
// [товары]: https://www.directelectric.ru/catalog/product/180142/
// [список товаров]: https://github.com/dotnow/moysklad/issues/3
func ParseAllObjectAndLoadToMoySklad() {
	// Получить все категории на директ электрике
	links := directelectric.ParseCatalogs()
	ParseLinks(links)
}

// Спасить [товары] по введённой ссылке и загрузить из в МойСклад в [список товаров].
//
// [товары]: https://www.directelectric.ru/catalog/product/180142/
// [список товаров]: https://github.com/dotnow/moysklad/issues/3
func ParseLinkAndLoadToMoySklad() {
	// Получить все категории на директ электрике
	//links := directelectric.ParseCatalogs()
	var errorReader error
	links := make([]string, 1)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Вставьте ссылку: ")
	links[0], errorReader = reader.ReadString('\n')
	if errorReader != nil {
		log.Fatalln(errorReader)
	}
	//links[0] = "https://www.directelectric.ru/catalog/?q=brite&s"
	links[0] = strings.ReplaceAll(links[0], directelectric.URL, "")
	links[0] = strings.TrimSpace(links[0])
	//links[0], _ = dataFile("link")
	//fmt.Println(">" + links[0] + "<")
	ParseLinks(links)
}

func Schuse() {
	fmt.Print("1. Пропарсить всё\n2. Пропарсить по ссылке\n > ")
	reader := bufio.NewReader(os.Stdin)
	text, errorReader := reader.ReadString('\n')
	if errorReader != nil {
		log.Fatalln(errorReader)
	}
	text = strings.TrimSpace(text)
	if text == "1" {
		ParseAllObjectAndLoadToMoySklad()
	} else if text == "2" {
		ParseLinkAndLoadToMoySklad()
	} else {
		fmt.Println("Неверный ввод. Перезапустите меня.")
	}

	fmt.Println("Press 'q' to quit")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		exit := scanner.Text()
		if exit == "q" {
			break
		} else {
			fmt.Println("Press 'q' to quit")
		}
	}
}
