<script lang="ts">
  import { onDestroy } from 'svelte';

  export let previewHtml: string = '';

  // Дебаунс перерисовки: набор текста в редакторе не должен перезаписывать DOM
  // тяжёлого предпросмотра (несколько листов А4) на каждое нажатие клавиши.
  const RENDER_DEBOUNCE_MS = 300;
  const ZOOM_STEP = 10;
  const ZOOM_MIN = 50;
  const ZOOM_MAX = 200;
  // Масштаб ленты листов передаётся CSS-переменной в документ предпросмотра
  const ZOOM_VAR = '--zoom-level';
  const SCALE_WRAPPER_CLASS = 'a4-scale-wrapper';

  let zoomLevel: number = 100;
  let iframeEl: HTMLIFrameElement;
  let renderTimer: number | null = null;
  let pendingHtml: string = '';
  let renderedHtml: string = '';
  // Постраничный показ: на экране виден ровно один лист (как в эталоне bid_app_gui.py)
  let currentPage: number = 1;
  let totalPages: number = 1;

  $: scale = zoomLevel / 100;
  $: canPrev = currentPage > 1;
  $: canNext = currentPage < totalPages;

  $: if (iframeEl && previewHtml && previewHtml !== renderedHtml) {
    scheduleRender(previewHtml);
  }

  // Зум меняет только CSS-переменную документа предпросмотра: перезаписи тяжёлого
  // документа и пересчёта стилей body не происходит
  $: if (iframeEl && scale) {
    applyZoom(scale);
  }

  function applyZoom(value: number) {
    const doc = iframeEl?.contentDocument;
    if (!doc || !doc.documentElement) return;
    doc.documentElement.style.setProperty(ZOOM_VAR, String(value));
  }

  // showPage оставляет видимым единственный лист: снимает active со всех страниц
  // документа предпросмотра, помечает нужную и синхронизирует счётчик в тулбаре
  function showPage(pageNum: number, resetScroll: boolean = true) {
    const doc = iframeEl?.contentDocument;
    if (!doc) return;

    const pages = Array.from(doc.querySelectorAll<HTMLElement>('.a4-page'));
    if (pages.length === 0) {
      totalPages = 1;
      currentPage = 1;
      return;
    }

    totalPages = pages.length;
    const target = Math.min(Math.max(Math.round(pageNum), 1), pages.length);

    pages.forEach((page) => page.classList.remove('active'));
    const page = doc.getElementById(`page-${target}`) ?? pages[target - 1];
    if (page) {
      page.classList.add('active');
    }

    currentPage = target;
    if (resetScroll) {
      scrollDocumentToTop(doc);
    }
  }

  // scrollDocumentToTop возвращает прокрутку к началу листа при смене страницы
  function scrollDocumentToTop(doc: Document) {
    doc.documentElement.scrollTop = 0;
    doc.documentElement.scrollLeft = 0;
    if (doc.body) {
      doc.body.scrollTop = 0;
      doc.body.scrollLeft = 0;
    }
  }

  function prevPage() {
    if (currentPage > 1) showPage(currentPage - 1);
  }

  function nextPage() {
    if (currentPage < totalPages) showPage(currentPage + 1);
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
    if (!doc) return;

    doc.open();
    doc.write(htmlContent);
    doc.close();
    renderedHtml = htmlContent;

    ensureScaleWrapper(doc);
    applyZoom(scale);
    // Число листов пересчитывается, активная страница сохраняется (или усекается,
    // если после правки листов стало меньше)
    showPage(currentPage, false);
    // Ctrl + колесо мыши над листом: масштаб без перезагрузки документа
    doc.addEventListener('wheel', handleInnerWheel, { passive: false });
  }

  // ensureScaleWrapper оборачивает ленту листов в единый масштабируемый контейнер,
  // не создавая второго скролла: скроллбаром владеет только body документа предпросмотра
  function ensureScaleWrapper(doc: Document): HTMLElement | null {
    const body = doc.body;
    if (!body) return null;

    const first = body.firstElementChild as HTMLElement | null;
    if (first && first.classList.contains(SCALE_WRAPPER_CLASS)) {
      return first;
    }

    const wrapper = doc.createElement('div');
    wrapper.className = SCALE_WRAPPER_CLASS;
    while (body.firstChild) {
      wrapper.appendChild(body.firstChild);
    }
    body.appendChild(wrapper);

    return wrapper;
  }

  function handleInnerWheel(e: WheelEvent) {
    if (!e.ctrlKey) return;
    e.preventDefault();
    if (e.deltaY < 0) {
      zoomIn();
    } else {
      zoomOut();
    }
  }

  function handleIframeLoad() {
    // После (пере)загрузки документа масштабируемая обёртка, масштаб и активный
    // лист восстанавливаются
    const doc = iframeEl?.contentDocument;
    if (!doc) return;
    ensureScaleWrapper(doc);
    applyZoom(scale);
    showPage(currentPage, false);
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
  <!-- Тулбар превью: постраничная навигация, масштаб и печать -->
  <div class="preview-toolbar">
    <div class="toolbar-group">
      <button
        class="tool-btn"
        title="Предыдущая страница"
        disabled={!canPrev}
        on:click={prevPage}>◀</button>
      <span class="page-indicator">Стр. {currentPage} из {totalPages}</span>
      <button
        class="tool-btn"
        title="Следующая страница"
        disabled={!canNext}
        on:click={nextPage}>▶</button>
    </div>

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

  <!-- Область предпросмотра занимает всю ширину и высоту панели: iframe растянут на 100%,
       а единственный скроллбар живёт на body документа внутри него (правый край панели) -->
  <div class="sheet-viewport">
    <iframe
      bind:this={iframeEl}
      on:load={handleIframeLoad}
      title="Предпросмотр заявки A4"
      class="a4-sheet"
    ></iframe>
  </div>
</div>

<style>
  .preview-panel {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    background-color: var(--ctp-crust);
    overflow: hidden;
    position: relative;
    user-select: none;
  }

  .preview-toolbar {
    flex: 0 0 auto;
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

  /* Счётчик страниц: «Стр. X из Y» (эталон bid_app_gui.py) */
  .page-indicator {
    font-size: 9.5pt;
    color: var(--ctp-subtext1);
    text-align: center;
    min-width: 92px;
    padding: 0 2px;
    white-space: nowrap;
  }

  .print-btn {
    padding: 2px 10px;
    gap: 4px;
  }

  .print-btn:hover {
    background-color: var(--ctp-surface1);
    border-color: var(--ctp-surface2);
  }

  /* Область предпросмотра: вся доступная ширина и высота правой панели.
     Здесь нет прокрутки — скроллбаром владеет документ внутри iframe,
     поэтому он всегда расположен у правого края панели. */
  .sheet-viewport {
    flex: 1 1 auto;
    width: 100%;
    height: 100%;
    min-height: 0;
    overflow: hidden;
    position: relative;
    background-color: #2b2d3a;
  }

  /* iframe растянут на всю область: колесо мыши над листом и над тёмными полями
     прокручивает один и тот же документ — без двойного скролла и заеданий */
  .a4-sheet {
    display: block;
    width: 100%;
    height: 100%;
    border: none;
    background-color: #2b2d3a;
  }
</style>
