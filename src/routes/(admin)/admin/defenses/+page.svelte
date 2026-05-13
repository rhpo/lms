<script lang="ts">
  import { goto } from "$app/navigation";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import DateInput from "$lib/components/ui/DateInput.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { defenses, dateFilter, statusFilter } = $derived(data);

  const STATUS_LABELS: Record<string, string> = {
    scheduled: "Planifiee",
    done: "Passee",
    postponed: "Reportee",
  };

  function statusVariant(
    status: string,
  ): "success" | "warning" | "danger" | "info" | "neutral" {
    switch (status) {
      case "scheduled":
        return "warning";
      case "done":
        return "success";
      case "postponed":
        return "danger";
      default:
        return "neutral";
    }
  }

  function applyFilter(key: string, value: string) {
    const params = new URLSearchParams(window.location.search);
    if (value) {
      params.set(key, value);
    } else {
      params.delete(key);
    }
    const qs = params.toString();
    goto(qs ? `/admin/defenses?${qs}` : "/admin/defenses");
  }

  function clearFilters() {
    goto("/admin/defenses");
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

<Page
  title="Soutenances"
  subtitle="Liste de toutes les soutenances planifiees et passees"
>
  {#snippet actions()}
    <a href="/admin/defenses/new">
      <Button variant="primary">Planifier une soutenance</Button>
    </a>
    <Button variant="ghost" onclick={clearFilters}>Reinitialiser filtres</Button
    >
  {/snippet}

  <div class="filters">
    <label>
      Date
      <DateInput
        value={dateFilter ?? ""}
        onchange={(e) =>
          applyFilter("date", (e.currentTarget as HTMLInputElement).value)}
      />
    </label>

    <label>
      Statut
      <select
        class="input"
        value={statusFilter ?? ""}
        onchange={(e) => applyFilter("status", e.currentTarget.value)}
      >
        <option value="">Tous</option>
        {#each Object.keys(STATUS_LABELS) as status}
          <option value={status}>{STATUS_LABELS[status]}</option>
        {/each}
      </select>
    </label>
  </div>

  <div class="table-wrapper">
    <table>
      <thead>
        <tr>
          <th>Date</th>
          <th>Salle</th>
          <th>Statut</th>
          <th>Note finale</th>
          <th>Resultat</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each defenses as defense}
          <tr>
            <td>{formatDate(defense.scheduled_at)}</td>
            <td>{defense.room}</td>
            <td>
              <Badge
                variant={statusVariant(defense.status)}
                label={STATUS_LABELS[defense.status] ?? defense.status}
              />
            </td>
            <td>
              {#if defense.final_grade !== null}
                {defense.final_grade}/20
              {:else}
                —
              {/if}
            </td>
            <td>
              {#if defense.result === "admitted"}
                <Badge variant="success" label="Admis" />
              {:else if defense.result === "corrections_required"}
                <Badge variant="warning" label="Corrections requises" />
              {:else if defense.result === "not_admitted"}
                <Badge variant="danger" label="Non admis" />
              {:else}
                —
              {/if}
            </td>
            <td>
              <a href="/admin/defenses/{defense.id}">
                <Button variant="ghost">Voir</Button>
              </a>
            </td>
          </tr>
        {:else}
          <tr>
            <td colspan="6" class="empty">Aucune soutenance trouvée</td>
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
    width: 200px;
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

  a {
    color: var(--color-accent);
    text-decoration: none;
  }

  a:hover {
    text-decoration: underline;
  }

  .empty {
    text-align: center;
    color: var(--color-text-muted);
    font-style: italic;
    padding: 2rem;
  }
</style>
