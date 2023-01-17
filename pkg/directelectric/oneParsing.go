package directelectric

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/gocolly/colly"
)

// Говно, надо отказаься от идеи использования глобальных переменных
var itemIndexGlobal int = 0

// Функция парсит каждую [страницу] товара
//
// [страницу]: https://www.directelectric.ru/catalog/product/180150/
func (items *DirectelEctricObjects) ParseAllItem() {
	c := colly.NewCollector()

	itemIndexGlobal = 0

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

	// Размерность
	c.OnHTML("div[class=buy__price] div span:last-of-type", func(e *colly.HTMLElement) {
		items.Data[itemIndexGlobal].Dimension = strings.ReplaceAll(e.DOM.Text(), "/", "")
	})

	// Фото imageLink
	c.OnHTML("div[class=product__slider-item] img", func(e *colly.HTMLElement) {
		items.Data[itemIndexGlobal].imageLink, _ = e.DOM.Attr("src")
		items.Data[itemIndexGlobal].imageLink = URL + items.Data[itemIndexGlobal].imageLink
	})

	// Описание товара
	// product-tabs__tab product-tabs__tab_active
	c.OnHTML("div[class*=product-tabs__tab_active]", func(e *colly.HTMLElement) {
		//fmt.Println(e.DOM.Text())
		items.Data[itemIndexGlobal].NameFull = e.DOM.Find("div[class=h3]").Text() // Полное название товара
		e.DOM.Find("div").Remove()
		items.Data[itemIndexGlobal].Description = strings.TrimSpace(e.DOM.Text())
	})

	// Вывод прогресса
	bar := pb.StartNew(len(items.Data))
	fmt.Println("Парсинг каждой карточки товара")

	for _, itemVal := range items.Data {
		bar.Increment() // Прибавляем 1 к отображению
		//fmt.Println("Parse Item: ", itemIndexGlobal+1, "/", len(items.Data))
		items.Data[itemIndexGlobal].Specifications = make(map[string]string) // Выделяем память в мапу
		c.Visit(URL + itemVal.Link)
		itemIndexGlobal++

	}

	bar.Finish()
	/*
		items.Data[itemIndexGlobal].Specifications = make(map[string]string) // Выделяем память в мапу
		c.Visit(URL + items.Data[0].Link)
	*/
}
