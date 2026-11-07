package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"tool-bid-app/internal/model"
	"tool-bid-app/internal/parser"
)

// App представляет собой сервис приложения Wails
type App struct {
	ctx context.Context
}

// NewApp создает новый экземпляр сервиса приложения
func NewApp() *App {
	return &App{}
}

// Startup сохраняет контекст runtime при запуске приложения и регистрирует OnFileDrop
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	runtime.OnFileDrop(ctx, a.HandleFileDrop)
}

// HandleFileDrop обрабатывает перетаскивание файлов из Проводника Windows
func (a *App) HandleFileDrop(x, y int, paths []string) {
	if a.ctx == nil {
		return
	}
	for _, path := range paths {
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".md" && ext != ".txt" {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			runtime.EventsEmit(a.ctx, "file-drop-error", fmt.Sprintf("Ошибка чтения файла '%s': %v", path, err))
			continue
		}
		doc, err := parser.ParseMarkdown(string(data))
		if err != nil {
			runtime.EventsEmit(a.ctx, "file-drop-error", fmt.Sprintf("Ошибка парсинга файла '%s': %v", filepath.Base(path), err))
			continue
		}
		doc.FilePath = path
		doc.IsModified = false
		runtime.EventsEmit(a.ctx, "file-dropped", doc)
	}
}

// OpenMarkdownDialog открывает системный диалог выбора файла .md/.txt и парсит его
func (a *App) OpenMarkdownDialog() (*model.Document, error) {
	if a.ctx == nil {
		return nil, fmt.Errorf("контекст приложения не инициализирован")
	}

	selectedPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Открыть заявку в формате Markdown",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Файлы Markdown (*.md;*.txt)",
				Pattern:     "*.md;*.txt",
			},
			{
				DisplayName: "Все файлы (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия диалога: %w", err)
	}
	if selectedPath == "" {
		return nil, nil // Пользователь отменил выбор
	}

	data, err := os.ReadFile(selectedPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл '%s': %w", selectedPath, err)
	}

	doc, err := parser.ParseMarkdown(string(data))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла '%s': %w", filepath.Base(selectedPath), err)
	}

	doc.FilePath = selectedPath
	doc.IsModified = false
	return doc, nil
}
