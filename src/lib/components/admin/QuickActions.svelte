<script lang="ts">
    import { CheckSquare, Users, Clock, ChevronRight } from "lucide-svelte";
    import type { Component } from "svelte";

    interface Action {
        icon: typeof CheckSquare;
        iconColor: string;
        title: string;
        description: string;
        href: string;
    }

    const actions: Action[] = [
        {
            icon: CheckSquare,
            iconColor: "#2563eb",
            title: "Validate Subjects",
            description: "Review 12 pending proposals",
            href: "/admin/subjects",
        },
        {
            icon: Users,
            iconColor: "#d97706",
            title: "Manage Users",
            description: "Add or edit student accounts",
            href: "/admin/users",
        },
        {
            icon: Clock,
            iconColor: "#7c3aed",
            title: "Set Deadlines",
            description: "Update PFE cycle timeline",
            href: "/admin/settings",
        },
    ];
</script>

<div class="qa-panel">
    <h2 class="qa-heading">Quick Actions</h2>
    <div class="qa-list">
        {#each actions as action}
            <a href={action.href} class="qa-row">
                <div class="qa-icon" style:--icon-color={action.iconColor}>
                    <action.icon size={18} color={action.iconColor} />
                </div>
                <div class="qa-text">
                    <p class="qa-title">{action.title}</p>
                    <p class="qa-desc">{action.description}</p>
                </div>
                <ChevronRight size={16} color="var(--color-text-muted)" />
            </a>
        {/each}
    </div>

    <div class="insight-card">
        <p class="insight-heading">Efficiency Insight</p>
        <p class="insight-body">
            Your validation time has improved by <strong>2.4 days</strong> compared
            to last semester.
        </p>
    </div>
</div>

<style>
    .qa-panel {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        box-shadow: var(--shadow-sm);
        display: flex;
        flex-direction: column;
        height: 100%;
    }
    .qa-heading {
        font-size: var(--text-lg);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text);
        padding: 1.25rem 1.5rem 1rem;
        margin: 0;
        border-bottom: 1px solid var(--color-border);
    }
    .qa-list {
        display: flex;
        flex-direction: column;
        padding: 0.5rem 0;
    }
    .qa-row {
        display: flex;
        align-items: center;
        gap: 0.85rem;
        padding: 0.9rem 1.25rem;
        text-decoration: none;
        transition: background var(--transition-fast);
        cursor: pointer;
    }
    .qa-row:hover {
        background: var(--color-background);
        text-decoration: none;
    }

    .qa-icon {
        width: 36px;
        height: 36px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        /* bg is derived from the icon color, adapts automatically in dark mode */
        background: color-mix(
            in srgb,
            var(--icon-color) 12%,
            var(--color-surface)
        );
    }
    .qa-text {
        flex: 1;
        min-width: 0;
    }
    .qa-title {
        font-size: var(--text-sm);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text);
        margin: 0 0 0.15rem;
    }
    .qa-desc {
        font-size: var(--text-xs);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin: 0;
    }
    .insight-card {
        margin: auto 1.25rem 1.25rem;
        background: color-mix(
            in srgb,
            var(--color-accent) 10%,
            var(--color-background-100)
        );
        border: 1px solid var(--color-border);
        border-radius: 10px;
        padding: 1.25rem;
        min-height: 100px;
    }
    .insight-heading {
        font-size: var(--text-md);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text);
        margin: 0 0 0.5rem;
    }
    .insight-body {
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin: 0;
        line-height: 1.5;
    }
    .insight-body strong {
        color: var(--color-text);
    }
</style>
