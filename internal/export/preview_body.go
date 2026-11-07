package export

const previewHTMLBody = `<body>
<div class="page-container">
  <!-- Верхний заголовок: Название заявки с подчеркиванием номера слева, блок подразделения справа -->
  <table class="top-header-table">
    <tr>
      <td style="width: 58%; vertical-align: bottom;">
        <table style="width: 100%; border-collapse: collapse;">
          <tr>
            <td class="bid-title">Заявка на закупку</td>
            <td style="width: 100%; vertical-align: bottom;">
              <span class="bid-number-underline">№{{.Meta.BidNum}} от {{.Meta.BidDate}} г.</span>
            </td>
          </tr>
        </table>
      </td>
      <td style="width: 42%; vertical-align: top; text-align: right;">
        <div class="dept-box-container">
          <div class="dept-box-header">ПОДРАЗДЕЛЕНИЕ</div>
          <div class="dept-box-value">{{if .Meta.Dept}}{{.Meta.Dept}}{{else}}&mdash;{{end}}</div>
        </div>
      </td>
    </tr>
  </table>

  <!-- Таблица подписантов (2 строки: Инициатор и Куратор) -->
  <table class="signatories-table">
    <tr>
      <td style="width: 28%;" class="sig-title">ИНИЦИАТОР ЗАКУПКИ</td>
      <td style="width: 28%;">{{if .Meta.InitRole}}{{.Meta.InitRole}}{{else}}&mdash;{{end}}</td>
      <td style="width: 24%;">{{if .Meta.InitFio}}{{.Meta.InitFio}}{{else}}&mdash;{{end}}</td>
      <td style="width: 10%;" class="sig-sign-col">Подп.</td>
      <td style="width: 10%;" class="sig-phone-col">Конт.тел.</td>
    </tr>
    <tr>
      <td class="sig-title">КУРАТОР ЗАКУПКИ</td>
      <td>{{if .Meta.CurRole}}{{.Meta.CurRole}}{{else}}&mdash;{{end}}</td>
      <td>{{if .Meta.CurFio}}{{.Meta.CurFio}}{{else}}&mdash;{{end}}</td>
      <td class="sig-sign-col">Подп.</td>
      <td class="sig-phone-col">Конт.тел.</td>
    </tr>
  </table>

  <!-- Раздел I: Общие сведения -->
  <div class="section-title">I. Общие сведения</div>
  <table class="section-table">
    <thead>
      <tr>
        <th class="sec1-col-num">№ п/п</th>
        <th class="sec1-col-name">Наименование пункта</th>
        <th class="sec1-col-desc">Характеристики товара/услуги</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td class="col-center">1</td>
        <td class="col-bold">{{if .Meta.ShopName}}{{.Meta.ShopName}}{{else}}&mdash;{{end}}</td>
        <td>&mdash;</td>
      </tr>
      <tr>
        <td class="col-center">2</td>
        <td>Желаемая дата поставки</td>
        <td>{{if .Meta.DeliveryDate}}{{.Meta.DeliveryDate}}{{else}}&mdash;{{end}}</td>
      </tr>
      <tr>
        <td class="col-center">3</td>
        <td>Требования к качеству товара, работ, услуг</td>
        <td>{{if .Meta.QualityReq}}{{.Meta.QualityReq}}{{else}}&mdash;{{end}}</td>
      </tr>
      <tr>
        <td class="col-center">4</td>
        <td>Цель закупки</td>
        <td>{{if .Meta.Purpose}}{{.Meta.Purpose}}{{else}}&mdash;{{end}}</td>
      </tr>
    </tbody>
  </table>

  <!-- Раздел II: Техническая часть документации -->
  <div class="section-title">II. Техническая часть документации</div>
  <table class="items-table">
    <thead>
      <tr>
        <th class="sec2-col-num">№ п/п</th>
        <th class="sec2-col-name">Наименование товара</th>
        <th class="sec2-col-desc">Характеристики товара</th>
        <th class="sec2-col-unit">Ед. изм.</th>
        <th class="sec2-col-qty">Кол-во</th>
      </tr>
    </thead>
    <tbody>
      {{if .Items}}
        {{range .Items}}
        <tr>
          <td class="sec2-col-num">{{.Num}}</td>
          <td class="sec2-col-name">{{.Name}}</td>
          <td class="sec2-col-desc">{{.Desc}}</td>
          <td class="sec2-col-unit">{{.Unit}}</td>
          <td class="sec2-col-qty">{{.Qty}}</td>
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
</div>
</body>
</html>`
