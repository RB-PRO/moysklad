package directelectric

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/gocolly/colly"
)

// Функция парсит все [ссылки] из массива и результат сохраняет в множество Xlsx
//
// [ссылки]: https://www.directelectric.ru/catalog/
func ParseItemsAndSaveAnotherXlsx(links []string) {
	for _, link := range links {
		var items DirectelEctricObjects
		fmt.Println("> Парсинг подкаталога", URL+link)
		items.ParseItem(link)
		items.ParseAllItem()
		link = strings.ReplaceAll(link, "catalog/", "")
		link = strings.ReplaceAll(link, "/", "")
		items.SaveXlsx(link)
	}
}

// Метод парсит все [ссылки] из массива и результат записывает в приемник
//
// [ссылки]: https://www.directelectric.ru/catalog/
func (items *DirectelEctricObjects) ParseItems(links []string) {
	for _, link := range links {
		fmt.Println("> Парсинг подкаталога", URL+link)
		items.ParseItem(link)
	}
}

// Метод парсит [страницу] по определённому по всем его возможным страницам
//
// [страницу]: https://www.directelectric.ru/catalog/rozetki-i-vyklyuchateli/filter/vendor_new-is-schneider%20electric/serial-is-atlasdesign/apply/?PAGEN_1=2&nal=y
func (items *DirectelEctricObjects) ParseItem(link string) {
	//fmt.Println("Parse", link)

	var pagesString string
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
		//fmt.Println(e.DOM.Find("div[class^=item__stock-title]").Text())
		if e.DOM.Find("div[class^=item__stock-title]").Text() == "В наличии" { // "В наличии" "Нет в наличии"
			item.Availability = true
		}

		items.Data = append(items.Data, item)
	})

	// Проверить можно ли дальше листать
	c.OnHTML("[class='pagination__next']", func(e *colly.HTMLElement) {
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
			if strings.Contains(hrefNext, "PAGEN_1") {
				pagesString = "PAGEN_1"
				if m["PAGEN_1"][0] == strconv.Itoa(schetchik) {

					next = false
				}
			}
			if strings.Contains(hrefNext, "PAGEN_2") {
				pagesString = "PAGEN_2"
				if m["PAGEN_2"][0] == strconv.Itoa(schetchik) {

					next = false
				}
			}

			// Если нет вообще ничего
			// https://www.directelectric.ru/catalog/elektricheskie-teplye-poly/nagrevatelnye-maty/filter/vendor_new-is-ekf/apply/
			if !strings.Contains(hrefNext, "PAGEN_1") && !strings.Contains(hrefNext, "PAGEN_2") {
				next = false
			}
		}
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Если вообще ничего не нашли
		if e.DOM.Find("[class='pagination__next']").Text() == "" {
			next = false
		}
	})
	/* // Одичночный парсинг
	fmt.Println(next, URL+link)
	linkPages, _ := MakeLinkWithPage(URL+link, schetchik)
	c.Visit(linkPages)
	*/

	fmt.Println("-> Парсинг всех страниц категории:")
	bar := pb.StartNew(1000)

	// Цикл по всем страницам
	for {
		//fmt.Println("--> Страница", schetchik)

		// Выход из цикла парсинга
		if next == false { //!next
			next = true
			break
		}

		// Делаем ссылку со страницей
		linkPages, _ := MakeLinkWithPage(URL+link, pagesString, schetchik)

		bar.Increment() // Прибавляем 1 к отображению

		// Парсим
		//fmt.Println(linkPages)

		c.Visit(linkPages)
		schetchik++

		//break

	}
	bar.Finish()

}
