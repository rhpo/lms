<script lang="ts">
  let {
    tabs,
    active,
  }: {
    tabs: { id: string; label: string; count?: number }[];
    active: string;
  } = $props();
</script>

<div class="tabs">
  {#each tabs as tab}
    <button
      class="tab"
      class:is-active={tab.id === active}
      onclick={() => (active = tab.id)}
    >
      {tab.label}
      {#if tab.count !== undefined}
        <span class="count">{tab.count}</span>
      {/if}
    </button>
  {/each}
</div>

<style>
  .tabs {
    display: flex;
    gap: 0;
    border-bottom: 2px solid var(--color-border);
    margin-bottom: var(--spacing-lg);
  }

  .tab {
    padding: var(--spacing-sm) var(--spacing-lg);
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-text-muted);
    background: none;
    border: none;
    border-bottom: 2px solid transparent;
    margin-bottom: -2px;
    cursor: pointer;
    transition: color var(--transition-fast), border-color var(--transition-fast);
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);

    &:hover {
      color: var(--color-text);
    }

    &.is-active {
      color: var(--color-accent);
      border-bottom-color: var(--color-accent);
    }
  }

  .count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 20px;
    height: 20px;
    padding: 0 6px;
    border-radius: 999px;
    background: color-mix(in srgb, var(--color-accent) 12%, transparent);
    color: var(--color-accent);
    font-size: var(--text-xs);
    font-weight: 600;
  }
</style>
