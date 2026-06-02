<script lang="ts">
    import { Eye } from "lucide-svelte";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { supervisedPfes: pfes } = $derived(data);

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
        const variants: Record<string, "info" | "warning" | "success" | "danger"> = {
            en_cours: "info",
            soutenance_planifiee: "warning",
            valide: "success",
            refuse: "danger",
        };
        return variants;
    });
</script>

<Page
    title="Encadrements"
    subtitle="Consultez les PFE que vous encadrez."
>
    {#if pfes.length === 0}
        <p class="empty">
            Aucun encadrement pour le moment.
        </p>
    {:else}
        <table>
            <thead>
                <tr>
                    <th>Sujet</th>
                    <th>Etudiant(s)</th>
                    <th>Co-promoteur</th>
                    <th>Statut</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each pfes as pfe}
                    <tr>
                        <td>
                            {pfe.subject?.title ?? pfe.subject_title ?? "PFE #" + pfe.id}
                        </td>
                        <td>
                            {pfe.student?.profile?.full_name ?? pfe.student_id}
                            {#if pfe.student2_id}
                                , {pfe.student2?.profile?.full_name ?? pfe.student2_id}
                            {/if}
                            {#if pfe.student3_id}
                                , {pfe.student3?.profile?.full_name ?? pfe.student3_id}
                            {/if}
                        </td>
                        <td>
                            {pfe.co_supervisor?.profile?.full_name ?? "-"}
                        </td>
                        <td>
                            <Badge
                                variant={pfeStatusVariant[pfe.status] ?? "info"}
                                label={pfeStatusLabel[pfe.status] ?? pfe.status}
                            />
                        </td>
                        <td>
                            <Button
                                variant="ghost"
                                size="sm"
                                href="/company/supervised-pfes/{pfe.id}"
                            >
                                <Eye size={14} />
                                Detail
                            </Button>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {/if}
</Page>

<style>
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
</style>
