package directelectric

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/gocolly/colly"
)

// Функция парсит каждую [страницу] товара
//
// [страницу]: https://www.directelectric.ru/catalog/product/180150/
func (prod *Product) SingleCart() {
	c := colly.NewCollector()

	// Характеристики
	c.OnHTML("div[class=characteristics__list] div[class^=characteristics__item]", func(e *colly.HTMLElement) {
		key := e.DOM.Find("div[class=characteristics__name] span").Text()
		value := e.DOM.Find("div[class=characteristics__value]").Text()
		switch key {
		case "":
			break
		case "Артикул":
			prod.Article = value // Заполняем артикул
		case "Код":
			prod.Code = value // Заполняем Код товара
		case "Вес":
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				prod.Weight = floatValue // Заполняем Вес
			}
		case "Объём":
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				prod.Volume = floatValue // Заполняем Объём
			}
		default:
			prod.Specifications[key] = value // Заполняем мапу характеристиками товара
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
			prod.Price = n
		}
	})

	// Размерность
	c.OnHTML("div[class=buy__price] div span:last-of-type", func(e *colly.HTMLElement) {
		prod.Dimension = strings.ReplaceAll(e.DOM.Text(), "/", "")
	})

	// Фото imageLink
	c.OnHTML("div[class=product__slider-item] img", func(e *colly.HTMLElement) {
		prod.ImageLink, _ = e.DOM.Attr("src")
		prod.ImageLink = URL + prod.ImageLink
	})

	// Описание товара
	// product-tabs__tab product-tabs__tab_active
	c.OnHTML("div[class*=product-tabs__tab_active]", func(e *colly.HTMLElement) {
		//fmt.Println(e.DOM.Text())
		prod.NameFull = e.DOM.Find("div[class=h3]").Text() // Полное название товара
		e.DOM.Find("div").Remove()
		prod.Description = strings.TrimSpace(e.DOM.Text())
	})

	c.Visit(URL + prod.Link)
	/*
		prod.Specifications = make(map[string]string) // Выделяем память в мапу
		c.Visit(URL + items.Data[0].Link)
	*/

}

func (items *DirectelEctricObjects) ParseAllItem() {
	// Вывод прогресса
	fmt.Println("--> Парсинг каждой карточки товара")
	bar := pb.StartNew(len(items.Data))
	for indexItem := range items.Data {
		bar.Increment() // Прибавляем 1 к отображению
		//fmt.Println("Parse Item: ", itemIndexGlobal+1, "/", len(items.Data))
		items.Data[indexItem].Specifications = make(map[string]string) // Выделяем память в мапу
		items.Data[indexItem].SingleCart()
	}
	bar.Finish()
}
