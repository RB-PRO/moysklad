package app

import (
	"fmt"
	"io"
	"os"

	"github.com/RB-PRO/moysklad/pkg/directelectric"
)

// Тестовый запуск по всем категориям
func RunAllCategory() {

	// Получить все категории на директ электрике
	links := directelectric.ParseCatalogs()

	fmt.Println("Ссылки для парсинга: ", links)

	// Определение структуры
	var items directelectric.DirectelEctricObjects

	// Пропарсить в подкатегории
	items.ParseItems(links)

	// Пропарсить карточки товаров
	items.ParseAllItem()

	items.SaveXlsx("directelectric")
}

// Спасить весь сайт и загрузить по разным Xlsx файлам
func ParseItemsAndSave() {
	// Получить все категории на директ электрике
	links := directelectric.ParseCatalogs()

	fmt.Println("Ссылки для парсинга: ", links)

	fmt.Println("len", len(links))

	links = links[9:]

	directelectric.ParseItemsAndSaveAnotherXlsx(links)
}

// Тестовый запуск по link розеток
func RunOneLink() {

	link, _ := dataFile("link")
	links := make([]string, 1)
	links[0] = link

	// Определение структуры
	var items directelectric.DirectelEctricObjects

	// Пропарсить в подкатегории
	items.ParseItems(links)

	// Пропарсить карточки товаров
	items.ParseAllItem()

	items.SaveXlsx("directelectric")
}

// Получение значение из файла
func dataFile(filename string) (string, error) {
	// Открыть файл
	fileToken, errorToken := os.Open(filename)
	if errorToken != nil {
		return "", errorToken
	}

	// Прочитать значение файла
	data := make([]byte, 512)
	n, err := fileToken.Read(data)
	if err == io.EOF { // если конец файла
		return "", errorToken
	}
	fileToken.Close() // Закрытие файла

	return string(data[:n]), nil
}
