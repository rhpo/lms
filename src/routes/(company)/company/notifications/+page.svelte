<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { Bell, CheckCheck, Mail, MailOpen } from "lucide-svelte";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { notifications } = $derived(data);

    const notificationTypeLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            affectation: "Affectation",
            validation: "Validation",
            soutenance: "Soutenance",
        };
        return labels;
    });

    const notificationTypeVariant = $derived.by(() => {
        const variants: Record<string, "info" | "success" | "warning"> = {
            affectation: "info",
            validation: "success",
            soutenance: "warning",
        };
        return variants;
    });

    const unreadCount = $derived(
        notifications.filter((n: { read_at: string | null }) => !n.read_at).length,
    );

    async function markReadAction(id: string) {
        try {
            const res = await fetch("/api/subjects", {
                method: "PATCH",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ notification_id: id }),
            });
            if (res.ok) {
                await invalidateAll();
            }
        } catch {
            // silent
        }
    }

    async function markAllReadAction() {
        try {
            const res = await fetch("/api/subjects", {
                method: "PATCH",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ all: true }),
            });
            if (res.ok) {
                await invalidateAll();
            }
        } catch {
            // silent
        }
    }
</script>

<Page
    title="Notifications"
    subtitle="Consultez toutes vos notifications."
>
    <div class="notif-header">
        <div class="notif-count">
            <Bell size={16} />
            <span>
                {notifications.length} notification{notifications.length !== 1 ? "s" : ""}
                {#if unreadCount > 0}
                    ({unreadCount} non lue{unreadCount !== 1 ? "s" : ""})
                {/if}
            </span>
        </div>

        {#if unreadCount > 0}
            <Button variant="ghost" size="sm" onclick={markAllReadAction}>
                <CheckCheck size={14} />
                Tout marquer comme lu
            </Button>
        {/if}
    </div>

    {#if notifications.length === 0}
        <p class="empty">
            Aucune notification pour le moment.
        </p>
    {:else}
        <div class="notif-list">
            {#each notifications as notif}
                <div
                    class="notif-card"
                    class:is-unread={!notif.read_at}
                    onclick={() => {
                        if (!notif.read_at) markReadAction(notif.id);
                    }}
                    role="button"
                    tabindex="0"
                >
                    <div class="notif-icon">
                        {#if notif.read_at}
                            <MailOpen size={16} />
                        {:else}
                            <Mail size={16} />
                        {/if}
                    </div>

                    <div class="notif-content">
                        <div class="notif-top">
                            <Badge
                                variant={notificationTypeVariant[notif.type] ?? "info"}
                                label={notificationTypeLabel[notif.type] ?? notif.type}
                            />
                            <span class="notif-date">
                                {new Date(
                                    notif.created_at,
                                ).toLocaleDateString("fr-FR", {
                                    day: "numeric",
                                    month: "short",
                                    hour: "2-digit",
                                    minute: "2-digit",
                                })}
                            </span>
                        </div>
                        <p class="notif-message">
                            {notif.message ?? "Notification"}
                        </p>
                    </div>

                    {#if !notif.read_at}
                        <Button
                            variant="ghost"
                            size="sm"
                            onclick={(e: Event) => {
                                e.stopPropagation();
                                markReadAction(notif.id);
                            }}
                        >
                            <CheckCheck size={14} />
                        </Button>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</Page>

<style>
    .notif-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: var(--spacing-sm);
        margin-bottom: var(--spacing-lg);
        flex-wrap: wrap;
    }

    .notif-count {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        font-style: italic;
        font-size: var(--text-sm);
        padding: 2rem;
        font-family: var(--font-sans);
    }

    .notif-list {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xs);
    }

    .notif-card {
        display: flex;
        align-items: flex-start;
        gap: var(--spacing-md);
        padding: var(--spacing-md);
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        cursor: pointer;
        transition: background var(--transition-fast);

        &:hover {
            background: var(--color-background-100);
        }
    }

    .notif-card.is-unread {
        border-left: 3px solid var(--color-accent);
        background: color-mix(
            in srgb,
            var(--color-accent) 3%,
            var(--color-surface)
        );
    }

    .notif-icon {
        flex-shrink: 0;
        color: var(--color-text-muted);
        padding-top: 0.15rem;
    }

    .notif-content {
        flex: 1;
        min-width: 0;
    }

    .notif-top {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        margin-bottom: 0.35rem;
        flex-wrap: wrap;
    }

    .notif-date {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    .notif-message {
        font-size: var(--text-sm);
        color: var(--color-text);
        font-family: var(--font-sans);
        margin: 0;
        line-height: 1.5;
    }
</style>
