<!--
  Modal.svelte - Reusable modal primitive.

  The modern replacement for prompt()/alert()/confirm().
  Renders a centered overlay with backdrop blur, escape-to-close,
  and click-outside-to-close.

  Usage:
    <Modal bind:open title="Pick a font">
      <YourContent />
    </Modal>
-->
<script lang="ts">
  import { X } from "lucide-svelte";
  import { fade, scale } from "svelte/transition";

  let {
    open = $bindable(false),
    title = "",
    description = "",
    width = "720px",
    onClose,
    children,
  }: {
    open?: boolean;
    title?: string;
    description?: string;
    width?: string;
    onClose?: () => void;
    children: any;
  } = $props();

  function close() {
    open = false;
    onClose?.();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") close();
  }

  function handleBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) close();
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="modal-backdrop"
    onclick={handleBackdropClick}
    transition:fade={{ duration: 150 }}
  >
    <div
      class="modal-container"
      style:max-width={width}
      transition:scale={{ duration: 200, start: 0.95 }}
    >
      <!-- Header -->
      {#if title}
        <header class="modal-header">
          <div class="modal-info">
            <h2 class="modal-title">{title}</h2>
            {#if description}
              <p class="modal-description">{description}</p>
            {/if}
          </div>
          <button class="modal-close" onclick={close} aria-label="Close modal">
            <X size={16} />
          </button>
        </header>
      {/if}

      <!-- Body -->
      <div class="modal-body">
        {@render children?.()}
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    inset: 0;
    z-index: 9999;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    padding: 24px;
  }

  .modal-container {
    width: 100%;
    max-height: calc(100vh - 48px);
    display: flex;
    flex-direction: column;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 16px;
    box-shadow: var(--shadow-popover);
    overflow: hidden;
  }

  .modal-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;
    padding: 20px 24px 16px;
    border-bottom: 1px solid var(--color-border);
    flex-shrink: 0;
  }

  .modal-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .modal-title {
    margin: 0;
    font-size: 16px;
    font-weight: 700;
    color: var(--color-text);
    letter-spacing: -0.01em;
  }

  .modal-description {
    margin: 0;
    font-size: 12px;
    color: var(--color-text-muted);
    line-height: 1.4;
  }

  .modal-close {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border: none;
    background: var(--color-background-100);
    color: var(--color-text-muted);
    border-radius: 8px;
    cursor: pointer;
    flex-shrink: 0;
    transition: all 0.15s ease;
  }
  .modal-close:hover {
    background: var(--color-background-200);
    color: var(--color-text);
  }

  .modal-body {
    flex: 1;
    overflow-y: auto;
    padding: 16px 24px 24px;

    /* Custom scrollbar */
    scrollbar-width: thin;
    scrollbar-color: var(--color-border) transparent;
  }
  .modal-body::-webkit-scrollbar {
    width: 6px;
  }
  .modal-body::-webkit-scrollbar-track {
    background: transparent;
  }
  .modal-body::-webkit-scrollbar-thumb {
    background: var(--color-border);
    border-radius: 3px;
  }
  .modal-body::-webkit-scrollbar-thumb:hover {
    background: var(--color-border-strong);
  }
</style>
