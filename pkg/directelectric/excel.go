package directelectric

import (
	"strconv"

	"github.com/xuri/excelize/v2"
)

// Создать книгу
func makeWorkBook() (*excelize.File, error) {
	// Создать книгу Excel
	f := excelize.NewFile()
	// Create a new sheet.
	_, err := f.NewSheet("main")
	if err != nil {
		return f, err
	}
	f.DeleteSheet("Sheet1")
	return f, nil
}

// Сохранить и закрыть файл
func closeXlsx(f *excelize.File, filename string) error {
	if err := f.SaveAs(filename + ".xlsx"); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

/*
func WriteOneLine(f *excelize.File, ssheet string, row int, SearchBasicRes SearchBasicResponse, SearchBasicIndex int, GetPartsRemainsByCodeRes GetPartsRemainsByCodeResponse, GetPartsRemainsByCodeIndex int) {
	// SearchBasic
	writeHeadOne(f, ssheet, 1, row, SearchBasicRes.Data.Items[SearchBasicIndex].Code, "")
}
*/

// Сохранить данные в Xlsx
func (items DirectelEctricObjects) SaveXlsx(filename string) error {
	// Создать книгу
	book, makeBookError := makeWorkBook()
	if makeBookError != nil {
		return makeBookError
	}

	wotkSheet := "main"
	setHead(book, wotkSheet, 1, "Ссылка на товар")            // Link
	setHead(book, wotkSheet, 2, "Фото")                       // imageLink
	setHead(book, wotkSheet, 3, "Полное Название товара")     // NameFull
	setHead(book, wotkSheet, 4, "Краткое Название товара")    // NameFew
	setHead(book, wotkSheet, 5, "Описание товара")            // Description
	setHead(book, wotkSheet, 6, "Артикул")                    // Article
	setHead(book, wotkSheet, 7, "Цена")                       // Price
	setHead(book, wotkSheet, 8, "Наличие")                    // Availability
	setHead(book, wotkSheet, 9, "Размерность(шт, кг, тонна)") // Dimension
	startIndexCollumn := 10

	// Создаём мапу, которая будет содержать значения номеров колонок
	colName := make(map[string]int)
	for indexItem, valItem := range items.Data {
		setCell(book, wotkSheet, indexItem, 1, URL+valItem.Link)     // Ссылка на товар
		setCell(book, wotkSheet, indexItem, 2, valItem.ImageLink)    // Фото
		setCell(book, wotkSheet, indexItem, 3, valItem.NameFull)     // Полное Название товар
		setCell(book, wotkSheet, indexItem, 4, valItem.NameFew)      // Краткое Название товара
		setCell(book, wotkSheet, indexItem, 5, valItem.Description)  // Описание товара
		setCell(book, wotkSheet, indexItem, 6, valItem.Article)      // Артикул
		setCell(book, wotkSheet, indexItem, 7, valItem.Price)        // Цена
		setCell(book, wotkSheet, indexItem, 8, valItem.Availability) // Наличие
		setCell(book, wotkSheet, indexItem, 9, valItem.Dimension)    // Размерность(шт, кг, тонна)

		for key, val := range valItem.Specifications {
			if _, ok := colName[key]; ok { // Если такое значение существует(т.е. существует колонка)
				//do something here
				setCell(book, wotkSheet, indexItem, colName[key], val)
			} else {
				colName[key] = startIndexCollumn
				setHead(book, wotkSheet, colName[key], key)
				setCell(book, wotkSheet, indexItem, colName[key], val)
				startIndexCollumn++
			}

		}
	}

	// Закрыть книгу
	closeBookError := closeXlsx(book, filename)
	if closeBookError != nil {
		return closeBookError
	}
	return nil
}

// Вписать значение в ячейку
func setCell(file *excelize.File, wotkSheet string, y, x int, value interface{}) {
	collumnStr, _ := excelize.ColumnNumberToName(x)
	file.SetCellValue(wotkSheet, collumnStr+strconv.Itoa(y+2), value)
}

// Вписать значение в название колонки
func setHead(file *excelize.File, wotkSheet string, x int, value interface{}) {
	collumnStr, _ := excelize.ColumnNumberToName(x)
	file.SetCellValue(wotkSheet, collumnStr+"1", value)
}
