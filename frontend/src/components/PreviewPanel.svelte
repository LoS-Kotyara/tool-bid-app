<script lang="ts">
  import { onDestroy } from 'svelte';

  export let previewHtml: string = '';

  // Дебаунс перерисовки: набор текста в редакторе не должен перезаписывать DOM
  // тяжёлого предпросмотра (несколько листов А4) на каждое нажатие клавиши.
  const RENDER_DEBOUNCE_MS = 300;
  const ZOOM_STEP = 10;
  const ZOOM_MIN = 50;
  const ZOOM_MAX = 200;
  // Высота одного листа А4 при 96 dpi — резерв до первой синхронизации размеров
  const A4_HEIGHT_PX = Math.round((297 / 25.4) * 96);
  const A4_WIDTH_MM = 210;

  let zoomLevel: number = 100;
  let iframeEl: HTMLIFrameElement;
  let renderTimer: number | null = null;
  let pendingHtml: string = '';
  let renderedHtml: string = '';
  let sheetHeight: number = 0;

  $: scale = zoomLevel / 100;
  $: stageHeight = Math.max(Math.round(sheetHeight * scale), Math.round(A4_HEIGHT_PX * scale));

  $: if (iframeEl && previewHtml && previewHtml !== renderedHtml) {
    scheduleRender(previewHtml);
  }

  // scheduleRender откладывает перезапись документа предпросмотра на RENDER_DEBOUNCE_MS,
  // схлопывая серию быстрых правок в одно обновление. Первый рендер выполняется сразу.
  function scheduleRender(htmlContent: string) {
    pendingHtml = htmlContent;
    if (renderTimer !== null) {
      window.clearTimeout(renderTimer);
      renderTimer = null;
    }
    if (renderedHtml === '') {
      updateIframeContent(pendingHtml);
      return;
    }
    renderTimer = window.setTimeout(() => {
      renderTimer = null;
      updateIframeContent(pendingHtml);
    }, RENDER_DEBOUNCE_MS);
  }

  function updateIframeContent(htmlContent: string) {
    if (!iframeEl) return;
    const doc = iframeEl.contentDocument || iframeEl.contentWindow?.document;
    if (doc) {
      doc.open();
      doc.write(htmlContent);
      doc.close();
      renderedHtml = htmlContent;
      // Документ рендерится дискретными листами А4: iframe растягивается на все листы,
      // а разбиение на страницы выполняет движок печати.
      window.setTimeout(fitIframeHeight, 60);
    }
  }

  // fitIframeHeight растягивает iframe по фактической высоте ленты листов А4
  function fitIframeHeight() {
    if (!iframeEl) return;
    const doc = iframeEl.contentDocument;
    if (!doc || !doc.body) return;
    const contentHeight = Math.max(doc.body.scrollHeight, doc.documentElement.scrollHeight);
    if (contentHeight > 0) {
      sheetHeight = contentHeight;
      iframeEl.style.height = `${contentHeight}px`;
    }
  }

  function zoomIn() {
    if (zoomLevel < ZOOM_MAX) zoomLevel += ZOOM_STEP;
  }

  function zoomOut() {
    if (zoomLevel > ZOOM_MIN) zoomLevel -= ZOOM_STEP;
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

  onDestroy(() => {
    if (renderTimer !== null) {
      window.clearTimeout(renderTimer);
      renderTimer = null;
    }
  });
</script>

<div class="preview-panel" on:wheel={handleWheel}>
  <!-- Тулбар превью: масштаб и печать (пагинация выполняется в самом документе) -->
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

  <!-- Область отображения ленты листов А4: масштаб применяется к внутреннему
       контейнеру страниц через transform: scale, стили документа не пересчитываются -->
  <div class="sheet-viewport">
    <div class="a4-zoom-box" style="width: calc({A4_WIDTH_MM}mm * {scale}); height: {stageHeight}px;">
      <div
        class="a4-sheet-container"
        style="transform: scale({scale}) translateZ(0); transform-origin: top left; will-change: transform;"
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
    /* Собственный слой композиции: плавный скролл длинной ленты листов */
    will-change: transform;
    transform: translateZ(0);
  }

  /* Коробка резервирует место под масштабированную ленту, чтобы скроллбар был корректным */
  .a4-zoom-box {
    flex: 0 0 auto;
    display: flex;
    justify-content: center;
  }

  .a4-sheet-container {
    transition: transform 0.1s ease-out;
    display: flex;
    justify-content: center;
    will-change: transform;
    backface-visibility: hidden;
  }

  /* Лента дискретных листов: ширина 210 мм, высота синхронизируется по содержимому */
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
