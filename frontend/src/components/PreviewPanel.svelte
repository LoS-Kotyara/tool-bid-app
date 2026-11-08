<script lang="ts">
  import { onMount } from 'svelte';

  export let previewHtml: string = '';

  let zoomLevel: number = 100;
  let iframeEl: HTMLIFrameElement;

  $: if (iframeEl && previewHtml) {
    updateIframeContent(previewHtml);
  }

  function updateIframeContent(htmlContent: string) {
    if (!iframeEl) return;
    const doc = iframeEl.contentDocument || iframeEl.contentWindow?.document;
    if (doc) {
      doc.open();
      doc.write(htmlContent);
      doc.close();
      // Документ рендерится единым непрерывным потоком: высота листа подстраивается под
      // содержимое, а разделение на страницы выполняет движок печати.
      window.setTimeout(fitIframeHeight, 60);
    }
  }

  // fitIframeHeight растягивает iframe по фактической высоте непрерывного листа А4
  function fitIframeHeight() {
    if (!iframeEl) return;
    const doc = iframeEl.contentDocument;
    if (!doc || !doc.body) return;
    const contentHeight = Math.max(doc.body.scrollHeight, doc.documentElement.scrollHeight);
    if (contentHeight > 0) {
      iframeEl.style.height = `${contentHeight}px`;
    }
  }

  function zoomIn() {
    if (zoomLevel < 200) zoomLevel += 10;
  }

  function zoomOut() {
    if (zoomLevel > 50) zoomLevel -= 10;
  }

  function zoomReset() {
    zoomLevel = 100;
  }

  function handleWheel(e: WheelEvent) {
    if (e.ctrlKey) {
      e.preventDefault();
      if (e.deltaY < 0) {
        zoomIn();
      } else {
        zoomOut();
      }
    }
  }

  function printPreview() {
    const printWindow = iframeEl?.contentWindow;
    if (!printWindow) return;
    printWindow.focus();
    printWindow.print();
  }

  onMount(() => {
    updateIframeContent(previewHtml);
  });
</script>

<div class="preview-panel" on:wheel={handleWheel}>
  <!-- Тулбар превью: масштаб и печать (пагинация не требуется — лист непрерывный) -->
  <div class="preview-toolbar">
    <div class="toolbar-group">
      <button class="tool-btn" title="Уменьшить (Ctrl + Колесо вниз)" on:click={zoomOut}>−</button>
      <button class="tool-btn zoom-text" title="Сбросить масштаб (100%)" on:click={zoomReset}>
        {zoomLevel}%
      </button>
      <button class="tool-btn" title="Увеличить (Ctrl + Колесо вверх)" on:click={zoomIn}>+</button>
      <button class="tool-btn" title="Масштаб 1:1" on:click={zoomReset}>1:1</button>
    </div>

    <div class="toolbar-group">
      <button class="tool-btn print-btn" title="Печать документа (Ctrl + P)" on:click={printPreview}>
        Печать
      </button>
    </div>
  </div>

  <!-- Область отображения непрерывного листа А4 с центрированием и зумом -->
  <div class="sheet-viewport">
    <div
      class="a4-sheet-container"
      style="transform: scale({zoomLevel / 100}); transform-origin: top center;"
    >
      <iframe
        bind:this={iframeEl}
        on:load={fitIframeHeight}
        title="Предпросмотр заявки A4"
        class="a4-sheet"
      ></iframe>
    </div>
  </div>
</div>

<style>
  .preview-panel {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: var(--ctp-crust);
    overflow: hidden;
    position: relative;
    user-select: none;
  }

  .preview-toolbar {
    height: 36px;
    background-color: var(--ctp-mantle);
    border-bottom: 1px solid var(--ctp-surface0);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 12px;
    z-index: 10;
  }

  .toolbar-group {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .tool-btn {
    background: transparent;
    border: 1px solid transparent;
    color: var(--ctp-text);
    padding: 2px 8px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 10pt;
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 24px;
    height: 24px;
  }

  .tool-btn:hover:not(:disabled) {
    background-color: var(--ctp-surface0);
    border-color: var(--ctp-surface1);
  }

  .tool-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .zoom-text {
    min-width: 50px;
    font-size: 9.5pt;
  }

  .print-btn {
    padding: 2px 10px;
    gap: 4px;
  }

  .print-btn:hover {
    background-color: var(--ctp-surface1);
    border-color: var(--ctp-surface2);
  }

  .sheet-viewport {
    flex: 1;
    overflow: auto;
    display: flex;
    justify-content: center;
    align-items: flex-start;
    padding: 24px 16px;
  }

  .a4-sheet-container {
    transition: transform 0.1s ease-out;
    display: flex;
    justify-content: center;
  }

  /* Непрерывный лист: ширина 210мм, высота подстраивается под содержимое документа */
  .a4-sheet {
    width: 210mm;
    min-height: 297mm;
    background-color: #2b2d3a;
    border: none;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.45);
    border-radius: 2px;
    display: block;
  }
</style>
