package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	headerRegex = regexp.MustCompile(`^#(?P<num>\S+)\s+от\s+(?P<date>\d{2}\.\d{2}\.\d{4})`)
	dateRegex   = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
)

// ValidateDateFormat проверяет соответствие даты формату ДД.ММ.ГГГГ и календарную валидность
func ValidateDateFormat(dateStr string) error {
	trimmed := strings.TrimSpace(dateStr)
	if trimmed == "" {
		return nil
	}
	if !dateRegex.MatchString(trimmed) {
		return fmt.Errorf("неверный формат даты '%s', ожидается ДД.ММ.ГГГГ", trimmed)
	}
	_, err := time.Parse("02.01.2006", trimmed)
	if err != nil {
		return fmt.Errorf("некорректное значение даты '%s': %w", trimmed, err)
	}
	return nil
}

// splitSignatory делит строку подписанта "Должность; ФИО"
func splitSignatory(value string) (role string, fio string) {
	val := strings.TrimRight(strings.TrimSpace(value), ";")
	val = strings.TrimSpace(val)

	parts := strings.SplitN(val, ";", 2)
	if len(parts) == 2 {
		r := strings.TrimRight(strings.TrimSpace(parts[0]), ";")
		f := strings.TrimRight(strings.TrimSpace(parts[1]), ";")
		for strings.HasSuffix(f, ";") || strings.HasSuffix(strings.TrimSpace(f), ";") {
			f = strings.TrimRight(strings.TrimSpace(f), ";")
		}
		return strings.TrimSpace(r), strings.TrimSpace(f)
	}
	return val, ""
}

// NormalizeQuantity проверяет и форматирует числовые значения количества
func NormalizeQuantity(qtyStr string) (string, error) {
	trimmed := strings.TrimSpace(qtyStr)
	if trimmed == "" {
		return "", errors.New("количество не указано")
	}
	normalized := strings.ReplaceAll(trimmed, ",", ".")
	val, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return "", fmt.Errorf("некорректное числовое количество: '%s'", qtyStr)
	}
	if val <= 0 {
		return "", fmt.Errorf("количество должно быть положительным числом, получено: %f", val)
	}
	if val == float64(int64(val)) {
		return strconv.FormatInt(int64(val), 10), nil
	}
	return strconv.FormatFloat(val, 'f', -1, 64), nil
}

