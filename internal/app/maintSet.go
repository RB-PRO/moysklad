package app

import (
	"fmt"
	"log"

	"github.com/RB-PRO/moysklad/pkg/directelectric"
	"github.com/cheggaaa/pb"
	"github.com/dotnow/moysklad"
)

// Спасить ВСЕ [товары] из ДиректЭлектрика и загрузить из в МойСклад в [список товаров].
//
// [товары]: https://www.directelectric.ru/catalog/product/180142/
// [список товаров]: https://github.com/dotnow/moysklad/issues/3
func ParseAllObjectAndLoadToMoySklad() {
	// Авторизация на сервисе МойСклад
	username, _ := dataFile("username") // Получить логин из файла
	password, _ := dataFile("password") // Получить пароль из файла

	ms := moysklad.NewClientWithBasicAuth(username, password)
	ms.PrettyPrintJson(true)       // Включить вывод форматированного JSON
	ms.DisableWebhookContent(true) // Отключить уведомление вебхуков в контексте данного запроса
	ms.SetAttempts(10)
	ms.SetTimeout(100)

	// *********************************************************************************************
	// Получить все категории на директ электрике
	links := directelectric.ParseCatalogs()

	fmt.Println("Найдено всего", len(links), "категорий")
	//links = links[:1]
	fmt.Println("links", links)

	// Определение структуры с данными
	var items directelectric.DirectelEctricObjects

	// Пропарсить в подкатегории
	items.ParseItems(links)

	// Получаем список дополнительных полей
	metaAttributes, _ := MetaAttr(ms)
	fmt.Println("Всего дополнительных полей:", len(metaAttributes))

	// Единицы измерения
	metaUom, errorUom := MetaUom(ms)
	if errorUom != nil {
		log.Println(errorUom)
	}

	// Пропарсить карточки товаров и добавить сразу на МойСклад // items.ParseAllItem()
	fmt.Println("-> Парсинг каждой карточки товара")
	bar := pb.StartNew(len(items.Data))
	for indexItem := range items.Data { // Индекс по всем карточкам товаров
		bar.Increment()                                                                                // Прибавляем 1 к отображению
		items.Data[indexItem].Specifications = make(map[string]string)                                 // Выделяем память в мапу
		items.Data[indexItem].SingleCart()                                                             // Пропарсить карточку товара
		items.Data[indexItem] = AddProductMoySklad(items.Data[indexItem], ms, metaAttributes, metaUom) // Добавить товар в корзину
	}
	bar.Finish()

	items.SaveXlsx("directelectric")
}
