package export

import (
	"fmt"
	"strings"
	"tool-bid-app/internal/model"
)

// appendODSRows9to19 генерирует строки 9-19 ODS (Раздел I, таблица сведений, Раздел II, шапка позиций)
func appendODSRows9to19(sb *strings.Builder, doc *model.Document) {
	// Строка 9: I. Общие сведения
	sb.WriteString(`      <table:table-row table:style-name="ro9">
        <table:table-cell table:style-name="ceSecTitle" office:value-type="string" table:number-columns-spanned="6"><text:p>I. Общие сведения</text:p></table:table-cell>
        <table:covered-table-cell/><table:covered-table-cell/><table:covered-table-cell/><table:covered-table-cell/><table:covered-table-cell/>
      </table:table-row>
`)

	// Строка 10: разделитель (4.5pt)
	sb.WriteString(`      <table:table-row table:style-name="ro10"><table:table-cell table:number-columns-repeated="6"/></table:table-row>
`)

	// Строка 11: Шапка раздела I
	sb.WriteString(`      <table:table-row table:style-name="ro11">
        <table:table-cell table:style-name="ceHeader" office:value-type="string"><text:p>№</text:p></table:table-cell>
        <table:table-cell table:style-name="ceHeader" office:value-type="string" table:number-columns-spanned="2"><text:p>Наименование</text:p></table:table-cell>
        <table:covered-table-cell/>
        <table:table-cell table:style-name="ceHeader" office:value-type="string" table:number-columns-spanned="3"><text:p>Содержание</text:p></table:table-cell>
        <table:covered-table-cell/><table:covered-table-cell/>
      </table:table-row>
`)

	// Пункты Раздела I: Строки 12-15
	sec1 := []struct {
		style string
		num   string
		label string
		val   string
	}{
		{"ro12", "1", "Служебная записка (Пункт)", doc.Meta.ShopName},
		{"ro13", "2", "Срок поставки (не позднее)", doc.Meta.DeliveryDate},
		{"ro14", "3", "Требования к качеству", doc.Meta.QualityReq},
		{"ro15", "4", "Цель приобретения", doc.Meta.Purpose},
	}

	for _, item := range sec1 {
		sb.WriteString(fmt.Sprintf(`      <table:table-row table:style-name="%s">
        <table:table-cell table:style-name="ceDataCenter" office:value-type="string"><text:p>%s</text:p></table:table-cell>
        <table:table-cell table:style-name="ceDataBold" office:value-type="string" table:number-columns-spanned="2"><text:p>%s</text:p></table:table-cell>
        <table:covered-table-cell/>
        <table:table-cell table:style-name="ceDataLeft" office:value-type="string" table:number-columns-spanned="3"><text:p>%s</text:p></table:table-cell>
        <table:covered-table-cell/><table:covered-table-cell/>
      </table:table-row>
`, item.style, escapeXML(item.num), escapeXML(item.label), escapeXML(strings.TrimSpace(item.val))))
	}

	// Строка 16: разделитель (3.0pt)
	sb.WriteString(`      <table:table-row table:style-name="ro16"><table:table-cell table:number-columns-repeated="6"/></table:table-row>
`)

	// Строка 17: II. Техническая часть документации
	sb.WriteString(`      <table:table-row table:style-name="ro17">
        <table:table-cell table:style-name="ceSecTitle" office:value-type="string" table:number-columns-spanned="6"><text:p>II. Техническая часть документации</text:p></table:table-cell>
        <table:covered-table-cell/><table:covered-table-cell/><table:covered-table-cell/><table:covered-table-cell/><table:covered-table-cell/>
      </table:table-row>
`)

	// Строка 18: разделитель (4.5pt)
	sb.WriteString(`      <table:table-row table:style-name="ro18"><table:table-cell table:number-columns-repeated="6"/></table:table-row>
`)

	// Строка 19: Шапка таблицы спецификации инструментов
	sb.WriteString(`      <table:table-row table:style-name="ro19">
        <table:table-cell table:style-name="ceHeader" office:value-type="string"><text:p>№ п/п</text:p></table:table-cell>
        <table:table-cell table:style-name="ceHeader" office:value-type="string"><text:p>Наименование товара</text:p></table:table-cell>
        <table:table-cell table:style-name="ceHeader" office:value-type="string" table:number-columns-spanned="2"><text:p>Характеристики товара</text:p></table:table-cell>
        <table:covered-table-cell/>
        <table:table-cell table:style-name="ceHeader" office:value-type="string"><text:p>Ед. изм.</text:p></table:table-cell>
        <table:table-cell table:style-name="ceHeader" office:value-type="string"><text:p>Кол-во</text:p></table:table-cell>
      </table:table-row>
`)
}
