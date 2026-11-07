package parser_test

import (
	"strings"
	"testing"

	"tool-bid-app/internal/model"
	"tool-bid-app/internal/parser"
)

const canonicalSample = `#52 от 11.06.2026
Подразделение: Цех 3
Пункт: Механосборочный цех № 3
Инициатор: Инженер-программист; Давиденко А.А.
Куратор: Начальник цеха №3; Рыбкин П.В.
Дата поставки: 15.06.2026
Требования: Без дефектов, заводской сертификат качества
Цель: Для изготовления КЭМЗ БК8.035.825

Сверло т/с; EDS08090-0, ExTex D9.0x8D; шт.; 2
Развертка т/с; D9.2H7; шт.; 2
Фреза концевая монолитная; HM60-S4-05130550, ExTex (D5); шт.; 4
`

func TestParseMarkdown_CanonicalSample(t *testing.T) {
	doc, err := parser.ParseMarkdown(canonicalSample)
	if err != nil {
		t.Fatalf("ParseMarkdown returned error: %v", err)
	}

	if doc.Meta.BidNum != "52" {
		t.Errorf("Expected BidNum '52', got '%s'", doc.Meta.BidNum)
	}
	if doc.Meta.BidDate != "11.06.2026" {
		t.Errorf("Expected BidDate '11.06.2026', got '%s'", doc.Meta.BidDate)
	}
	if doc.Meta.Dept != "Цех 3" {
		t.Errorf("Expected Dept 'Цех 3', got '%s'", doc.Meta.Dept)
	}
	if doc.Meta.ShopName != "Механосборочный цех № 3" {
		t.Errorf("Expected ShopName 'Механосборочный цех № 3', got '%s'", doc.Meta.ShopName)
	}
	if doc.Meta.InitRole != "Инженер-программист" {
		t.Errorf("Expected InitRole 'Инженер-программист', got '%s'", doc.Meta.InitRole)
	}
	if doc.Meta.InitFio != "Давиденко А.А." {
		t.Errorf("Expected InitFio 'Давиденко А.А.', got '%s'", doc.Meta.InitFio)
	}
	if doc.Meta.CurRole != "Начальник цеха №3" {
		t.Errorf("Expected CurRole 'Начальник цеха №3', got '%s'", doc.Meta.CurRole)
	}
	if doc.Meta.CurFio != "Рыбкин П.В." {
		t.Errorf("Expected CurFio 'Рыбкин П.В.', got '%s'", doc.Meta.CurFio)
	}
	if doc.Meta.DeliveryDate != "15.06.2026" {
		t.Errorf("Expected DeliveryDate '15.06.2026', got '%s'", doc.Meta.DeliveryDate)
	}
	if doc.Meta.QualityReq != "Без дефектов, заводской сертификат качества" {
		t.Errorf("Expected QualityReq 'Без дефектов, заводской сертификат качества', got '%s'", doc.Meta.QualityReq)
	}
	if doc.Meta.Purpose != "Для изготовления КЭМЗ БК8.035.825" {
		t.Errorf("Expected Purpose 'Для изготовления КЭМЗ БК8.035.825', got '%s'", doc.Meta.Purpose)
	}

	if len(doc.Items) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(doc.Items))
	}

	expectedItems := []model.Item{
		{Num: 1, Name: "Сверло т/с", Desc: "EDS08090-0, ExTex D9.0x8D", Unit: "шт.", Qty: "2"},
		{Num: 2, Name: "Развертка т/с", Desc: "D9.2H7", Unit: "шт.", Qty: "2"},
		{Num: 3, Name: "Фреза концевая монолитная", Desc: "HM60-S4-05130550, ExTex (D5)", Unit: "шт.", Qty: "4"},
	}

	for i, expected := range expectedItems {
		item := doc.Items[i]
		if item.Num != expected.Num || item.Name != expected.Name || item.Desc != expected.Desc || item.Unit != expected.Unit || item.Qty != expected.Qty {
			t.Errorf("Item %d mismatch. Expected %+v, got %+v", i+1, expected, item)
		}
	}
}


