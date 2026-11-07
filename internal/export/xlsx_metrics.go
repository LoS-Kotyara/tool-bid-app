package export

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// setupXLSXMetrics настраивает ширины колонок и базовые высоты строк по ТЗ
func setupXLSXMetrics(f *excelize.File, sheetName string) error {
	colWidths := map[string]float64{
		"A": 4.14,
		"B": 29.57,
		"C": 20.85,
		"D": 21.57,
		"E": 8.43,
		"F": 9.29,
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return fmt.Errorf("ошибка установки ширины колонки %s: %w", col, err)
		}
	}

	rowHeights := map[int]float64{
		1:  15.0,
		2:  5.25,
		3:  15.0,
		4:  4.5,
		5:  27.0,
		6:  35.45,
		7:  33.6,
		8:  4.5,
		9:  16.5,
		10: 4.5,
		11: 32.25,
		12: 54.75,
		13: 15.0,
		14: 33.0,
		15: 54.0,
		16: 3.0,
		17: 16.5,
		18: 4.5,
		19: 36.0,
	}
	for r, h := range rowHeights {
		if err := f.SetRowHeight(sheetName, r, h); err != nil {
			return fmt.Errorf("ошибка установки высоты строки %d: %w", r, err)
		}
	}
	return nil
}
