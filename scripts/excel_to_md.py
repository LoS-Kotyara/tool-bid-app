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

def get_sheet_data(ws):
    max_row = ws.max_row
    max_col = max(ws.max_column, 6)
    data = [[None] * max_col for _ in range(max_row)]
    
    for r in range(1, max_row + 1):
        for c in range(1, max_col + 1):
            data[r - 1][c - 1] = ws.cell(r, c).value
            
    for rng in ws.merged_cells.ranges:
        val = ws.cell(rng.min_row, rng.min_col).value
        for r in range(rng.min_row, rng.max_row + 1):
            for c in range(rng.min_col, min(rng.max_col, max_col) + 1):
                data[r - 1][c - 1] = val
                
    return data

def find_row_containing(data, keyword, max_row_limit=None):
    kw = keyword.lower()
    limit = max_row_limit or len(data)
    for i, row in enumerate(data[:limit]):
        for cell in row[:6]:
            if cell and kw in str(cell).lower():
                return i
    return -1

def get_person_info(data, keyword):
    idx = find_row_containing(data, keyword, max_row_limit=10)
    if idx != -1:
        row_slice = data[idx][:6]
        vals = []
        for cell in row_slice[2:]:
            if cell is not None and str(cell).strip():
                c_str = str(cell).strip()
                if keyword.lower() not in c_str.lower() and c_str not in vals:
                    vals.append(c_str)
        if len(vals) >= 2:
            return clean_signatory_part(vals[0]), clean_signatory_part(vals[1])
        elif len(vals) == 1:
            return clean_signatory_part(vals[0]), ""
    return "", ""

def get_section1_row_value(data, row_idx, keyword):
    if row_idx == -1 or row_idx >= len(data):
        return ""
    row_slice = data[row_idx][:6]
    for cell in row_slice[1:]:
        if cell is not None and str(cell).strip():
            c_str = str(cell).strip()
            if keyword and keyword.lower() in c_str.lower():
                continue
            if c_str.isdigit():
                continue
            if "служебная записка" in c_str.lower() or "срок поставки" in c_str.lower():
                continue
            if "требования к качеству" in c_str.lower() or "цель приобретения" in c_str.lower():
                continue
            return cell
    return ""

