package directelectric

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

const URL string = "https://www.directelectric.ru"

// Весь массив с данными товаров
type DirectelEctricObjects struct {
	Data []Product
}
type Product struct {
	Link           string            // Ссылка на товар
	NameFull       string            // Полное Название товара
	NameFew        string            // Краткое Название товара
	Article        string            // Артикул
	Price          float64           // Цена
	Specifications map[string]string // Остальные характеристики
}

func MakeLinkWithPage(link string, page int) (string, error) {
	urlA, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	values := urlA.Query()
	values.Set("PAGEN_1", strconv.Itoa(page))

	urlA.RawQuery = values.Encode()

	return urlA.String(), nil
}

// Функция парсит [страницу] по определённому page
//
// [страницу]: https://www.directelectric.ru/catalog/rozetki-i-vyklyuchateli/filter/vendor_new-is-schneider%20electric/serial-is-atlasdesign/apply/?PAGEN_1=2&nal=y
func (items *DirectelEctricObjects) ParseItems(link string) {
	fmt.Println("Parse", link)

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("div[class=item]", func(e *colly.HTMLElement) {
		var item Product
		item.NameFew = e.DOM.Find("a[class^=item__title]").Text()
		item.Link, _ = e.DOM.Find("a[class^=item__title]").Attr("href")
		items.Data = append(items.Data, item)
	})

	c.Visit(link)

}

// Функция парсит каждую [страницу] товаров
//
// [страницу]: https://www.directelectric.ru/catalog/product/180150/
func (items *DirectelEctricObjects) ParseAllItem(link string) {
	fmt.Println("Parse Item", link)

	c := colly.NewCollector()
	var itemIndex int = 0

	// Характеристики
	c.OnHTML("div[class=characteristics__list] div[class^=characteristics__item]", func(e *colly.HTMLElement) {
		key := e.DOM.Find("div[class=characteristics__name] span").Text()
		value := e.DOM.Find("div[class=characteristics__value]").Text()
		switch key {
		case "":
			break
		case "Артикул":
			items.Data[itemIndex].Article = value // Заполняем артикул
		default:
			items.Data[itemIndex].Specifications[key] = value // Заполняем мапу характеристиками товара
		}
	})

	// Цена
	c.OnHTML("div[class=buy__price] span div span:nth-of-type(1)", func(e *colly.HTMLElement) {
		cost := e.DOM.Text()
		fmt.Println("cost", cost)
		reg := regexp.MustCompile("0123456789,")
		replaceStr := reg.ReplaceAllString(cost, "")
		fmt.Println("replaceStr", replaceStr)
	})

	items.Data[itemIndex].Specifications = make(map[string]string) // Выделяем память в мапу
	c.Visit(URL + items.Data[0].Link)

	/*
		for _, itemVal := range items.Data {
			c.Visit(URL + itemVal.Link)
			itemIndex++
		}
	*/

}
