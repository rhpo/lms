<script lang="ts">
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { academicYears } = $derived(data);

  let selectedYear = $state("");
  let exportType = $state("assignments");
  let format = $state("xlsx");
  let loading = $state(false);
  let error = $state("");

  async function handleExport(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    try {
      const params = new URLSearchParams();
      if (selectedYear) params.set("academic_year_id", selectedYear);
      params.set("type", exportType);
      params.set("format", format);

      const res = await fetch(`/api/exports?${params.toString()}`);

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.message || "Erreur lors de l'export");
      }

      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `export-${exportType}.${format}`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      loading = false;
    }
  }
</script>

<Page title="Exports" subtitle="Exportez les donnees en Excel ou PDF">
  <form onsubmit={handleExport}>
    {#if error}
      <div class="error">{error}</div>
    {/if}

    <div class="field">
      <label for="year">Annee universitaire</label>
      <select id="year" bind:value={selectedYear}>
        <option value="">Toutes les annees</option>
        {#each academicYears as year}
          <option value={year.id}>{year.label}</option>
        {/each}
      </select>
    </div>

    <div class="field">
      <label for="type">Type d'export</label>
      <select id="type" bind:value={exportType}>
        <option value="assignments">Affectations</option>
        <option value="defenses">Plannings de soutenance</option>
        <option value="statistics">Statistiques</option>
      </select>
    </div>

    <div class="field">
      <label for="format">Format</label>
      <select id="format" bind:value={format}>
        <option value="xlsx">Excel (.xlsx)</option>
        <option value="pdf">PDF</option>
      </select>
    </div>

    <div class="actions">
      <Button variant="primary" type="submit" disabled={loading}>
        {loading ? "Export en cours..." : "Exporter"}
      </Button>
    </div>
  </form>
</Page>

<style>
  form {
    max-width: 480px;
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .error {
    padding: 0.75rem 1rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    color: var(--color-danger);
    border-radius: 8px;
    font-size: var(--text-sm);
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  label {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
  }

  select {
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-surface);
    color: var(--color-text);
  }

  .actions {
    display: flex;
    gap: var(--spacing-sm);
    align-items: center;
  }
</style>
