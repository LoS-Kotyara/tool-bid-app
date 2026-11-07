package export

import (
	"fmt"
	"strings"
	"tool-bid-app/internal/model"
)

// appendODSRows1to8 генерирует строки 1-8 ODS (Шапка, подразделение, подписанты)
func appendODSRows1to8(sb *strings.Builder, doc *model.Document) {
	// Строка 1: B1: ЗАЯВКА №, C1:F1: [num] от [date] г.
	headerText := fmt.Sprintf("%s от %s г.", strings.TrimSpace(doc.Meta.BidNum), strings.TrimSpace(doc.Meta.BidDate))
	sb.WriteString(`      <table:table-row table:style-name="ro1">
        <table:table-cell table:style-name="ceDefault"/>
        <table:table-cell table:style-name="ceDefault" office:value-type="string"><text:p>ЗАЯВКА №</text:p></table:table-cell>
        <table:table-cell table:style-name="ceUnderline" office:value-type="string" table:number-columns-spanned="4"><text:p>` + escapeXML(headerText) + `</text:p></table:table-cell>
        <table:covered-table-cell/><table:covered-table-cell/><table:covered-table-cell/>
      </table:table-row>
`)

	// Строка 2: разделитель (5.25pt)
	sb.WriteString(`      <table:table-row table:style-name="ro2">
        <table:table-cell table:number-columns-repeated="6"/>
      </table:table-row>
`)

	// Строка 3: D3:F3 - Блок подразделения
	deptText := strings.TrimSpace(doc.Meta.Dept)
	sb.WriteString(`      <table:table-row table:style-name="ro3">
        <table:table-cell table:style-name="ceDefault" table:number-columns-repeated="3"/>
        <table:table-cell table:style-name="ceDept" office:value-type="string" table:number-columns-spanned="3"><text:p>` + escapeXML(deptText) + `</text:p></table:table-cell>
        <table:covered-table-cell/><table:covered-table-cell/>
      </table:table-row>
`)

	// Строки 4-5: разделители
	sb.WriteString(`      <table:table-row table:style-name="ro4"><table:table-cell table:number-columns-repeated="6"/></table:table-row>
      <table:table-row table:style-name="ro5"><table:table-cell table:number-columns-repeated="6"/></table:table-row>
`)

	// Строка 6: Инициатор закупки (A6:B6, C6, D6 пустая, E6:F6 ФИО)
	initRole := strings.TrimSpace(doc.Meta.InitRole)
	initFio := strings.TrimSpace(doc.Meta.InitFio)
	sb.WriteString(`      <table:table-row table:style-name="ro6">
        <table:table-cell table:style-name="ceSigTitle" office:value-type="string" table:number-columns-spanned="2"><text:p>ИНИЦИАТОР ЗАКУПКИ</text:p></table:table-cell>
        <table:covered-table-cell/>
        <table:table-cell table:style-name="ceSigRole" office:value-type="string"><text:p>` + escapeXML(initRole) + `</text:p></table:table-cell>
        <table:table-cell table:style-name="ceDefault"/>
        <table:table-cell table:style-name="ceSigFio" office:value-type="string" table:number-columns-spanned="2"><text:p>` + escapeXML(initFio) + `</text:p></table:table-cell>
        <table:covered-table-cell/>
      </table:table-row>
`)

	// Строка 7: Куратор закупки (A7:B7, C7, D7 пустая, E7:F7 ФИО)
	curRole := strings.TrimSpace(doc.Meta.CurRole)
	curFio := strings.TrimSpace(doc.Meta.CurFio)
	sb.WriteString(`      <table:table-row table:style-name="ro7">
        <table:table-cell table:style-name="ceSigTitle" office:value-type="string" table:number-columns-spanned="2"><text:p>КУРАТОР ЗАКУПКИ</text:p></table:table-cell>
        <table:covered-table-cell/>
        <table:table-cell table:style-name="ceSigRole" office:value-type="string"><text:p>` + escapeXML(curRole) + `</text:p></table:table-cell>
        <table:table-cell table:style-name="ceDefault"/>
        <table:table-cell table:style-name="ceSigFio" office:value-type="string" table:number-columns-spanned="2"><text:p>` + escapeXML(curFio) + `</text:p></table:table-cell>
        <table:covered-table-cell/>
      </table:table-row>
`)

	// Строка 8: разделитель (4.5pt)
	sb.WriteString(`      <table:table-row table:style-name="ro8"><table:table-cell table:number-columns-repeated="6"/></table:table-row>
`)
}
