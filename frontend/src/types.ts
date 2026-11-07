export interface DocumentTab {
  id: string;
  filePath: string;
  isModified: boolean;
  rawText: string;
  meta: MetaInfoState;
  items: ItemState[];
  // История Undo/Redo
  history: HistorySnapshot[];
  historyIndex: number;
}

export interface MetaInfoState {
  bidNum: string;
  bidDate: string;
  dept: string;
  shopName: string;
  initRole: string;
  initFio: string;
  curRole: string;
  curFio: string;
  deliveryDate: string;
  qualityReq: string;
  purpose: string;
}

export interface ItemState {
  num: number;
  name: string;
  desc: string;
  unit: string;
  qty: string;
}

export interface HistorySnapshot {
  meta: MetaInfoState;
  rawText: string;
}

export type ViewMode = 'split' | 'edit' | 'preview';
export type CatppuccinTheme = 'mocha' | 'macchiato' | 'frappe' | 'latte';

export function getEmptyMeta(): MetaInfoState {
  const now = new Date();
  const d = String(now.getDate()).padStart(2, '0');
  const m = String(now.getMonth() + 1).padStart(2, '0');
  const y = now.getFullYear();
  const todayStr = `${d}.${m}.${y}`;

  return {
    bidNum: '',
    bidDate: todayStr,
    dept: '',
    shopName: '',
    initRole: '',
    initFio: '',
    curRole: '',
    curFio: '',
    deliveryDate: todayStr,
    qualityReq: '',
    purpose: '',
  };
}

export function parseItemsFromText(text: string): ItemState[] {
  const lines = text.split('\n');
  const items: ItemState[] = [];
  let itemNum = 1;

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed) continue;

    const parts = trimmed.split(';');
    if (parts.length >= 4) {
      const name = parts[0].trim();
      const qty = parts[parts.length - 1].trim();
      const unit = parts[parts.length - 2].trim();
      const desc = parts.slice(1, parts.length - 2).join(';').trim();

      if (name) {
        items.push({
          num: itemNum++,
          name,
          desc,
          unit,
          qty,
        });
      }
    }
  }

  return items;
}

export function itemsToText(items: ItemState[]): string {
  return items.map(it => `${it.name}; ${it.desc}; ${it.unit}; ${it.qty}`).join('\n');
}
