<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { DocumentTab } from '../types';

  export let tabs: DocumentTab[] = [];
  export let activeTabId: string = '';

  const dispatch = createEventDispatcher<{
    selectTab: string;
    closeTab: string;
    newTab: void;
  }>();

  function formatTabName(tab: DocumentTab): string {
    const rawNum = tab.meta.bidNum.trim();
    if (!rawNum) {
      return 'Новая';
    }
    if (rawNum.length <= 5) {
      return rawNum;
    }
    return rawNum.slice(0, 4) + '…';
  }
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <span class="sidebar-title">Заявки ({tabs.length})</span>
    <button class="add-tab-btn" title="Новая заявка (Ctrl+N)" on:click={() => dispatch('newTab')}>
      +
    </button>
  </div>

  <div class="tabs-list">
    {#each tabs as tab (tab.id)}
      <div
        class="tab-item"
        class:active={tab.id === activeTabId}
        on:click={() => dispatch('selectTab', tab.id)}
        role="button"
        tabindex="0"
        on:keydown={(e) => e.key === 'Enter' && dispatch('selectTab', tab.id)}
      >
        <span class="tab-indicator" class:modified={tab.isModified} title={tab.isModified ? 'Несохраненные изменения' : 'Сохранено'}>
          {tab.isModified ? '○' : ''}
        </span>

        <span class="tab-label" title={tab.meta.bidNum ? `Заявка №${tab.meta.bidNum}` : 'Новая заявка'}>
          {formatTabName(tab)}
        </span>

        <button
          class="close-tab-btn"
          title="Закрыть вкладку"
          on:click|stopPropagation={() => dispatch('closeTab', tab.id)}
        >
          ×
        </button>
      </div>
    {/each}
  </div>
</aside>

<style>
  .sidebar {
    width: 140px;
    min-width: 120px;
    background-color: var(--ctp-mantle);
    border-right: 1px solid var(--ctp-surface0);
    display: flex;
    flex-direction: column;
    height: 100%;
    user-select: none;
  }

  .sidebar-header {
    height: 36px;
    padding: 0 10px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid var(--ctp-surface0);
  }

  .sidebar-title {
    font-size: 9pt;
    text-transform: uppercase;
    color: var(--ctp-subtext0);
    letter-spacing: 0.5px;
  }

  .add-tab-btn {
    background: transparent;
    border: none;
    color: var(--ctp-text);
    font-size: 14pt;
    line-height: 1;
    width: 24px;
    height: 24px;
    border-radius: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .add-tab-btn:hover {
    background-color: var(--ctp-surface0);
    color: var(--ctp-blue);
  }

  .tabs-list {
    flex: 1;
    overflow-y: auto;
    padding: 6px;
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  .tab-item {
    display: flex;
    align-items: center;
    padding: 6px 8px;
    border-radius: 6px;
    background-color: transparent;
    cursor: pointer;
    font-size: 10.5pt;
    color: var(--ctp-subtext1);
    transition: all 0.15s ease;
    outline: none;
  }

  .tab-item:hover {
    background-color: var(--ctp-surface0);
    color: var(--ctp-text);
  }

  .tab-item.active {
    background-color: var(--ctp-surface1);
    color: var(--ctp-text);
    font-weight: 500;
  }

  .tab-indicator {
    width: 14px;
    font-size: 11pt;
    text-align: center;
    margin-right: 4px;
    color: var(--ctp-peach);
  }

  .tab-indicator.modified {
    color: #ffffff;
    font-weight: bold;
  }

  .tab-label {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .close-tab-btn {
    background: transparent;
    border: none;
    color: var(--ctp-overlay1);
    font-size: 13pt;
    line-height: 1;
    width: 18px;
    height: 18px;
    border-radius: 3px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0.6;
    transition: opacity 0.15s;
  }

  .tab-item:hover .close-tab-btn {
    opacity: 1;
  }

  .close-tab-btn:hover {
    background-color: var(--ctp-surface2);
    color: var(--ctp-red);
  }
</style>
