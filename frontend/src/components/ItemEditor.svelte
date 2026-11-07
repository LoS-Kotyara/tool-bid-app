<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';

  export let text: string = '';
  export let editorFont: string = 'Cascadia Code';

  const dispatch = createEventDispatcher<{
    change: string;
    undo: void;
    redo: void;
  }>();

  let textareaEl: HTMLTextAreaElement;
  let lines: string[] = [''];
  let activeLineIndex: number = 0;
  let scrollTop: number = 0;

  $: {
    lines = text.split('\n');
    if (lines.length === 0) lines = [''];
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
    const pos = textareaEl.selectionStart;
    const textBefore = textareaEl.value.slice(0, pos);
    activeLineIndex = textBefore.split('\n').length - 1;
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
  }

  onMount(() => {
    updateCursorLine();
  });
</script>

<div class="editor-container" style="--editor-font: '{editorFont}', Consolas, monospace;">
  <div class="editor-gutter" style="transform: translateY(-{scrollTop}px);">
    {#each lines as _, idx}
      <div class="gutter-num" class:active-line={idx === activeLineIndex}>
        {idx + 1}
      </div>
    {/each}
  </div>

  <div class="editor-body">
    <textarea
      bind:this={textareaEl}
      value={text}
      on:input={handleInput}
      on:scroll={handleScroll}
      on:click={updateCursorLine}
      on:keyup={updateCursorLine}
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
    line-height: 20px;
  }

  .editor-gutter {
    width: 48px;
    background-color: var(--ctp-mantle);
    border-right: 1px solid var(--ctp-surface0);
    padding: 8px 0;
    user-select: none;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    will-change: transform;
  }

  .gutter-num {
    height: 20px;
    padding-right: 10px;
    text-align: right;
    color: var(--ctp-overlay0);
    font-size: 9.5pt;
    line-height: 20px;
  }

  .gutter-num.active-line {
    color: var(--ctp-text);
    font-weight: bold;
    background-color: var(--ctp-surface0);
  }

  .editor-body {
    flex: 1;
    position: relative;
    display: flex;
  }

  textarea {
    flex: 1;
    height: 100%;
    padding: 8px 12px;
    margin: 0;
    border: none;
    background: transparent;
    color: var(--ctp-text);
    font-family: inherit;
    font-size: inherit;
    line-height: 20px;
    resize: none;
    outline: none;
    white-space: pre-wrap;
    word-break: break-word;
    overflow-y: scroll;
    box-sizing: border-box;
  }

  textarea::placeholder {
    color: var(--ctp-overlay0);
    font-style: italic;
    font-size: 10.5pt;
  }
</style>
