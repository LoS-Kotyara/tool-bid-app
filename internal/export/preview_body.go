package export

const previewHTMLBody = `<body>
<div class="page-container">
  <table class="header-grid">
    <tr>
      <td style="width: 60%;">
        <div class="bid-title">ЗАЯВКА НА ЗАКУПКУ ИНСТРУМЕНТА</div>
        <div style="margin-top: 3px;">
          № <span class="bid-number-underline">{{.Meta.BidNum}}</span> от <span>{{.Meta.BidDate}}</span> г.
        </div>
      </td>
      <td style="width: 40%; text-align: right;">
        <div class="dept-box">{{if .Meta.Dept}}{{.Meta.Dept}}{{else}}&mdash;{{end}}</div>
      </td>
    </tr>
  </table>
  <div class="section-title">I. Общие сведения</div>
  <table class="meta-table">
    <tr>
      <td class="meta-label">Пункт назначения:</td>
      <td class="meta-value-line">{{if .Meta.ShopName}}{{.Meta.ShopName}}{{else}}&mdash;{{end}}</td>
    </tr>
    <tr>
      <td class="meta-label">Инициатор:</td>
      <td class="meta-value-line">
        {{if .Meta.InitRole}}{{.Meta.InitRole}}{{end}}{{if and .Meta.InitRole .Meta.InitFio}}; {{end}}{{if .Meta.InitFio}}{{.Meta.InitFio}}{{end}}
        {{if and (not .Meta.InitRole) (not .Meta.InitFio)}}&mdash;{{end}}
      </td>
    </tr>
    <tr>
      <td class="meta-label">Куратор:</td>
      <td class="meta-value-line">
        {{if .Meta.CurRole}}{{.Meta.CurRole}}{{end}}{{if and .Meta.CurRole .Meta.CurFio}}; {{end}}{{if .Meta.CurFio}}{{.Meta.CurFio}}{{end}}
        {{if and (not .Meta.CurRole) (not .Meta.CurFio)}}&mdash;{{end}}
      </td>
    </tr>
    <tr>
      <td class="meta-label">Желаемая дата:</td>
      <td class="meta-value-line">{{if .Meta.DeliveryDate}}{{.Meta.DeliveryDate}}{{else}}&mdash;{{end}}</td>
    </tr>
    <tr>
      <td class="meta-label">Цель закупки:</td>
      <td class="meta-value-line">{{if .Meta.Purpose}}{{.Meta.Purpose}}{{else}}&mdash;{{end}}</td>
    </tr>
    <tr>
      <td class="meta-label">Требования к кач.:</td>
      <td class="meta-value-line">{{if .Meta.QualityReq}}{{.Meta.QualityReq}}{{else}}&mdash;{{end}}</td>
    </tr>
  </table>
  <div class="section-title">II. Техническая часть документации</div>
  <table class="items-table">
    <thead>
      <tr>
        <th class="col-num">№ п/п</th>
        <th class="col-name">Наименование инструмента</th>
        <th class="col-desc">Характеристики товара</th>
        <th class="col-unit">Ед. изм.</th>
        <th class="col-qty">Кол-во</th>
      </tr>
    </thead>
    <tbody>
      {{if .Items}}
        {{range .Items}}
        <tr>
          <td class="col-num">{{.Num}}</td>
          <td class="col-name">{{.Name}}</td>
          <td class="col-desc">{{.Desc}}</td>
          <td class="col-unit">{{.Unit}}</td>
          <td class="col-qty">{{.Qty}}</td>
        </tr>
        {{end}}
      {{else}}
        <tr>
          <td colspan="5" style="text-align: center; padding: 12px; color: #777777; font-style: italic;">
            Позиции спецификации не добавлены
          </td>
        </tr>
      {{end}}
    </tbody>
  </table>
  <table class="sign-footer">
    <tr>
      <td style="width: 18%; font-weight: bold;">Инициатор:</td>
      <td style="width: 32%;">
        <span class="sign-line"></span>
        <div style="font-size: 8pt; text-align: center; color: #555555;">(подпись, дата)</div>
      </td>
      <td style="width: 18%; font-weight: bold; text-align: right;">Куратор:</td>
      <td style="width: 32%;">
        <span class="sign-line"></span>
        <div style="font-size: 8pt; text-align: center; color: #555555;">(подпись, дата)</div>
      </td>
    </tr>
  </table>
</div>
</body>
</html>`
