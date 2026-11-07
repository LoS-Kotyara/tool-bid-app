import type { DocumentTab, HistorySnapshot } from './types';
import { getEmptyMeta, itemsToText, parseItemsFromText } from './types';
import { model } from '../wailsjs/go/models';

export function createNewTab(initialDoc?: model.Document): DocumentTab {
  const id = 'tab_' + Date.now() + '_' + Math.random().toString(36).substring(2, 6);
  const meta = initialDoc ? {
    bidNum: initialDoc.meta?.bidNum || '',
    bidDate: initialDoc.meta?.bidDate || '',
    dept: initialDoc.meta?.dept || '',
    shopName: initialDoc.meta?.shopName || '',
    initRole: initialDoc.meta?.initRole || '',
    initFio: initialDoc.meta?.initFio || '',
    curRole: initialDoc.meta?.curRole || '',
    curFio: initialDoc.meta?.curFio || '',
    deliveryDate: initialDoc.meta?.deliveryDate || '',
    qualityReq: initialDoc.meta?.qualityReq || '',
    purpose: initialDoc.meta?.purpose || '',
  } : getEmptyMeta();

  const items = initialDoc?.items ? initialDoc.items.map(it => ({
    num: it.num,
    name: it.name,
    desc: it.desc,
    unit: it.unit,
    qty: it.qty,
  })) : [];

  const rawText = itemsToText(items);
  const initialSnapshot: HistorySnapshot = {
    meta: JSON.parse(JSON.stringify(meta)),
    rawText,
  };

  return {
    id,
    filePath: initialDoc?.filePath || '',
    isModified: false,
    rawText,
    meta,
    items,
    history: [initialSnapshot],
    historyIndex: 0,
  };
}

export function pushTabHistory(tab: DocumentTab) {
  tab.isModified = true;
  const snapshot: HistorySnapshot = {
    meta: JSON.parse(JSON.stringify(tab.meta)),
    rawText: tab.rawText,
  };

  tab.history = tab.history.slice(0, tab.historyIndex + 1);
  tab.history.push(snapshot);
  if (tab.history.length > 50) {
    tab.history.shift();
  }
  tab.historyIndex = tab.history.length - 1;
}

export function buildDocumentModel(tab: DocumentTab): model.Document {
  const items = parseItemsFromText(tab.rawText).map(it => new model.Item({
    num: it.num,
    name: it.name,
    desc: it.desc,
    unit: it.unit,
    qty: it.qty,
  }));

  return new model.Document({
    id: tab.id,
    filePath: tab.filePath,
    isModified: tab.isModified,
    meta: new model.MetaInfo(tab.meta),
    items: items,
  });
}
