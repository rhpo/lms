<script lang="ts">
    import { onMount } from "svelte";
    import { Check, CheckCheck } from "lucide-svelte";
    import { notifications as notifApi } from "$lib/api";
    import { notificationStore } from "$lib/stores/notifications";
    import type { Notification } from "$lib/types";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";
    import { NOTIFICATION_TYPE_LABELS } from "$lib/constants/labels";

    let notifs = $state<Notification[]>([]);
    let loading = $state(true);

    const NOTIFICATION_TYPE_VARIANTS: Record<
        string,
        "info" | "warning" | "success" | "danger" | "neutral"
    > = {
        validation_requise: "warning",
        affectation: "info",
        jury: "neutral",
        disponibilite: "info",
    };

    function parseMessage(notif: Notification): string {
        if (notif.message) return notif.message;
        if (notif.payload) {
            try {
                return JSON.parse(notif.payload).message;
            } catch {
                return notif.payload;
            }
        }
        return notif.type;
    }

    async function loadNotifications() {
        loading = true;

        try {
            notifs = (await notifApi.list()) ?? [];
        } catch {
            notifs = [];
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        loadNotifications();
    });

    async function markRead(id: number) {
        try {
            await notifApi.markRead(id);
            notificationStore.decrement();

            notifs = notifs.map((n) =>
                n.id === id ? { ...n, read_at: new Date().toISOString() } : n,
            );
        } catch {

        }
    }

    async function markAllRead() {
        try {
            await notifApi.markAllRead();
            notificationStore.clear();
            notifs = notifs.map((n) => ({
                ...n,
                read_at: n.read_at ?? new Date().toISOString(),
            }));
        } catch {

        }
    }

    let unreadCount = $derived(notifs.filter((n) => !n.read_at).length);
</script>

<Page
    title="Notifications"
    subtitle="Toutes vos notifications.{unreadCount > 0
        ? ` (${unreadCount} non lue${unreadCount > 1 ? 's' : ''})`
        : ''}"
>
    {#if loading}
        <p class="empty">Chargement...</p>
    {:else if notifs.length === 0}
        <div class="empty">
            <p>Aucune notification.</p>
        </div>
    {:else}
        <div class="toolbar">
            {#if unreadCount > 0}
                <Button variant="ghost" Icon={CheckCheck} onclick={markAllRead}>
                    Tout marquer comme lu
                </Button>
            {/if}
        </div>
        <div class="list">
            {#each notifs as notif (notif.id)}
                <div class="item" class:unread={!notif.read_at}>
                    <div class="item-header">
                        <div class="item-header-left">
                            <Badge
                                variant={NOTIFICATION_TYPE_VARIANTS[
                                    notif.type
                                ] || "info"}
                                label={NOTIFICATION_TYPE_LABELS[notif.type] ||
                                    notif.type}
                            />
                            {#if !notif.read_at}
                                <span class="unread-dot" title="Non lu"></span>
                            {/if}
                        </div>
                        <span class="date">
                            {new Date(notif.created_at).toLocaleDateString(
                                "fr-FR",
                                {
                                    day: "numeric",
                                    month: "short",
                                    year: "numeric",
                                    hour: "2-digit",
                                    minute: "2-digit",
                                },
                            )}
                        </span>
                    </div>
                    <p class="message">{parseMessage(notif)}</p>
                    {#if !notif.read_at}
                        <div class="item-actions">
                            <Button
                                variant="ghost"
                                Icon={Check}
                                onclick={() => markRead(notif.id)}
                            >
                                Marquer comme lu
                            </Button>
                        </div>
                    {:else}
                        <span class="read-label">Lu</span>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</Page>

<style>
    .toolbar {
        margin-bottom: var(--spacing-md);
        display: flex;
        justify-content: flex-end;
    }

    .empty {
        text-align: center;
        padding: 3rem 1rem;
        color: var(--color-text-muted);
    }

    .empty p {
        font-size: var(--text-sm);
        margin: 0;
    }

    .list {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-sm);
    }

    .item {
        padding: var(--spacing-md);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-surface);
        transition: background var(--transition-fast);
    }

    .item.unread {
        background: var(--color-accent-50);
        border-color: var(--color-accent-200);
    }

    .item-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 0.5rem;
    }

    .item-header-left {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .unread-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: var(--color-accent);
        display: inline-block;
        flex-shrink: 0;
    }

    .date {
        font-size: 0.75rem;
        color: var(--color-text-muted);
    }

    .message {
        margin: 0;
        font-size: var(--text-sm);
        color: var(--color-text);
        line-height: 1.5;
    }

    .item-actions {
        margin-top: 0.5rem;
        display: flex;
        justify-content: flex-end;
    }

    .read-label {
        display: block;
        margin-top: 0.5rem;
        font-size: 0.75rem;
        color: var(--color-text-muted);
        text-align: right;
        font-style: italic;
    }
</style>
