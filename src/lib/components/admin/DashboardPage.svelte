<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    title,
    subtitle,
    actions,
    children,
  }: {
    title: string;
    subtitle?: string;
    actions?: Snippet;
    children?: Snippet;
  } = $props();
</script>

<div class="page">
  <div class="page-header">
    <div class="page-header-text">
      <h1>{title}</h1>
      {#if subtitle}
        <p>{subtitle}</p>
      {/if}
    </div>

    {#if actions}
      <div class="page-header-actions">
        {@render actions()}
      </div>
    {/if}
  </div>

  {@render children?.()}
</div>

<style>
  .page {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    flex-wrap: wrap;
  }

  h1 {
    font-size: var(--text-3xl);
    font-weight: var(--font-bold);
    font-family: var(--font-sans);
    color: var(--color-text);
    line-height: var(--leading-tight);
    margin: 0 0 0.35rem;
  }

  p {
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text-muted);
    margin: 0;
  }

  .page-header-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }

  @media (max-width: 480px) {
    .page-header {
      flex-direction: column;
    }
  }
</style>
