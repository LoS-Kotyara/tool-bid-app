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
  .top-header-table {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 2mm;
  }
  .top-header-table td {
    vertical-align: bottom;
  }
  .bid-title {
    font-size: 10pt;
    font-weight: normal;
    text-transform: uppercase;
    white-space: nowrap;
    padding-right: 8px;
    padding-bottom: 2px;
  }
  .bid-number-underline {
    border-bottom: 1px solid #000000;
    display: inline-block;
    width: 100%;
    padding-bottom: 2px;
    font-weight: bold;
    font-size: 10pt;
    text-align: left;
  }
  .dept-box-container {
    width: 75mm;
    border: 1.2pt solid #000000;
    text-align: center;
    margin-left: auto;
  }
  .dept-box-header {
    background-color: #d9d9d9;
    font-size: 9pt;
    font-weight: bold;
    padding: 3px 6px;
    border-bottom: 1.2pt solid #000000;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .dept-box-value {
    font-size: 10.5pt;
    font-weight: bold;
    padding: 6px 8px;
    min-height: 24px;
  }
  .signatories-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 3mm;
    margin-bottom: 3.5mm;
  }
  .signatories-table td {
    border: 1px solid #000000;
    padding: 4px 6px;
    font-size: 9pt;
    vertical-align: middle;
  }
  .sig-title {
    font-weight: bold;
    background-color: #f2f2f2;
    text-transform: uppercase;
    white-space: nowrap;
  }
  .sig-sign-col, .sig-phone-col {
    color: #333333;
    font-size: 8.5pt;
    text-align: center;
    white-space: nowrap;
  }
  .section-title {
    font-size: 10pt;
    font-weight: bold;
    margin-top: 2.5mm;
    margin-bottom: 1.5mm;
  }
  table.section-table, table.items-table {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
    margin-top: 1mm;
    margin-bottom: 2.5mm;
  }
  table.section-table th, table.items-table th {
    background-color: #d9d9d9;
    border: 1px solid #000000;
    padding: 4px 4px;
    font-size: 9pt;
    font-weight: bold;
    text-align: center;
    vertical-align: middle;
    word-break: break-word;
    overflow-wrap: anywhere;
  }
  table.section-table td, table.items-table td {
    border: 1px solid #000000;
    padding: 4px 5px;
    font-size: 9pt;
    vertical-align: middle;
    word-break: break-word;
    overflow-wrap: anywhere;
  }
  thead {
    display: table-header-group;
  }
  tr {
    page-break-inside: avoid;
  }
  .col-center {
    text-align: center;
  }
  .col-bold {
    font-weight: bold;
  }
  .sec1-col-num { width: 6%; text-align: center; }
  .sec1-col-name { width: 35%; font-weight: bold; }
  .sec1-col-desc { width: 59%; }
  .sec2-col-num { width: 6%; text-align: center; }
  .sec2-col-name { width: 32%; text-align: left; }
  .sec2-col-desc { width: 44%; text-align: left; }
  .sec2-col-unit { width: 9%; text-align: center; }
  .sec2-col-qty { width: 9%; text-align: center; }
</style>
</head>
`
