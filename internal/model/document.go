package model

import "time"

// Document представляет каноничную структуру заводской заявки
type Document struct {
	ID         string   `json:"id"`
	FilePath   string   `json:"filePath"`
	IsModified bool     `json:"isModified"`
	Meta       MetaInfo `json:"meta"`
	Items      []Item   `json:"items"`
}

// MetaInfo содержит реквизиты, подписантов и общие сведения по заявке
type MetaInfo struct {
	BidNum       string `json:"bidNum"`
	BidDate      string `json:"bidDate"`
	Dept         string `json:"dept"`         // "Подразделение" в шапке (напр. Цех 3)
	ShopName     string `json:"shopName"`     // "Пункт" (напр. Механосборочный цех № 3)
	InitRole     string `json:"initRole"`     // Должность инициатора
	InitFio      string `json:"initFio"`      // ФИО инициатора
	CurRole      string `json:"curRole"`      // Должность куратора
	CurFio       string `json:"curFio"`       // ФИО куратора
	DeliveryDate string `json:"deliveryDate"` // Желаемая дата поставки
	QualityReq   string `json:"qualityReq"`   // Требования к качеству
	Purpose      string `json:"purpose"`      // Цель закупки
}

// Item представляет позицию инструмента в спецификации
type Item struct {
	Num  int    `json:"num"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Unit string `json:"unit"`
	Qty  string `json:"qty"`
}

// NewEmptyDocument создает новый пустой документ с базовыми значениями по умолчанию
func NewEmptyDocument() *Document {
	today := time.Now().Format("02.01.2006")
	return &Document{
		ID:         "",
		FilePath:   "",
		IsModified: false,
		Meta: MetaInfo{
			BidNum:       "",
			BidDate:      today,
			Dept:         "",
			ShopName:     "",
			InitRole:     "",
			InitFio:      "",
			CurRole:      "",
			CurFio:       "",
			DeliveryDate: today,
			QualityReq:   "",
			Purpose:      "",
		},
		Items: make([]Item, 0),
	}
}
