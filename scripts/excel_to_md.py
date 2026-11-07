import sys
import os
import re
from datetime import datetime, date
import openpyxl

def normalize_date(val):
    if val is None:
        return ""
    if isinstance(val, (datetime, date)):
        return val.strftime("%d.%m.%Y")
    
    s = str(val).strip()
    if not s:
        return ""
    
    m1 = re.match(r"^(\d{1,2})[./\-](\d{1,2})[./\-](\d{4})", s)
    if m1:
        d, m, y = int(m1.group(1)), int(m1.group(2)), int(m1.group(3))
        return f"{d:02d}.{m:02d}.{y:04d}"
    
    m2 = re.match(r"^(\d{4})[./\-](\d{1,2})[./\-](\d{1,2})", s)
    if m2:
        y, m, d = int(m2.group(1)), int(m2.group(2)), int(m2.group(3))
        return f"{d:02d}.{m:02d}.{y:04d}"
    
    m3 = re.match(r"^(\d{1,2})[./\-](\d{1,2})[./\-](\d{2})$", s)
    if m3:
        d, m, y = int(m3.group(1)), int(m3.group(2)), int(m3.group(3))
        y += 2000 if y < 50 else 1900
        return f"{d:02d}.{m:02d}.{y:04d}"

    return s

def clean_signatory_part(s):
    if not s:
        return ""
    s = str(s).strip()
    s = re.sub(r"[;\s]+$", "", s)
    return s.strip()

def get_cell_val(ws, row, col):
    cell = ws.cell(row=row, column=col)
    if cell.coordinate in ws.merged_cells:
        for rng in ws.merged_cells.ranges:
            if cell.coordinate in rng:
                top_left = ws.cell(row=rng.min_row, column=rng.min_col)
                return top_left.value
    return cell.value

def parse_excel_bid(filepath):
    wb = openpyxl.load_workbook(filepath, data_only=True)
    ws = wb.active

    header_val = ""
    for r in range(1, 4):
        for c in range(2, 6):
            v = get_cell_val(ws, r, c)
            if v and any(k in str(v).lower() for k in ["от", "заявка", "№"]):
                header_val = str(v).strip()
                break
        if header_val:
            break

    bid_num = "1"
    bid_date = datetime.now().strftime("%d.%m.%Y")
    
    m_num = re.search(r"№?\s*([A-Za-zА-Яа-я0-9\-_]+)\s+от\s+(\d{1,2}[./\-]\d{1,2}[./\-]\d{2,4})", header_val)
    if m_num:
        bid_num = m_num.group(1).lstrip("№").strip()
        bid_date = normalize_date(m_num.group(2))
    else:
        c1 = get_cell_val(ws, 1, 3)
        if c1:
            m_c1 = re.search(r"([A-Za-zА-Яа-я0-9\-_]+)\s+от\s+(\d{1,2}[./\-]\d{1,2}[./\-]\d{2,4})", str(c1))
            if m_c1:
                bid_num = m_c1.group(1).lstrip("№").strip()
                bid_date = normalize_date(m_c1.group(2))

    dept = ""
    v_dept = get_cell_val(ws, 3, 4)
    if v_dept:
        dept = str(v_dept).strip()

    init_role = clean_signatory_part(get_cell_val(ws, 6, 3))
    init_fio = clean_signatory_part(get_cell_val(ws, 6, 5))
    cur_role = clean_signatory_part(get_cell_val(ws, 7, 3))
    cur_fio = clean_signatory_part(get_cell_val(ws, 7, 5))

    shop_name = str(get_cell_val(ws, 12, 4) or "").strip()
    raw_delivery = get_cell_val(ws, 13, 4)
    delivery_date = normalize_date(raw_delivery)
    quality_req = str(get_cell_val(ws, 14, 4) or "").strip()
    purpose = str(get_cell_val(ws, 15, 4) or "").strip()

    start_row = 20
    for r in range(16, 25):
        val_a = str(get_cell_val(ws, r, 1) or "").lower()
        val_b = str(get_cell_val(ws, r, 2) or "").lower()
        if "наименование" in val_b or "№" in val_a:
            start_row = r + 1
            break

    items = []
    for r in range(start_row, ws.max_row + 1):
        name_v = get_cell_val(ws, r, 2)
        desc_v = get_cell_val(ws, r, 3)
        unit_v = get_cell_val(ws, r, 5)
        qty_v = get_cell_val(ws, r, 6)

        if not name_v and not desc_v and not qty_v:
            continue

        name_str = str(name_v or "").strip().replace(";", ",")
        desc_str = str(desc_v or "").strip().replace(";", ",")
        unit_str = str(unit_v or "шт.").strip().replace(";", "")
        
        qty_str = ""
        if qty_v is not None:
            if isinstance(qty_v, float) and qty_v.is_integer():
                qty_str = str(int(qty_v))
            else:
                qty_str = str(qty_v).strip().replace(";", "")
        
        if name_str or desc_str:
            items.append((name_str, desc_str, unit_str, qty_str))

    return {
        "bid_num": bid_num,
        "bid_date": bid_date,
        "dept": dept,
        "shop_name": shop_name,
        "init_role": init_role,
        "init_fio": init_fio,
        "cur_role": cur_role,
        "cur_fio": cur_fio,
        "delivery_date": delivery_date,
        "quality_req": quality_req,
        "purpose": purpose,
        "items": items,
    }


