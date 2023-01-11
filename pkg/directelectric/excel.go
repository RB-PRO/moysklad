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
func closeXlsx(f *excelize.File) error {
	if err := f.SaveAs("directelectric.xlsx"); err != nil {
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
func (items DirectelEctricObjects) SaveXlsx() error {
	// Создать книгу
	book, makeBookError := makeWorkBook()
	if makeBookError != nil {
		return makeBookError
	}

	wotkSheet := "main"
	setHead(book, wotkSheet, 1, "Ссылка на товар")
	setHead(book, wotkSheet, 2, "Фото")
	setHead(book, wotkSheet, 3, "Полное Название товара")
	setHead(book, wotkSheet, 4, "Краткое Название товара")
	setHead(book, wotkSheet, 5, "Артикул")
	setHead(book, wotkSheet, 6, "Цена")
	startIndexCollumn := 7

	// Создаём мапу, которая будет содержать значения номеров колонок
	colName := make(map[string]int)
	for indexItem, valItem := range items.Data {
		setCell(book, wotkSheet, indexItem, 1, URL+valItem.Link)
		setCell(book, wotkSheet, indexItem, 2, valItem.imageLink)
		setCell(book, wotkSheet, indexItem, 3, valItem.NameFull)
		setCell(book, wotkSheet, indexItem, 4, valItem.NameFew)
		setCell(book, wotkSheet, indexItem, 5, valItem.Article)
		setCell(book, wotkSheet, indexItem, 6, valItem.Price)

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
	closeBookError := closeXlsx(book)
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
