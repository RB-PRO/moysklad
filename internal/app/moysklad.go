package app

import "github.com/dotnow/moysklad"

func AddProduct() {

	username, _ := dataFile("username")
	password, _ := dataFile("password")

	ms := moysklad.NewClientWithBasicAuth(username, password)

	ms.PrettyPrintJson(true) //Включить вывод форматированного JSON

	ms.DisableWebhookContent(true) //Отключить уведомление вебхуков в контексте данного запроса

}
