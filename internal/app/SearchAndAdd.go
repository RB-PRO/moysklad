package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/RB-PRO/moysklad/pkg/directelectric"
	"github.com/dotnow/moysklad"
	"github.com/dotnow/moysklad/client"
	"github.com/dotnow/moysklad/entity"
)

// Тестовый парсинг [товара] и добавление его в [список товаров] на сервисе мой склад
//
// [товара]: https://www.directelectric.ru/catalog/product/180142/
// [список товаров]: https://github.com/dotnow/moysklad/issues/3
func Search_add() {
	// Авторизация на сервисе МойСклад
	username, _ := dataFile("username") // Получить логин из файла
	password, _ := dataFile("password") // Получить пароль из файла

	ms := moysklad.NewClientWithBasicAuth(username, password)
	ms.PrettyPrintJson(true)       // Включить вывод форматированного JSON
	ms.DisableWebhookContent(true) // Отключить уведомление вебхуков в контексте данного запроса
	ms.SetAttempts(10)
	ms.SetTimeout(100)

	// *********************************************************************************************

	// Создаём массив для тестового овара
	items := directelectric.DirectelEctricObjects{[]directelectric.Product{directelectric.Product{
		Link: "/catalog/product/180142/",
	}}}
	items.ParseAllItem() // Парсим товар и обновляем массив данных

	fmt.Printf("Название товара: %#v\n", items.Data[0].NameFull)

	// Получаем список дополнительных полей
	metaAttributes, _ := MetaAttr(ms)
	fmt.Println("Всего дополнительных полей:", len(metaAttributes))

	// ************************************************************************

	// Создаём товар
	product := new(entity.Product)

	// Обязательные поля
	product.Name = items.Data[0].NameFull           // Название товара
	product.Code = items.Data[0].Code               // Код товара
	product.Article = items.Data[0].Article         // Артикул
	product.Description = items.Data[0].Description // Описание товара
	product.Weight = items.Data[0].Weight           // Артикул
	product.Description = items.Data[0].Description // Описание товара

	// Получить мета данные рубля
	metaRuble, errorRuble := MetaRuble(ms)
	if errorRuble != nil {
		log.Fatal(errorRuble)
	}
	product.BuyPrice = &entity.BuyPrice{Value: items.Data[0].Price, Currency: &metaRuble} // Закупочная цена

	// Единицы измерения
	metaUom, errorUom := MetaUom(ms)
	if errorUom != nil {
		log.Fatal(errorUom)
	}
	product.Uom = &metaUom

	//product.Images.Meta.Href = items.Data[0].imageLink // Ссылка на картинку

	// Массив дополнительных полей, которые мы будет прикреплять к запросу на создание товара
	attrs := []entity.Attribute{}

	for keyAttr, valAttr := range items.Data[0].Specifications { // Цикл по дополнительным полям товара
		if _, ok := metaAttributes[keyAttr]; ok { // Если Дополнительное поле товара имеет аналог в МоёмСкладе
			// Добавляем соответствующее дополнительное поле
			attrs = append(attrs, entity.Attribute{
				Name:  keyAttr,                                // Название дополнительного поля
				Meta:  metaAttributes[keyAttr],                // Meta дополнительного поля
				Value: items.Data[0].Specifications[keyAttr]}) // Значение, которое необходимо записать
		} else { // В ином случае дополняем эту информацию в описание товара
			product.Description += "\n" + keyAttr + " - " + valAttr
		}
	}
	product.Attributes = &attrs // Записываем дополнительные поля в структуру запроса

	// Отправляем запрос за создание товара
	result, response := ms.Entity().Product().Create(product)
	fmt.Println("result", (result))
	fmt.Println("response", (response.Errors.Merge().Error()))
}

// Получить meta Рубля
func MetaRuble(ms *client.JSONApiClient) (entity.Currency, error) {
	currency, _ := ms.Entity().Currency().Get()
	//currency, _ := ms.Context().C.Entity().Currency().Get()
	//currency, _ := ms.Report().C.Entity().Currency().Get()
	for _, val := range currency.Rows {
		if val.FullName == "Российский рубль" { // Если это рубль
			return val, nil
		}
	}
	return entity.Currency{}, errors.New("не смог найти Рубль")
}

// Получить Единицы измерения
func MetaUom(ms *client.JSONApiClient) (entity.Uom, error) {
	currency, _ := ms.Entity().Uom().Get()
	for _, val := range currency.Rows {
		if val.Name == "шт" { // Если это рубль
			return val, nil
		}
	}
	return entity.Uom{}, errors.New("не смог найти шт")
}
