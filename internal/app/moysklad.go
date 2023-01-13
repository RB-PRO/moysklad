package app

import (
	"fmt"

	"github.com/dotnow/moysklad"
	"github.com/dotnow/moysklad/params"
)

func AddProduct() {

	username, _ := dataFile("username")
	password, _ := dataFile("password")

	ms := moysklad.NewClientWithBasicAuth(username, password)

	// Параметры запроса (фильтрация, сортировка, лимит, expand...)
	p := params.NewParams().Limit(20).OrderAsc("name")

	// Получаем структуру ответа 'ResultResponse[T]', с полем Result, содержащим ответ на запрос
	resultResponse := ms.Product().WithParams(p).Get()

	// Выведем кол-во секунд, за которое запрос был выполнен
	fmt.Printf("Запрос выполнен за: %f сек.\n", resultResponse.TimeElapsed)

	// Ответ может содержать ошибки (например, ошибки десериализации)
	// или ошибки ответа от сервиса moysklad
	if resultResponse.HasErrors() {
		fmt.Println(resultResponse.Error)

		for _, clientError := range resultResponse.ClientErrors {
			fmt.Println(clientError.Error, clientError.Code, clientError.MoreInfo)
		}
	}

	for _, product := range resultResponse.Result.Rows {
		fmt.Println(product.Name)
	}
}
