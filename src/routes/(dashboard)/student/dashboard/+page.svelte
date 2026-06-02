<script lang="ts">
    import { goto } from "$app/navigation";
    import { Heart, FileText, Bell } from "lucide-svelte";
    import { formatDate, PFE_ASSIGNMENT_STATUS_LABELS, PFE_ASSIGNMENT_STATUS_VARIANTS } from "$lib/constants/labels";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { currentPfe, wishes, notifications, yearId } = $derived(data);

    const wishCount = $derived(wishes.length);
    const unreadCount = $derived(
        notifications.filter((n) => !n.read_at).length,
    );

    const pfeStatusLabel = $derived.by(() => {
        if (!currentPfe) return "Pas encore affecte";
        return PFE_ASSIGNMENT_STATUS_LABELS[currentPfe.status] ?? "En cours";
    });

    const pfeStatusVariant = $derived.by<
        "success" | "danger" | "warning" | "info" | "neutral"
    >(() => {
        if (!currentPfe) return "neutral";
        return PFE_ASSIGNMENT_STATUS_VARIANTS[currentPfe.status] ?? "neutral";
    });
</script>

<Page title="Tableau de bord" subtitle="Bienvenue sur votre espace etudiant.">
    <div class="stat-grid">
        <div
            class="stat-card"
            onclick={() => goto("/student/my-pfe")}
            role="button"
            tabindex="0"
        >
            <div class="stat-icon">
                <FileText size={20} />
            </div>
            <p class="stat-label">Statut PFE</p>
            <Badge variant={pfeStatusVariant} label={pfeStatusLabel} />
        </div>

        <div
            class="stat-card"
            onclick={() => goto("/student/voeux")}
            role="button"
            tabindex="0"
        >
            <div class="stat-icon">
                <Heart size={20} />
            </div>
            <p class="stat-label">Voeux soumis</p>
            <p class="stat-value">{wishCount}</p>
        </div>

        <div
            class="stat-card"
            onclick={() => goto("/student/notifications")}
            role="button"
            tabindex="0"
        >
            <div class="stat-icon">
                <Bell size={20} />
            </div>
            <p class="stat-label">Notifications non lues</p>
            <p class="stat-value">{unreadCount}</p>
        </div>
    </div>

    <div class="dashboard-grid">
        <section>
            <h2>Notifications</h2>
            {#if unreadCount > 0}
                <div class="unread-banner">
                    <Bell size={16} />
                    <span>{unreadCount} notification(s) non lue(s)</span>
                </div>
            {:else}
                <p class="empty">Aucune notification non lue</p>
            {/if}
            <div class="section-actions">
                <Button
                    variant="ghost"
                    onclick={() => goto("/student/notifications")}
                >
                    Voir toutes les notifications
                </Button>
            </div>
        </section>

        <section>
            <h2>Dernieres notifications</h2>
            {#if notifications.length > 0}
                <div class="list">
                    {#each notifications.slice(0, 3) as notif}
                        <div class="list-item">
                            <div class="list-item-icon">
                                <Bell size={16} />
                            </div>
                            <div>
                                <p class="list-item-title">{notif.message ?? "Notification"}</p>
                                <p class="list-item-subtitle">
                                    {formatDate(notif.created_at)}
                                </p>
                            </div>
                        </div>
                    {/each}
                </div>
            {:else}
                <p class="empty">Aucune notification</p>
            {/if}
        </section>
    </div>

    <div class="quick-actions">
        <h2>Raccourcis</h2>
        <div class="actions-grid">
            <Button
                variant="primary"
                onclick={() => goto("/student/catalogue")}
            >
                Consulter le catalogue
            </Button>
            <Button variant="primary" onclick={() => goto("/student/voeux")}>
                Mes voeux
            </Button>
            <Button variant="primary" onclick={() => goto("/student/my-pfe")}>
                Mon PFE
            </Button>
        </div>
    </div>
</Page>

<style>
    .stat-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 1rem;
        margin-bottom: var(--spacing-xl);

        @media screen and (max-width: 768px) {
            & {
                grid-template-columns: repeat(2, 1fr);
            }
        }

        @media screen and (max-width: 480px) {
            & {
                grid-template-columns: 1fr;
            }
        }
    }

    .stat-card {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: 1.25rem;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        cursor: pointer;
        transition: box-shadow var(--transition-normal);

        &:hover {
            box-shadow: var(--shadow-md);
        }
    }

    .stat-icon {
        width: 38px;
        height: 38px;
        border-radius: 8px;
        background: var(--color-background);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-text-muted);
    }

    .stat-label {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        margin: 0;
    }

    .stat-value {
        font-size: 1.9rem;
        font-weight: 700;
        color: var(--color-text);
        font-family: var(--font-sans);
        margin: 0;
        line-height: 1;
    }

    .dashboard-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: var(--spacing-lg);
        margin-bottom: var(--spacing-xl);

        @media screen and (max-width: 768px) {
            & {
                grid-template-columns: 1fr;
            }
        }
    }

    section {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
    }

    h2 {
        font-size: var(--text-lg);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text);
        margin: 0 0 var(--spacing-md);
    }

    .list {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .list-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem;
        border-radius: 8px;
        cursor: pointer;
        transition: background var(--transition-fast);

        &:hover {
            background: var(--color-background);
        }
    }

    .list-item-icon {
        width: 32px;
        height: 32px;
        border-radius: 6px;
        background: var(--color-background);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-accent);
        flex-shrink: 0;
    }

    .list-item-title {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
        margin: 0;
    }

    .list-item-subtitle {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        margin: 0.15rem 0 0;
    }

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        font-style: italic;
        font-size: var(--text-sm);
        padding: 1rem;
    }

    .unread-banner {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem;
        background: color-mix(
            in srgb,
            var(--color-accent) 10%,
            var(--color-surface)
        );
        border-radius: 8px;
        color: var(--color-accent);
        font-size: var(--text-sm);
        font-weight: 600;
    }

    .section-actions {
        margin-top: var(--spacing-md);
    }

    .quick-actions {
        margin-bottom: var(--spacing-xl);
    }

    .actions-grid {
        display: flex;
        flex-wrap: wrap;
        gap: var(--spacing-sm);
    }
</style>
