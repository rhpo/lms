<script lang="ts">
  let {
    variant = "default",
    label = "",
  }: {
    variant?: "default" | "success" | "warning" | "danger" | "info" | "neutral";
    label?: string;
  } = $props();
</script>

<span class="badge" data-variant={variant}>
  {label}
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    padding: 0.2rem 0.6rem;
    border-radius: 999px;
    font-size: 0.75rem;
    font-weight: 600;
    font-family: var(--font-sans);
    text-transform: uppercase;
    letter-spacing: 0.03em;
    white-space: nowrap;
    line-height: 1.4;

    width: fit-content;
  }

  .badge[data-variant="default"] {
    background: var(--color-background-200);
    color: var(--color-text);
  }

  .badge[data-variant="success"] {
    background: color-mix(in srgb, var(--color-accent) 15%, transparent);
    color: var(--color-accent);
  }

  .badge[data-variant="warning"] {
    background: color-mix(in srgb, #f59e0b 15%, transparent);
    color: #f59e0b;
  }

  .badge[data-variant="danger"] {
    background: color-mix(in srgb, #ef4444 15%, transparent);
    color: #ef4444;
  }

  .badge[data-variant="info"] {
    background: color-mix(in srgb, #3b82f6 15%, transparent);
    color: #3b82f6;
  }

  .badge[data-variant="neutral"] {
    background: color-mix(in srgb, var(--color-text-muted) 12%, transparent);
    color: var(--color-text-muted);
  }
</style>
