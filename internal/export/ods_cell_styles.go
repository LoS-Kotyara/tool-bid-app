package export

func getODSCellStyles() string {
	return `    <style:style style:name="ceDefault" style:family="table-cell">
      <style:table-cell-properties fo:wrap-option="wrap" style:vertical-align="middle"/>
      <style:text-properties fo:font-size="11pt"/>
    </style:style>
    <style:style style:name="ceUnderline" style:family="table-cell">
      <style:table-cell-properties fo:border-bottom="0.8pt solid #000000" style:vertical-align="bottom"/>
      <style:text-properties fo:font-size="11pt" fo:font-weight="bold"/>
    </style:style>
    <style:style style:name="ceDept" style:family="table-cell">
      <style:table-cell-properties fo:border="1.2pt solid #000000" fo:padding="2pt" style:vertical-align="middle" fo:wrap-option="wrap"/>
      <style:text-properties fo:font-size="11pt" fo:font-weight="bold"/>
      <style:paragraph-properties fo:text-align="center"/>
    </style:style>
    <style:style style:name="ceSecTitle" style:family="table-cell">
      <style:table-cell-properties style:vertical-align="middle"/>
      <style:text-properties fo:font-size="11pt" fo:font-weight="bold"/>
    </style:style>
    <style:style style:name="ceSigTitle" style:family="table-cell">
      <style:table-cell-properties fo:wrap-option="wrap" style:vertical-align="middle"/>
      <style:text-properties fo:font-size="10pt" fo:font-weight="bold"/>
    </style:style>
    <style:style style:name="ceSigRole" style:family="table-cell">
      <style:table-cell-properties fo:wrap-option="wrap" style:vertical-align="middle"/>
      <style:text-properties fo:font-size="11pt"/>
    </style:style>
    <style:style style:name="ceSigFio" style:family="table-cell">
      <style:table-cell-properties fo:wrap-option="wrap" style:vertical-align="middle"/>
      <style:text-properties fo:font-size="11pt" fo:font-weight="bold"/>
      <style:paragraph-properties fo:text-align="center"/>
    </style:style>
    <style:style style:name="ceHeader" style:family="table-cell">
      <style:table-cell-properties fo:background-color="#f2f2f2" fo:border="0.5pt solid #000000" fo:padding="2pt" fo:wrap-option="wrap" style:vertical-align="middle"/>
      <style:text-properties fo:font-size="10pt" fo:font-weight="bold"/>
      <style:paragraph-properties fo:text-align="center"/>
    </style:style>
    <style:style style:name="ceDataLeft" style:family="table-cell">
      <style:table-cell-properties fo:border="0.5pt solid #000000" fo:padding="2pt" fo:wrap-option="wrap" style:vertical-align="middle"/>
      <style:text-properties fo:font-size="10pt"/>
      <style:paragraph-properties fo:text-align="left"/>
    </style:style>
    <style:style style:name="ceDataCenter" style:family="table-cell">
      <style:table-cell-properties fo:border="0.5pt solid #000000" fo:padding="2pt" fo:wrap-option="wrap" style:vertical-align="middle"/>
      <style:text-properties fo:font-size="10pt"/>
      <style:paragraph-properties fo:text-align="center"/>
    </style:style>
    <style:style style:name="ceDataBold" style:family="table-cell">
      <style:table-cell-properties fo:border="0.5pt solid #000000" fo:padding="2pt" fo:wrap-option="wrap" style:vertical-align="middle"/>
      <style:text-properties fo:font-size="10pt" fo:font-weight="bold"/>
      <style:paragraph-properties fo:text-align="left"/>
    </style:style>
`
}

func buildODSAutomaticStyles() string {
	return "  <office:automatic-styles>\n" +
		getODSColStyles() +
		getODSRowStyles() +
		getODSCellStyles() +
		"  </office:automatic-styles>\n"
}
