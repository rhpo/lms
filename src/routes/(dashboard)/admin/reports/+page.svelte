<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { CheckCircle, XCircle } from "lucide-svelte";
    import { admin } from "$lib/api";
    import {
        COMPANY_REPORT_STATUS_LABELS,
        formatDate,
    } from "$lib/constants/labels";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();
    const { reports } = $derived(data);

    async function resolveReport(reportId: number) {
        try {
            await admin.reportAction(reportId, "resolve");
            await invalidateAll();
        } catch {}
    }

    async function rejectReport(reportId: number) {
        try {
            await admin.reportAction(reportId, "reject");
            await invalidateAll();
        } catch {}
    }

    const variantMap: Record<
        string,
        "success" | "warning" | "danger" | "neutral"
    > = {
        en_attente: "warning",
        resolu: "success",
        rejete: "danger",
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
                            <td class="company-cell"
                                >{report.company?.company_name ?? "-"}</td
                            >
                            <td>{report.submitted_by}</td>
                            <td>{report.correction_type}</td>
                            <td class="desc-cell">{report.description}</td>
                            <td>
                                <Badge
                                    variant={variantMap[report.status] ??
                                        "neutral"}
                                    label={COMPANY_REPORT_STATUS_LABELS[
                                        report.status
                                    ] ?? report.status}
                                />
                            </td>
                            <td class="date-cell">
                                <div>{formatDate(report.created_at)}</div>
                                {#if report.resolved_at}
                                    <div class="resolved-date">
                                        {report.status === "resolu"
                                            ? "Résolu le "
                                            : "Rejeté le "}
                                        {formatDate(report.resolved_at)}
                                    </div>
                                {/if}
                            </td>
                            <td>
                                {#if report.status === "en_attente"}
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
