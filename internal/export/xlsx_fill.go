package export

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
	"tool-bid-app/internal/model"
)

// fillXLSXContent заполняет ячейки книги метаданными и спецификацией позиций
func fillXLSXContent(f *excelize.File, sheetName string, doc *model.Document, s *xlsxStyles) error {
	// Строка 1: Заголовок заявки
	f.SetCellValue(sheetName, "B1", "ЗАЯВКА №")
	f.SetCellStyle(sheetName, "B1", "B1", s.normal)

	headerText := fmt.Sprintf("%s от %s г.", strings.TrimSpace(doc.Meta.BidNum), strings.TrimSpace(doc.Meta.BidDate))
	f.SetCellValue(sheetName, "C1", headerText)
	f.MergeCell(sheetName, "C1", "F1")
	f.SetCellStyle(sheetName, "C1", "F1", s.headerUnderline)

	// Строка 3: Блок подразделения
	f.SetCellValue(sheetName, "D3", strings.TrimSpace(doc.Meta.Dept))
	f.MergeCell(sheetName, "D3", "F3")
	f.SetCellStyle(sheetName, "D3", "F3", s.deptBox)

	// Строки 6-7: Инициатор и Куратор
	f.SetCellValue(sheetName, "A6", "ИНИЦИАТОР ЗАКУПКИ")
	f.MergeCell(sheetName, "A6", "B6")
	f.SetCellStyle(sheetName, "A6", "B6", s.signatoryTitle)

	f.SetCellValue(sheetName, "C6", strings.TrimSpace(doc.Meta.InitRole))
	f.SetCellStyle(sheetName, "C6", "C6", s.signatoryRole)

	f.SetCellValue(sheetName, "E6", strings.TrimSpace(doc.Meta.InitFio))
	f.MergeCell(sheetName, "E6", "F6")
	f.SetCellStyle(sheetName, "E6", "F6", s.signatoryFIO)

	f.SetCellValue(sheetName, "A7", "КУРАТОР ЗАКУПКИ")
	f.MergeCell(sheetName, "A7", "B7")
	f.SetCellStyle(sheetName, "A7", "B7", s.signatoryTitle)

	f.SetCellValue(sheetName, "C7", strings.TrimSpace(doc.Meta.CurRole))
	f.SetCellStyle(sheetName, "C7", "C7", s.signatoryRole)

	f.SetCellValue(sheetName, "E7", strings.TrimSpace(doc.Meta.CurFio))
	f.MergeCell(sheetName, "E7", "F7")
	f.SetCellStyle(sheetName, "E7", "F7", s.signatoryFIO)

	// Строка 9: Раздел I
	f.SetCellValue(sheetName, "A9", "I. Общие сведения")
	f.MergeCell(sheetName, "A9", "F9")
	f.SetCellStyle(sheetName, "A9", "F9", s.sectionHeader)

	// Строка 11: Шапка Раздела I
	f.SetCellValue(sheetName, "A11", "№")
	f.SetCellStyle(sheetName, "A11", "A11", s.tableHeader)

	f.SetCellValue(sheetName, "B11", "Наименование")
	f.MergeCell(sheetName, "B11", "C11")
	f.SetCellStyle(sheetName, "B11", "C11", s.tableHeader)

	f.SetCellValue(sheetName, "D11", "Содержание")
	f.MergeCell(sheetName, "D11", "F11")
	f.SetCellStyle(sheetName, "D11", "F11", s.tableHeader)

	// Строки 12-15: Пункты Раздела I
	sec1 := []struct {
		num   string
		label string
		val   string
	}{
		{"1", "Служебная записка (Пункт)", doc.Meta.ShopName},
		{"2", "Срок поставки (не позднее)", doc.Meta.DeliveryDate},
		{"3", "Требования к качеству", doc.Meta.QualityReq},
		{"4", "Цель приобретения", doc.Meta.Purpose},
	}

	for idx, item := range sec1 {
		r := 12 + idx
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", r), item.num)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", r), fmt.Sprintf("A%d", r), s.tableCellCenter)

		f.SetCellValue(sheetName, fmt.Sprintf("B%d", r), item.label)
		f.MergeCell(sheetName, fmt.Sprintf("B%d", r), fmt.Sprintf("C%d", r))
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", r), fmt.Sprintf("C%d", r), s.tableLabelBold)

		f.SetCellValue(sheetName, fmt.Sprintf("D%d", r), strings.TrimSpace(item.val))
		f.MergeCell(sheetName, fmt.Sprintf("D%d", r), fmt.Sprintf("F%d", r))
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", r), fmt.Sprintf("F%d", r), s.tableCellLeft)
	}

	// Строка 17: Раздел II
	f.SetCellValue(sheetName, "A17", "II. Техническая часть документации")
	f.MergeCell(sheetName, "A17", "F17")
	f.SetCellStyle(sheetName, "A17", "F17", s.sectionHeader)

	// Строка 19: Шапка таблицы позиций
	f.SetCellValue(sheetName, "A19", "№\nп/п")
	f.SetCellStyle(sheetName, "A19", "A19", s.tableHeader)

	f.SetCellValue(sheetName, "B19", "Наименование товара")
	f.SetCellStyle(sheetName, "B19", "B19", s.tableHeader)

	f.SetCellValue(sheetName, "C19", "Характеристики товара")
	f.MergeCell(sheetName, "C19", "D19")
	f.SetCellStyle(sheetName, "C19", "D19", s.tableHeader)

	f.SetCellValue(sheetName, "E19", "Ед.\nизм.")
	f.SetCellStyle(sheetName, "E19", "E19", s.tableHeader)

	f.SetCellValue(sheetName, "F19", "Кол-во")
	f.SetCellStyle(sheetName, "F19", "F19", s.tableHeader)

	// Строки 20+: Позиции инструмента
	currentRow := 20
	for _, it := range doc.Items {
		f.SetRowHeight(sheetName, currentRow, 30.0)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), it.Num)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), s.tableCellCenter)

		f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), strings.TrimSpace(it.Name))
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", currentRow), fmt.Sprintf("B%d", currentRow), s.tableCellLeft)

		f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), strings.TrimSpace(it.Desc))
		f.MergeCell(sheetName, fmt.Sprintf("C%d", currentRow), fmt.Sprintf("D%d", currentRow))
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", currentRow), fmt.Sprintf("D%d", currentRow), s.tableCellLeft)

		f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), strings.TrimSpace(it.Unit))
		f.SetCellStyle(sheetName, fmt.Sprintf("E%d", currentRow), fmt.Sprintf("E%d", currentRow), s.tableCellCenter)

		f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), strings.TrimSpace(it.Qty))
		f.SetCellStyle(sheetName, fmt.Sprintf("F%d", currentRow), fmt.Sprintf("F%d", currentRow), s.tableCellCenter)

		currentRow++
	}

	return nil
}
