package export

import (
	"strings"
	"testing"

	"tool-bid-app/internal/model"
)

// newPreviewTestDocument собирает заполненный документ для проверки рендеринга печатной формы
func newPreviewTestDocument() *model.Document {
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

	return doc
}

// renderPreview рендерит HTML-предпросмотр и падает при ошибке шаблона
func renderPreview(t *testing.T, doc *model.Document) string {
	t.Helper()

	htmlStr, err := GenerateDocumentPreviewHTML(doc)
	if err != nil {
		t.Fatalf("GenerateDocumentPreviewHTML returned error: %v", err)
	}

	return htmlStr
}

func assertPreviewContains(t *testing.T, htmlStr string, substrings []string) {
	t.Helper()

	for _, sub := range substrings {
		if !strings.Contains(htmlStr, sub) {
			t.Errorf("HTML preview expected to contain %q, but it was missing", sub)
		}
	}
}

func assertPreviewNotContains(t *testing.T, htmlStr string, substrings []string) {
	t.Helper()

	for _, sub := range substrings {
		if strings.Contains(htmlStr, sub) {
			t.Errorf("HTML preview must not contain %q", sub)
		}
	}
}

// assertPreviewOrder проверяет, что фрагменты идут в документе в указанном порядке
func assertPreviewOrder(t *testing.T, htmlStr string, fragments []string) {
	t.Helper()

	prevIndex := -1
	for _, fragment := range fragments {
		index := strings.Index(htmlStr, fragment)
		if index < 0 {
			t.Fatalf("HTML preview expected to contain %q, but it was missing", fragment)
		}
		if index <= prevIndex {
			t.Errorf("HTML preview fragment %q is out of expected order", fragment)
		}
		prevIndex = index
	}
}

