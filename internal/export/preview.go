package export

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"strings"

	"tool-bid-app/internal/model"
)

// GenerateDocumentPreviewHTML формирует автономный HTML-документ листа А4 для печати и предпросмотра
func GenerateDocumentPreviewHTML(doc *model.Document) (string, error) {
	if doc == nil {
		doc = model.NewEmptyDocument()
	}

	fullTemplate := previewHTMLHeader + previewHTMLBody

	tmpl, err := template.New("preview").Funcs(template.FuncMap{
		"escape": html.EscapeString,
		"trim":   strings.TrimSpace,
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
