package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"tool-bid-app/internal/export"
	"tool-bid-app/internal/parser"
)

func TestParseRealConvertedFiles(t *testing.T) {
	convertedDir := filepath.Join("..", "..", "testdata", "converted_md")
	outputDir := filepath.Join("..", "..", "testdata", "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	entries, err := os.ReadDir(convertedDir)
	if err != nil {
		t.Fatalf("Failed to read converted_md dir: %v", err)
	}

	mdFiles := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			mdFiles = append(mdFiles, filepath.Join(convertedDir, entry.Name()))
		}
	}

	if len(mdFiles) == 0 {
		t.Fatalf("No .md files found in %s", convertedDir)
	}

	for _, mdFile := range mdFiles {
		t.Run(filepath.Base(mdFile), func(t *testing.T) {
			content, err := os.ReadFile(mdFile)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", mdFile, err)
			}

			doc, err := parser.ParseMarkdown(string(content))
			if err != nil {
				t.Fatalf("ParseMarkdown failed for %s: %v", mdFile, err)
			}

			// Проверка базовых метаданных
			if doc.Meta.BidNum == "" {
				t.Errorf("%s: BidNum is empty", mdFile)
			}
			if doc.Meta.BidDate == "" {
				t.Errorf("%s: BidDate is empty", mdFile)
			}
			if doc.Meta.Dept == "" {
				t.Errorf("%s: Dept is empty", mdFile)
			}
			if doc.Meta.ShopName == "" {
				t.Errorf("%s: ShopName is empty", mdFile)
			}
			if doc.Meta.InitRole == "" || doc.Meta.InitFio == "" {
				t.Errorf("%s: Incomplete Initiator: role='%s', fio='%s'", mdFile, doc.Meta.InitRole, doc.Meta.InitFio)
			}
			if doc.Meta.CurRole == "" || doc.Meta.CurFio == "" {
				t.Errorf("%s: Incomplete Curator: role='%s', fio='%s'", mdFile, doc.Meta.CurRole, doc.Meta.CurFio)
			}
			if doc.Meta.DeliveryDate == "" {
				t.Errorf("%s: DeliveryDate is empty", mdFile)
			}
			if len(doc.Items) == 0 {
				t.Errorf("%s: Items list is empty", mdFile)
			}

			// Проверка каждой позиции
			for i, item := range doc.Items {
				if item.Num != i+1 {
					t.Errorf("%s: Item %d has unexpected Num %d", mdFile, i, item.Num)
				}
				if item.Name == "" {
					t.Errorf("%s: Item %d has empty Name", mdFile, i)
				}
				if item.Unit == "" {
					t.Errorf("%s: Item %d has empty Unit", mdFile, i)
				}
				if item.Qty == "" {
					t.Errorf("%s: Item %d has empty Qty", mdFile, i)
				}
			}

			// Проверка экспорта (Round-trip)
			baseName := filepath.Base(mdFile)
			baseWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]

			// 1. XLSX экспорт
			xlsxPath := filepath.Join(outputDir, baseWithoutExt+".xlsx")
			if err := export.ExportXLSX(doc, xlsxPath); err != nil {
				t.Fatalf("%s: ExportXLSX failed: %v", mdFile, err)
			}
			info, err := os.Stat(xlsxPath)
			if err != nil || info.Size() == 0 {
				t.Fatalf("%s: ExportXLSX generated invalid file", mdFile)
			}

			// 2. ODS экспорт
			odsPath := filepath.Join(outputDir, baseWithoutExt+".ods")
			if err := export.ExportODS(doc, odsPath); err != nil {
				t.Fatalf("%s: ExportODS failed: %v", mdFile, err)
			}
			infoOds, err := os.Stat(odsPath)
			if err != nil || infoOds.Size() == 0 {
				t.Fatalf("%s: ExportODS generated invalid file", mdFile)
			}
		})
	}
}
