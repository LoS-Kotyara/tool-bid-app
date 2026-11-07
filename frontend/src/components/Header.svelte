<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { ViewMode } from '../types';

  export let viewMode: ViewMode = 'split';
  export let sidebarOpen: boolean = true;

  const dispatch = createEventDispatcher<{
    newTab: void;
    open: void;
    save: void;
    saveAs: void;
    exportXlsx: void;
    exportOds: void;
    openSettings: void;
    toggleSidebar: void;
    changeViewMode: ViewMode;
  }>();

  let fileMenuOpen = false;

  function toggleFileMenu(e: MouseEvent) {
    e.stopPropagation();
    fileMenuOpen = !fileMenuOpen;
  }

  function closeMenu() {
    fileMenuOpen = false;
  }
</script>

<svelte:window on:click={closeMenu} />

<header class="header">
  <div class="left-section">
    <button
      class="sidebar-toggle-btn"
      title={sidebarOpen ? "Свернуть панель вкладок" : "Развернуть панель вкладок"}
      on:click={() => dispatch('toggleSidebar')}
    >
      {sidebarOpen ? '◀ Вкладки' : '▶ Вкладки'}
    </button>

    <div class="menu-container">
      <button class="menu-btn" on:click={toggleFileMenu}>
        Файл ▾
      </button>
      {#if fileMenuOpen}
        <div class="dropdown-menu">
          <button class="dropdown-item" on:click={() => { closeMenu(); dispatch('newTab'); }}>
            <span class="shortcut">Ctrl+N</span> Создать новую заявку
          </button>
          <button class="dropdown-item" on:click={() => { closeMenu(); dispatch('open'); }}>
            <span class="shortcut">Ctrl+O</span> Открыть Markdown...
          </button>
          <div class="divider"></div>
          <button class="dropdown-item" on:click={() => { closeMenu(); dispatch('save'); }}>
            <span class="shortcut">Ctrl+S</span> Сохранить
          </button>
          <button class="dropdown-item" on:click={() => { closeMenu(); dispatch('saveAs'); }}>
            Сохранить как...
          </button>
          <div class="divider"></div>
          <button class="dropdown-item" on:click={() => { closeMenu(); dispatch('exportXlsx'); }}>
            Экспорт в Excel (.xlsx)
          </button>
          <button class="dropdown-item" on:click={() => { closeMenu(); dispatch('exportOds'); }}>
            Экспорт в LibreOffice (.ods)
          </button>
        </div>
      {/if}
    </div>
  </div>

  <div class="center-section">
    <div class="mode-switch">
      <button
        class="mode-btn"
        class:active={viewMode === 'split'}
        on:click={() => dispatch('changeViewMode', 'split')}
        title="Сплит-режим 50/50"
      >
        ⚖ Сплит
      </button>
      <button
        class="mode-btn"
        class:active={viewMode === 'edit'}
        on:click={() => dispatch('changeViewMode', 'edit')}
        title="100% редактор"
      >
        📝 Ввод
      </button>
      <button
        class="mode-btn"
        class:active={viewMode === 'preview'}
        on:click={() => dispatch('changeViewMode', 'preview')}
        title="100% превью"
      >
        🔍 Превью
      </button>
    </div>
  </div>

  <div class="right-section">
    <button
      class="settings-btn"
      title="Настройки темы и шрифта"
      on:click={() => dispatch('openSettings')}
    >
      ⚙ Настройки
    </button>
  </div>
</header>

<style>
  .header {
    height: 42px;
    background-color: var(--ctp-mantle);
    border-bottom: 1px solid var(--ctp-surface0);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 12px;
    z-index: 100;
  }

  .left-section, .right-section, .center-section {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .sidebar-toggle-btn, .menu-btn, .settings-btn {
    background: transparent;
    border: 1px solid transparent;
    color: var(--ctp-text);
    padding: 4px 10px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 11pt;
    transition: background-color 0.15s, border-color 0.15s;
  }

  .sidebar-toggle-btn:hover, .menu-btn:hover, .settings-btn:hover {
    background-color: var(--ctp-surface0);
    border-color: var(--ctp-surface1);
  }

  .menu-container {
    position: relative;
  }

  .dropdown-menu {
    position: absolute;
    top: 32px;
    left: 0;
    min-width: 220px;
    background-color: var(--ctp-surface0);
    border: 1px solid var(--ctp-surface1);
    border-radius: 6px;
    box-shadow: 0 8px 18px rgba(0, 0, 0, 0.4);
    padding: 4px 0;
    z-index: 1000;
  }

  .dropdown-item {
    width: 100%;
    text-align: left;
    background: transparent;
    border: none;
    color: var(--ctp-text);
    padding: 7px 12px;
    cursor: pointer;
    font-size: 10.5pt;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .dropdown-item:hover {
    background-color: var(--ctp-surface1);
    color: var(--ctp-blue);
  }

  .shortcut {
    font-size: 9pt;
    color: var(--ctp-subtext0);
    margin-right: 12px;
  }

  .divider {
    height: 1px;
    background-color: var(--ctp-surface1);
    margin: 4px 0;
  }

  .mode-switch {
    display: flex;
    background-color: var(--ctp-crust);
    border: 1px solid var(--ctp-surface0);
    border-radius: 6px;
    padding: 2px;
    gap: 2px;
  }

  .mode-btn {
    background: transparent;
    border: none;
    color: var(--ctp-subtext0);
    padding: 3px 12px;
    border-radius: 4px;
    font-size: 10.5pt;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .mode-btn:hover {
    color: var(--ctp-text);
  }

  .mode-btn.active {
    background-color: var(--ctp-surface0);
    color: var(--ctp-blue);
    font-weight: 500;
  }
</style>
