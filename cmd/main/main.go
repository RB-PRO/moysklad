// # Парсер [DirectElectric] с загрузкой на сервий [MoySklad]
//
// # Данные проект используется специально для единоразовой загрузки товаров, однако устроен так, что его с лёгкостью можно приспособить для обновления актуальность товаров в CRM МойСклад
//
// В качестве бибилиотеки используется модуль [moysklad] от open source разработчика [dotnow]. Выражаю безмерное спасибо за модуль и я рад, что принял участие в коммите [v0.0.2-beta3]
//
// Данная программа работает совместно с версией github.com/dotnow/moyskladv0.0.2-beta3
//
// Проект разделён на 3 слоя:
// - cmd - запускает функции из internal;
// - internal - бизнес-логика, само приложение, правильно использует методы парсинга и загрузки данных, оперирует внутренними структурами, предоставляемыми внешними пакетами;
// - pkg - код, который может быть переиспользован в остальных проектах. В данном случае пакет directelectric сипользуется для парсинга соответствующего сайта.
//
// [moysklad]: https://github.com/dotnow/moysklad
// [dotnow]: https://github.com/dotnow/moysklad/commits?author=dotnow
// [v0.0.2-beta3]: https://github.com/dotnow/moysklad/commit/3dff569944480d9b8c639ba6a152f6aaf6e995d7
//
// [DirectElectric]: https://www.directelectric.ru
// [MoySklad]: https://www.moysklad.ru/
package main

import "github.com/RB-PRO/moysklad/internal/app"

func main() {
	//app.RunAllCategory() // Запуск по всем категориям
	//app.RunOneLink() // Запуск по тестовой ссылке
	//app.AddProduct() // Работа с сервисом moysklad
	//app.ParseItemsAndSave() // Сохранить весь сайт в Xlsx
	//app.Search_add()                      // Добавление товара
	//app.ParseAllObjectAndLoadToMoySklad() // Спарсить все товары и добавить в МойСклад
	//app.ParseLinkAndLoadToMoySklad() // спарсить по ссылке
	app.Schuse()
}
