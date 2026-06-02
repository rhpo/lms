<script lang="ts">
    import { fly, scale } from "svelte/transition";
    import { ChevronDown } from "lucide-svelte";
    import type { Component } from "svelte";

    import Avatar from "./Avatar.svelte";

    interface Link {
        href?: string;
        label: string;
        icon: any;
        count?: number;
        action?: () => void;
    }

    interface Props {
        user: { full_name: string; role: string; avatar_url?: string | null };
        links: Link[];
    }

    let { user, links = $bindable() }: Props = $props();

    let totalCount = $derived(
        links.reduce((sum, link) => sum + (link.count ?? 0), 0),
    );

    let open = $state(false);
    let menuEl: HTMLDivElement | undefined = $state();

    function toggle() {
        open = !open;
    }

    function close() {
        open = false;
    }

    $effect(() => {
        if (!open) return;
        function handleClick(e: MouseEvent) {
            if (menuEl && !menuEl.contains(e.target as Node)) {
                close();
            }
        }
        function handleKey(e: KeyboardEvent) {
            if (e.key === "Escape") close();
        }
        requestAnimationFrame(() => {
            document.addEventListener("click", handleClick);
            document.addEventListener("keydown", handleKey);
        });
        return () => {
            document.removeEventListener("click", handleClick);
            document.removeEventListener("keydown", handleKey);
        };
    });
</script>

<div class="profile-wrapper">
    <!-- svelte-ignore a11y_no_noninteractive_tabindex -->
    <button
        class="trigger"
        onclick={toggle}
        aria-haspopup="true"
        aria-expanded={open}
    >
        <Avatar {user} size={32} />
        <span class="trigger-name">{user.full_name}</span>
        {#if totalCount > 0}
            <span class="trigger-badge"
                >{totalCount > 99 ? "99+" : totalCount}</span
            >
        {/if}
        <div class="chevron {open ? 'open' : ''}">
            <ChevronDown size={14} />
        </div>
    </button>

    {#if open}
        <div
            class="dropdown"
            bind:this={menuEl}
            transition:fly={{ y: -8, duration: 150, opacity: 0 }}
        >
            <div class="dropdown-header">
                <span class="dropdown-name">{user.full_name}</span>
                <span class="dropdown-role">{user.role}</span>
            </div>
            <div class="dropdown-divider"></div>
            {#each links as link}
                {#if link.action}
                    <button
                        class="dropdown-item"
                        onclick={() => {
                            close();
                            link.action!();
                        }}
                    >
                        <link.icon size={16} />
                        <span class="dropdown-item-label">{link.label}</span>
                    </button>
                {:else}
                    <a href={link.href} class="dropdown-item" onclick={close}>
                        <link.icon size={16} />
                        <span class="dropdown-item-label">{link.label}</span>
                        {#if link.count && link.count > 0}
                            <span class="dropdown-item-badge"
                                >{link.count > 99 ? "99+" : link.count}</span
                            >
                        {/if}
                    </a>
                {/if}
            {/each}
        </div>
    {/if}
</div>

<style>
    .profile-wrapper {
        position: relative;
    }

    .trigger {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: rgba(255, 255, 255, 0.1);
        border: 1px solid rgba(255, 255, 255, 0.15);
        border-radius: 999px;
        padding: 0.25rem 0.75rem 0.25rem 0.25rem;
        cursor: pointer;
        color: #fff;
        font-family: var(--font-sans);
        font-size: var(--text-sm);
        transition:
            background 0.15s,
            transform 0.15s,
            box-shadow 0.15s;

        &:hover {
            background: rgba(255, 255, 255, 0.18);
            transform: scale(1.03);
            box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
        }

        &:active {
            transform: scale(0.95);
        }
    }

    .trigger-name {
        font-weight: 500;
        white-space: nowrap;
    }

    .trigger-badge {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 18px;
        height: 18px;
        padding: 0 5px;
        border-radius: 999px;
        background: #ef4444;
        color: #fff;
        font-size: 11px;
        font-weight: 700;
        line-height: 1;
    }

    .dropdown-item-label {
        flex: 1;
    }

    .dropdown-item-badge {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 20px;
        height: 20px;
        padding: 0 6px;
        border-radius: 999px;
        background: #ef4444;
        color: #fff;
        font-size: 11px;
        font-weight: 700;
        line-height: 1;
    }

    .chevron {
        transition: transform 0.2s;
        &.open {
            transform: rotate(180deg);
        }
    }

    .dropdown {
        position: absolute;
        right: 0;
        top: calc(100% + 0.5rem);
        min-width: 200px;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
        overflow: hidden;
        z-index: 100;
    }

    .dropdown-header {
        padding: 0.75rem 1rem;
        display: flex;
        flex-direction: column;
        gap: 0.15rem;
    }

    .dropdown-name {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
    }

    .dropdown-role {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .dropdown-divider {
        height: 1px;
        background: var(--color-border);
        margin: 0;
    }

    .dropdown-item {
        display: flex;
        align-items: center;
        gap: 0.65rem;
        padding: 0.65rem 1rem;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        text-decoration: none;
        background: none;
        border: none;
        width: 100%;
        text-align: left;
        cursor: pointer;
        transition:
            background 0.12s,
            transform 0.12s,
            padding-left 0.12s;

        margin: 0.25rem;
        border-radius: 0.5rem;

        &,
        & :global(*) {
            color: var(--color-text);
        }

        &:hover {
            background: var(--color-background-200);
            padding-left: 1.25rem;
        }

        &:active {
            transform: scale(0.97);
        }
    }
</style>
