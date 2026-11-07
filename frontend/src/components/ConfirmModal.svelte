<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let title: string = 'Подтверждение';
  export let message: string = '';
  export let confirmText: string = 'Продолжить';
  export let cancelText: string = 'Отмена';
  export let danger: boolean = false;

  const dispatch = createEventDispatcher<{
    confirm: void;
    cancel: void;
  }>();

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      dispatch('cancel');
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="modal-backdrop" on:click={() => dispatch('cancel')}>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="modal-window" on:click|stopPropagation>
    <div class="modal-header">
      <span class="modal-title">{title}</span>
    </div>

    <div class="modal-body">
      <p>{message}</p>
    </div>

    <div class="modal-footer">
      <button class="btn btn-secondary" on:click={() => dispatch('cancel')}>
        {cancelText}
      </button>
      <button
        class="btn"
        class:btn-danger={danger}
        class:btn-primary={!danger}
        on:click={() => dispatch('confirm')}
      >
        {confirmText}
      </button>
    </div>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2100;
  }

  .modal-window {
    background-color: var(--ctp-surface0);
    border: 1px solid var(--ctp-surface1);
    border-radius: 8px;
    width: 400px;
    box-shadow: 0 16px 36px rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .modal-header {
    height: 40px;
    padding: 0 16px;
    background-color: var(--ctp-mantle);
    border-bottom: 1px solid var(--ctp-surface1);
    display: flex;
    align-items: center;
  }

  .modal-title {
    font-size: 11pt;
    font-weight: 600;
    color: var(--ctp-text);
  }

  .modal-body {
    padding: 16px;
    font-size: 10.5pt;
    line-height: 1.4;
    color: var(--ctp-text);
  }

  .modal-footer {
    padding: 10px 16px;
    background-color: var(--ctp-mantle);
    border-top: 1px solid var(--ctp-surface1);
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }

  .btn {
    padding: 6px 14px;
    border-radius: 4px;
    font-size: 10pt;
    cursor: pointer;
    border: none;
  }

  .btn-secondary {
    background-color: var(--ctp-surface1);
    color: var(--ctp-text);
  }

  .btn-secondary:hover {
    background-color: var(--ctp-surface2);
  }

  .btn-primary {
    background-color: var(--ctp-blue);
    color: var(--ctp-crust);
    font-weight: 600;
  }

  .btn-danger {
    background-color: var(--ctp-red);
    color: var(--ctp-crust);
    font-weight: 600;
  }
</style>

