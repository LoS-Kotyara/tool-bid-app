package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
	"tool-bid-app/internal/export"
	"tool-bid-app/internal/model"
)

func createSampleDoc() *model.Document {
	return &model.Document{
		ID:         "test-doc-52",
		FilePath:   "sample.md",
		IsModified: false,
		Meta: model.MetaInfo{
			BidNum:       "52",
			BidDate:      "11.06.2026",
			Dept:         "Цех 3",
			ShopName:     "Механосборочный цех № 3",
			InitRole:     "Инженер-программист",
			InitFio:      "Давиденко А.А.",
			CurRole:      "Начальник цеха №3",
			CurFio:       "Рыбкин П.В.",
			DeliveryDate: "15.06.2026",
			QualityReq:   "Без дефектов, заводской сертификат качества",
			Purpose:      "Для изготовления КЭМЗ БК8.035.825",
		},
		Items: []model.Item{
			{
				Num:  1,
				Name: "Сверло т/с",
				Desc: "EDS08090-0, ExTex D9.0x8D",
				Unit: "шт.",
				Qty:  "2",
			},
			{
				Num:  2,
				Name: "Развертка т/с",
				Desc: "D9.2H7",
				Unit: "шт.",
				Qty:  "2",
			},
			{
				Num:  3,
				Name: "Фреза концевая монолитная",
				Desc: "HM60-S4-05130550, ExTex (D5)",
				Unit: "шт.",
				Qty:  "4",
			},
		},
	}
}

func TestExportXLSX(t *testing.T) {
	doc := createSampleDoc()
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "bid_52.xlsx")

	err := export.ExportXLSX(doc, outPath)
	if err != nil {
		t.Fatalf("ExportXLSX failed: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("Stat failed on generated xlsx: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("Generated xlsx file is empty")
	}

	f, err := excelize.OpenFile(outPath)
	if err != nil {
		t.Fatalf("OpenFile failed on generated xlsx: %v", err)
	}
	defer f.Close()

	sheet := "Заявка"

	b1, _ := f.GetCellValue(sheet, "B1")
	if b1 != "ЗАЯВКА №" {
		t.Errorf("expected B1 'ЗАЯВКА №', got '%s'", b1)
	}

	c1, _ := f.GetCellValue(sheet, "C1")
	if !strings.Contains(c1, "52") || !strings.Contains(c1, "11.06.2026") {
		t.Errorf("expected C1 to contain '52' and '11.06.2026', got '%s'", c1)
	}

	d3, _ := f.GetCellValue(sheet, "D3")
	if d3 != "Цех 3" {
		t.Errorf("expected D3 'Цех 3', got '%s'", d3)
	}

	c6, _ := f.GetCellValue(sheet, "C6")
	if c6 != "Инженер-программист" {
		t.Errorf("expected C6 'Инженер-программист', got '%s'", c6)
	}
	e6, _ := f.GetCellValue(sheet, "E6")
	if e6 != "Давиденко А.А." {
		t.Errorf("expected E6 'Давиденко А.А.', got '%s'", e6)
	}

	d12, _ := f.GetCellValue(sheet, "D12")
	if d12 != "Механосборочный цех № 3" {
		t.Errorf("expected D12 'Механосборочный цех № 3', got '%s'", d12)
	}

	b20, _ := f.GetCellValue(sheet, "B20")
	if b20 != "Сверло т/с" {
		t.Errorf("expected B20 'Сверло т/с', got '%s'", b20)
	}

	definedNames := f.GetDefinedName()
	foundPrintTitles := false
	for _, dn := range definedNames {
		if dn.Name == "_xlnm.Print_Titles" {
			foundPrintTitles = true
			if !strings.Contains(dn.RefersTo, "$19:$19") {
				t.Errorf("expected Print_Titles to refer to row 19, got: %s", dn.RefersTo)
			}
			break
		}
	}
	if !foundPrintTitles {
		t.Errorf("expected defined name _xlnm.Print_Titles not found")
	}

	merges, err := f.GetMergeCells(sheet)
	if err != nil {
		t.Fatalf("GetMergeCells failed: %v", err)
	}
	hasMergeC1F1 := false
	hasMergeD3F3 := false
	for _, m := range merges {
		if m.GetStartAxis() == "C1" && m.GetEndAxis() == "F1" {
			hasMergeC1F1 = true
		}
		if m.GetStartAxis() == "D3" && m.GetEndAxis() == "F3" {
			hasMergeD3F3 = true
		}
	}
	if !hasMergeC1F1 {
		t.Errorf("Merge C1:F1 not found")
	}
	if !hasMergeD3F3 {
		t.Errorf("Merge D3:F3 not found")
	}
}
