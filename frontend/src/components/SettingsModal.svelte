<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import type { CatppuccinTheme } from '../types';
  import { GetInstalledFonts } from '../../wailsjs/go/app/App';

  export let currentTheme: CatppuccinTheme = 'mocha';
  export let currentFont: string = 'Cascadia Code';

  const dispatch = createEventDispatcher<{
    close: void;
    save: { theme: CatppuccinTheme; font: string };
  }>();

  let selectedTheme: CatppuccinTheme = currentTheme;
  let selectedFont: string = currentFont;
  let installedFonts: string[] = ['Cascadia Code', 'Consolas', 'Courier New'];
  let loadingFonts = true;

  onMount(async () => {
    try {
      const fonts = await GetInstalledFonts();
      if (fonts && fonts.length > 0) {
        installedFonts = fonts;
      }
    } catch (e) {
      console.error('Failed to load system fonts:', e);
    } finally {
      loadingFonts = false;
    }
  });

  function saveSettings() {
    dispatch('save', { theme: selectedTheme, font: selectedFont });
  }

  function closeDialog() {
    dispatch('close');
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      closeDialog();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="modal-backdrop" on:click={closeDialog}>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="modal-window" on:click|stopPropagation>
    <div class="modal-header">
      <span class="modal-title">⚙ Настройки интерфейса</span>
      <button class="close-btn" on:click={closeDialog}>×</button>
    </div>

    <div class="modal-body">
      <!-- Выбор темы Catppuccin -->
      <div class="setting-item">
        <label for="theme-select">Цветовая тема (Catppuccin):</label>
        <select id="theme-select" bind:value={selectedTheme}>
          <option value="mocha">Catppuccin Mocha (Темная по умолчанию)</option>
          <option value="macchiato">Catppuccin Macchiato (Темная)</option>
          <option value="frappe">Catppuccin Frappé (Темная)</option>
          <option value="latte">Catppuccin Latte (Светлая)</option>
        </select>
      </div>

      <!-- Выбор системного шрифта -->
      <div class="setting-item">
        <label for="font-select">
          Шрифт редактора:
          {#if loadingFonts}
            <span class="loading-label">(чтение реестра...)</span>
          {/if}
        </label>
        <select id="font-select" bind:value={selectedFont}>
          {#each installedFonts as font}
            <option value={font}>{font}</option>
          {/each}
        </select>
        <div class="font-preview" style="font-family: '{selectedFont}', monospace;">
          Образец: Фреза концевая КЭМЗ; D12; шт.; 4
        </div>
      </div>
    </div>

    <div class="modal-footer">
      <button class="btn btn-secondary" on:click={closeDialog}>Отмена</button>
      <button class="btn btn-primary" on:click={saveSettings}>Применить</button>
    </div>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
  }

  .modal-window {
    background-color: var(--ctp-surface0);
    border: 1px solid var(--ctp-surface1);
    border-radius: 8px;
    width: 440px;
    box-shadow: 0 16px 36px rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .modal-header {
    height: 44px;
    padding: 0 16px;
    background-color: var(--ctp-mantle);
    border-bottom: 1px solid var(--ctp-surface1);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .modal-title {
    font-size: 11pt;
    font-weight: 600;
    color: var(--ctp-text);
  }

  .close-btn {
    background: transparent;
    border: none;
    color: var(--ctp-subtext0);
    font-size: 16pt;
    cursor: pointer;
    line-height: 1;
  }

  .close-btn:hover {
    color: var(--ctp-red);
  }

  .modal-body {
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .setting-item {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  label {
    font-size: 9pt;
    color: var(--ctp-subtext1);
    font-weight: 500;
  }

  .loading-label {
    font-size: 8pt;
    color: var(--ctp-overlay1);
  }

  select {
    height: 32px;
    background-color: var(--ctp-base);
    border: 1px solid var(--ctp-surface1);
    border-radius: 4px;
    color: var(--ctp-text);
    padding: 0 8px;
    outline: none;
    font-size: 10pt;
  }

  select:focus {
    border-color: var(--ctp-blue);
  }

  .font-preview {
    margin-top: 4px;
    padding: 8px 10px;
    background-color: var(--ctp-base);
    border: 1px solid var(--ctp-surface1);
    border-radius: 4px;
    font-size: 10pt;
    color: var(--ctp-subtext0);
  }

  .modal-footer {
    padding: 12px 16px;
    background-color: var(--ctp-mantle);
    border-top: 1px solid var(--ctp-surface1);
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }

  .btn {
    padding: 6px 14px;
    border-radius: 4px;
    font-size: 10pt;
    cursor: pointer;
    border: none;
  }

  .btn-secondary {
    background-color: var(--ctp-surface1);
    color: var(--ctp-text);
  }

  .btn-secondary:hover {
    background-color: var(--ctp-surface2);
  }

  .btn-primary {
    background-color: var(--ctp-blue);
    color: var(--ctp-crust);
    font-weight: 600;
  }

  .btn-primary:hover {
    opacity: 0.9;
  }
</style>

