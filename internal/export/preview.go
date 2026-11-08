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

// GenerateDocumentPreviewHTML формирует автономный HTML-документ листа А4 для печати и предпросмотра
func GenerateDocumentPreviewHTML(doc *model.Document) (string, error) {
	if doc == nil {
		doc = model.NewEmptyDocument()
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
	if err := tmpl.Execute(&buf, doc); err != nil {
		return "", fmt.Errorf("ошибка рендеринга HTML-шаблона предпросмотра: %w", err)
	}

	return buf.String(), nil
}
