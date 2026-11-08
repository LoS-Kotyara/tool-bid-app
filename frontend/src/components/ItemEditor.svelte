<script lang="ts">
  import { createEventDispatcher, onMount, onDestroy, tick } from 'svelte';

  export let text: string = '';
  export let editorFont: string = 'Cascadia Code';

  const dispatch = createEventDispatcher<{
    change: string;
    undo: void;
    redo: void;
  }>();

  const lineHeight: number = 22;

  let textareaEl: HTMLTextAreaElement;
  let mirrorEl: HTMLDivElement;
  let lines: string[] = [''];
  let lineHeights: number[] = [];
  let activeLine: number = 1;
  let scrollTop: number = 0;
  let textareaWidth: number = 0;
  let resizeObserver: ResizeObserver | null = null;
  let prevText: string = '';

  $: if (text !== prevText) {
    prevText = text;
    lines = text.split('\n');
    if (lines.length === 0) lines = [''];
    if (activeLine > lines.length) {
      activeLine = Math.max(1, lines.length);
    }
    updateLineHeights();
  }

  let activeLineTop: number = 0;
  let activeLineHeight: number = lineHeight;

  $: {
    let top = 0;
    const count = Math.min(activeLine - 1, lines.length);
    for (let i = 0; i < count; i++) {
      top += (lineHeights[i] || lineHeight);
    }
    activeLineTop = top - scrollTop;
    activeLineHeight = lineHeights[activeLine - 1] || lineHeight;
  }

  async function updateLineHeights() {
    if (!textareaEl) return;
    const newWidth = textareaEl.clientWidth || textareaEl.getBoundingClientRect().width || 500;
    if (newWidth > 0 && newWidth !== textareaWidth) {
      textareaWidth = newWidth;
    }
    await tick();
    if (!mirrorEl) return;
    const children = mirrorEl.children;
    const len = children.length;
    const newHeights = new Array(len);
    for (let i = 0; i < len; i++) {
      const h = (children[i] as HTMLElement).offsetHeight;
      newHeights[i] = h >= lineHeight ? h : lineHeight;
    }
    lineHeights = newHeights;
  }

  function handleInput(e: Event) {
    const target = e.target as HTMLTextAreaElement;
    text = target.value;
    updateCursorLine();
    dispatch('change', text);
    updateLineHeights();
  }

  function handleScroll(e: Event) {
    const target = e.target as HTMLTextAreaElement;
    scrollTop = target.scrollTop;
  }

  function updateCursorLine() {
    if (!textareaEl) return;
    const pos = textareaEl.selectionStart ?? 0;
    const textBefore = textareaEl.value.slice(0, pos);
    let newlineCount = 0;
    for (let i = 0; i < textBefore.length; i++) {
      if (textBefore[i] === '\n') {
        newlineCount++;
      }
    }
    activeLine = newlineCount + 1;
  }

  function handleKeydown(e: KeyboardEvent) {
    // Поддержка хоткеев Undo/Redo внутри textarea
    if (e.ctrlKey && (e.code === 'KeyZ')) {
      if (e.shiftKey) {
        e.preventDefault();
        dispatch('redo');
      } else {
        e.preventDefault();
        dispatch('undo');
      }
    } else if (e.ctrlKey && e.code === 'KeyY') {
      e.preventDefault();
      dispatch('redo');
    }

    requestAnimationFrame(updateCursorLine);
  }

  onMount(() => {
    updateCursorLine();
    if (textareaEl) {
      textareaWidth = textareaEl.clientWidth;
      updateLineHeights();
    }
    if (typeof ResizeObserver !== 'undefined' && textareaEl) {
      resizeObserver = new ResizeObserver(() => {
        updateLineHeights();
      });
      resizeObserver.observe(textareaEl);
    }
  });

  onDestroy(() => {
    if (resizeObserver) {
      resizeObserver.disconnect();
      resizeObserver = null;
    }
  });
</script>

