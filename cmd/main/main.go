package main

import "github.com/RB-PRO/moysklad/internal/app"

func main() {
	//app.RunAllCategory() // Запуск по всем категориям
	//app.RunOneLink() // Запуск по тестовой ссылке
	//app.AddProduct() // Работа с сервисом moysklad
	//app.ParseItemsAndSave() // Сохранить весь сайт в Xlsx
	//app.Search_add()                      // Добавление товара
	app.ParseAllObjectAndLoadToMoySklad() // Спарсить все товары и добавить в МойСклад
}
