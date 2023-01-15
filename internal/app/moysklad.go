package app

import (
	"fmt"

	"github.com/dotnow/moysklad"
	"github.com/dotnow/moysklad/client"
	"github.com/dotnow/moysklad/params"
)

func AddProduct() {

	username, _ := dataFile("username")
	password, _ := dataFile("password")

	ms := moysklad.NewClientWithBasicAuth(username, password)
	ms.PrettyPrintJson(true) // Включить вывод форматированного JSON

	ms.DisableWebhookContent(true) // Отключить уведомление вебхуков в контексте данного запроса
	ms.SetAttempts(10)
	ms.SetTimeout(100)

	// ***

	productClient := client.ProductClient(ms)

	// Задаём параметр
	p := params.NewParams()
	p.Offset(0)
	p.Limit(50)
	p.OrderAsc("name")
	productClient.WithParams(p)

	// Выполняем запрос
	productClientGet := productClient.Get()

	fmt.Println(len(productClientGet.Result.Rows))

	for _, valls := range productClientGet.Result.Rows {
		fmt.Println(valls.Article, valls.Name)
	}

	fmt.Println(productClientGet.Result.Meta)

	// ***

	metadataClient := client.MetadataClient(ms)

	metadataClientGet := metadataClient.Get()

	fmt.Println(metadataClientGet.Result.CustomerOrder.States[0].Meta.Href)
	fmt.Println()

	fmt.Printf("%#+v\n", metadataClientGet.Result.CustomerOrder.States[0])
	fmt.Println()
	fmt.Println(metadataClientGet.Result.CustomerOrder.States[0])

}

// Структура, которая содержит ссылки на доп. поля и их названия
type metas struct {
	met []met
}
type met struct {
	id   string // ID дополнительного поля
	name string // Название поля
}

// Метод, который заполняет массив дополнительных полей
func (met *metas) metaData() {

}
