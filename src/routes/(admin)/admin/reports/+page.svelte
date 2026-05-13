<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { FileWarning, CheckCircle, XCircle } from "lucide-svelte";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();
    const { reports } = $derived(data);

    async function resolveReport(reportId: string) {
        try {
            const res = await fetch(`/api/reports/${reportId}/resolve`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ status: "resolved" }),
            });
            if (res.ok) await invalidateAll();
        } catch {
            // ignore
        }
    }

    async function rejectReport(reportId: string) {
        try {
            const res = await fetch(`/api/reports/${reportId}/resolve`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ status: "rejected" }),
            });
            if (res.ok) await invalidateAll();
        } catch {
            // ignore
        }
    }

    const variantMap: Record<
        string,
        "success" | "warning" | "danger" | "neutral"
    > = {
        resolved: "success",
        pending: "warning",
        rejected: "danger",
    };
</script>

<Page
    title="Reports"
    subtitle="Demandes de correction soumises par les entreprises."
>
    {#if reports.length === 0}
        <div class="empty-state">
            <CheckCircle size={48} />
            <h2>Aucun report</h2>
            <p>Tous les reports ont ete traites. Rien en attente.</p>
        </div>
    {:else}
        <div class="table-container">
            <table>
                <thead>
                    <tr>
                        <th>Entreprise</th>
                        <th>Demandeur</th>
                        <th>Type</th>
                        <th>Description</th>
                        <th>Statut</th>
                        <th>Date</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each reports as report (report.id)}
                        <tr>
                            <td class="company-cell">{report.company_name}</td>
                            <td>{report.reporter_name}</td>
                            <td>{report.report_type_label}</td>
                            <td class="desc-cell">{report.description}</td>
                            <td>
                                <Badge
                                    variant={variantMap[report.status] ??
                                        "neutral"}
                                    label={report.status_label}
                                />
                            </td>
                            <td class="date-cell">
                                <div>{report.created_at_formatted}</div>
                                {#if report.resolved_at_formatted && report.resolved_at_formatted !== "—"}
                                    <div class="resolved-date">
                                        {report.status === "resolved"
                                            ? "Resolu le "
                                            : "Rejete le "}
                                        {report.resolved_at_formatted}
                                    </div>
                                {/if}
                            </td>
                            <td>
                                {#if report.status === "pending"}
                                    <div class="actions">
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            Icon={CheckCircle}
                                            onclick={() =>
                                                resolveReport(report.id)}
                                            title="Marquer comme resolu"
                                        />
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            Icon={XCircle}
                                            onclick={() =>
                                                rejectReport(report.id)}
                                            title="Rejeter"
                                        />
                                    </div>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</Page>

<style>
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 1rem;
        padding: 3rem;
        text-align: center;
        color: var(--color-text-muted);

        h2 {
            font-size: var(--text-lg);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
        }

        p {
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            margin: 0;
        }
    }

    .table-container {
        overflow-x: auto;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
    }

    table {
        width: 100%;
        border-collapse: collapse;

        thead {
            th {
                padding: var(--spacing-sm) var(--spacing-md);
                font-size: var(--text-xs);
                font-weight: 600;
                text-transform: uppercase;
                letter-spacing: 0.05em;
                color: var(--color-text-muted);
                font-family: var(--font-sans);
                text-align: left;
                border-bottom: 1px solid var(--color-border);
                background: var(--color-background);
                white-space: nowrap;
            }
        }

        tbody {
            tr {
                border-bottom: 1px solid var(--color-border);
                transition: background var(--transition-fast);

                &:last-child {
                    border-bottom: none;
                }

                &:hover {
                    background: var(--color-background-100);
                }
            }

            td {
                padding: var(--spacing-sm) var(--spacing-md);
                font-size: var(--text-sm);
                font-family: var(--font-sans);
                color: var(--color-text);
                vertical-align: middle;
            }
        }
    }

    .company-cell {
        font-weight: 600;
    }

    .desc-cell {
        max-width: 280px;
        white-space: normal;
        line-height: 1.4;
    }

    .date-cell {
        font-size: var(--text-xs);
    }

    .resolved-date {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        margin-top: 0.25rem;
    }

    .actions {
        display: flex;
        gap: var(--spacing-xs);
    }
</style>
