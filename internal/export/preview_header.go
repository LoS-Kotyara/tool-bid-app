package export

const previewHTMLHeader = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<style>
  @page {
    size: A4 portrait;
    margin: 10mm;
  }
  * {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }
  body {
    background-color: #ffffff;
    color: #000000;
    font-family: "Liberation Sans", "DejaVu Sans", "Arial", sans-serif;
    font-size: 10pt;
    line-height: 1.25;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
  .page-container {
    width: 190mm;
    margin: 0 auto;
    padding: 0;
    background: #ffffff;
  }
  .header-grid {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 2mm;
  }
  .header-grid td {
    vertical-align: top;
    padding: 2px 4px;
  }
  .bid-title {
    font-size: 12pt;
    font-weight: bold;
    text-transform: uppercase;
  }
  .bid-number-underline {
    border-bottom: 1.5px solid #000000;
    display: inline-block;
    min-width: 120px;
    padding-bottom: 1px;
    font-weight: bold;
    font-size: 11pt;
  }
  .dept-box {
    border: 2px solid #000000;
    padding: 4px 8px;
    font-weight: bold;
    text-align: center;
    font-size: 10.5pt;
    min-width: 180px;
  }
  .meta-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 2mm;
    margin-bottom: 3mm;
  }
  .meta-table td {
    padding: 3px 4px;
    vertical-align: top;
    font-size: 9.5pt;
  }
  .meta-label {
    font-weight: bold;
    white-space: nowrap;
    width: 22%;
  }
  .meta-value-line {
    border-bottom: 1px solid #777777;
    width: 78%;
  }
  .section-title {
    font-size: 10.5pt;
    font-weight: bold;
    background-color: #e9e9e9;
    padding: 3px 6px;
    margin-top: 3mm;
    margin-bottom: 1.5mm;
    border: 1px solid #cccccc;
    text-transform: uppercase;
  }
  table.items-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 2mm;
  }
  thead {
    display: table-header-group;
  }
  tr {
    page-break-inside: avoid;
  }
  table.items-table th {
    background-color: #f2f2f2;
    border: 1px solid #000000;
    padding: 5px 4px;
    font-size: 9.5pt;
    font-weight: bold;
    text-align: center;
    vertical-align: middle;
  }
  table.items-table td {
    border: 1px solid #000000;
    padding: 4px 5px;
    font-size: 9.5pt;
    vertical-align: middle;
  }
  .col-num { width: 6%; text-align: center; }
  .col-name { width: 34%; text-align: left; }
  .col-desc { width: 44%; text-align: left; }
  .col-unit { width: 8%; text-align: center; }
  .col-qty { width: 8%; text-align: right; }
  .sign-footer {
    margin-top: 8mm;
    width: 100%;
    border-collapse: collapse;
    page-break-inside: avoid;
  }
  .sign-footer td {
    padding: 4px 6px;
    font-size: 9.5pt;
    vertical-align: bottom;
  }
  .sign-line {
    border-bottom: 1px solid #000000;
    display: block;
    height: 14px;
  }
</style>
</head>
`
