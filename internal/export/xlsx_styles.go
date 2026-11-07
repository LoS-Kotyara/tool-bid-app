package export

import "github.com/xuri/excelize/v2"

type xlsxStyles struct {
	normal          int
	headerUnderline int
	deptBox         int
	sectionHeader   int
	signatoryTitle  int
	signatoryRole   int
	signatoryFIO    int
	tableHeader     int
	tableCellLeft   int
	tableCellCenter int
	tableLabelBold  int
}

func initXLSXStyles(f *excelize.File) (*xlsxStyles, error) {
	fontName := "Calibri"

	stNormal, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 11},
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
	})

	stHeaderUnderline, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Vertical: "bottom"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	stDeptBox, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
	})

	stSectionHeader, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})

	stSignatoryTitle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 10, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
	})
	stSignatoryRole, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
	})
	stSignatoryFIO, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})

	stTableHeader, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 10, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"F2F2F2"}, Pattern: 1},
	})

	stTableCellLeft, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	stTableCellCenter, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	stTableLabelBold, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: fontName, Size: 10, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	return &xlsxStyles{
		normal:          stNormal,
		headerUnderline: stHeaderUnderline,
		deptBox:         stDeptBox,
		sectionHeader:   stSectionHeader,
		signatoryTitle:  stSignatoryTitle,
		signatoryRole:   stSignatoryRole,
		signatoryFIO:    stSignatoryFIO,
		tableHeader:     stTableHeader,
		tableCellLeft:   stTableCellLeft,
		tableCellCenter: stTableCellCenter,
		tableLabelBold:  stTableLabelBold,
	}, nil
}