<div class="editor-container" style="--editor-font: '{editorFont}', Consolas, monospace;">
  <div class="editor-gutter" style="transform: translateY(-{scrollTop}px);">
    {#each lines as _, idx}
      <div
        class="gutter-num"
        class:active-line={idx + 1 === activeLine}
        style="height: {lineHeights[idx] || lineHeight}px; line-height: {lineHeight}px;"
      >
        {idx + 1}
      </div>
    {/each}
  </div>

  <div class="editor-body">
    <div
      class="mirror-container"
      bind:this={mirrorEl}
      style="width: {textareaWidth ? `${textareaWidth}px` : '100%'};"
    >
      {#each lines as line}
        <div class="mirror-line">{line === '' ? '\u200B' : line}</div>
      {/each}
    </div>

    {#if activeLine >= 1 && activeLine <= lines.length}
      <div
        class="active-line-backdrop"
        style="top: {activeLineTop}px; height: {activeLineHeight}px;"
      ></div>
    {/if}
    <textarea
      bind:this={textareaEl}
      value={text}
      on:input={handleInput}
      on:scroll={handleScroll}
      on:click={updateCursorLine}
      on:keyup={updateCursorLine}
      on:select={updateCursorLine}
      on:mouseup={updateCursorLine}
      on:focus={updateCursorLine}
      on:keydown={handleKeydown}
      spellcheck="false"
      placeholder="Наименование; Характеристики; Ед. изм.; Количество (например: Сверло т/с; EDS08090-0; шт.; 2)"
    ></textarea>
  </div>
</div>

<style>
  .editor-container {
    flex: 1;
    display: flex;
    overflow: hidden;
    position: relative;
    background-color: var(--ctp-base);
    font-family: var(--editor-font);
    font-size: 12pt;
    line-height: 22px;
  }

  .editor-gutter {
    width: 46px;
    min-width: 46px;
    max-width: 46px;
    background-color: var(--ctp-mantle);
    border-right: 1px solid var(--ctp-surface0);
    padding: 0;
    box-sizing: border-box;
    overflow: hidden;
    user-select: none;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    will-change: transform;
  }

  .gutter-num {
    height: 22px;
    padding-right: 10px;
    text-align: right;
    color: var(--ctp-overlay0);
    font-size: 12pt;
    line-height: 22px;
    box-sizing: border-box;
    overflow: visible;
    white-space: nowrap;
  }

  .gutter-num.active-line {
    color: var(--ctp-lavender, var(--ctp-text));
    font-weight: bold;
    background-color: var(--ctp-surface0, rgba(255, 255, 255, 0.05));
  }

  .editor-body {
    flex: 1;
    position: relative;
    display: flex;
    overflow: hidden;
  }

  .mirror-container {
    position: absolute;
    top: 0;
    left: 0;
    visibility: hidden;
    pointer-events: none;
    z-index: -1;
    padding: 0 12px;
    margin: 0;
    border: none;
    box-sizing: border-box;
    font-family: inherit;
    font-size: 12pt;
    line-height: 22px;
    white-space: pre-wrap;
    word-break: break-word;
    overflow-wrap: break-word;
    overflow: hidden;
  }

  .mirror-line {
    min-height: 22px;
    line-height: 22px;
    word-break: break-word;
    overflow-wrap: break-word;
    white-space: pre-wrap;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  .active-line-backdrop {
    position: absolute;
    left: 0;
    width: 100%;
    background-color: var(--ctp-surface0, rgba(255, 255, 255, 0.05));
    pointer-events: none;
    z-index: 0;
  }

  textarea {
    flex: 1;
    height: 100%;
    padding: 0 12px;
    margin: 0;
    border: none;
    background: transparent;
    color: var(--ctp-text);
    font-family: inherit;
    font-size: 12pt;
    line-height: 22px;
    resize: none;
    outline: none;
    white-space: pre-wrap;
    word-break: break-word;
    overflow-wrap: break-word;
    overflow-x: hidden;
    overflow-y: scroll;
    box-sizing: border-box;
    position: relative;
    z-index: 1;
  }

  textarea::placeholder {
    color: var(--ctp-overlay0);
    font-style: italic;
    font-size: 10.5pt;
    line-height: 22px;
  }
</style>
