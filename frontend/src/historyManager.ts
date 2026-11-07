import type { DocumentTab, HistorySnapshot } from './types';

export function executeUndo(tab: DocumentTab, parseItems: (text: string) => any[]): void {
  if (tab.historyIndex <= 0) return;
  tab.historyIndex--;
  const state = tab.history[tab.historyIndex];
  tab.meta = JSON.parse(JSON.stringify(state.meta));
  tab.rawText = state.rawText;
  tab.items = parseItems(state.rawText);
  tab.isModified = true;
}

export function executeRedo(tab: DocumentTab, parseItems: (text: string) => any[]): void {
  if (tab.historyIndex >= tab.history.length - 1) return;
  tab.historyIndex++;
  const state = tab.history[tab.historyIndex];
  tab.meta = JSON.parse(JSON.stringify(state.meta));
  tab.rawText = state.rawText;
  tab.items = parseItems(state.rawText);
  tab.isModified = true;
}
