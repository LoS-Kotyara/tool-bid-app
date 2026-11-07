# Office Export Engine (Go XLSX, ODS & PDF)

Специализируется на низкоуровневой генерации файлов электронных таблиц и документов печати в Go Core.

## Зоны ответственности:
1. Генерация XLSX (через github.com/xuri/excelize/v2):
   - Точные ширины колонок (A: 4.14, B: 29.57, C: 20.85, D: 21.57, E: 8.43, F: 9.29)[cite: 1].
   - Матрица высот строк: R1: 15.0, R2: 5.25, R3: 15.0, R4: 4.5, R5: 27.0, R6: 35.45, R7: 33.6, R8: 4.5, R9: 16.5, R10: 4.5, R11: 32.25, R12: 54.75, R13: 15.0, R14: 33.0, R15: 54.0, R16: 3.0, R17: 16.5, R18: 4.5, R19: 36.0, R20+: 30.0[cite: 1].
   - Объединения: C1:F1 (подчеркивание), D3:F3 (контур подразделения medium), A6:B6, A7:B7, A9:F9, B11:C11, D11:F11, C19:D19[cite: 1].
   - Печать: Portrait, FitToWidth=1, сквозной заголовок строк 19:19[cite: 1].
2. Автономная генерация ODS (OpenDocument Spreadsheet):
   - Использование исключительно стандартных пакетов Go `archive/zip` и `encoding/xml`[cite: 1].
   - Структура пакета: mimetype (STORED / без сжатия, application/vnd.oasis.opendocument.spreadsheet), META-INF/manifest.xml, styles.xml, content.xml[cite: 1].
   - Метрики в сантиметрах: A: 1.0cm, B: 6.0cm, C: 4.3cm, D: 4.3cm, E: 1.7cm, F: 1.9cm[cite: 1].
   - Полная валидность и совместимость с Apache OpenOffice Calc и LibreOffice Calc[cite: 1].
3. PDF Рендеринг:
   - A4 печатный HTML/CSS шаблон с полями 10 мм, `thead { display: table-header-group; }` для авто-переноса шапки[cite: 1].