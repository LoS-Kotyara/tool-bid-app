<script lang="ts">
  import { onMount } from 'svelte';
  import './appLayout.css';
  import Header from './components/Header.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import Workspace from './components/Workspace.svelte';
  import SettingsModal from './components/SettingsModal.svelte';
  import ConfirmModal from './components/ConfirmModal.svelte';

  import type { DocumentTab, ViewMode, CatppuccinTheme } from './types';
  import { parseItemsFromText } from './types';
  import { createNewTab, pushTabHistory } from './tabManager';
  import { executeUndo, executeRedo } from './historyManager';
  import { setupHotkeys, subscribeFileDrop } from './appEvents';
  import {
    openMarkdownFile, saveMarkdown, reloadMarkdown,
    exportToXlsx, exportToOds, getPreviewHTML
  } from './fileOperations';
  import type { model } from '../wailsjs/go/models';

  let tabs: DocumentTab[] = [];
  let activeTabId = '';
  let viewMode: ViewMode = 'split';
  let sidebarOpen = true;
  let currentTheme: CatppuccinTheme = 'mocha';
  let currentFont = 'Cascadia Code';
  let settingsOpen = false;
  let previewHtml = '';

  let confirmDialog = {
    open: false,
    title: '',
    message: '',
    danger: false,
    confirmAction: () => {},
  };

  $: activeTab = tabs.find(t => t.id === activeTabId);

  let previewTimer: any;
  $: if (activeTab) {
    clearTimeout(previewTimer);
    previewTimer = setTimeout(async () => {
      if (activeTab) {
        try { previewHtml = await getPreviewHTML(activeTab); } catch (e) { console.error(e); }
      }
    }, 200);
  }

  function addNewTab(doc?: model.Document) {
    const newTab = createNewTab(doc);
    tabs = [...tabs, newTab];
    activeTabId = newTab.id;
  }

  function handleCloseTab(id: string) {
    const tab = tabs.find(t => t.id === id);
    if (!tab) return;
    if (tab.isModified) {
      confirmDialog = {
        open: true,
        title: 'Несохраненные изменения',
        message: `В заявке №${tab.meta.bidNum || 'Новая'} есть изменения. Закрыть без сохранения?`,
        danger: true,
        confirmAction: () => closeTabDirectly(id)
      };
    } else {
      closeTabDirectly(id);
    }
  }

  function closeTabDirectly(id: string) {
    const idx = tabs.findIndex(t => t.id === id);
    tabs = tabs.filter(t => t.id !== id);
    if (activeTabId === id) {
      activeTabId = tabs.length > 0 ? tabs[Math.max(0, idx - 1)].id : '';
      if (!activeTabId) addNewTab();
    }
  }

  async function handleOpen() {
    try {
      const tab = await openMarkdownFile();
      if (!tab) return;
      tabs = [...tabs, tab];
      activeTabId = tab.id;
    } catch (e) { alert(`Ошибка открытия: ${e}`); }
  }

  async function handleSave(silent = false) {
    if (!activeTab) return;
    try {
      const path = await saveMarkdown(activeTab, silent);
      if (path) { activeTab.filePath = path; activeTab.isModified = false; tabs = tabs; }
    } catch (e) { alert(`Ошибка сохранения: ${e}`); }
  }

  async function handleReload() {
    if (!activeTab || !activeTab.filePath) return;
    if (activeTab.isModified) {
      confirmDialog = {
        open: true,
        title: 'Перезагрузка файла с диска',
        message: 'Файл изменен. Перезагрузить с диска и потерять изменения?',
        danger: true,
        confirmAction: () => reloadDirectly()
      };
    } else {
      await reloadDirectly();
    }
  }

  async function reloadDirectly() {
    if (!activeTab || !activeTab.filePath) return;
    try {
      const tab = await reloadMarkdown(activeTab.filePath);
      const idx = tabs.findIndex(t => t.id === activeTabId);
      if (idx !== -1) { tab.id = activeTabId; tabs[idx] = tab; tabs = tabs; }
    } catch (e) { alert(`Ошибка перезагрузки: ${e}`); }
  }

  const handleGlobalKeyDown = setupHotkeys({
    onSave: () => handleSave(true),
    onReload: handleReload,
    onNewTab: () => addNewTab(),
    onOpen: handleOpen,
    onUndo: () => { if (activeTab) { executeUndo(activeTab, parseItemsFromText); tabs = tabs; } },
    onRedo: () => { if (activeTab) { executeRedo(activeTab, parseItemsFromText); tabs = tabs; } }
  });

  onMount(() => {
    const th = localStorage.getItem('tool-bid-theme') as CatppuccinTheme;
    if (th) currentTheme = th;
    const fn = localStorage.getItem('tool-bid-font');
    if (fn) currentFont = fn;
    document.documentElement.setAttribute('data-theme', currentTheme);
    addNewTab();
    subscribeFileDrop((doc) => addNewTab(doc), (err) => alert(err));
  });

  function applySettings(e: CustomEvent<{ theme: CatppuccinTheme; font: string }>) {
    currentTheme = e.detail.theme;
    currentFont = e.detail.font;
    localStorage.setItem('tool-bid-theme', currentTheme);
    localStorage.setItem('tool-bid-font', currentFont);
    document.documentElement.setAttribute('data-theme', currentTheme);
    settingsOpen = false;
  }
</script>

<svelte:window on:keydown={handleGlobalKeyDown} />

<div class="app-layout" data-theme={currentTheme}>
  <Header
    {viewMode}
    {sidebarOpen}
    on:newTab={() => addNewTab()}
    on:open={handleOpen}
    on:save={() => handleSave(true)}
    on:saveAs={() => handleSave(false)}
    on:exportXlsx={() => activeTab && exportToXlsx(activeTab)}
    on:exportOds={() => activeTab && exportToOds(activeTab)}
    on:openSettings={() => settingsOpen = true}
    on:toggleSidebar={() => sidebarOpen = !sidebarOpen}
    on:changeViewMode={(e) => viewMode = e.detail}
  />

  <main class="main-body">
    {#if sidebarOpen}
      <Sidebar
        {tabs}
        {activeTabId}
        on:selectTab={(e) => activeTabId = e.detail}
        on:closeTab={(e) => handleCloseTab(e.detail)}
        on:newTab={() => addNewTab()}
      />
    {/if}

    {#if activeTab}
      <Workspace
        {activeTab}
        {viewMode}
        {currentFont}
        {previewHtml}
        on:metaChange={() => { pushTabHistory(activeTab); tabs = tabs; }}
        on:textChange={(e) => {
          activeTab.rawText = e.detail;
          activeTab.items = parseItemsFromText(e.detail);
          pushTabHistory(activeTab);
          tabs = tabs;
        }}
        on:undo={() => { executeUndo(activeTab, parseItemsFromText); tabs = tabs; }}
        on:redo={() => { executeRedo(activeTab, parseItemsFromText); tabs = tabs; }}
      />
    {/if}
  </main>

  {#if settingsOpen}
    <SettingsModal
      {currentTheme}
      {currentFont}
      on:close={() => settingsOpen = false}
      on:save={applySettings}
    />
  {/if}

  {#if confirmDialog.open}
    <ConfirmModal
      title={confirmDialog.title}
      message={confirmDialog.message}
      danger={confirmDialog.danger}
      on:cancel={() => confirmDialog.open = false}
      on:confirm={() => { confirmDialog.open = false; confirmDialog.confirmAction(); }}
    />
  {/if}
</div>
