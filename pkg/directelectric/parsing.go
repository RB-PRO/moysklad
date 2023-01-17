package directelectric

import (
	"net/url"
	"strconv"

	"github.com/gocolly/colly"
)

const URL string = "https://www.directelectric.ru"

// Весь массив с данными товаров
type DirectelEctricObjects struct {
	Data []Product
}

// Структура карточки товара
//
// Главная особенность - гарантированные данные
type Product struct {
	Link           string            // Ссылка на товар
	imageLink      string            // Ссылка на фото товара
	NameFull       string            // Полное Название товара
	NameFew        string            // Краткое Название товара
	Description    string            // Описание товара
	Article        string            // Артикул
	Price          float64           // Цена
	Availability   bool              // Наличие
	Dimension      string            // Размерность(шт, кг, тонна)
	Specifications map[string]string // Остальные характеристики
}

// Изменить ссылку с параметром запрашиваемой страницы
func MakeLinkWithPage(link string, page int) (string, error) {
	urlA, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	values := urlA.Query()
	values.Set("PAGEN_1", strconv.Itoa(page)) // Страница
	values.Set("nal", "n")                    // В наличии = false(Значит будут показываеться товары как в наличии, так и нет)
	values.Set("limit", "16")                 // Показывать по 16 товаров на странице. Это необходимо для корректного отображения всех товаров. ПРи 24 отображается всего 20

	urlA.RawQuery = values.Encode()

	return urlA.String(), nil
}

// Функция парсит [каталог] с ссылками на все подкатегории
//
// [каталог]: https://www.directelectric.ru/catalog/
func ParseCatalogs() []string {
	c := colly.NewCollector() // Создаём экземпляр для библиотеки gocolly

	var links []string

	// Идём по категориям товара
	c.OnHTML("a[class^=catalog-top__main-section]", func(e *colly.HTMLElement) {
		link, linksExit := e.DOM.Attr("href")
		//fmt.Println(link, linksExit)
		if linksExit {
			links = append(links, link)
		}
	})

	c.Visit(URL + "/catalog/")

	return links
}
