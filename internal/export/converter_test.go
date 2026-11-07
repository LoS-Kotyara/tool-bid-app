package export_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"tool-bid-app/internal/export"
	"tool-bid-app/internal/parser"
)

func TestPythonConverterEndToEnd(t *testing.T) {
	// Проверяем наличие python
	_, err := exec.LookPath("python")
	if err != nil {
		t.Skip("python not available in environment, skipping python converter test")
	}

	doc := createSampleDoc()
	tmpDir := t.TempDir()
	xlsxPath := filepath.Join(tmpDir, "bid_52.xlsx")

	// 1. Генерируем XLSX через ExportXLSX
	if err := export.ExportXLSX(doc, xlsxPath); err != nil {
		t.Fatalf("ExportXLSX failed: %v", err)
	}

	// 2. Запускаем scripts/excel_to_md.py
	scriptPath := filepath.Join("..", "..", "scripts", "excel_to_md.py")
	absScript, err := filepath.Abs(scriptPath)
	if err != nil {
		t.Fatalf("failed to resolve script path: %v", err)
	}

	cmd := exec.Command("python", absScript, xlsxPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("python converter failed: %v\nOutput: %s", err, string(output))
	}

	// 3. Проверяем сгенерированный .md
	mdPath := filepath.Join(tmpDir, "bid_52.md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("failed to read converted markdown file: %v", err)
	}

	mdContent := string(data)
	if !strings.HasPrefix(mdContent, "#52 от 11.06.2026") {
		t.Errorf("markdown header mismatch: %s", mdContent)
	}

	// 4. Парсим получившийся .md нашим Go-парсером
	parsedDoc, err := parser.ParseMarkdown(mdContent)
	if err != nil {
		t.Fatalf("parser.ParseMarkdown failed on converted markdown: %v", err)
	}

	if parsedDoc.Meta.BidNum != "52" {
		t.Errorf("expected BidNum '52', got '%s'", parsedDoc.Meta.BidNum)
	}
	if parsedDoc.Meta.Dept != "Цех 3" {
		t.Errorf("expected Dept 'Цех 3', got '%s'", parsedDoc.Meta.Dept)
	}
	if parsedDoc.Meta.InitRole != "Инженер-программист" {
		t.Errorf("expected InitRole 'Инженер-программист', got '%s'", parsedDoc.Meta.InitRole)
	}
	if parsedDoc.Meta.InitFio != "Давиденко А.А." {
		t.Errorf("expected InitFio 'Давиденко А.А.', got '%s'", parsedDoc.Meta.InitFio)
	}
	if len(parsedDoc.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(parsedDoc.Items))
	}
	if parsedDoc.Items[0].Name != "Сверло т/с" {
		t.Errorf("expected item 0 Name 'Сверло т/с', got '%s'", parsedDoc.Items[0].Name)
	}
	if parsedDoc.Items[0].Qty != "2" {
		t.Errorf("expected item 0 Qty '2', got '%s'", parsedDoc.Items[0].Qty)
	}
}
