package export

// previewHTMLBody — тело печатной формы: дискретные листы А4 по эталону 85.pdf.
// Лист 1 содержит шапку заявки, подписантов, разделы I и II и первые 27 строк
// спецификации; каждый следующий лист начинается сразу с таблицы спецификации
// (шапка таблицы повторяется) и содержит до 42 строк. Колонтитулов нет.
// Все таблицы построены на строгой 6-колоночной сетке A-F (см. preview_header.go).
const previewHTMLBody = `<body>
{{range .Pages}}
{{$meta := $.Meta}}
<div class="a4-page" id="page-{{.PageNum}}">
{{if .IsFirst}}

  <!-- ===== Шапка документа ===== -->
  <!-- C1: наименование документа (единственный жирный элемент) | D1:F1: номер и дата заявки -->
  <!-- C3: графа ПОДРАЗДЕЛЕНИЕ (без рамок)                        | D3:F3: блок цеха (рамка 1.5pt) -->
  <table class="grid">
    <colgroup>
      <col class="col-a"><col class="col-b"><col class="col-c"><col class="col-d"><col class="col-e"><col class="col-f">
    </colgroup>
    <tbody>
      <tr>
        <td colspan="2"></td>
        <td class="doc-title">Заявка на закупку</td>
        <td class="bid-number" colspan="3">№{{$meta.BidNum}} от {{$meta.BidDate}}</td>
      </tr>
      <tr class="spacer-row">
        <td colspan="6"></td>
      </tr>
      <tr>
        <td colspan="2"></td>
        <td class="dept-label">ПОДРАЗДЕЛЕНИЕ</td>
        <td class="dept-box" colspan="3">{{dash $meta.Dept}}</td>
      </tr>
    </tbody>
  </table>

  <!-- ===== Сетка подписантов: строка заголовков граф, инициатор, куратор ===== -->
  <table class="grid sign-table">
    <colgroup>
      <col class="col-a"><col class="col-b"><col class="col-c"><col class="col-d"><col class="col-e"><col class="col-f">
    </colgroup>
    <tbody>
      <tr>
        <td colspan="2"></td>
        <td class="boxed centered">Должность</td>
        <td class="boxed centered">Ф.И.О.</td>
        <td class="boxed centered">Подп.</td>
        <td class="boxed centered">Конт.тел.</td>
      </tr>
      <tr class="sign-row">
        <td class="boxed" colspan="2">ИНИЦИАТОР ЗАКУПКИ</td>
        <td class="boxed">{{dash $meta.InitRole}}</td>
        <td class="boxed">{{dash $meta.InitFio}}</td>
        <td class="boxed"></td>
        <td class="boxed"></td>
      </tr>
      <tr class="sign-row">
        <td class="boxed" colspan="2">КУРАТОР ЗАКУПКИ</td>
        <td class="boxed">{{dash $meta.CurRole}}</td>
        <td class="boxed">{{dash $meta.CurFio}}</td>
        <td class="boxed"></td>
        <td class="boxed"></td>
      </tr>
    </tbody>
  </table>

  <!-- ===== Раздел I. Общие сведения: A — № п/п, B:C — наименование пункта, D:F — характеристики ===== -->
  <div class="section-title">I. Общие сведения</div>
  <table class="grid">
    <colgroup>
      <col class="col-a"><col class="col-b"><col class="col-c"><col class="col-d"><col class="col-e"><col class="col-f">
    </colgroup>
    <thead>
      <tr>
        <th class="boxed centered">№ п/п</th>
        <th class="boxed centered" colspan="2">Наименование пункта</th>
        <th class="boxed centered" colspan="3">Характеристики товара/услуги</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td class="boxed centered">1</td>
        <td class="boxed" colspan="2">{{dash $meta.ShopName}}</td>
        <td class="boxed" colspan="3">—</td>
      </tr>
      <tr>
        <td class="boxed centered">2</td>
        <td class="boxed" colspan="2">Желаемая дата поставки</td>
        <td class="boxed" colspan="3">{{dash $meta.DeliveryDate}}</td>
      </tr>
      <tr>
        <td class="boxed centered">3</td>
        <td class="boxed" colspan="2">Требования к качеству товара, работ, услуг</td>
        <td class="boxed" colspan="3">{{dash $meta.QualityReq}}</td>
      </tr>
      <tr>
        <td class="boxed centered">4</td>
        <td class="boxed" colspan="2">Цель закупки</td>
        <td class="boxed" colspan="3">{{dash $meta.Purpose}}</td>
      </tr>
    </tbody>
  </table>

  <!-- ===== Раздел II. Техническая часть документации ===== -->
  <div class="section-title">II. Техническая часть документации</div>
{{end}}

  <!-- Спецификация: 5 колонок 4.41% + 31.51% + 45.20% (C+D) + 8.98% + 9.89% = 100% -->
  <table class="items-table">
    <colgroup>
      <col class="col-num"><col class="col-name"><col class="col-desc"><col class="col-unit"><col class="col-qty">
    </colgroup>
    <thead>
      <tr>
        <th class="col-num">№<br>п/п</th>
        <th class="col-name">Наименов. товара</th>
        <th class="col-desc">Характеристики товара</th>
        <th class="col-unit">Ед. изм</th>
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
          <td class="empty-row" colspan="5">Позиции спецификации не добавлены</td>
        </tr>
      {{end}}
    </tbody>
  </table>

</div>
{{end}}
</body>
</html>`
