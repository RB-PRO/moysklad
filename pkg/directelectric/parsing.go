package directelectric

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const URL string = "https://www.directelectric.ru"

// Весь массив с данными товаров
type DirectelEctricObjects struct {
	Data []Product // Массив с товарами, который используется для загрузки товаров на МойСклад
}

// Структура карточки товара
//
// Главная особенность - гарантированные данные
type Product struct {
	Link           string            // Ссылка на товар
	ImageLink      string            // Ссылка на фото товара
	NameFull       string            // Полное Название товара
	Code           string            // Код
	NameFew        string            // Краткое Название товара
	Description    string            // Описание товара
	Article        string            // Артикул
	Price          float64           // Цена
	Availability   bool              // Наличие
	Dimension      string            // Размерность(шт, кг, тонна)
	Weight         float64           // Вес
	Volume         float64           // Объём
	Specifications map[string]string // Остальные характеристики
}

// Изменить ссылку с параметром запрашиваемой страницы
func MakeLinkWithPage(link string, pagesString string, page int) (string, error) {
	urlA, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	values := urlA.Query()
	if strings.Contains(link, "PAGEN_1") {
		pagesString = "PAGEN_1"
	}
	if strings.Contains(link, "PAGEN_2") {
		pagesString = "PAGEN_2"
	}
	if pagesString != "" {
		values.Set(pagesString, strconv.Itoa(page)) // Страница
	}

	//values.Set("nal", "n")                    // В наличии = false(Значит будут показываеться товары как в наличии, так и нет)
	values.Set("limit", "16") // Показывать по 16 товаров на странице. Это необходимо для корректного отображения всех товаров. ПРи 24 отображается всего 20

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
