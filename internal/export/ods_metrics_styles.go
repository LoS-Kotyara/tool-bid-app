package export

import (
	"bytes"
	"encoding/xml"
)

func escapeXML(s string) string {
	var buf bytes.Buffer
	xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

func getODSColStyles() string {
	return `    <style:style style:name="coA" style:family="table-column"><style:table-column-properties style:column-width="1.0cm"/></style:style>
    <style:style style:name="coB" style:family="table-column"><style:table-column-properties style:column-width="6.0cm"/></style:style>
    <style:style style:name="coC" style:family="table-column"><style:table-column-properties style:column-width="4.3cm"/></style:style>
    <style:style style:name="coD" style:family="table-column"><style:table-column-properties style:column-width="4.3cm"/></style:style>
    <style:style style:name="coE" style:family="table-column"><style:table-column-properties style:column-width="1.7cm"/></style:style>
    <style:style style:name="coF" style:family="table-column"><style:table-column-properties style:column-width="1.9cm"/></style:style>
`
}

func getODSRowStyles() string {
	return `    <style:style style:name="ro1" style:family="table-row"><style:table-row-properties style:row-height="15.0pt"/></style:style>
    <style:style style:name="ro2" style:family="table-row"><style:table-row-properties style:row-height="5.25pt"/></style:style>
    <style:style style:name="ro3" style:family="table-row"><style:table-row-properties style:row-height="15.0pt"/></style:style>
    <style:style style:name="ro4" style:family="table-row"><style:table-row-properties style:row-height="4.5pt"/></style:style>
    <style:style style:name="ro5" style:family="table-row"><style:table-row-properties style:row-height="27.0pt"/></style:style>
    <style:style style:name="ro6" style:family="table-row"><style:table-row-properties style:row-height="35.45pt"/></style:style>
    <style:style style:name="ro7" style:family="table-row"><style:table-row-properties style:row-height="33.6pt"/></style:style>
    <style:style style:name="ro8" style:family="table-row"><style:table-row-properties style:row-height="4.5pt"/></style:style>
    <style:style style:name="ro9" style:family="table-row"><style:table-row-properties style:row-height="16.5pt"/></style:style>
    <style:style style:name="ro10" style:family="table-row"><style:table-row-properties style:row-height="4.5pt"/></style:style>
    <style:style style:name="ro11" style:family="table-row"><style:table-row-properties style:row-height="32.25pt"/></style:style>
    <style:style style:name="ro12" style:family="table-row"><style:table-row-properties style:row-height="54.75pt"/></style:style>
    <style:style style:name="ro13" style:family="table-row"><style:table-row-properties style:row-height="15.0pt"/></style:style>
    <style:style style:name="ro14" style:family="table-row"><style:table-row-properties style:row-height="33.0pt"/></style:style>
    <style:style style:name="ro15" style:family="table-row"><style:table-row-properties style:row-height="54.0pt"/></style:style>
    <style:style style:name="ro16" style:family="table-row"><style:table-row-properties style:row-height="3.0pt"/></style:style>
    <style:style style:name="ro17" style:family="table-row"><style:table-row-properties style:row-height="16.5pt"/></style:style>
    <style:style style:name="ro18" style:family="table-row"><style:table-row-properties style:row-height="4.5pt"/></style:style>
    <style:style style:name="ro19" style:family="table-row"><style:table-row-properties style:row-height="36.0pt"/></style:style>
    <style:style style:name="roItem" style:family="table-row"><style:table-row-properties style:row-height="30.0pt"/></style:style>
`
}
