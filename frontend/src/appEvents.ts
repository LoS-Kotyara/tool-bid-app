import { EventsOn } from '../wailsjs/runtime/runtime';
import { model } from '../wailsjs/go/models';

export function setupHotkeys(callbacks: {
  onSave: () => void;
  onReload: () => void;
  onNewTab: () => void;
  onOpen: () => void;
  onUndo: () => void;
  onRedo: () => void;
}) {
  return function handleKeyDown(e: KeyboardEvent) {
    if (e.ctrlKey && e.code === 'KeyS') {
      e.preventDefault();
      callbacks.onSave();
    } else if (e.code === 'F5') {
      e.preventDefault();
      callbacks.onReload();
    } else if (e.ctrlKey && e.code === 'KeyN') {
      e.preventDefault();
      callbacks.onNewTab();
    } else if (e.ctrlKey && e.code === 'KeyO') {
      e.preventDefault();
      callbacks.onOpen();
    } else if (e.ctrlKey && e.code === 'KeyZ') {
      e.preventDefault();
      if (e.shiftKey) callbacks.onRedo(); else callbacks.onUndo();
    } else if (e.ctrlKey && e.code === 'KeyY') {
      e.preventDefault();
      callbacks.onRedo();
    }
  };
}

export function subscribeFileDrop(
  onDrop: (doc: model.Document) => void,
  onError: (err: string) => void
) {
  EventsOn('file-dropped', onDrop);
  EventsOn('file-drop-error', onError);
}
