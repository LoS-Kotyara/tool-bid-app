package app

import (
	"testing"
)

func TestGetInstalledFonts(t *testing.T) {
	app := NewApp()
	fonts, err := app.GetInstalledFonts()
	if err != nil {
		t.Fatalf("GetInstalledFonts error: %v", err)
	}

	if len(fonts) == 0 {
		t.Errorf("Expected at least one installed font, got empty list")
	}

	foundMonospace := false
	for _, font := range fonts {
		if font == "Consolas" || font == "Courier New" || font == "Arial" || font == "Segoe UI" {
			foundMonospace = true
			break
		}
	}

	if !foundMonospace {
		t.Errorf("Expected to find standard Windows font in list, got: %v", fonts[:min(5, len(fonts))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
