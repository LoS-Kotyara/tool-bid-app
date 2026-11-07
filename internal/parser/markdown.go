package parser

import (
	"errors"
	"fmt"
	"strings"

	"tool-bid-app/internal/model"
)

// ParseMarkdown разбирает каноничный Markdown заявки в структуру model.Document
func ParseMarkdown(content string) (*model.Document, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return nil, errors.New("пустой документ")
	}

	startIdx := 0
	for startIdx < len(lines) && strings.TrimSpace(lines[startIdx]) == "" {
		startIdx++
	}

	if startIdx >= len(lines) {
		return nil, errors.New("документ не содержит данных")
	}

	headerLine := strings.TrimSpace(lines[startIdx])
	matches := headerRegex.FindStringSubmatch(headerLine)
	if matches == nil {
		return nil, fmt.Errorf("первая строка обязана содержать номер и дату в формате '#[номер] от [ДД.ММ.ГГГГ]', получено: '%s'", headerLine)
	}

	bidNum := matches[headerRegex.SubexpIndex("num")]
	bidDate := matches[headerRegex.SubexpIndex("date")]

	if err := ValidateDateFormat(bidDate); err != nil {
		return nil, fmt.Errorf("ошибка валидации даты заявки: %w", err)
	}

	doc := model.NewEmptyDocument()
	doc.Meta.BidNum = bidNum
	doc.Meta.BidDate = bidDate

	currentLine := startIdx + 1
	hasEmptyLineSeparator := false

	for currentLine < len(lines) {
		line := lines[currentLine]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			hasEmptyLineSeparator = true
			currentLine++
			break
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) < 2 {
			// Если строка содержит разделитель спецификации ';' - это признак того, что пользователь
			// начал спецификацию без обязательной пустой строки
			if strings.Contains(trimmed, ";") {
				return nil, errors.New("отсутствует обязательная пустая строка между блоком метаданных и спецификацией")
			}
			return nil, fmt.Errorf("строка %d: неверный формат метаданных, ожидается 'Ключ: Значение': '%s'", currentLine+1, trimmed)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "Подразделение":
			doc.Meta.Dept = val
		case "Пункт":
			doc.Meta.ShopName = val
		case "Инициатор":
			role, fio := splitSignatory(val)
			doc.Meta.InitRole = role
			doc.Meta.InitFio = fio
		case "Куратор":
			role, fio := splitSignatory(val)
			doc.Meta.CurRole = role
			doc.Meta.CurFio = fio
		case "Дата поставки":
			if val != "" {
				if err := ValidateDateFormat(val); err != nil {
					return nil, fmt.Errorf("строка %d: ошибка даты поставки: %w", currentLine+1, err)
				}
			}
			doc.Meta.DeliveryDate = val
		case "Требования":
			doc.Meta.QualityReq = val
		case "Цель":
			doc.Meta.Purpose = val
		default:
			return nil, fmt.Errorf("строка %d: неизвестный префикс метаданных '%s:'", currentLine+1, key)
		}

		currentLine++
	}

	if !hasEmptyLineSeparator && currentLine < len(lines) {
		return nil, errors.New("отсутствует обязательная пустая строка между блоком метаданных и спецификацией")
	}

	itemNum := 1
	for currentLine < len(lines) {
		line := lines[currentLine]
		trimmed := strings.TrimSpace(line)
		currentLine++

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "+ ") {
			return nil, fmt.Errorf("строка %d: спецификация инструмента должна быть без маркеров списка ('-', '*'): '%s'", currentLine, trimmed)
		}

		parts := strings.Split(trimmed, ";")
		if len(parts) < 4 {
			return nil, fmt.Errorf("строка %d: неверный формат позиции спецификации (ожидается 'Наименование; Характеристики; Ед. изм.; Количество', найдено %d колонок): '%s'", currentLine, len(parts), trimmed)
		}

		name := strings.TrimSpace(parts[0])
		qty := strings.TrimSpace(parts[len(parts)-1])
		unit := strings.TrimSpace(parts[len(parts)-2])
		descParts := make([]string, 0, len(parts)-3)
		for _, dp := range parts[1 : len(parts)-2] {
			descParts = append(descParts, strings.TrimSpace(dp))
		}
		desc := strings.Join(descParts, "; ")

		if name == "" {
			return nil, fmt.Errorf("строка %d: наименование инструмента не может быть пустым", currentLine)
		}

		item := model.Item{
			Num:  itemNum,
			Name: name,
			Desc: desc,
			Unit: unit,
			Qty:  qty,
		}
		doc.Items = append(doc.Items, item)
		itemNum++
	}

	if len(doc.Items) > 0 && !hasEmptyLineSeparator {
		return nil, errors.New("отсутствует обязательная пустая строка перед спецификацией")
	}

	return doc, nil
}
