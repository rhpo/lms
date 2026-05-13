<script lang="ts">
    import { goto } from "$app/navigation";
    import { GraduationCap, FileText } from "lucide-svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();
    const { assignments } = $derived(data);

    const statusVariant: Record<
        string,
        "success" | "warning" | "danger" | "info" | "neutral"
    > = {
        en_cours: "info",
        soutenance_planifiee: "warning",
        valide: "success",
        refuse: "danger",
    };
</script>

<Page title="PFE" subtitle="Tous les PFE actifs et archives">
    {#snippet actions()}
        <Button onclick={() => goto("/admin/defenses/new")}>
            <GraduationCap size={16} />
            Planifier une soutenance
        </Button>
    {/snippet}

    {#if assignments && assignments.length > 0}
        <div class="table-wrapper">
            <table>
                <thead>
                    <tr>
                        <th>Code PFE</th>
                        <th>Sujet</th>
                        <th>Étudiants</th>
                        <th>Encadrant</th>
                        <th>Co-encadrant</th>
                        <th>Spécialité</th>
                        <th>Statut</th>
                        <th>Mémoire</th>
                    </tr>
                </thead>
                <tbody>
                    {#each assignments as a}
                        <tr>
                            <td class="code">{a.pfe_code}</td>
                            <td>{a.subject_title}</td>
                            <td>
                                {#if a.students.length > 0}
                                    <ul class="student-list">
                                        {#each a.students as s}
                                            <li>{s.full_name}</li>
                                        {/each}
                                    </ul>
                                {:else}
                                    —
                                {/if}
                            </td>
                            <td>{a.supervisor_name}</td>
                            <td>{a.co_supervisor_name ?? "—"}</td>
                            <td>{a.specialty}</td>
                            <td>
                                <Badge
                                    variant={statusVariant[a.status] ??
                                        "neutral"}
                                    label={a.status_label}
                                />
                            </td>
                            <td>
                                {#if a.memoire_url}
                                    <a
                                        href={a.memoire_url}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                    >
                                        <FileText size={16} />
                                        Voir
                                    </a>
                                {:else}
                                    —
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {:else}
        <p class="empty">Aucun PFE enregistré.</p>
    {/if}
</Page>

<style>
    .table-wrapper {
        overflow-x: auto;
    }

    table {
        width: 100%;
        border-collapse: collapse;

        thead {
            th {
                text-align: left;
                padding: var(--spacing-sm) var(--spacing-md);
                font-size: var(--text-sm);
                font-weight: 600;
                color: var(--color-text-muted);
                border-bottom: 2px solid var(--color-border);
                white-space: nowrap;
            }
        }

        tbody {
            tr {
                &:hover {
                    background: var(--color-background-100);
                }

                td {
                    padding: var(--spacing-sm) var(--spacing-md);
                    font-size: var(--text-sm);
                    border-bottom: 1px solid var(--color-border);
                    vertical-align: top;

                    &.code {
                        font-family: monospace;
                        font-weight: 600;
                    }
                }
            }
        }
    }

    .student-list {
        margin: 0;
        padding: 0;
        list-style: none;

        li {
            font-size: var(--text-sm);
            line-height: 1.5;
        }
    }

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        padding: var(--spacing-xl);
    }
</style>
