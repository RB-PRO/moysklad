package directelectric

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

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

// Функция парсит [страницу] по определённому page
//
// [страницу]: https://www.directelectric.ru/catalog/rozetki-i-vyklyuchateli/filter/vendor_new-is-schneider%20electric/serial-is-atlasdesign/apply/?PAGEN_1=2&nal=y
func (items *DirectelEctricObjects) ParseItems(links []string) {
	//fmt.Println("Parse", link)

	var next bool = true
	var schetchik int = 1

	c := colly.NewCollector()

	// Карточки товара
	c.OnHTML("div[class=item]", func(e *colly.HTMLElement) {

		// Создаём экземпляр товара
		var item Product

		// Краткое название
		item.NameFew = e.DOM.Find("a[class^=item__title]").Text()
		item.NameFew = strings.TrimSpace(item.NameFew)

		// Ссылка на товар
		item.Link, _ = e.DOM.Find("a[class^=item__title]").Attr("href")

		// Наличие
		fmt.Println(e.DOM.Find("a[class^=item__stock-title]").Text())
		if e.DOM.Find("a[class^=item__stock-title]").Text() == "В наличии" { // "В наличии" "Нет в наличии"
			item.Availability = true
		}

		items.Data = append(items.Data, item)
	})

	// Проверить можно ли дальше листать
	c.OnHTML("[class=pagination__next]", func(e *colly.HTMLElement) {
		//fmt.Println(e.DOM.Text())
		hrefNext, hrefNextIsExit := e.DOM.Attr("href")
		//fmt.Println(">>", hrefNext, hrefNextIsExit)
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

	///*
	fmt.Println(next, URL+links[0])
	linkPages, _ := MakeLinkWithPage(URL+links[0], schetchik)
	c.Visit(linkPages)
	//*/

	/*
		for _, link := range links {
			fmt.Println("> Парсинг подкаталога", URL+link)
			for {
				fmt.Println("--> Страница", schetchik)

				// Выход из цикла парсинга
				if !next {
					next = true
					break
				}

				// Делаем ссылку со страницей
				linkPages, _ := MakeLinkWithPage(URL+link, schetchik)

				// Парсим
				c.Visit(linkPages)
				schetchik++
			}
		}
	*/
}

// Говно, надо отказаься от идеи использования глобальных переменных
var itemIndexGlobal int = 0

// Функция парсит каждую [страницу] товара
//
// [страницу]: https://www.directelectric.ru/catalog/product/180150/
func (items *DirectelEctricObjects) ParseAllItem() {
	c := colly.NewCollector()

	// Характеристики
	c.OnHTML("div[class=characteristics__list] div[class^=characteristics__item]", func(e *colly.HTMLElement) {
		key := e.DOM.Find("div[class=characteristics__name] span").Text()
		value := e.DOM.Find("div[class=characteristics__value]").Text()
		switch key {
		case "":
			break
		case "Артикул":
			items.Data[itemIndexGlobal].Article = value // Заполняем артикул
		default:
			items.Data[itemIndexGlobal].Specifications[key] = value // Заполняем мапу характеристиками товара
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
			items.Data[itemIndexGlobal].Price = n
		}
	})

	// Фото imageLink
	c.OnHTML("div[class=product__slider-item] img", func(e *colly.HTMLElement) {
		items.Data[itemIndexGlobal].imageLink, _ = e.DOM.Attr("src")
		items.Data[itemIndexGlobal].imageLink = URL + items.Data[itemIndexGlobal].imageLink
	})

	for _, itemVal := range items.Data {
		fmt.Println("Parse Item: ", itemIndexGlobal+1, "/", len(items.Data))
		items.Data[itemIndexGlobal].Specifications = make(map[string]string) // Выделяем память в мапу
		c.Visit(URL + itemVal.Link)
		itemIndexGlobal++

	}

	/*
		items.Data[itemIndexGlobal].Specifications = make(map[string]string) // Выделяем память в мапу
		c.Visit(URL + items.Data[0].Link)
	*/
}
