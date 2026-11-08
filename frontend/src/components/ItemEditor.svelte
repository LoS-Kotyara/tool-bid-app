<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';

  export let text: string = '';
  export let editorFont: string = 'Cascadia Code';

  const dispatch = createEventDispatcher<{
    change: string;
    undo: void;
    redo: void;
  }>();

  const lineHeight: number = 22;

  let textareaEl: HTMLTextAreaElement;
  let lines: string[] = [''];
  let activeLine: number = 1;
  let scrollTop: number = 0;

  $: {
    lines = text.split('\n');
    if (lines.length === 0) lines = [''];
    if (activeLine > lines.length) {
      activeLine = Math.max(1, lines.length);
    }
  }

  function handleInput(e: Event) {
    const target = e.target as HTMLTextAreaElement;
    text = target.value;
    updateCursorLine();
    dispatch('change', text);
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
  });
</script>

<div class="editor-container" style="--editor-font: '{editorFont}', Consolas, monospace;">
  <div class="editor-gutter" style="transform: translateY(-{scrollTop}px);">
    {#each lines as _, idx}
      <div class="gutter-num" class:active-line={idx + 1 === activeLine}>
        {idx + 1}
      </div>
    {/each}
  </div>

  <div class="editor-body">
    {#if activeLine >= 1 && activeLine <= lines.length}
      <div
        class="active-line-backdrop"
        style="top: {(activeLine - 1) * lineHeight - scrollTop}px; height: {lineHeight}px;"
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
    white-space: pre;
    word-wrap: normal;
  }

  .editor-gutter {
    width: 48px;
    background-color: var(--ctp-mantle);
    border-right: 1px solid var(--ctp-surface0);
    padding: 0;
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

  .active-line-backdrop {
    position: absolute;
    left: 0;
    width: 100%;
    height: 22px;
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
    white-space: pre;
    overflow-x: auto;
    overflow-y: scroll;
    word-wrap: normal;
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
