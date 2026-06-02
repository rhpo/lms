<script lang="ts">
    import type { IconComponent } from "$lib/types/domain";

    interface Props {
        title: string;
        value: string | number;
        badge?: string;
        badgeVariant?: "green" | "orange" | "purple" | "blue";
        icon: IconComponent;
        barColor?: string;
    }

    let {
        title,
        value,
        badge,
        badgeVariant = "blue",
        icon: Icon,
        barColor = "#3b82f6",
    }: Props = $props();
</script>

<main data-aos="fade-up">
    <div class="stat-card">
        <div class="card-head">
            <div class="card-icon">
                <Icon size={20} />
            </div>
            {#if badge}
                <span class="card-badge {badgeVariant}">{badge}</span>
            {/if}
        </div>
        <div class="card-body">
            <p class="card-label">{title}</p>
            <p class="card-value">{value}</p>
        </div>
        <div class="card-bar" style:background={barColor}></div>
    </div>
</main>

<style>
    .stat-card {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: 1.25rem 1.25rem 1rem;
        display: flex;
        flex-direction: column;
        gap: 0.85rem;
        flex: 1;
        min-width: 0;
        box-shadow: var(--shadow-sm);
        transition: box-shadow var(--transition-normal);
    }
    .stat-card:hover {
        box-shadow: var(--shadow-md);
    }

    .card-head {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
    .card-icon {
        width: 38px;
        height: 38px;
        border-radius: 8px;
        background: var(--color-background-100);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-text-muted);
    }
    .card-badge {
        font-size: 0.72rem;
        font-weight: 600;
        font-family: var(--font-sans);
        padding: 0.2rem 0.55rem;
        border-radius: 100px;
    }
    .card-badge.green {
        background: color-mix(in srgb, #10b981 15%, var(--color-surface));
        color: #065f46;
    }
    .card-badge.orange {
        background: color-mix(in srgb, #f59e0b 15%, var(--color-surface));
        color: #c2410c;
    }
    .card-badge.purple {
        background: color-mix(in srgb, #8b5cf6 15%, var(--color-surface));
        color: #7c3aed;
    }
    .card-badge.blue {
        background: color-mix(in srgb, #3b82f6 15%, var(--color-surface));
        color: #1d4ed8;
    }

    :global(html.dark) .card-badge.green {
        color: #6ee7b7;
    }
    :global(html.dark) .card-badge.orange {
        color: #fcd34d;
    }
    :global(html.dark) .card-badge.purple {
        color: #c4b5fd;
    }
    :global(html.dark) .card-badge.blue {
        color: #93c5fd;
    }

    .card-label {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        margin: 0 0 0.35rem;
    }
    .card-value {
        font-size: 1.9rem;
        font-weight: 700;
        color: var(--color-text);
        font-family: var(--font-sans);
        margin: 0;
        line-height: 1;
    }
    .card-bar {
        height: 3px;
        border-radius: 100px;
        width: 45%;
    }
</style>
