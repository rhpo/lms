<script lang="ts">
    import { Eye, Plus, RefreshCw, Pencil } from "lucide-svelte";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { subjects } = $derived(data);

    const subjectStatusLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            en_attente: "En attente",
            valide: "Validé",
            accepte_sous_reserve: "Sous réserve",
            refuse: "Refusé",
            expire: "Expiré",
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
            accepte_sous_reserve: "warning",
            refuse: "danger",
            expire: "info",
        };
        return variants;
    });
</script>

<Page
    title="Mes sujets"
    subtitle="Consultez et Gérer les sujets que vous avez proposes."
>
    <div class="page-actions">
        <Button variant="primary" href="/company/propose-subject">
            <Plus size={14} />
            Proposer un sujet
        </Button>
    </div>

    {#if subjects.length === 0}
        <p class="empty">Vous n'avez propose aucun sujet pour le moment.</p>
    {:else}
        <table>
            <thead>
                <tr>
                    <th>Titre</th>
                    <th>Domaines</th>
                    <th>Groupe</th>
                    <th>Statut</th>
                    <th>Date</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each subjects as subject}
                    <tr>
                        <td>
                            <a
                                href="/company/my-subjects/{subject.id}/candidats"
                            >
                                {subject.title}
                            </a>
                        </td>
                        <td
                            >{subject.domains?.map((d) => d.name).join(", ") ??
                                "-"}</td
                        >
                        <td>{subject.group_type}</td>
                        <td>
                            <Badge
                                variant={subjectStatusVariant[subject.status] ??
                                    "info"}
                                label={subjectStatusLabel[subject.status] ??
                                    subject.status}
                            />
                        </td>
                        <td>
                            {new Date(subject.created_at).toLocaleDateString(
                                "fr-FR",
                            )}
                        </td>
                        <td>
                            <div class="actions-cell">
                                {#if subject.status === "valide"}
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        href="/company/my-subjects/{subject.id}/candidats"
                                    >
                                        <Eye size={14} />
                                        Candidats
                                    </Button>
                                {/if}
                                {#if subject.status === "accepte_sous_reserve"}
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        href="/company/my-subjects/{subject.id}/edit"
                                    >
                                        <Pencil size={14} />
                                        Modifier
                                    </Button>
                                {/if}
                                {#if subject.status === "refuse"}
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        href="/company/propose-subject"
                                    >
                                        <RefreshCw size={14} />
                                        Proposer un nouveau
                                    </Button>
                                {/if}
                            </div>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {/if}
</Page>

<style>
    .page-actions {
        display: flex;
        justify-content: flex-end;
        margin-bottom: var(--spacing-lg);
    }

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        font-style: italic;
        font-size: var(--text-sm);
        padding: 2rem;
        font-family: var(--font-sans);
    }

    table {
        width: 100%;
        border-collapse: collapse;
        font-family: var(--font-sans);
        font-size: var(--text-sm);
    }

    thead th {
        text-align: left;
        padding: var(--spacing-sm) var(--spacing-md);
        font-weight: 600;
        color: var(--color-text-muted);
        border-bottom: 2px solid var(--color-border);
        white-space: nowrap;
    }

    tbody tr {
        border-bottom: 1px solid var(--color-border);
    }

    tbody tr:hover {
        background: var(--color-background-100);
    }

    tbody td {
        padding: var(--spacing-sm) var(--spacing-md);
        color: var(--color-text);
    }

    tbody td a {
        color: var(--color-accent);
        text-decoration: none;
        font-weight: 500;
    }

    tbody td a:hover {
        text-decoration: underline;
    }

    .actions-cell {
        display: flex;
        gap: var(--spacing-xs);
    }
</style>
