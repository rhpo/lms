<script lang="ts">
  import { goto } from "$app/navigation";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import DateInput from "$lib/components/ui/DateInput.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { logs, actionTypeFilter, actorFilter, dateFrom, dateTo, search } =
    $derived(data);

  const ACTION_LABELS: Record<string, string> = {
    SUBJECT_SUBMITTED: "Sujet soumis",
    SUBJECT_VALIDATED: "Sujet valide",
    SUBJECT_REFUSED: "Sujet refuse",
    STUDENT_ASSIGNED: "Etudiant affecte",
    JURY_ASSIGNED: "Jury designe",
    DEFENSE_SCHEDULED: "Soutenance planifiee",
    GRADE_SUBMITTED: "Note saisie",
    GRADE_VALIDATED: "Note validee",
    COMPANY_REGISTERED: "Entreprise inscrite",
    WISH_SUBMITTED: "Voeu soumis",
    NOTIFICATION_SENT: "Notification envoyee",
    CSV_IMPORTED: "CSV importe",
    ROLE_TRANSFERRED: "Role transfere",
  };

  function actionVariant(
    action: string,
  ): "success" | "warning" | "danger" | "info" | "neutral" {
    if (action.includes("VALIDATED") || action.includes("ASSIGNED"))
      return "success";
    if (action.includes("REFUSED") || action.includes("SCHEDULED"))
      return "warning";
    if (action.includes("SUBMITTED") || action.includes("SENT")) return "info";
    return "neutral";
  }

  function applyFilter(key: string, value: string) {
    const params = new URLSearchParams(window.location.search);
    if (value) {
      params.set(key, value);
    } else {
      params.delete(key);
    }
    const qs = params.toString();
    goto(qs ? `/admin/audit-log?${qs}` : "/admin/audit-log");
  }

  function clearFilters() {
    goto("/admin/audit-log");
  }

  function formatDate(iso: string): string {
    return new Date(iso).toLocaleDateString("fr-FR", {
      day: "numeric",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
</script>

<Page title="Journal d'activité" subtitle="Audit des actions significatives">
  {#snippet actions()}
    <Button variant="ghost" onclick={clearFilters}>Reinitialiser filtres</Button
    >
  {/snippet}

  <div class="filters">
    <label>
      Type d'action
      <select
        class="input"
        value={actionTypeFilter ?? ""}
        onchange={(e) => applyFilter("action_type", e.currentTarget.value)}
      >
        <option value="">Tous</option>
        {#each Object.keys(ACTION_LABELS) as action}
          <option value={action}>{ACTION_LABELS[action]}</option>
        {/each}
      </select>
    </label>

    <label>
      Acteur
      <input
        type="text"
        class="input"
        value={actorFilter ?? ""}
        placeholder="ID acteur"
        onchange={(e) => applyFilter("actor", e.currentTarget.value)}
      />
    </label>

    <label>
      Du
      <DateInput
        value={dateFrom ?? ""}
        onchange={(e) =>
          applyFilter("date_from", (e.currentTarget as HTMLInputElement).value)}
      />
    </label>

    <label>
      Au
      <DateInput
        value={dateTo ?? ""}
        onchange={(e) =>
          applyFilter("date_to", (e.currentTarget as HTMLInputElement).value)}
      />
    </label>

    <label>
      Recherche
      <input
        type="text"
        class="input"
        value={search ?? ""}
        placeholder="Rechercher..."
        onchange={(e) => applyFilter("search", e.currentTarget.value)}
      />
    </label>
  </div>

  <div class="table-wrapper">
    <table>
      <thead>
        <tr>
          <th>Date</th>
          <th>Action</th>
          <th>Entite</th>
          <th>Acteur</th>
        </tr>
      </thead>
      <tbody>
        {#each logs as log}
          <tr>
            <td>{formatDate(log.created_at)}</td>
            <td>
              <Badge
                variant={actionVariant(log.action)}
                label={ACTION_LABELS[log.action] ?? log.action}
              />
            </td>
            <td>{log.entity} #{log.entity_id}</td>
            <td>{log.actor_id}</td>
          </tr>
        {:else}
          <tr>
            <td colspan="4" class="empty">Aucune entrée trouvée</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</Page>

<style>
  .filters {
    display: flex;
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-lg);
    flex-wrap: wrap;
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
  }

  .input {
    width: 180px;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-surface);
    color: var(--color-text);
  }

  .table-wrapper {
    width: 100%;
    overflow-x: auto;
    border-radius: 12px;
    border: 1px solid var(--color-border);
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-sans);
    font-size: 0.9rem;
  }

  thead {
    background: var(--color-background-100);
  }

  th {
    text-align: left;
    padding: 0.75rem 1rem;
    font-weight: 600;
    color: var(--color-text);
    font-size: 0.8rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    border-bottom: 1px solid var(--color-border);
  }

  td {
    padding: 0.75rem 1rem;
    color: var(--color-text);
    border-bottom: 1px solid var(--color-border);
  }

  tr:last-child td {
    border-bottom: none;
  }

  tbody tr:hover {
    background: var(--color-background-100);
  }

  .empty {
    text-align: center;
    color: var(--color-text-muted);
    font-style: italic;
    padding: 2rem;
  }
</style>
