<script lang="ts">
  import { goto } from "$app/navigation";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const {
    academicYears,
    yearFilter,
    totalPfes,
    totalTeachers,
    totalStudents,
    totalCompanies,
    totalDefenses,
    validatedSubjects,
    pendingSubjects,
    refusedSubjects,
    activePfes,
    defendedPfes,
    scheduledDefenses,
    doneDefenses,
  } = $derived(data);

  function handleYearChange(e: Event) {
    const value = (e.currentTarget as HTMLSelectElement).value;
    const params = new URLSearchParams(window.location.search);
    if (value) {
      params.set("academic_year_id", value);
    } else {
      params.delete("academic_year_id");
    }
    const qs = params.toString();
    goto(qs ? `/admin/statistics?${qs}` : "/admin/statistics");
  }
</script>

<Page title="Statistiques" subtitle="Tableau de bord analytique">
  <div class="year-filter">
    <label>
      Annee universitaire
      <select value={yearFilter ?? ""} onchange={handleYearChange}>
        <option value="">Toutes les annees</option>
        {#each academicYears as year}
          <option value={year.id}>{year.label}</option>
        {/each}
      </select>
    </label>
  </div>

  <div class="cards">
    <div class="card">
      <span class="card-value">{totalPfes}</span>
      <span class="card-label">PFE</span>
    </div>
    <div class="card">
      <span class="card-value">{totalTeachers}</span>
      <span class="card-label">Enseignants</span>
    </div>
    <div class="card">
      <span class="card-value">{totalStudents}</span>
      <span class="card-label">Etudiants</span>
    </div>
    <div class="card">
      <span class="card-value">{totalCompanies}</span>
      <span class="card-label">Entreprises</span>
    </div>
    <div class="card">
      <span class="card-value">{totalDefenses}</span>
      <span class="card-label">Soutenances</span>
    </div>
  </div>

  <div class="grid">
    <div class="section">
      <h3>Sujets</h3>
      <table>
        <tbody>
          <tr>
            <td>Valides</td>
            <td><Badge variant="success" label={String(validatedSubjects)} /></td>
          </tr>
          <tr>
            <td>En attente</td>
            <td><Badge variant="warning" label={String(pendingSubjects)} /></td>
          </tr>
          <tr>
            <td>Refuses</td>
            <td><Badge variant="danger" label={String(refusedSubjects)} /></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="section">
      <h3>PFE</h3>
      <table>
        <tbody>
          <tr>
            <td>En cours</td>
            <td><Badge variant="info" label={String(activePfes)} /></td>
          </tr>
          <tr>
            <td>Soutenus / valides</td>
            <td><Badge variant="success" label={String(defendedPfes)} /></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="section">
      <h3>Soutenances</h3>
      <table>
        <tbody>
          <tr>
            <td>Planifiees</td>
            <td><Badge variant="warning" label={String(scheduledDefenses)} /></td>
          </tr>
          <tr>
            <td>Realisees</td>
            <td><Badge variant="success" label={String(doneDefenses)} /></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</Page>

<style>
  .year-filter {
    margin-bottom: var(--spacing-lg);
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

  select {
    width: 240px;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-surface);
    color: var(--color-text);
  }

  .cards {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-lg);
  }

  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-lg);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
  }

  .card-value {
    font-size: 2rem;
    font-weight: 700;
    color: var(--color-accent);
  }

  .card-label {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: var(--spacing-lg);
  }

  .section {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-lg);
  }

  h3 {
    margin: 0 0 var(--spacing-md);
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-sans);
    font-size: 0.9rem;
  }

  td {
    padding: 0.5rem 0;
    color: var(--color-text);
    border-bottom: 1px solid var(--color-border);
  }

  td:last-child {
    text-align: right;
  }

  tr:last-child td {
    border-bottom: none;
  }
</style>
