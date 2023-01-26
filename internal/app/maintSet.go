package app

import (
	"fmt"

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
	//var errorReader error
	links := make([]string, 1)

	/*
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		links[0], errorReader = reader.ReadString('\n')
		if errorReader != nil {
			log.Fatalln(errorReader)
		}
		links[0] = strings.ReplaceAll(links[0], directelectric.URL, "")

	*/
	links[0], _ = dataFile("link")
	fmt.Println(">" + links[0] + "<")
	ParseLinks(links)
}
