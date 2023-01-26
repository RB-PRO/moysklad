// Модуль для загрузки товаров в МойСклад
//
// Данный модуль использует библиотеку github.com/dotnow/moysklad@v0.0.2-beta3
package app

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/RB-PRO/moysklad/pkg/directelectric"
	"github.com/dotnow/moysklad/client"
	"github.com/dotnow/moysklad/entity"
)

// Метод загрузки товара на МойСклад
//
//	prod directelectric.Product - Сам товар для загрузки
//	ms *client.JSONApiClient - Авторизация клиента
//	metaAttributes map[string]*entity.Meta - Валюта
//	metaUom entity.Uom - Единицы измерения
func AddProductMoySklad(prod directelectric.Product, ms *client.JSONApiClient, metaAttributes map[string]*entity.Meta, metaUom entity.Uom, valCurrency entity.Currency) directelectric.Product {
	product := new(entity.Product) // Создаём товар

	// Обязательные поля
	product.Name = prod.NameFull           // Название товара
	product.Code = prod.Code               // Код товара
	product.Article = prod.Article         // Артикул
	product.Description = prod.Description // Описание товара
	product.Weight = prod.Weight           // Артикул
	product.Description = prod.Description // Описание товара

	// Цена
	product.BuyPrice = &entity.BuyPrice{Value: entity.FloatPrice(prod.Price), Currency: &valCurrency} // Закупочная цена prod.Price

	// Единицы измерения
	product.Uom = &metaUom

	// Загрузить создать объект картинки для загрузки
	if prod.ImageLink != directelectric.URL {
		errorDownload := downloadFile(prod.ImageLink, "main.jpeg") // Загрузить изображение
		if errorDownload == nil {                                  // Если картинка загружена
			img, imgError := ImageBase("main.jpeg") // Получить экземпляр изображения
			if imgError != nil {
				log.Println(imgError)
			}
			product.Images.Add(img) // Добавить изображение в товар
		}
	}

	// Массив дополнительных полей, которые мы будет прикреплять к запросу на создание товара
	attrs := []entity.Attribute{}
	for keyAttr, valAttr := range prod.Specifications { // Цикл по дополнительным полям товара
		if _, ok := metaAttributes[keyAttr]; ok { // Если Дополнительное поле товара имеет аналог в МоёмСкладе
			// Добавляем соответствующее дополнительное поле
			attrs = append(attrs, entity.Attribute{
				Name:  keyAttr,                       // Название дополнительного поля
				Meta:  metaAttributes[keyAttr],       // Meta дополнительного поля
				Value: prod.Specifications[keyAttr]}) // Значение, которое необходимо записать
		} else { // В ином случае дополняем эту информацию в описание товара
			product.Description += "\n" + "<strong>" + keyAttr + " - " + valAttr + "</strong>"
		}
	}
	product.Attributes = attrs // Записываем дополнительные поля в структуру запроса

	// Отправляем запрос за создание товара
	_, response := ms.Entity().Product().Create(*product)
	//fmt.Println("result", (result))
	//fmt.Println("response", (response.GetErrorsInline()))
	if response.GetErrorsInline() != nil {
		log.Print(response.GetErrorsInline().Error() + "; ")
		if strings.Contains(response.GetErrorsInline().Error(), "enexpected end of JSON input") {
			log.Fatalln("Не смог загрузить товар. Паникую!")
		}
	}
	return prod
}

// Получить рубль
func EntityRuble(ms *client.JSONApiClient) (entity.Currency, error) {
	currency, _ := ms.Entity().Currency().Get() // Получить информацию по валюте
	for _, val := range currency.Rows {         // Цикл по всем валютам
		if val.FullName == "Российский рубль" { // Если это рубль
			return val, nil
		}
	}
	return entity.Currency{}, errors.New("нет Рублей")
}

// Скачать картинку
func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// Получить мапу дополнительных полей
//
// map[string]string - map["Название дополнительного поля"] = "Ссылка на поле"
func MetaAttr(ms *client.JSONApiClient) (map[string]*entity.Meta, error) {
	attributes := make(map[string]*entity.Meta)                          // Выделяем память в структуру, которая хранит данные о дополнительных полях
	MetadataAttr, response := ms.Entity().Product().MetadataAttributes() // Выполнить запрос дополнительных полей
	if response.HasErrors() {                                            // Проверяем на наличие ошибки в запросе
		return nil, response.GetErrorsInline()
	}
	for _, val := range MetadataAttr.Rows { // Цикл по результатам запроса
		attributes[val.Name] = val.Meta // Заполняем map
	}
	return attributes, nil
}

// Создать экземпляр загрузки картинки
func ImageBase(filename string) (entity.Image, error) {
	var img entity.Image
	f, fOpenError := os.Open(filename) // Открыть картинку
	if fOpenError != nil {
		return entity.Image{}, fOpenError
	}
	reader := bufio.NewReader(f)                 // Сканируем содержимое
	content, fContentError := io.ReadAll(reader) // Читаем содержимое
	if fContentError != nil {
		return entity.Image{}, fContentError
	}
	encoded := base64.StdEncoding.EncodeToString(content) // кодируем в base64
	img.Filename = filename                               // Вставляем название картинки
	img.Content = encoded                                 // Вставляем содержимое base64
	return img, nil
}

// Получить meta Рубля
func MetaRuble(ms *client.JSONApiClient) (entity.Currency, error) {
	//currency := ms.Entity().Currency()
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