func TestGenerateDocumentPreviewHTML(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	requiredSubstrings := []string{
		"<!DOCTYPE html>",
		"@page {",
		"size: A4 portrait;",
		"margin: 10mm;",
		"display: table-header-group;",
		"page-break-inside: avoid;",
		"Заявка на закупку",
		"ПОДРАЗДЕЛЕНИЕ",
		"Должность",
		"Ф.И.О.",
		"Подп.",
		"Конт.тел.",
		"ИНИЦИАТОР ЗАКУПКИ",
		"КУРАТОР ЗАКУПКИ",
		"I. Общие сведения",
		"№ п/п",
		"Наименование пункта",
		"Характеристики товара/услуги",
		"II. Техническая часть документации",
		"№52 от 11.06.2026",
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

	assertPreviewContains(t, htmlStr, requiredSubstrings)
}

func TestPreviewHTMLBuildsExcelSixColumnGrid(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	columnWidths := []string{
		".col-a { width: 4.41%; }",
		".col-b { width: 31.51%; }",
		".col-c { width: 22.22%; }",
		".col-d { width: 22.98%; }",
		".col-e { width: 8.98%; }",
		".col-f { width: 9.89%; }",
	}
	assertPreviewContains(t, htmlStr, columnWidths)

	// Шапка, сетка подписантов и раздел I используют полную сетку A-F (6 колонок)
	sixColumnColgroup := `<col class="col-a"><col class="col-b"><col class="col-c">` +
		`<col class="col-d"><col class="col-e"><col class="col-f">`
	if count := strings.Count(htmlStr, sixColumnColgroup); count != 3 {
		t.Errorf("expected 3 six-column colgroups (header, signatories, section I), got %d", count)
	}

	// Раздел II: характеристики объединяют колонки C + D -> 45.20%
	assertPreviewContains(t, htmlStr, []string{
		`<col class="col-num"><col class="col-name"><col class="col-desc"><col class="col-unit"><col class="col-qty">`,
		".col-num { width: 4.41%; text-align: center; vertical-align: middle; }",
		".col-name { width: 31.51%; text-align: left; vertical-align: middle; }",
		".col-desc { width: 45.20%; text-align: left; vertical-align: top; }",
		".col-unit { width: 8.98%; text-align: center; vertical-align: middle; }",
		".col-qty { width: 9.89%; text-align: center; vertical-align: middle; }",
	})
}

func TestPreviewHTMLHeaderCellsUseGridPlacement(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	assertPreviewContains(t, htmlStr, []string{
		`<td class="doc-title">Заявка на закупку</td>`,
		`<td class="bid-number" colspan="3">№52 от 11.06.2026</td>`,
		`<td class="dept-label">ПОДРАЗДЕЛЕНИЕ</td>`,
		`<td class="dept-box" colspan="3">Цех 3</td>`,
		`<td class="boxed" colspan="2">ИНИЦИАТОР ЗАКУПКИ</td>`,
		`<td class="boxed" colspan="2">КУРАТОР ЗАКУПКИ</td>`,
	})

	// Подчеркивание номера заявки и рамка блока цеха 1.5pt
	assertPreviewContains(t, htmlStr, []string{
		".bid-number {\n    border-bottom: 1px solid #000000;",
		".dept-box {\n    border: 1.5pt solid #000000;",
	})
}


func TestPreviewHTMLHasNoPageFooters(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	forbidden := []string{
		"page-footer",
		"Лист ",
		"Стр. ",
		"page-number",
		"counter(page",
		"counter(pages",
		"Страница ",
		"page-break-after",
	}
	assertPreviewNotContains(t, htmlStr, forbidden)

	// Документ рендерится единым листом/потоком
	if count := strings.Count(htmlStr, `class="page-sheet"`); count != 1 {
		t.Errorf("expected exactly 1 page-sheet container, got %d", count)
	}
}

func TestPreviewHTMLUsesWhiteBackgroundOnly(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	// Серые плашки эталонной формы 1С отсутствуют
	assertPreviewNotContains(t, htmlStr, []string{
		"#d9d9d9",
		"#f2f2f2",
		"#e5e7eb",
		"#777777",
		"background-color: #d",
		"background-color: #e",
		"background-color: #f",
	})

	assertPreviewContains(t, htmlStr, []string{
		"background: #ffffff;",
		"color: #000000;",
	})

	// Тёмный фон предназначен только для экранного просмотра и снимается при печати
	assertPreviewContains(t, htmlStr, []string{
		"background-color: #2b2d3a;",
		".page-sheet { width: 100%; min-height: auto; padding: 0; box-shadow: none; }",
	})
}

func TestPreviewHTMLBoldOnlyForDocumentTitle(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	if count := strings.Count(htmlStr, "font-weight: bold"); count != 1 {
		t.Fatalf("expected exactly 1 bold declaration (document title), got %d", count)
	}

	assertPreviewContains(t, htmlStr, []string{
		".doc-title {\n    font-weight: bold;",
		".items-table th {\n    text-align: center;\n    vertical-align: middle;\n    font-weight: normal;",
		".section-title {\n    margin: 3mm 0 1.5mm 0;\n    text-align: center;\n    font-size: 12pt;\n    font-weight: normal;",
	})

	// Названия разделов и графы таблиц не выделены жирным
	assertPreviewNotContains(t, htmlStr, []string{
		"<b>",
		"<strong>",
		"font-weight: 700",
		"font-weight:bold",
	})
}

func TestPreviewHTMLKeepsFontAndBorderTokens(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	assertPreviewContains(t, htmlStr, []string{
		"font-family: Calibri, Arial, sans-serif;",
		"font-size: 12pt;",
		"border: 1px solid #000000;",
		"table-layout: fixed;",
		"border-collapse: collapse;",
		// Строки спецификации набираются кеглем 11pt
		".items-table th, .items-table td {\n    border: 1px solid #000000;\n    padding: 2px 4px;\n    font-size: 11pt;",
	})
}

func TestPreviewHTMLDocumentStructureOrder(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	// Порядок проверяется по разметке тела документа, без учета <head>
	bodyStart := strings.Index(htmlStr, "<body>")
	if bodyStart < 0 {
		t.Fatalf("HTML preview expected to contain <body>, but it was missing")
	}
	body := htmlStr[bodyStart:]

	assertPreviewOrder(t, body, []string{
		`<div class="page-sheet">`,
		`<td class="doc-title">Заявка на закупку</td>`,
		`ПОДРАЗДЕЛЕНИЕ`,
		`Должность`,
		`Ф.И.О.`,
		`Подп.`,
		`Конт.тел.`,
		`ИНИЦИАТОР ЗАКУПКИ`,
		`КУРАТОР ЗАКУПКИ`,
		`I. Общие сведения`,
		`Наименование пункта`,
		`Характеристики товара/услуги`,
		`II. Техническая часть документации`,
		`<th class="col-name">Наименов. товара</th>`,
		`<th class="col-desc">Характеристики товара</th>`,
		`<th class="col-unit">Ед. изм</th>`,
		`<th class="col-qty">Кол-во</th>`,
		`Сверло т/с`,
		`Паста алмазная`,
	})
}


func TestPreviewHTMLRendersItemsInContinuousTable(t *testing.T) {
	htmlStr := renderPreview(t, newPreviewTestDocument())

	assertPreviewContains(t, htmlStr, []string{
		`<td class="col-num">1</td>`,
		`<td class="col-name">Сверло т/с</td>`,
		`<td class="col-desc">EDS08090-0, ExTex D9.0x8D</td>`,
		`<td class="col-unit">шт.</td>`,
		`<td class="col-qty">2</td>`,
		`<td class="col-num">2</td>`,
		`<td class="col-name">Паста алмазная</td>`,
		`<td class="col-desc">АСМ 20/14; 16/12</td>`,
		`<td class="col-qty">1</td>`,
	})

	// Нумерация арабскими цифрами без точек
	assertPreviewNotContains(t, htmlStr, []string{
		`<td class="col-num">1.</td>`,
		`<td class="col-num">2.</td>`,
		"Позиции спецификации не добавлены",
	})

	// Шапка таблицы повторяется на каждой странице печати, строки не разрываются
	assertPreviewContains(t, htmlStr, []string{
		".items-table thead {\n    display: table-header-group;",
		".items-table tr {\n    page-break-inside: avoid;\n    break-inside: avoid;",
	})
}

func TestPreviewHTMLEmptyDocumentRendersDashes(t *testing.T) {
	htmlStr := renderPreview(t, nil)

	assertPreviewContains(t, htmlStr, []string{
		`<td class="dept-box" colspan="3">—</td>`,
		"Позиции спецификации не добавлены",
	})

	if count := strings.Count(htmlStr, ">—<"); count < 5 {
		t.Errorf("expected at least 5 dash placeholders for an empty document, got %d", count)
	}
}

func TestPreviewDashHelper(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  string
	}{
		{name: "empty string", value: "", want: "—"},
		{name: "spaces only", value: "   ", want: "—"},
		{name: "nil value", value: nil, want: "—"},
		{name: "filled string", value: "  Цех 3  ", want: "Цех 3"},
		{name: "integer value", value: 42, want: "42"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := previewDash(tc.value); got != tc.want {
				t.Errorf("previewDash(%#v) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}

