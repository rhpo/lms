<script lang="ts">
    import { FileText, Users, Bell, Plus, Eye } from "lucide-svelte";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { subjects, pfes, notifications } = $derived(data);

    const subjectStatusLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            en_attente: "En attente",
            valide: "Valide",
            refuse: "Refuse",
            expire: "Expire",
        };
        return labels;
    });

    const subjectStatusVariant = $derived.by(() => {
        const variants: Record<
            string,
            "warning" | "success" | "danger" | "info"
        > = {
            en_attente: "warning",
            valide: "success",
            refuse: "danger",
            expire: "info",
        };
        return variants;
    });

    const pfeStatusLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            en_cours: "En cours",
            soutenance_planifiee: "Soutenance planifiee",
            valide: "Valide",
            refuse: "Refuse",
        };
        return labels;
    });

    const pfeStatusVariant = $derived.by(() => {
        const variants: Record<
            string,
            "info" | "warning" | "success" | "danger"
        > = {
            en_cours: "info",
            soutenance_planifiee: "warning",
            valide: "success",
            refuse: "danger",
        };
        return variants;
    });

    const recentNotifications = $derived(notifications.slice(0, 3));
</script>

<Page
    title="Tableau de bord"
    subtitle="Vue d'ensemble de votre activite entreprise."
>
    <div class="stats-grid">
        <div class="stat-card">
            <FileText size={24} />
            <div class="stat-info">
                <span class="stat-value">{subjects.length}</span>
                <span class="stat-label">Sujets proposes</span>
            </div>
        </div>
        <div class="stat-card">
            <Users size={24} />
            <div class="stat-info">
                <span class="stat-value">{pfes.length}</span>
                <span class="stat-label">Etudiants encadres</span>
            </div>
        </div>
        <div class="stat-card">
            <Bell size={24} />
            <div class="stat-info">
                <span class="stat-value">
                    {notifications.filter(
                        (n: { read_at: string | null }) => !n.read_at,
                    ).length}
                </span>
                <span class="stat-label">Notifications non lues</span>
            </div>
        </div>
    </div>

    <div class="section">
        <div class="section-header">
            <h2>Mes sujets</h2>
            <Button variant="primary" href="/company/propose-subject">
                <Plus size={14} />
                Proposer un sujet
            </Button>
        </div>

        {#if subjects.length === 0}
            <p class="empty">Aucun sujet propose pour le moment.</p>
        {:else}
            <div class="subject-list">
                {#each subjects as subject}
                    <div class="subject-row">
                        <div class="subject-info">
                            <span class="subject-title">{subject.title}</span>
                            <Badge
                                variant={subjectStatusVariant[subject.status] ??
                                    "info"}
                                label={subjectStatusLabel[subject.status] ??
                                    subject.status}
                            />
                        </div>
                        <Button
                            variant="ghost"
                            size="sm"
                            Icon={Eye}
                            href="/company/my-subjects"
                        >
                            Voir
                        </Button>
                    </div>
                {/each}
            </div>
        {/if}
    </div>

    <div class="section">
        <div class="section-header">
            <h2>Encadrements en cours</h2>
            <Button variant="ghost" href="/company/supervised-pfes">
                Tout voir
            </Button>
        </div>

        {#if pfes.length === 0}
            <p class="empty">Aucun encadrement pour le moment.</p>
        {:else}
            <div class="pfe-list">
                {#each pfes.slice(0, 5) as pfe}
                    <div class="pfe-row">
                        <div class="pfe-info">
                            <span class="pfe-title">
                                {pfe.subject?.title ?? pfe.subject_title ?? "-"}
                            </span>
                            <Badge
                                variant={pfeStatusVariant[pfe.status] ?? "info"}
                                label={pfeStatusLabel[pfe.status] ?? pfe.status}
                            />
                        </div>
                        <Button
                            variant="ghost"
                            size="sm"
                            href="/company/supervised-pfes/{pfe.id}"
                        >
                            <Eye size={14} />
                            Detail
                        </Button>
                    </div>
                {/each}
            </div>
        {/if}
    </div>

    {#if recentNotifications.length > 0}
        <div class="section">
            <div class="section-header">
                <h2>Notifications recentes</h2>
                <Button variant="ghost" href="/company/notifications">
                    Tout voir
                </Button>
            </div>
            <div class="notif-list">
                {#each recentNotifications as notif}
                    <div class="notif-row" class:is-unread={!notif.read_at}>
                        <span class="notif-message">
                            {notif.message ?? "Notification"}
                        </span>
                        <span class="notif-date">
                            {new Date(notif.created_at).toLocaleDateString(
                                "fr-FR",
                            )}
                        </span>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</Page>

<style>
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: var(--spacing-md);
        margin-bottom: var(--spacing-lg);
    }

    .stat-card {
        display: flex;
        align-items: center;
        gap: var(--spacing-md);
        padding: var(--spacing-lg);
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        color: var(--color-accent);
    }

    .stat-info {
        display: flex;
        flex-direction: column;
    }

    .stat-value {
        font-size: 1.5rem;
        font-weight: 700;
        font-family: var(--font-sans);
        color: var(--color-text);
    }

    .stat-label {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    .section {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        margin-bottom: var(--spacing-lg);
    }

    .section-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: var(--spacing-md);

        h2 {
            font-size: var(--text-lg);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
        }
    }

    .empty {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-style: italic;
        font-family: var(--font-sans);
        text-align: center;
        padding: 1rem;
    }

    .subject-list,
    .pfe-list,
    .notif-list {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xs);
    }

    .subject-row,
    .pfe-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: var(--spacing-sm);
        padding: var(--spacing-sm) var(--spacing-md);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-background);
    }

    .subject-info,
    .pfe-info {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        flex: 1;
        min-width: 0;
    }

    .subject-title,
    .pfe-title {
        font-size: var(--text-sm);
        font-weight: 500;
        color: var(--color-text);
        font-family: var(--font-sans);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .notif-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: var(--spacing-sm);
        padding: var(--spacing-sm) var(--spacing-md);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-background);
    }

    .notif-row.is-unread {
        border-left: 3px solid var(--color-accent);
    }

    .notif-message {
        font-size: var(--text-sm);
        color: var(--color-text);
        font-family: var(--font-sans);
        flex: 1;
        min-width: 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .notif-date {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        white-space: nowrap;
    }
</style>
