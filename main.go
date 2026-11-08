package main

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"tool-bid-app/internal/app"
)

var GitCommit = "dev"

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Создаем экземпляр сервиса приложения
	appService := app.NewApp()

	appTitle := "Управление заявками на инструмент"
	if GitCommit != "" && GitCommit != "dev" {
		appTitle = fmt.Sprintf("Управление заявками на инструмент [%s]", GitCommit)
	}

	// Запуск приложения с окном 1440x900, минимумом 1100x700 и цветом Catppuccin Mocha (#1e1e2e)
	err := wails.Run(&options.App{
		Title:             appTitle,
		Width:             1440,
		Height:            900,
		MinWidth:          1100,
		MinHeight:         700,
		BackgroundColour:  &options.RGBA{R: 30, G: 30, B: 46, A: 255}, // #1e1e2e (Catppuccin Mocha)
		OnStartup:         appService.Startup,
		Bind: []interface{}{
			appService,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop:     true,
			DisableWebViewDrop: false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
