package export

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"strings"

	"tool-bid-app/internal/model"
)

// previewPlaceholder — прочерк, выводимый в незаполненных графах печатной формы
const previewPlaceholder = "—"

// Стандарт нарезки спецификации по эталону 85.pdf: на первом листе помещается
// меньше строк из-за шапки заявки, разделов I и II; последующие листы — плотнее.
const (
	previewFirstPageRows = 27
	previewNextPageRows  = 42
)

// PageView — один дискретный лист А4 предпросмотра с собственной порцией позиций
type PageView struct {
	PageNum int
	IsFirst bool
	Items   []model.Item
}

// PreviewData — данные рендеринга HTML-предпросмотра: реквизиты заявки и набор листов
type PreviewData struct {
	Meta  model.MetaInfo
	Pages []PageView
}

// previewDash подставляет прочерк вместо пустого или состоящего из пробелов значения графы
func previewDash(value any) string {
	var text string
	switch v := value.(type) {
	case nil:
		text = ""
	case string:
		text = v
	default:
		text = fmt.Sprint(v)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return previewPlaceholder
	}
	return text
}

// buildPreviewPages нарезает позиции спецификации на дискретные листы А4:
// лист 1 — до 27 позиций, каждый последующий — строго порциями до 42 позиций.
func buildPreviewPages(items []model.Item) []PageView {
	pages := make([]PageView, 0, 1+len(items)/previewNextPageRows)

	appendPage := func(chunk []model.Item) {
		pages = append(pages, PageView{
			PageNum: len(pages) + 1,
			IsFirst: len(pages) == 0,
			Items:   chunk,
		})
	}

	rest := items
	if len(rest) > previewFirstPageRows {
		appendPage(rest[:previewFirstPageRows])
		rest = rest[previewFirstPageRows:]
	} else {
		appendPage(rest)
		rest = nil
	}

	for len(rest) > 0 {
		size := previewNextPageRows
		if len(rest) < size {
			size = len(rest)
		}
		appendPage(rest[:size])
		rest = rest[size:]
	}

	return pages
}

// GenerateDocumentPreviewHTML формирует автономный HTML-документ дискретных листов А4
// для печати и предпросмотра
func GenerateDocumentPreviewHTML(doc *model.Document) (string, error) {
	if doc == nil {
		doc = model.NewEmptyDocument()
	}

	data := PreviewData{
		Meta:  doc.Meta,
		Pages: buildPreviewPages(doc.Items),
	}

	fullTemplate := previewHTMLHeader + previewHTMLBody

	tmpl, err := template.New("preview").Funcs(template.FuncMap{
		"escape": html.EscapeString,
		"trim":   strings.TrimSpace,
		"dash":   previewDash,
	}).Parse(fullTemplate)
	if err != nil {
		return "", fmt.Errorf("ошибка компиляции HTML-шаблона предпросмотра: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("ошибка рендеринга HTML-шаблона предпросмотра: %w", err)
	}

	return buf.String(), nil
}
