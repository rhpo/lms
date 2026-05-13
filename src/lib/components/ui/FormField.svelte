<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    label,
    error,
    required = false,
    children,
  }: {
    label: string;
    error?: string;
    required?: boolean;
    children?: Snippet;
  } = $props();
</script>

<div class="field">
  <label>
    {label}
    {#if required}
      <span class="required">*</span>
    {/if}
  </label>
  {@render children?.()}
  {#if error}
    <p class="error">{error}</p>
  {/if}
</div>

<style>
  .field {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  label {
    font-size: var(--text-sm);
    font-weight: 600;
    font-family: var(--font-sans);
    color: var(--color-text);
  }

  .required {
    color: var(--color-danger);
  }

  .error {
    font-size: var(--text-xs);
    font-family: var(--font-sans);
    color: var(--color-danger);
    margin: 0;
  }

  :global(.input) {
    width: 100%;
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-surface);
    color: var(--color-text);
    transition: border-color 0.15s, box-shadow 0.15s;
    box-sizing: border-box;
  }

  :global(.input:focus) {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-accent) 15%, transparent);
  }

  :global(.input.has-error) {
    border-color: var(--color-danger);
  }

  :global(textarea.input) {
    min-height: 100px;
    resize: vertical;
  }

  :global(select.input) {
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    padding-right: 2rem;
  }

  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .form-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-top: 0.5rem;
  }

  @media (max-width: 600px) {
    .form-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
