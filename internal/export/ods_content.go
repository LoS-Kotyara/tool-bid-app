package export

import (
	"fmt"
	"strings"
	"tool-bid-app/internal/model"
)

// appendODSItems генерирует строки таблицы спецификации (позиции инструмента R20+)
func appendODSItems(sb *strings.Builder, items []model.Item) {
	for _, it := range items {
		sb.WriteString(fmt.Sprintf(`      <table:table-row table:style-name="roItem">
        <table:table-cell table:style-name="ceDataCenter" office:value-type="float" office:value="%d"><text:p>%d</text:p></table:table-cell>
        <table:table-cell table:style-name="ceDataLeft" office:value-type="string"><text:p>%s</text:p></table:table-cell>
        <table:table-cell table:style-name="ceDataLeft" office:value-type="string" table:number-columns-spanned="2"><text:p>%s</text:p></table:table-cell>
        <table:covered-table-cell/>
        <table:table-cell table:style-name="ceDataCenter" office:value-type="string"><text:p>%s</text:p></table:table-cell>
        <table:table-cell table:style-name="ceDataCenter" office:value-type="string"><text:p>%s</text:p></table:table-cell>
      </table:table-row>
`, it.Num, it.Num, escapeXML(strings.TrimSpace(it.Name)), escapeXML(strings.TrimSpace(it.Desc)), escapeXML(strings.TrimSpace(it.Unit)), escapeXML(strings.TrimSpace(it.Qty))))
	}
}

// buildODSContentXML формирует полный XML-документ content.xml для ODS-пакета
func buildODSContentXML(doc *model.Document) []byte {
	var sb strings.Builder

	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<office:document-content xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0"
  xmlns:style="urn:oasis:names:tc:opendocument:xmlns:style:1.0"
  xmlns:text="urn:oasis:names:tc:opendocument:xmlns:text:1.0"
  xmlns:table="urn:oasis:names:tc:opendocument:xmlns:table:1.0"
  xmlns:fo="urn:oasis:names:tc:opendocument:xmlns:xsl-fo-compatible:1.0"
  office:version="1.2">
  <office:scripts/>
  <office:font-face-decls>
    <style:font-face style:name="Calibri" svg:font-family="Calibri, 'Liberation Sans', sans-serif" xmlns:svg="urn:oasis:names:tc:opendocument:xmlns:svg-compatible:1.0"/>
  </office:font-face-decls>
`)

	sb.WriteString(buildODSAutomaticStyles())

	sb.WriteString(`  <office:body>
    <office:spreadsheet>
      <table:table table:name="Заявка">
        <!-- Определение 6 колонок A-F с заданными ширинами в см -->
        <table:table-column table:style-name="coA"/>
        <table:table-column table:style-name="coB"/>
        <table:table-column table:style-name="coC"/>
        <table:table-column table:style-name="coD"/>
        <table:table-column table:style-name="coE"/>
        <table:table-column table:style-name="coF"/>
`)

	appendODSRows1to8(&sb, doc)
	appendODSRows9to19(&sb, doc)
	appendODSItems(&sb, doc.Items)

	sb.WriteString(`      </table:table>
    </office:spreadsheet>
  </office:body>
</office:document-content>`)

	return []byte(sb.String())
}
