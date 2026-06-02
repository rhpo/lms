<script lang="ts">
  import type { Snippet } from "svelte";
  import PageMeta from "./PageMeta.svelte";

  let {
    children,
    title,
    subtitle,
    description = subtitle,
    actions,
  }: {
    children: Snippet;
    title: string;
    subtitle?: string;
    description?: string;
    actions?: Snippet;
  } = $props();
</script>

<PageMeta {title} description={description || "Welcome!"}>
  <main>
    <header>
      {#if description}
        <div class="header-info">
          <h1>{title}</h1>
          <p class="description">{description}</p>
        </div>
      {:else}
        <h1>{title}</h1>
      {/if}
      {#if actions}
        <div class="actions">
          {@render actions()}
        </div>
      {/if}
    </header>

    {@render children()}
  </main>
</PageMeta>

<style>
  main {
    padding: var(--spacing-lg);
    max-width: 1200px;
    margin: 0 auto;
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: var(--spacing-xl);

    h1 {
      margin: 0;
      font-size: var(--text-2xl);
      font-weight: 700;
    }

    .header-info {
      .description {
        margin: var(--spacing-xs) 0 0;
        color: var(--color-text-muted);
        font-size: var(--text-base);
      }
    }

    .actions {
      display: flex;
      gap: var(--spacing-sm);
      flex-shrink: 0;
    }
  }
</style>
