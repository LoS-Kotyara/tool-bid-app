package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"tool-bid-app/internal/export"
	"tool-bid-app/internal/model"
	"tool-bid-app/internal/parser"
)

// SaveMarkdownFile сериализует заявку в Markdown и сохраняет на диск.
// Если path пустой — вызывает системный диалог "Сохранить как...".
func (a *App) SaveMarkdownFile(path string, doc *model.Document) (string, error) {
	if doc == nil {
		return "", fmt.Errorf("документ не передан")
	}

	targetPath := path
	if strings.TrimSpace(targetPath) == "" {
		if a.ctx == nil {
			return "", fmt.Errorf("контекст приложения не инициализирован")
		}
		defaultFilename := "заявка.md"
		if doc.Meta.BidNum != "" {
			defaultFilename = fmt.Sprintf("заявка №%s от %s.md", doc.Meta.BidNum, doc.Meta.BidDate)
		}

		selectedPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			Title:           "Сохранить заявку как Markdown",
			DefaultFilename: defaultFilename,
			Filters: []runtime.FileFilter{
				{
					DisplayName: "Файлы Markdown (*.md)",
					Pattern:     "*.md",
				},
			},
		})
		if err != nil {
			return "", fmt.Errorf("ошибка диалога сохранения: %w", err)
		}
		if selectedPath == "" {
			return "", nil // Отмена сохранения пользователем
		}
		targetPath = selectedPath
	}

	content, err := parser.SerializeToMarkdown(doc)
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации в Markdown: %w", err)
	}

	cleanPath := filepath.Clean(targetPath)
	if err := os.WriteFile(cleanPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("не удалось записать файл '%s': %w", cleanPath, err)
	}

	return cleanPath, nil
}

// ReloadMarkdownFile перезагружает файл с диска по указанному пути
func (a *App) ReloadMarkdownFile(path string) (*model.Document, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("путь к файлу не указан")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл '%s': %w", path, err)
	}

	doc, err := parser.ParseMarkdown(string(data))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла: %w", err)
	}

	doc.FilePath = path
	doc.IsModified = false
	return doc, nil
}

// ExportXLSXDialog вызывает диалог сохранения *.xlsx и экспортирует файл
func (a *App) ExportXLSXDialog(doc *model.Document) error {
	if a.ctx == nil {
		return fmt.Errorf("контекст приложения не инициализирован")
	}
	if doc == nil {
		return fmt.Errorf("документ не передан")
	}

	defaultFilename := "заявка.xlsx"
	if doc.Meta.BidNum != "" {
		defaultFilename = fmt.Sprintf("заявка №%s от %s.xlsx", doc.Meta.BidNum, doc.Meta.BidDate)
	}

	selectedPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Экспорт заявки в Excel (XLSX)",
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Книга Excel (*.xlsx)",
				Pattern:     "*.xlsx",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка диалога сохранения XLSX: %w", err)
	}
	if selectedPath == "" {
		return nil
	}

	if err := export.ExportXLSX(doc, selectedPath); err != nil {
		return fmt.Errorf("ошибка экспорта в XLSX: %w", err)
	}

	return nil
}

// ExportODSDialog вызывает диалог сохранения *.ods и экспортирует файл
func (a *App) ExportODSDialog(doc *model.Document) error {
	if a.ctx == nil {
		return fmt.Errorf("контекст приложения не инициализирован")
	}
	if doc == nil {
		return fmt.Errorf("документ не передан")
	}

	defaultFilename := "заявка.ods"
	if doc.Meta.BidNum != "" {
		defaultFilename = fmt.Sprintf("заявка №%s от %s.ods", doc.Meta.BidNum, doc.Meta.BidDate)
	}

	selectedPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Экспорт заявки в OpenDocument Spreadsheet (ODS)",
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Таблица OpenDocument (*.ods)",
				Pattern:     "*.ods",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка диалога сохранения ODS: %w", err)
	}
	if selectedPath == "" {
		return nil
	}

	if err := export.ExportODS(doc, selectedPath); err != nil {
		return fmt.Errorf("ошибка экспорта в ODS: %w", err)
	}

	return nil
}

// GetDocumentPreviewHTML возвращает HTML-код печатного листа A4 заявки
func (a *App) GetDocumentPreviewHTML(doc *model.Document) (string, error) {
	return export.GenerateDocumentPreviewHTML(doc)
}
