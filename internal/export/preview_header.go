package export

const previewHTMLHeader = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Заявка на закупку</title>
<style>
  /* ===== Печатный лист А4: поля 10 мм, без колонтитулов и номеров страниц ===== */
  @page {
    size: A4 portrait;
    margin: 10mm;
  }
  * {
    box-sizing: border-box;
  }
  html {
    margin: 0;
    padding: 0;
  }
  body {
    background-color: #2b2d3a;
    margin: 0;
    padding: 20px 0;
    display: flex;
    justify-content: center;
    color: #000000;
    font-family: Calibri, Arial, sans-serif;
    font-size: 12pt;
    line-height: 1.15;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
  .page-sheet {
    width: 210mm;
    min-height: 297mm;
    padding: 10mm;
    box-sizing: border-box;
    background: #ffffff;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
  }

  /* ===== Строгая 6-колоночная сетка 1С/Excel (A-F), базовая ширина 190 мм ===== */
  .col-a { width: 4.41%; }
  .col-b { width: 31.51%; }
  .col-c { width: 22.22%; }
  .col-d { width: 22.98%; }
  .col-e { width: 8.98%; }
  .col-f { width: 9.89%; }

  table.grid {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
  }
  table.grid th, table.grid td {
    border: none;
    padding: 2px 4px;
    font-size: 12pt;
    font-weight: normal;
    line-height: 1.15;
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
  .spacer-row td {
    height: 3mm;
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

  /* ===== Сетка подписантов ===== */
  .sign-table {
    margin-top: 4mm;
  }
  .sign-table tr.sign-row td {
    height: 8mm;
  }

  /* ===== Заголовки разделов ===== */
  .section-title {
    margin: 3mm 0 1.5mm 0;
    text-align: center;
    font-size: 12pt;
    font-weight: normal;
    white-space: nowrap;
  }

  /* ===== Раздел II. Непрерывный табличный поток с повтором шапки ===== */
  .items-table {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
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
    padding: 2px 4px;
    font-size: 11pt;
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

  @media print {
    body { background: transparent; padding: 0; display: block; }
    .page-sheet { width: 100%; min-height: auto; padding: 0; box-shadow: none; }
  }
</style>
</head>
`
