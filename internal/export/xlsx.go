package export

import (
	"fmt"
	"path/filepath"

	"github.com/xuri/excelize/v2"
	"tool-bid-app/internal/model"
)

// ExportXLSX генерирует Excel (.xlsx) файл заявки в точном соответствии с заводским шаблоном
func ExportXLSX(doc *model.Document, outputPath string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Заявка"
	defaultSheet := f.GetSheetName(0)
	f.SetSheetName(defaultSheet, sheetName)

	if err := setupXLSXMetrics(f, sheetName); err != nil {
		return err
	}

	styles, err := initXLSXStyles(f)
	if err != nil {
		return err
	}

	if err := fillXLSXContent(f, sheetName, doc, styles); err != nil {
		return err
	}

	// Настройки страницы: Fit to page (1 страница по ширине)
	fitToPage := true
	if err := f.SetSheetProps(sheetName, &excelize.SheetPropsOptions{
		FitToPage: &fitToPage,
	}); err != nil {
		return fmt.Errorf("ошибка настройки масштабирования листа: %w", err)
	}

	// Настройки страницы и печати: ориентация Portrait, A4 (Size: 9), FitToWidth = 1, FitToHeight = 0
	paperSize := 9 // A4
	orientation := "portrait"
	fitWidth := 1
	fitHeight := 0
	if err := f.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Orientation: &orientation,
		Size:        &paperSize,
		FitToWidth:  &fitWidth,
		FitToHeight: &fitHeight,
	}); err != nil {
		return fmt.Errorf("ошибка настройки страницы: %w", err)
	}

	// Сквозной повтор заголовка таблицы (строка 19) на каждой странице
	if err := f.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Titles",
		RefersTo: fmt.Sprintf("'%s'!$19:$19", sheetName),
		Scope:    sheetName,
	}); err != nil {
		return fmt.Errorf("ошибка установки сквозного заголовка: %w", err)
	}

	cleanPath := filepath.Clean(outputPath)
	if err := f.SaveAs(cleanPath); err != nil {
		return fmt.Errorf("ошибка сохранения файла XLSX '%s': %w", cleanPath, err)
	}

	return nil
}
