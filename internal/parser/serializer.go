package parser

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"tool-bid-app/internal/model"
)

// SerializeToMarkdown преобразует структуру model.Document в каноничный Markdown
func SerializeToMarkdown(doc *model.Document) (string, error) {
	if doc == nil {
		return "", errors.New("документ равен nil")
	}

	var sb strings.Builder

	bidNum := strings.TrimSpace(doc.Meta.BidNum)
	if bidNum == "" {
		bidNum = "0"
	}
	bidDate := strings.TrimSpace(doc.Meta.BidDate)
	if bidDate == "" {
		bidDate = time.Now().Format("02.01.2006")
	}

	sb.WriteString(fmt.Sprintf("#%s от %s\n", bidNum, bidDate))

	if doc.Meta.Dept != "" {
		sb.WriteString(fmt.Sprintf("Подразделение: %s\n", strings.TrimSpace(doc.Meta.Dept)))
	}
	if doc.Meta.ShopName != "" {
		sb.WriteString(fmt.Sprintf("Пункт: %s\n", strings.TrimSpace(doc.Meta.ShopName)))
	}

	initRole := strings.TrimSpace(doc.Meta.InitRole)
	initFio := strings.TrimSpace(doc.Meta.InitFio)
	if initRole != "" || initFio != "" {
		if initFio != "" {
			sb.WriteString(fmt.Sprintf("Инициатор: %s; %s\n", initRole, initFio))
		} else {
			sb.WriteString(fmt.Sprintf("Инициатор: %s\n", initRole))
		}
	}

	curRole := strings.TrimSpace(doc.Meta.CurRole)
	curFio := strings.TrimSpace(doc.Meta.CurFio)
	if curRole != "" || curFio != "" {
		if curFio != "" {
			sb.WriteString(fmt.Sprintf("Куратор: %s; %s\n", curRole, curFio))
		} else {
			sb.WriteString(fmt.Sprintf("Куратор: %s\n", curRole))
		}
	}

	if doc.Meta.DeliveryDate != "" {
		sb.WriteString(fmt.Sprintf("Дата поставки: %s\n", strings.TrimSpace(doc.Meta.DeliveryDate)))
	}
	if doc.Meta.QualityReq != "" {
		sb.WriteString(fmt.Sprintf("Требования: %s\n", strings.TrimSpace(doc.Meta.QualityReq)))
	}
	if doc.Meta.Purpose != "" {
		sb.WriteString(fmt.Sprintf("Цель: %s\n", strings.TrimSpace(doc.Meta.Purpose)))
	}

	// Обязательная пустая строка между блоком метаданных и спецификацией
	sb.WriteString("\n")

	for _, item := range doc.Items {
		sb.WriteString(fmt.Sprintf("%s; %s; %s; %s\n",
			strings.TrimSpace(item.Name),
			strings.TrimSpace(item.Desc),
			strings.TrimSpace(item.Unit),
			strings.TrimSpace(item.Qty),
		))
	}

	return sb.String(), nil
}
