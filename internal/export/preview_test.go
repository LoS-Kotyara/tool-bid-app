package export

import (
	"strings"
	"testing"

	"tool-bid-app/internal/model"
)

func TestGenerateDocumentPreviewHTML(t *testing.T) {
	doc := model.NewEmptyDocument()
	doc.Meta.BidNum = "52"
	doc.Meta.BidDate = "11.06.2026"
	doc.Meta.Dept = "Цех 3"
	doc.Meta.ShopName = "Механосборочный цех № 3"
	doc.Meta.InitRole = "Инженер-программист"
	doc.Meta.InitFio = "Давиденко А.А."
	doc.Meta.CurRole = "Начальник цеха №3"
	doc.Meta.CurFio = "Рыбкин П.В."
	doc.Meta.DeliveryDate = "15.06.2026"
	doc.Meta.QualityReq = "Без дефектов, заводской сертификат качества"
	doc.Meta.Purpose = "Для изготовления КЭМЗ БК8.035.825"

	doc.Items = []model.Item{
		{
			Num:  1,
			Name: "Сверло т/с",
			Desc: "EDS08090-0, ExTex D9.0x8D",
			Unit: "шт.",
			Qty:  "2",
		},
		{
			Num:  2,
			Name: "Паста алмазная",
			Desc: "АСМ 20/14; 16/12",
			Unit: "шт.",
			Qty:  "1",
		},
	}

	htmlStr, err := GenerateDocumentPreviewHTML(doc)
	if err != nil {
		t.Fatalf("GenerateDocumentPreviewHTML returned error: %v", err)
	}

	requiredSubstrings := []string{
		"@page {",
		"size: A4 portrait;",
		"margin: 10mm;",
		"display: table-header-group;",
		"page-break-inside: avoid;",
		"Заявка на закупку",
		"ПОДРАЗДЕЛЕНИЕ",
		"ИНИЦИАТОР ЗАКУПКИ",
		"КУРАТОР ЗАКУПКИ",
		"Подп.",
		"Конт.тел.",
		"I. Общие сведения",
		"Наименование пункта",
		"Характеристики товара/услуги",
		"II. Техническая часть документации",
		"№52 от 11.06.2026 г.",
		"Цех 3",
		"Механосборочный цех № 3",
		"Инженер-программист",
		"Давиденко А.А.",
		"Начальник цеха №3",
		"Рыбкин П.В.",
		"Желаемая дата поставки",
		"15.06.2026",
		"Требования к качеству товара, работ, услуг",
		"Без дефектов, заводской сертификат качества",
		"Цель закупки",
		"Для изготовления КЭМЗ БК8.035.825",
		"Сверло т/с",
		"EDS08090-0, ExTex D9.0x8D",
		"Паста алмазная",
		"АСМ 20/14; 16/12",
	}

	for _, sub := range requiredSubstrings {
		if !strings.Contains(htmlStr, sub) {
			t.Errorf("HTML preview expected to contain '%s', but it was missing", sub)
		}
	}
}
