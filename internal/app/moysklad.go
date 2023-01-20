package app

import (
	"fmt"
	"log"

	"github.com/dotnow/moysklad"
	"github.com/dotnow/moysklad/client"
	"github.com/dotnow/moysklad/entity"
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

	/*
		// Получить все дополнительные поля
		MetadataAttr := ms.Entity().Product().MetadataAttributes()
		fmt.Println("Всего дополнительных полей", len(MetadataAttr.Result.Rows))
		fmt.Println(MetadataAttr.Result.Rows[0].Name)
		fmt.Println(MetadataAttr.Result.Rows[1].Name)
		fmt.Println(MetadataAttr.Result.Rows[2].Name)
		fmt.Println(MetadataAttr.Result.Rows[2].AttributeEntityType)
	*/

	// Получаем все доп поля сущности 'Товар'
	attributes, response := ms.Entity().Product().MetadataAttributes()
	if response.HasErrors() {
		log.Fatalln(response.Errors)
	}
	fmt.Println("Всего доп. полей", len(attributes.Rows))
	fmt.Println(attributes.Rows[0].Name)
	fmt.Println(attributes.Rows[1].Name)
	fmt.Println(attributes.Rows[2].Name)
	fmt.Println()

	metaAttributes, _ := MetaAttr(ms)
	for key, val := range metaAttributes {
		fmt.Println(key, val)
	}

	// *********************************************************************************************

	//p := params.NewParams()

	//p := params.NewParams()

	//productsWithParams := ms.Entity().Product().Create(p)

	//CreateMetadataAttribute
	//asd := client.ProductClient(ms)

	//var zxc entity.Attribute

	var gf entity.Product

	gf.Code = "123"
	gf.Article = "123"
	gf.Name = "123"

	//zxc := asd.Create(gf)

	//zxc := ms.Entity().Product().Create(gf)

	//fmt.Println("zxc", zxc)
	//fmt.Printf("->>> %+#v\n\n", zxc)
	//fmt.Println("asd", asd)
	//fmt.Println("gf", gf)

}

// Получить мапу дополнительных полей
//
// map[string]string - map["Название дополнительного поля"] = "Ссылка на поле"
func MetaAttr(ms *client.JSONApiClient) (map[string]*entity.Meta, error) {
	attributes := make(map[string]*entity.Meta)                          // Выделяем память в структуру, которая хранит данные о дополнительных полях
	MetadataAttr, response := ms.Entity().Product().MetadataAttributes() // Выполнить запрос дополнительных полей
	if response.HasErrors() {                                            // Проверяем на наличие ошибки в запросе
		return nil, response.Errors.Merge()
	}
	for _, val := range MetadataAttr.Rows { // Цикл по результатам запроса
		attributes[val.Name] = val.Meta // Заполняем map
	}
	return attributes, nil
}