func TestRoundTrip_CanonicalSample(t *testing.T) {
	doc, err := parser.ParseMarkdown(canonicalSample)
	if err != nil {
		t.Fatalf("ParseMarkdown returned error: %v", err)
	}

	serialized, err := parser.SerializeToMarkdown(doc)
	if err != nil {
		t.Fatalf("SerializeToMarkdown returned error: %v", err)
	}

	normOrig := strings.ReplaceAll(canonicalSample, "\r\n", "\n")
	normSer := strings.ReplaceAll(serialized, "\r\n", "\n")

	if normOrig != normSer {
		t.Errorf("Roundtrip mismatch:\nExpected:\n%s\nGot:\n%s", normOrig, normSer)
	}

	doc2, err := parser.ParseMarkdown(serialized)
	if err != nil {
		t.Fatalf("Re-parsing serialized markdown returned error: %v", err)
	}

	if doc.Meta != doc2.Meta {
		t.Errorf("Metadata mismatch after re-parsing:\n%+v vs %+v", doc.Meta, doc2.Meta)
	}
	if len(doc.Items) != len(doc2.Items) {
		t.Fatalf("Items length mismatch: %d vs %d", len(doc.Items), len(doc2.Items))
	}
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "Missing header",
			input:   "Подразделение: Цех 3\n\nСверло; D10; шт.; 1",
			wantErr: "первая строка обязана содержать номер и дату",
		},
		{
			name:    "Invalid date format",
			input:   "#52 от 2026-06-11\n\nСверло; D10; шт.; 1",
			wantErr: "первая строка обязана содержать номер и дату",
		},
		{
			name:    "Invalid calendar date",
			input:   "#52 от 32.01.2026\n\nСверло; D10; шт.; 1",
			wantErr: "ошибка валидации даты заявки",
		},
		{
			name: "Missing empty line between meta and items",
			input: `#52 от 11.06.2026
Подразделение: Цех 3
Сверло; D10; шт.; 1`,
			wantErr: "отсутствует обязательная пустая строка",
		},
		{
			name: "List marker forbidden",
			input: `#52 от 11.06.2026
Подразделение: Цех 3

- Сверло; D10; шт.; 1`,
			wantErr: "без маркеров списка",
		},
		{
			name: "Invalid item columns",
			input: `#52 от 11.06.2026
Подразделение: Цех 3

Сверло; D10; шт.`,
			wantErr: "неверный формат позиции спецификации",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.ParseMarkdown(tt.input)
			if err == nil {
				t.Fatalf("Expected error containing '%s', got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing '%s', got '%v'", tt.wantErr, err)
			}
		})
	}
}

func TestSignatoryTrailingSemicolons(t *testing.T) {
	input := `#10 от 01.01.2026
Инициатор: Инженер ; Иванов И.И. ; ;

Фреза; D5; шт.; 2
`
	doc, err := parser.ParseMarkdown(input)
	if err != nil {
		t.Fatalf("ParseMarkdown returned error: %v", err)
	}

	if doc.Meta.InitRole != "Инженер" {
		t.Errorf("Expected InitRole 'Инженер', got '%s'", doc.Meta.InitRole)
	}
	if doc.Meta.InitFio != "Иванов И.И." {
		t.Errorf("Expected InitFio 'Иванов И.И.', got '%s'", doc.Meta.InitFio)
	}
}

func TestNormalizeQuantity(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"2", "2", false},
		{"2.0", "2", false},
		{"2,5", "2.5", false},
		{"0", "", true},
		{"-5", "", true},
		{"abc", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		got, err := parser.NormalizeQuantity(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("NormalizeQuantity(%q) err = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("NormalizeQuantity(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
