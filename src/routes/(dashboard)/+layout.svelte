<script lang="ts">
    import "$lib/styles/app.css";
    import { onMount, onDestroy } from "svelte";
    import { notificationStore } from "$lib/stores/notifications";
    import { atomic } from "$lib/stores/atomic.svelte";

    let { children } = $props();

    onMount(async () => {
        try {
            await atomic.load();
        } catch {

        }
        notificationStore.startPolling(30_000);
    });

    onDestroy(() => {
        notificationStore.stopPolling();
    });
</script>

{#if atomic.error}
    <div class="fatal-gate">
        <div class="fatal-card">
            <span class="fatal-icon">⚠</span>
            <h1>Erreur critique</h1>
            <p>{atomic.error}</p>
            <button onclick={() => location.reload()}>Réessayer</button>
        </div>
    </div>
{:else if atomic.ready}
    {@render children()}
{:else}
    <div class="fatal-gate">
        <div class="spinner"></div>
        <p>Chargement…</p>
    </div>
{/if}

<style>
    :global(*) {
        transition: all var(--transition-duration);
    }

    .fatal-gate {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 100vh;
        gap: 1rem;
        font-family: var(--font-sans, system-ui, sans-serif);
        color: var(--text-primary, #1a1a2e);
        background: var(--bg-primary, #f8f9fc);
    }

    .fatal-card {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.75rem;
        padding: 2.5rem;
        border-radius: 1rem;
        background: var(--bg-surface, #fff);
        box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
        max-width: 420px;
        text-align: center;
    }

    .fatal-icon {
        font-size: 2.5rem;
    }

    .fatal-card h1 {
        font-size: 1.25rem;
        font-weight: 700;
        margin: 0;
    }

    .fatal-card p {
        font-size: 0.875rem;
        color: var(--text-secondary, #64748b);
        margin: 0;
    }

    .fatal-card button {
        margin-top: 0.5rem;
        padding: 0.5rem 1.5rem;
        border: none;
        border-radius: 0.5rem;
        background: var(--accent, #6366f1);
        color: #fff;
        font-weight: 600;
        cursor: pointer;
    }

    .fatal-card button:hover {
        opacity: 0.9;
    }

    .spinner {
        width: 2rem;
        height: 2rem;
        border: 3px solid var(--border, #e2e8f0);
        border-top-color: var(--accent, #6366f1);
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }
</style>
