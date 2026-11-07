import {
  OpenMarkdownDialog,
  SaveMarkdownFile,
  ReloadMarkdownFile,
  ExportXLSXDialog,
  ExportODSDialog,
  GetDocumentPreviewHTML
} from '../wailsjs/go/app/App';
import type { DocumentTab } from './types';
import { buildDocumentModel, createNewTab } from './tabManager';

export async function openMarkdownFile(): Promise<DocumentTab | null> {
  const doc = await OpenMarkdownDialog();
  if (!doc) return null;
  return createNewTab(doc);
}

export async function saveMarkdown(tab: DocumentTab, silent: boolean = false): Promise<string> {
  const doc = buildDocumentModel(tab);
  const targetPath = silent ? tab.filePath : '';
  return await SaveMarkdownFile(targetPath, doc);
}

export async function reloadMarkdown(path: string): Promise<DocumentTab> {
  const doc = await ReloadMarkdownFile(path);
  return createNewTab(doc);
}

export async function exportToXlsx(tab: DocumentTab): Promise<void> {
  await ExportXLSXDialog(buildDocumentModel(tab));
}

export async function exportToOds(tab: DocumentTab): Promise<void> {
  await ExportODSDialog(buildDocumentModel(tab));
}

export async function getPreviewHTML(tab: DocumentTab): Promise<string> {
  return await GetDocumentPreviewHTML(buildDocumentModel(tab));
}
