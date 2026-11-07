package export_test

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"tool-bid-app/internal/export"
)

func TestExportODS(t *testing.T) {
	doc := createSampleDoc()
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "bid_52.ods")

	err := export.ExportODS(doc, outPath)
	if err != nil {
		t.Fatalf("ExportODS failed: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("Stat failed on generated ods: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("Generated ods file is empty")
	}

	zr, err := zip.OpenReader(outPath)
	if err != nil {
		t.Fatalf("Failed to open zip reader on ods: %v", err)
	}
	defer zr.Close()

	if len(zr.File) < 4 {
		t.Fatalf("expected at least 4 files in ods package, got %d", len(zr.File))
	}

	// 1. mimetype: первый файл, без сжатия (Store)
	firstFile := zr.File[0]
	if firstFile.Name != "mimetype" {
		t.Errorf("expected first file to be 'mimetype', got '%s'", firstFile.Name)
	}
	if firstFile.Method != zip.Store {
		t.Errorf("expected mimetype to be stored without compression (Store), got method %d", firstFile.Method)
	}

	rc, err := firstFile.Open()
	if err != nil {
		t.Fatalf("failed to read mimetype: %v", err)
	}
	content, _ := io.ReadAll(rc)
	rc.Close()
	expectedMime := "application/vnd.oasis.opendocument.spreadsheet"
	if string(content) != expectedMime {
		t.Errorf("expected mimetype content '%s', got '%s'", expectedMime, string(content))
	}

	// 2. Наличие необходимых файлов в пакете
	fileNames := make(map[string]bool)
	for _, f := range zr.File {
		fileNames[f.Name] = true
	}

	for _, req := range []string{"META-INF/manifest.xml", "styles.xml", "content.xml"} {
		if !fileNames[req] {
			t.Errorf("required file '%s' missing in ODS package", req)
		}
	}

	// 3. Проверка содержимого content.xml
	var contentFile *zip.File
	for _, f := range zr.File {
		if f.Name == "content.xml" {
			contentFile = f
			break
		}
	}
	if contentFile != nil {
		rcContent, err := contentFile.Open()
		if err != nil {
			t.Fatalf("failed to open content.xml: %v", err)
		}
		rawContent, _ := io.ReadAll(rcContent)
		rcContent.Close()
		contentStr := string(rawContent)

		expectedSubstrings := []string{
			"52 от 11.06.2026 г.",
			"Цех 3",
			"Механосборочный цех № 3",
			"Инженер-программист",
			"Давиденко А.А.",
			"Начальник цеха №3",
			"Рыбкин П.В.",
			"Сверло т/с",
			"EDS08090-0, ExTex D9.0x8D",
			"table:number-columns-spanned=\"4\"",
			"fo:border=\"1.2pt solid #000000\"",
			"fo:border-bottom=\"0.8pt solid #000000\"",
		}
		for _, sub := range expectedSubstrings {
			if !strings.Contains(contentStr, sub) {
				t.Errorf("content.xml does not contain expected substring: %s", sub)
			}
		}
	}
}
