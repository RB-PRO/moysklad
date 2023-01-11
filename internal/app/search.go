package app

import (
	"fmt"
	"io"
	"os"

	"github.com/RB-PRO/moysklad/pkg/directelectric"
)

func Run() {

	link, _ := dataFile("link")

	linkPages, _ := directelectric.MakeLinkWithPage(link, 2)

	var items directelectric.DirectelEctricObjects

	items.ParseItems(linkPages)

	items.ParseAllItem(linkPages)

	fmt.Println(items.Data[0].Specifications)
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
