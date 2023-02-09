package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/RB-PRO/moysklad/pkg/directelectric"
	"github.com/cheggaaa/pb"
	"github.com/dotnow/moysklad"
)

func ParseLinks(links []string) {
	// Авторизация на сервисе МойСклад
	username, _ := dataFile("username") // Получить логин из файла
	password, _ := dataFile("password") // Получить пароль из файла
	//username := User() // Получить логин
	//password := Pass() // Получить пароль

	ms := moysklad.NewClientWithBasicAuth(username, password)
	ms.PrettyPrintJson(true)       // Включить вывод форматированного JSON
	ms.DisableWebhookContent(true) // Отключить уведомление вебхуков в контексте данного запроса
	ms.PricePrecision(true)

	ms.SetAttempts(0)
	ms.SetTimeout(100)

	// *********************************************************************************************

	fmt.Println("Найдено всего", len(links), "категорий")
	//links = links[:1]
	fmt.Println("Ссылки на товары:", links)

	// Получаем список дополнительных полей
	metaAttributes, metaAttributesError := MetaAttr(ms)
	if metaAttributesError != nil {
		log.Fatalln(metaAttributesError)
	}
	fmt.Println("Всего дополнительных полей:", len(metaAttributes))

	// Единицы измерения
	metaUom, errorUom := MetaUom(ms)
	if errorUom != nil {
		log.Fatalln(errorUom)
	}
	// цена закупочная
	valCurrency, currencyError := EntityRuble(ms)
	if currencyError != nil {
		log.Fatalln(currencyError)
	}

	// *********************** START ***********************
	for _, link := range links { // Пропарсить в подкатегории // цикл по всем link, вместо items.ParseItems(links)
		// Определение структуры с данными
		var items directelectric.DirectelEctricObjects
		fmt.Println("> Парсинг подкаталога", directelectric.URL+link)
		items.ParseItem(link)

		// Пропарсить карточки товаров и добавить сразу на МойСклад // items.ParseAllItem()
		fmt.Println("-> Парсинг каждой карточки товара")
		bar := pb.StartNew(len(items.Data))
		for indexItem := range items.Data { // Индекс по всем карточкам товаров
			bar.Increment()                                                                                             // Прибавляем 1 к отображению
			items.Data[indexItem].Specifications = make(map[string]string)                                              // Выделяем память в мапу
			items.Data[indexItem].SingleCart()                                                                          // Пропарсить карточку товара
			items.Data[indexItem] = AddProductMoySklad(items.Data[indexItem], ms, metaAttributes, metaUom, valCurrency) // Добавить товар в корзину
		}
		bar.Finish()

		items.SaveXlsx(strings.ReplaceAll(link, "/", "-")) // Сохранить в XLSX
	}
}
