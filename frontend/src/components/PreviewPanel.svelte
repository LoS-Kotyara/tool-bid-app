<script lang="ts">
  import { onDestroy, onMount } from 'svelte';

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
  // Селектор листов и класс активной страницы в документе предпросмотра
  const PAGE_SELECTOR = '.a4-page';
  const ACTIVE_PAGE_CLASS = 'active';

  let zoomLevel: number = 100;
  let iframeEl: HTMLIFrameElement;
  let renderTimer: number | null = null;
  let pendingHtml: string = previewHtml;
  let renderedHtml: string = '';
  // Постраничный показ: на экране виден ровно один лист (как в эталоне bid_app_gui.py)
  let currentPage: number = 1;
  let totalPages: number = 1;
  // Защита от рекурсии: запись документа сама вызывает load у iframe, а восстановление
  // содержимого вызывается из обработчиков DOM предпросмотра
  let isWriting: boolean = false;
  let isRecovering: boolean = false;
  let resizeObserver: ResizeObserver | null = null;

  $: scale = zoomLevel / 100;

  // Число листов считается из входного HTML-пропса (источник истины), а не из DOM
  // iframe: при переключении режима «Сплит» -> «Ввод» -> «Сплит» iframe сбрасывается
  // в about:blank, его DOM пуст, и кнопки пагинации «залипали» в disabled.
  $: if (previewHtml) {
    totalPages = extractTotalPages(previewHtml);
    if (currentPage > totalPages) currentPage = totalPages;
    if (currentPage < 1) currentPage = 1;
  } else if (totalPages !== 1 || currentPage !== 1) {
    totalPages = 1;
    currentPage = 1;
  }

  // Любое обновление пропса HTML гарантированно доходит до iframe: если документ
  // предпросмотра пуст (панель была скрыта/пересоздана), запись идёт немедленно,
  // иначе — с дебаунсом, чтобы набор текста не перезаписывал DOM на каждый символ.
  $: if (iframeEl && previewHtml && previewHtml !== renderedHtml) {
    scheduleRender(previewHtml);
  }

  // Зум меняет только CSS-переменную документа предпросмотра: перезаписи тяжёлого
  // документа и пересчёта стилей body не происходит
  $: if (iframeEl && scale) {
    applyZoom(scale);
  }

  function applyZoom(value: number) {
    const doc = getIframeDocument();
    if (!doc || !doc.documentElement) return;
    doc.documentElement.style.setProperty(ZOOM_VAR, String(value));
  }

  // getIframeDocument возвращает документ предпросмотра независимо от способа доступа
  function getIframeDocument(): Document | null {
    if (!iframeEl) return null;
    return iframeEl.contentDocument || iframeEl.contentWindow?.document || null;
  }

  // hasPageNodes проверяет, что документ предпросмотра реально содержит листы А4
  function hasPageNodes(doc: Document | null): boolean {
    if (!doc) return false;
    return doc.querySelectorAll(PAGE_SELECTOR).length > 0;
  }

  // extractTotalPages считает листы по HTML-пропсу, не завися от готовности iframe
  function extractTotalPages(htmlContent: string): number {
    if (!htmlContent) return 1;
    const matches = htmlContent.match(/class=["'][^"']*a4-page[^"']*["']/g);
    return matches && matches.length > 0 ? matches.length : 1;
  }

  // showPage переключает активную страницу и сразу применяет её к DOM предпросмотра
  function showPage(targetPage: number | string) {
    const pageNum = Number(targetPage);
    if (isNaN(pageNum) || pageNum < 1 || pageNum > totalPages) return;
    currentPage = pageNum;
    applyPageToDOM(currentPage);
  }

  // applyPageToDOM оставляет видимым единственный лист: снимает active со всех страниц
  // документа предпросмотра, помечает нужную и возвращает прокрутку к началу листа
  function applyPageToDOM(pageNum: number) {
    if (!iframeEl) return;
    const doc = getIframeDocument();
    if (!doc) return;

    const pages = Array.from(doc.querySelectorAll<HTMLElement>(PAGE_SELECTOR));

    // iframe потерял содержимое (переключение режимов/пересоздание) — заполняем заново
    if (pages.length === 0) {
      recoverDocument();
      return;
    }

    pages.forEach((page, idx) => {
      if (idx + 1 === pageNum) {
        page.classList.add(ACTIVE_PAGE_CLASS);
      } else {
        page.classList.remove(ACTIVE_PAGE_CLASS);
      }
    });

    scrollDocumentToTop(doc);
  }

  // recoverDocument заполняет iframe заново, если он сбросился в about:blank
  function recoverDocument() {
    if (isRecovering || isWriting) return;
    if (!pendingHtml && !previewHtml) return;
    isRecovering = true;
    try {
      renderDocument();
    } finally {
      isRecovering = false;
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

  // scheduleRender откладывает перезапись документа предпросмотра на RENDER_DEBOUNCE_MS,
  // схлопывая серию быстрых правок в одно обновление. Немедленная запись выполняется,
  // когда документ предпросмотра пуст (первый рендер или сброс iframe при смене режима).
  function scheduleRender(htmlContent: string) {
    pendingHtml = htmlContent;
    if (renderTimer !== null) {
      window.clearTimeout(renderTimer);
      renderTimer = null;
    }
    if (renderedHtml === '' || !hasPageNodes(getIframeDocument())) {
      renderDocument();
      return;
    }
    renderTimer = window.setTimeout(() => {
      renderTimer = null;
      renderDocument();
    }, RENDER_DEBOUNCE_MS);
  }

  // renderDocument — единственная точка записи HTML в iframe. Вызывается при
  // монтировании, при обновлении пропса и при восстановлении сброшенного документа.
  // Идемпотентен: на неизменном документе лишь синхронизирует масштаб и активный лист,
  // поэтому кнопки пагинации всегда работают с актуальным DOM.
  function renderDocument() {
    if (!iframeEl) return;
    const doc = getIframeDocument();
    if (!doc) return;

    const htmlContent = pendingHtml || previewHtml;
    if (!htmlContent) return;

    if (htmlContent === renderedHtml && hasPageNodes(doc)) {
      applyZoom(scale);
      applyPageToDOM(currentPage);
      return;
    }

    isWriting = true;
    try {
      doc.open();
      doc.write(htmlContent);
      doc.close();
      renderedHtml = htmlContent;

      ensureScaleWrapper(doc);
      applyZoom(scale);
      // Применяем активную страницу к свежезаписанному DOM
      applyPageToDOM(currentPage);
      // Ctrl + колесо мыши над листом: масштаб без перезагрузки документа
      doc.addEventListener('wheel', handleInnerWheel, { passive: false });
    } finally {
      // load-событие от собственной записи гасится флагом isWriting
      window.setTimeout(() => {
        isWriting = false;
      }, 0);
    }
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
    // Собственная запись в документ тоже вызывает load — не перезаписываем сами себя
    if (isWriting) return;
    const doc = getIframeDocument();
    if (!doc) return;
    ensureScaleWrapper(doc);
    applyZoom(scale);
    // iframe потерял документ (перемонтирование/сброс в about:blank) — заполняем заново,
    // иначе клики по [◀]/[▶] не найдут элементов .a4-page
    if (!hasPageNodes(doc)) {
      recoverDocument();
      return;
    }
    applyPageToDOM(currentPage);
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

  // Гарантированная запись при монтировании: переход «Ввод» -> «Сплит» пересоздаёт
  // компонент, и предпросмотр обязан сразу содержать листы
  onMount(() => {
    renderDocument();
    if (typeof ResizeObserver !== 'undefined' && iframeEl) {
      // Страховка на смену видимости панели: если iframe при показе потерял документ,
      // содержимое восстанавливается. Прокрутка и активный лист не сбрасываются,
      // пока документ на месте (иначе перетаскивание разделителя дёргало бы лист).
      resizeObserver = new ResizeObserver(() => {
        if (!hasPageNodes(getIframeDocument())) {
          recoverDocument();
        }
      });
      resizeObserver.observe(iframeEl);
    }
  });

  onDestroy(() => {
    if (renderTimer !== null) {
      window.clearTimeout(renderTimer);
      renderTimer = null;
    }
    if (resizeObserver) {
      resizeObserver.disconnect();
      resizeObserver = null;
    }
  });
</script>

<div class="preview-panel" on:wheel={handleWheel}>
  <!-- Тулбар превью: постраничная навигация, масштаб и печать -->
  <div class="preview-toolbar">
    <div class="toolbar-group">
      <button
        class="tool-btn page-btn"
        type="button"
        title="Предыдущая страница"
        disabled={currentPage <= 1}
        on:click={() => showPage(currentPage - 1)}>◀</button>
      <span class="page-indicator">Стр. {currentPage} из {totalPages}</span>
      <button
        class="tool-btn page-btn"
        type="button"
        title="Следующая страница"
        disabled={currentPage >= totalPages}
        on:click={() => showPage(currentPage + 1)}>▶</button>
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

  /* Кнопки постраничной навигации: блокируются только по границам диапазона
     (currentPage/totalPages), без посторонних флагов загрузки/рендеринга */
  .page-btn {
    min-width: 26px;
    font-size: 9pt;
    line-height: 1;
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