def format_as_canonical_md(data):
    lines = []
    lines.append(f"#{data['bid_num']} от {data['bid_date']}")
    
    if data["dept"]:
        lines.append(f"Подразделение: {data['dept']}")
    if data["shop_name"]:
        lines.append(f"Пункт: {data['shop_name']}")
    
    init_parts = []
    if data["init_role"]:
        init_parts.append(data["init_role"])
    if data["init_fio"]:
        init_parts.append(data["init_fio"])
    if init_parts:
        lines.append(f"Инициатор: {'; '.join(init_parts)}")
    
    cur_parts = []
    if data["cur_role"]:
        cur_parts.append(data["cur_role"])
    if data["cur_fio"]:
        cur_parts.append(data["cur_fio"])
    if cur_parts:
        lines.append(f"Куратор: {'; '.join(cur_parts)}")
    
    if data["delivery_date"]:
        lines.append(f"Дата поставки: {data['delivery_date']}")
    if data["quality_req"]:
        lines.append(f"Требования: {data['quality_req']}")
    if data["purpose"]:
        lines.append(f"Цель: {data['purpose']}")
    
    lines.append("")
    
    for item in data["items"]:
        name, desc, unit, qty = item
        lines.append(f"{name}; {desc}; {unit}; {qty}")
        
    return "\n".join(lines) + "\n"

def convert_file(file_path):
    if not os.path.isfile(file_path):
        print(f"[!] Файл не найден: {file_path}")
        return False
    
    ext = os.path.splitext(file_path)[1].lower()
    if ext not in [".xlsx", ".xlsm"]:
        print(f"[!] Пропуск не-Excel файла: {file_path}")
        return False

    print(f"[*] Обработка: {file_path} ...")
    try:
        data = parse_excel_bid(file_path)
        md_text = format_as_canonical_md(data)
        
        md_path = os.path.splitext(file_path)[0] + ".md"
        with open(md_path, "w", encoding="utf-8") as f:
            f.write(md_text)
            
        print(f"[+] Успешно сохранено в Markdown: {md_path}")
        return True
    except Exception as e:
        print(f"[ERROR] Ошибка конвертации {file_path}: {e}")
        return False

def main():
    if len(sys.argv) < 2:
        print("Использование: excel_to_md.py <файл1.xlsx> [<файл2.xlsx> ...]")
        print("Или перетащите Excel-файл(ы) на батник run_converter.bat.")
        input("\nНажмите Enter для выхода...")
        return

    success_count = 0
    for arg in sys.argv[1:]:
        if convert_file(arg):
            success_count += 1

    print(f"\nЗавершено. Успешно сконвертировано: {success_count} файл(ов).")
    if sys.stdout.isatty():
        input("\nНажмите Enter для выхода...")

if __name__ == "__main__":
    main()

