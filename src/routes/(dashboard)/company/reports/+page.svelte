<script lang="ts">
  import { goto } from "$app/navigation";
  import { FileWarning } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import {
    COMPANY_REPORT_STATUS_LABELS,
    formatDate,
  } from "$lib/constants/labels";
  import type { CompanyReportStatus } from "$lib/types";

  let { data } = $props();
  const { reports } = $derived(data);

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
  title="Demandes de correction"
  subtitle="Historique de vos demandes de correction de profil."
>
  <div class="actions-bar">
    <Button
      variant="primary"
      Icon={FileWarning}
      onclick={() => goto("/company/reports/new")}
    >
      Nouvelle demande
    </Button>
  </div>
  {#if reports.length === 0}
    <div class="empty-state">
      <FileWarning size={48} />
      <h2>Aucune demande</h2>
      <p>Vous n'avez soumis aucune demande de correction pour le moment.</p>
    </div>
  {:else}
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Type</th>
            <th>Description</th>
            <th>Statut</th>
            <th>Date</th>
            <th>Reponse admin</th>
          </tr>
        </thead>
        <tbody>
          {#each reports as report (report.id)}
            <tr>
              <td>{report.correction_type}</td>
              <td class="desc-cell">{report.description}</td>
              <td>
                <Badge
                  variant={variantMap[report.status] ?? "neutral"}
                  label={COMPANY_REPORT_STATUS_LABELS[
                    report.status as CompanyReportStatus
                  ] ?? report.status}
                />
              </td>
              <td class="date-cell">
                <div>{formatDate(report.created_at)}</div>
                {#if report.resolved_at}
                  <div class="resolved-date">
                    {report.status === "resolu" ? "Resolu le " : "Rejete le "}
                    {formatDate(report.resolved_at)}
                  </div>
                {/if}
              </td>
              <td class="desc-cell">{report.requested_value ?? "-"}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Page>

<style>
  .actions-bar {
    display: flex;
    justify-content: flex-end;
    margin-bottom: var(--spacing-md);
  }

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
</style>
