package export

// previewHTMLHeader — стили автономного предпросмотра: дискретные листы А4 (210x297 мм)
// на экране и полное совпадение с печатью (@page A4, поля 8/10 мм, без колонтитулов).
const previewHTMLHeader = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Заявка на закупку</title>
<style>
  /* ===== Печатный лист А4: поля по эталону ReportLab (верх/низ 8 мм, право/лево 10 мм) ===== */
  @page {
    size: A4 portrait;
    margin: 8mm 10mm;
  }
  * {
    box-sizing: border-box;
  }
  html {
    margin: 0;
    padding: 0;
    height: 100%;
    /* Масштаб ленты листов: значение задаёт панель предпросмотра */
    --zoom-level: 1;
  }
  /* ===== Экранный режим ===== */
  /* overflow-y на body наследуется вьюпортом iframe: документ прокручивается
     единственным скроллбаром у правого края панели, двойного скролла нет */
  body {
    background-color: #2b2d3a;
    margin: 0;
    padding: 24px 0;
    width: 100%;
    min-height: 100%;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
    overflow-y: auto;
    color: #000000;
    font-family: Calibri, Arial, sans-serif;
    font-size: 12pt;
    line-height: 1.15;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }

  /* ===== Масштабирование ленты листов ===== */
  /* Обёртка задаёт фактическую (масштабированную) ширину ленты, поэтому листы
     остаются строго центрированными, а скроллбар — у правого края панели.
     margin: 0 auto центрирует ленту, пока хватает места, и прижимает её влево,
     когда панель уже масштабированного листа: левый край не обрезается. */
  .a4-scale-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: 0 0 auto;
    min-width: 100%;
    width: calc(210mm * var(--zoom-level));
    /* auto по вертикали: активный лист центрируется по центру экрана, пока он
       меньше видимой области, и остаётся доступным сверху при увеличении */
    margin: auto;
    padding: 24px 0;
    box-sizing: border-box;
  }

  /* ===== Дискретный лист А4: на экране виден только активный лист ===== */
  .a4-page {
    display: none;
    width: 210mm;
    height: 297mm;
    max-height: 297mm;
    /* Поля листа по эталону ReportLab: верх/низ 8 мм, право/лево 10 мм */
    padding: 8mm 10mm;
    background: #ffffff;
    box-shadow: 0 4px 18px rgba(0, 0, 0, 0.45);
    box-sizing: border-box;
    overflow: hidden;
    page-break-after: always;
    break-after: page;
    contain: layout paint;
    /* Масштаб приходит CSS-переменной от панели предпросмотра. transform не влияет
       на раскладку, поэтому нижний отступ компенсирует высоту листа до
       масштабированной (297 мм * zoom) — прокрутка знает реальный размер листа. */
    transform-origin: top center;
    transform: scale(var(--zoom-level));
    margin-bottom: calc(297mm * (var(--zoom-level) - 1));
  }
  /* В режиме экрана отображается только активная страница */
  .a4-page.active {
    display: block !important;
  }

  /* ===== Строгая 6-колоночная сетка 1С/Excel (A-F), базовая ширина 190 мм ===== */
  .col-a { width: 4.41%; }
  .col-b { width: 31.51%; }
  .col-c { width: 22.22%; }
  .col-d { width: 22.98%; }
  .col-e { width: 8.98%; }
  .col-f { width: 9.89%; }

  /* ===== Общие стили таблиц ===== */
  table {
    border-collapse: collapse;
    width: 100%;
  }
  table.grid {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
  }
  table.grid th, table.grid td {
    border: none;
    font-weight: normal;
    vertical-align: middle;
    text-align: left;
    word-break: break-word;
    overflow-wrap: anywhere;
  }
  table.grid .boxed {
    border: 1px solid #000000;
  }
  table.grid .centered {
    text-align: center;
  }

  /* ===== Верхний блок (шапка документа) ===== */
  .top-header-table td {
    padding: 1px 2px;
    font-size: 11pt;
    line-height: 1.15;
  }

  /* ===== Сетка подписантов ===== */
  .signatories-table {
    margin-top: 3mm;
  }
  .signatories-table td {
    padding: 1.5px 3px;
    font-size: 10.5pt;
    line-height: 1.15;
  }
  /* Рамки рисуются только у граф с классом boxed: пустая ячейка A:B остаётся без рамки */
  .signatories-table td.boxed {
    border: 1px solid #000000;
  }
  .signatories-table tr.sign-row td {
    height: 8mm;
  }
  .spacer-row td {
    height: 2mm;
    padding: 0;
  }

  /* ===== Шапка документа ===== */
  .doc-title {
    font-weight: bold;
    white-space: nowrap;
    vertical-align: bottom;
    padding: 0 2px 1px 0;
  }
  .bid-number {
    border-bottom: 1px solid #000000;
    vertical-align: bottom;
    padding: 0 0 1px 4px;
  }
  .dept-label {
    text-align: center;
    vertical-align: middle;
    text-transform: uppercase;
    letter-spacing: 0.3px;
  }
  .dept-box {
    border: 1.5pt solid #000000;
    text-align: center;
    vertical-align: middle;
    padding: 4px 6px;
  }
  /* Подчеркивание номера заявки (D1:F1) */
  .bid-number-underline {
    border-bottom: 1px solid #000000 !important;
  }
  /* Рамка вокруг наименования подразделения (D3:F3) */
  .dept-box-border {
    border: 1.5pt solid #000000 !important;
  }

  /* ===== Заголовки разделов I и II ===== */
  .section-title {
    margin: 1.8mm 0 1mm 0;
    text-align: center;
    font-size: 11pt;
    font-weight: normal;
    line-height: 1.1;
    white-space: nowrap;
  }

  /* ===== Таблица Раздела I (Общие сведения) ===== */
  .section-table {
    margin-top: 0;
  }
  .section-table td, .section-table th {
    border: 1px solid #000000;
    padding: 1.5px 3px;
    font-size: 10.5pt;
    line-height: 1.15;
  }

  /* ===== Раздел II. Таблица спецификации: шапка повторяется на каждом листе ===== */
  .items-table {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
    margin-top: 1mm;
    font-family: Calibri, Arial, sans-serif;
  }
  .items-table thead {
    display: table-header-group;
  }
  .items-table tr {
    page-break-inside: avoid;
    break-inside: avoid;
  }
  .items-table th, .items-table td {
    border: 1px solid #000000;
    padding: 1.5px 3px;
    font-size: 10.5pt;
    line-height: 1.15;
    word-break: break-word;
    overflow-wrap: anywhere;
  }
  .items-table th {
    text-align: center;
    vertical-align: middle;
    font-weight: normal;
  }
  .col-num { width: 4.41%; text-align: center; vertical-align: middle; }
  .col-name { width: 31.51%; text-align: left; vertical-align: middle; }
  .col-desc { width: 45.20%; text-align: left; vertical-align: top; }
  .col-unit { width: 8.98%; text-align: center; vertical-align: middle; }
  .col-qty { width: 9.89%; text-align: center; vertical-align: middle; }
  .empty-row {
    text-align: center;
    vertical-align: middle;
    font-style: italic;
    padding: 8px 4px;
  }

  /* ===== Печать: показываются ВСЕ листы подряд без экранной косметики ===== */
  @media print {
    html {
      height: auto;
      --zoom-level: 1;
    }
    body {
      background: transparent !important;
      padding: 0 !important;
      margin: 0 !important;
      gap: 0;
      width: auto;
      min-height: 0;
      display: block;
      overflow: visible !important;
    }
    /* Масштаб и центрирующая обёртка при печати отключаются */
    .a4-scale-wrapper {
      display: block;
      width: auto;
      min-width: 0;
      margin: 0;
      padding: 0;
      transform: none;
    }
    .a4-page {
      display: block !important;
      box-shadow: none !important;
      width: 100% !important;
      height: 100% !important;
      max-height: none !important;
      padding: 0 !important;
      margin: 0 !important;
      overflow: visible;
      page-break-after: always !important;
      break-after: page !important;
      /* Экранный масштаб (inline-переменная) при печати не применяется */
      transform: none;
    }
    /* Последний лист не должен порождать пустую страницу */
    .a4-page:last-child {
      page-break-after: auto !important;
      break-after: auto !important;
    }
  }
</style>
</head>
`
