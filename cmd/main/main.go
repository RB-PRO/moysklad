package main

import (
	"github.com/RB-PRO/moysklad/internal/app"
)

func main() {
	//app.RunAllCategory() // Запуск по всем категориям
	//app.RunOneLink() // Запуск по тестовой ссылке
	app.AddProduct() // Работа с сервисом moysklad
}
