<script lang="ts">
  import { onMount } from 'svelte';

  export let previewHtml: string = '';

  let zoomLevel: number = 100;
  let iframeEl: HTMLIFrameElement;
  let currentPage: number = 1;
  let totalPages: number = 1;

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
      // Вычисляем приблизительное количество страниц (высота листа A4 ~ 1122px при 96 DPI)
      setTimeout(() => {
        const bodyHeight = doc.body.scrollHeight;
        totalPages = Math.max(1, Math.ceil(bodyHeight / 1120));
      }, 100);
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

  function nextPage() {
    if (currentPage < totalPages) {
      currentPage++;
      scrollToPage(currentPage);
    }
  }

  function prevPage() {
    if (currentPage > 1) {
      currentPage--;
      scrollToPage(currentPage);
    }
  }

  function scrollToPage(p: number) {
    if (!iframeEl || !iframeEl.contentWindow) return;
    const offset = (p - 1) * 1120;
    iframeEl.contentWindow.scrollTo({ top: offset, behavior: 'smooth' });
  }

  onMount(() => {
    updateIframeContent(previewHtml);
  });
</script>

<div class="preview-panel" on:wheel={handleWheel}>
  <!-- Тулбар превью: Зум и пагинация -->
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
      <button class="tool-btn" disabled={currentPage <= 1} on:click={prevPage} title="Предыдущая страница">
        ◀
      </button>
      <span class="page-indicator">Стр. {currentPage} из {totalPages}</span>
      <button class="tool-btn" disabled={currentPage >= totalPages} on:click={nextPage} title="Следующая страница">
        ▶
      </button>
    </div>
  </div>

  <!-- Область отображения листа А4 с центрированием и зумом -->
  <div class="sheet-viewport">
    <div
      class="a4-sheet-container"
      style="transform: scale({zoomLevel / 100}); transform-origin: top center;"
    >
      <iframe
        bind:this={iframeEl}
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

  .page-indicator {
    font-size: 9pt;
    color: var(--ctp-subtext0);
    padding: 0 6px;
  }

  .sheet-viewport {
    flex: 1;
    overflow: auto;
    display: flex;
    justify-content: center;
    padding: 24px 16px;
  }

  .a4-sheet-container {
    transition: transform 0.1s ease-out;
    display: flex;
    justify-content: center;
  }

  /* Пропорции формата А4 (210мм x 297мм) в пикселях: 794px x 1123px */
  .a4-sheet {
    width: 210mm;
    min-height: 297mm;
    background-color: #ffffff;
    border: none;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.45);
    border-radius: 2px;
    display: block;
  }
</style>
