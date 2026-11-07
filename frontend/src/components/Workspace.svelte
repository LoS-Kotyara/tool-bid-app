<script lang="ts">
  import { createEventDispatcher } from 'svelte';
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
</script>

<div class="workspace-area">
  {#if viewMode === 'split' || viewMode === 'edit'}
    <div class="editor-pane" class:full-width={viewMode === 'edit'}>
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

  {#if viewMode === 'split' || viewMode === 'preview'}
    <div class="preview-pane" class:full-width={viewMode === 'preview'}>
      <PreviewPanel {previewHtml} />
    </div>
  {/if}
</div>
