<script lang="ts">
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import MetaForm from './MetaForm.svelte';
  import ItemEditor from './ItemEditor.svelte';
  import PreviewPanel from './PreviewPanel.svelte';
  import type { DocumentTab, ViewMode } from '../types';

  export let activeTab: DocumentTab;
  export let viewMode: ViewMode;
  export let currentFont: string;
  export let previewHtml: string;

  const dispatch = createEventDispatcher<{
    metaChange: void;
    textChange: string;
    undo: void;
    redo: void;
  }>();

  const SPLIT_STORAGE_KEY = 'tool_bid_splitter_ratio';
  let splitRatio: number = 50; // Процент ширины левой панели (25% - 75%)
  let isDragging: boolean = false;
  let workspaceEl: HTMLElement;

  onMount(() => {
    const saved = localStorage.getItem(SPLIT_STORAGE_KEY);
    if (saved) {
      const parsed = parseFloat(saved);
      if (!isNaN(parsed) && parsed >= 25 && parsed <= 75) {
        splitRatio = parsed;
      }
    }
  });

  function startDragging(e: MouseEvent) {
    if (viewMode !== 'split') return;
    e.preventDefault();
    isDragging = true;
    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', stopDragging);
  }

  function onMouseMove(e: MouseEvent) {
    if (!isDragging || !workspaceEl) return;
    const rect = workspaceEl.getBoundingClientRect();
    if (rect.width <= 0) return;

    const offsetX = e.clientX - rect.left;
    let ratio = (offsetX / rect.width) * 100;
    if (ratio < 25) ratio = 25;
    if (ratio > 75) ratio = 75;

    splitRatio = ratio;
  }

  function stopDragging() {
    if (!isDragging) return;
    isDragging = false;
    window.removeEventListener('mousemove', onMouseMove);
    window.removeEventListener('mouseup', stopDragging);
    localStorage.setItem(SPLIT_STORAGE_KEY, splitRatio.toString());
  }

  onDestroy(() => {
    window.removeEventListener('mousemove', onMouseMove);
    window.removeEventListener('mouseup', stopDragging);
  });
</script>

<div class="workspace-area" bind:this={workspaceEl} class:is-dragging={isDragging}>
  {#if viewMode === 'split' || viewMode === 'edit'}
    <div
      class="editor-pane"
      class:full-width={viewMode === 'edit'}
      style={viewMode === 'split' ? `width: ${splitRatio}%; flex: none;` : ''}
    >
      <MetaForm
        bind:meta={activeTab.meta}
        on:change={() => dispatch('metaChange')}
      />
      <ItemEditor
        text={activeTab.rawText}
        editorFont={currentFont}
        on:change={(e) => dispatch('textChange', e.detail)}
        on:undo={() => dispatch('undo')}
        on:redo={() => dispatch('redo')}
      />
    </div>
  {/if}

  {#if viewMode === 'split'}
    <!-- Интерактивный сплиттер между редактором и предпросмотром -->
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <div
      class="workspace-splitter"
      class:dragging={isDragging}
      on:mousedown={startDragging}
      title="Потяните для изменения пропорций сплита (25% - 75%)"
      role="separator"
      tabindex="-1"
    ></div>
  {/if}

  {#if viewMode === 'split' || viewMode === 'preview'}
    <div
      class="preview-pane"
      class:full-width={viewMode === 'preview'}
      class:disable-pointer={isDragging}
      style={viewMode === 'split' ? `width: calc(${100 - splitRatio}% - 5px); flex: 1;` : ''}
    >
      <PreviewPanel {previewHtml} />
    </div>
  {/if}
</div>
