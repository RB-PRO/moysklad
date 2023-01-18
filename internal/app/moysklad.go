package app

import (
	"fmt"

	"github.com/dotnow/moysklad"
	"github.com/dotnow/moysklad/params"
)

func AddProduct() {

	username, _ := dataFile("username") // Получить логин из файла
	password, _ := dataFile("password") // Получить пароль из файла

	ms := moysklad.NewClientWithBasicAuth(username, password)
	ms.PrettyPrintJson(true)       // Включить вывод форматированного JSON
	ms.DisableWebhookContent(true) // Отключить уведомление вебхуков в контексте данного запроса
	ms.SetAttempts(10)
	ms.SetTimeout(100)

	// *********************************************************************************************

	// Получить все дополнительные поля
	MetadataAttr := ms.Entity().Product().MetadataAttributes()
	fmt.Println(MetadataAttr.Result.Rows[0].Name)
	fmt.Println(MetadataAttr.Result.Rows[1].Name)
	fmt.Println(MetadataAttr.Result.Rows[2].Name)

	// *********************************************************************************************

	//p := params.NewParams()

	p := params.NewParams()

	productsWithParams := ms.Entity().Product().Create(p)

	fmt.Println(string(productsWithParams.Body))
}
