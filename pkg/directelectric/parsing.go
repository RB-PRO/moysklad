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

	var next bool = true
	var schetchik int = 1

	c := colly.NewCollector()

	// Карточки товара
	c.OnHTML("div[class=item]", func(e *colly.HTMLElement) {
		var item Product
		item.NameFew = e.DOM.Find("a[class^=item__title]").Text()
		item.Link, _ = e.DOM.Find("a[class^=item__title]").Attr("href")
		items.Data = append(items.Data, item)
	})

	// Проверить можно ли дальше листать
	c.OnHTML("[class=pagination__next]", func(e *colly.HTMLElement) {
		//fmt.Println(e.DOM.Text())
		hrefNext, hrefNextIsExit := e.DOM.Attr("href")
		fmt.Println(">>", hrefNext, hrefNextIsExit)
		if !hrefNextIsExit {
			next = false
		} else {
			u, err := url.Parse(URL + hrefNext)
			if err != nil {
				panic(err)
			}

			m, _ := url.ParseQuery(u.RawQuery)

			//fmt.Println("m[PAGEN_1][0]", m["PAGEN_1"][0], "schetchik", schetchik)
			if m["PAGEN_1"][0] == strconv.Itoa(schetchik) {
				next = false
			}
		}
	})

	for {
		if !next {
			break
		}
		// Делаем ссылку со страницей
		linkPages, _ := MakeLinkWithPage(link, schetchik)

		// Парсим
		c.Visit(linkPages)
		schetchik++
	}

}

// Функция парсит каждую [страницу] товара
//
// [страницу]: https://www.directelectric.ru/catalog/product/180150/
func (items *DirectelEctricObjects) ParseAllItem(link string) {
	//fmt.Println("Parse Item", link)

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
	c.OnHTML("div[class=buy__price] div span:first-of-type", func(e *colly.HTMLElement) {
		cost := e.DOM.Text()
		//fmt.Println("cost", cost)
		reg := regexp.MustCompile("[^0-9.]+")
		replaceStr := reg.ReplaceAllString(cost, "")
		//fmt.Println("replaceStr", replaceStr)
		if n, err := strconv.ParseFloat(replaceStr, 64); err == nil {
			items.Data[itemIndex].Price = n
		}
	})

	for itemIndex, itemVal := range items.Data {
		fmt.Println("Parse Item: ", itemIndex, "/", len(items.Data))
		items.Data[itemIndex].Specifications = make(map[string]string) // Выделяем память в мапу
		c.Visit(URL + itemVal.Link)
		itemIndex++
	}

	/*
		items.Data[itemIndex].Specifications = make(map[string]string) // Выделяем память в мапу
		c.Visit(URL + items.Data[0].Link)
	*/
}