def parse_excel_bid(filepath):
    wb = openpyxl.load_workbook(filepath, data_only=True)
    ws = wb.active
    data = get_sheet_data(ws)

    bid_num = "1"
    bid_date = datetime.now().strftime("%d.%m.%Y")
    found_hdr = False

    for row in data[:4]:
        for cell in row[:6]:
            if cell and ("№" in str(cell) or "от" in str(cell)):
                m = re.search(r"№?\s*([A-Za-zА-Яа-я0-9\-_]+)\s+от\s+(\d{1,2}[./\-]\d{1,2}[./\-]\d{2,4})", str(cell))
                if m:
                    bid_num = m.group(1).lstrip("№").strip()
                    bid_date = normalize_date(m.group(2))
                    found_hdr = True
                    break
        if found_hdr:
            break

    dept = ""
    dept_idx = find_row_containing(data, "ПОДРАЗДЕЛЕНИЕ", max_row_limit=5)
    if dept_idx != -1:
        for cell in data[dept_idx][:6]:
            if cell and "подразделение" not in str(cell).lower() and str(cell).strip():
                dept = str(cell).strip()
                break
    if not dept and len(data) >= 3:
        val = data[2][3]
        if val and str(val).strip():
            dept = str(val).strip()

    init_role, init_fio = get_person_info(data, "ИНИЦИАТОР ЗАКУПКИ")
    cur_role, cur_fio = get_person_info(data, "КУРАТОР ЗАКУПКИ")

    shop_name = ""
    p_idx = find_row_containing(data, "Служебная записка", max_row_limit=17)
    if p_idx == -1:
        p_idx = find_row_containing(data, "Механосборочный", max_row_limit=17)
    if p_idx == -1:
        p_idx = find_row_containing(data, "Пункт", max_row_limit=17)
    if p_idx != -1:
        shop_name = str(get_section1_row_value(data, p_idx, "Служебная записка") or "").strip()
    elif len(data) >= 12:
        shop_name = str(data[11][3] or "").strip()

    delivery_date = ""
    d_idx = find_row_containing(data, "дата поставки", max_row_limit=17)
    if d_idx == -1:
        d_idx = find_row_containing(data, "Срок поставки", max_row_limit=17)
    if d_idx != -1:
        raw_delivery = get_section1_row_value(data, d_idx, "поставки")
        delivery_date = normalize_date(raw_delivery)
    elif len(data) >= 13:
        delivery_date = normalize_date(data[12][3])

    quality_req = ""
    q_idx = find_row_containing(data, "Требования к качеству", max_row_limit=17)
    if q_idx != -1:
        quality_req = str(get_section1_row_value(data, q_idx, "Требования к качеству") or "").strip()
    elif len(data) >= 14:
        quality_req = str(data[13][3] or "").strip()

    purpose = ""
    c_idx = find_row_containing(data, "Цель", max_row_limit=17)
    if c_idx != -1:
        purpose = str(get_section1_row_value(data, c_idx, "Цель") or "").strip()
    elif len(data) >= 15:
        purpose = str(data[14][3] or "").strip()

    start_row = 20
    for r in range(16, min(25, len(data))):
        val_a = str(data[r - 1][0] or "").lower()
        val_b = str(data[r - 1][1] or "").lower()
        if "наименование" in val_b or "наименов" in val_b or "№" in val_a:
            start_row = r + 1
            break

    items = []
    for r in range(start_row, len(data) + 1):
        num_v = data[r - 1][0]
        name_v = data[r - 1][1]
        desc_v = data[r - 1][2]
        unit_v = data[r - 1][4]
        qty_v = data[r - 1][5]

        if not name_v and not desc_v and not qty_v:
            continue

        if not num_v or not str(num_v).strip().isdigit():
            if not qty_v:
                continue

        name_str = str(name_v or "").strip()
        #name_str = re.sub(r";\s*", ", ", name_str)

        desc_str = str(desc_v or "").strip()
        #desc_str = re.sub(r";\s*", ", ", desc_str)

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

def convert_file(file_path, output_dir=None):
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
        
        base_name = os.path.splitext(os.path.basename(file_path))[0] + ".md"
        if output_dir:
            os.makedirs(output_dir, exist_ok=True)
            md_path = os.path.join(output_dir, base_name)
        else:
            md_path = os.path.splitext(file_path)[0] + ".md"

        with open(md_path, "w", encoding="utf-8") as f:
            f.write(md_text)
            
        print(f"[+] Успешно сохранено в Markdown ({len(data['items'])} поз.): {md_path}")
        return True
    except Exception as e:
        print(f"[ERROR] Ошибка конвертации {file_path}: {e}")
        return False

def main():
    if len(sys.argv) < 2:
        print("Использование: excel_to_md.py <файл1.xlsx> [<файл2.xlsx> ...] [--out <dir>]")
        print("Или перетащите Excel-файл(ы) на батник run_converter.bat.")
        input("\nНажмите Enter для выхода...")
        return

    output_dir = None
    args = []
    i = 1
    while i < len(sys.argv):
        if sys.argv[i] == "--out" and i + 1 < len(sys.argv):
            output_dir = sys.argv[i + 1]
            i += 2
        else:
            args.append(sys.argv[i])
            i += 1

    success_count = 0
    for arg in args:
        if convert_file(arg, output_dir=output_dir):
            success_count += 1

    print(f"\nЗавершено. Успешно сконвертировано: {success_count} файл(ов).")
    if sys.stdout.isatty():
        input("\nНажмите Enter для выхода...")

if __name__ == "__main__":
    main()

