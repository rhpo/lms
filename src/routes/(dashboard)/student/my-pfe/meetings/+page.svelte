<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { Plus, FileText } from "lucide-svelte";
  import { student } from "$lib/api";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";

  let { data } = $props();
  const { meetings, pfe } = $derived(data);

  let showForm = $state(false);
  let submitting = $state(false);
  let submitError = $state("");

  let form = $state({
    meeting_date: "",
    duration: 60,
    meeting_type: "presentiel" as "presentiel" | "visio",
    topics: "",
    status: "a_faire" as "a_faire" | "en_cours" | "termine",
    observation: "",
  });

  const DURATION_OPTIONS = [
    { value: 30, label: "30 min" },
    { value: 60, label: "1 h" },
    { value: 90, label: "1 h 30" },
    { value: 120, label: "2 h" },
    { value: 150, label: "2 h 30" },
    { value: 180, label: "3 h" },
  ];

  const MEETING_TYPE_LABELS: Record<string, string> = {
    presentiel: "Présentiel",
    visio: "Visio",
  };

  const STATUS_LABELS: Record<string, string> = {
    a_faire: "À faire",
    en_cours: "En cours",
    termine: "Terminé",
  };

  let updatingStatusId = $state<number | null>(null);

  async function updateStatus(reportId: number, newStatus: string) {
    updatingStatusId = reportId;
    try {
      await student.updateMyMeeting(reportId, { status: newStatus });
      await invalidateAll();
    } catch {
    } finally {
      updatingStatusId = null;
    }
  }

  function formatDuration(minutes: number): string {
    if (minutes < 60) return `${minutes} min`;
    const h = Math.floor(minutes / 60);
    const m = minutes % 60;
    return m > 0 ? `${h} h ${m}` : `${h} h`;
  }

  function resetForm() {
    form = {
      meeting_date: "",
      duration: 60,
      meeting_type: "presentiel",
      topics: "",
      status: "a_faire",
      observation: "",
    };
    submitError = "";
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!form.meeting_date) {
      submitError = "La date est obligatoire.";
      return;
    }
    if (!form.topics.trim()) {
      submitError = "Les sujets abordés sont obligatoires.";
      return;
    }
    if (form.topics.trim().length < 5) {
      submitError =
        "Les sujets abordés doivent contenir au moins 5 caractères.";
      return;
    }
    submitting = true;
    submitError = "";
    try {
      await student.addMyMeeting({
        meeting_date: form.meeting_date,
        duration: form.duration,
        meeting_type: form.meeting_type,
        topics: form.topics,
        status: form.status,
        observation: form.observation || undefined,
      });
      resetForm();
      showForm = false;
      await invalidateAll();
    } catch (err: unknown) {
      submitError = err instanceof Error ? err.message : "Erreur réseau.";
    } finally {
      submitting = false;
    }
  }
</script>

<Page title="Suivi / Réunions" subtitle="Journal de suivi de votre PFE.">
  {#snippet actions()}
    {#if pfe}
      <Button
        variant="primary"
        Icon={Plus}
        onclick={() => (showForm = !showForm)}
      >
        {showForm ? "Annuler" : "Ajouter un suivi"}
      </Button>
    {/if}
  {/snippet}

  {#if !pfe}
    <div class="empty-state">
      <FileText size={48} />
      <h2>Aucun PFE en cours</h2>
      <p>Vous n'avez pas encore de PFE affecté.</p>
      <Button variant="primary" href="/student/my-pfe">Voir mon PFE</Button>
    </div>
  {:else}
    {#if showForm}
      <div class="form-card">
        <h3>Nouvelle entrée de suivi</h3>
        {#if submitError}
          <div class="error-banner">{submitError}</div>
        {/if}
        <form onsubmit={handleSubmit}>
          <div class="form-grid">
            <FormField label="Date" required>
              <input
                type="date"
                bind:value={form.meeting_date}
                required
                class="input"
              />
            </FormField>

            <FormField label="Durée" required>
              <select bind:value={form.duration} class="input">
                {#each DURATION_OPTIONS as opt}
                  <option value={opt.value}>{opt.label}</option>
                {/each}
              </select>
            </FormField>

            <FormField label="Type de réunion" required>
              <select bind:value={form.meeting_type} class="input">
                <option value="presentiel">Présentiel</option>
                <option value="visio">Visio</option>
              </select>
            </FormField>

            <FormField label="État" required>
              <select bind:value={form.status} class="input">
                <option value="a_faire">À faire</option>
                <option value="en_cours">En cours</option>
                <option value="termine">Terminé</option>
              </select>
            </FormField>
          </div>

          <FormField label="Points discutés">
            <textarea
              bind:value={form.topics}
              placeholder="Résumé des points abordés lors de la réunion..."
              rows={3}
              class="input"
            ></textarea>
          </FormField>

          <FormField label="Observation">
            <textarea
              bind:value={form.observation}
              placeholder="Remarques supplémentaires (optionnel)..."
              rows={2}
              class="input"
            ></textarea>
          </FormField>

          <div class="form-actions">
            <Button
              variant="ghost"
              type="button"
              onclick={() => {
                showForm = false;
                resetForm();
              }}
            >
              Annuler
            </Button>
            <Button variant="primary" type="submit" disabled={submitting}>
              {submitting ? "Enregistrement..." : "Enregistrer"}
            </Button>
          </div>
        </form>
      </div>
    {/if}

    {#if meetings.length === 0}
      <div class="empty-reports">
        <p>Aucune entrée de suivi pour le moment.</p>
      </div>
    {:else}
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>N° réunion</th>
              <th>Date</th>
              <th>Durée</th>
              <th>Type de réunion</th>
              <th>Points discutés</th>
              <th>État</th>
              <th>Observation</th>
            </tr>
          </thead>
          <tbody>
            {#each [...meetings].reverse() as report, i (report.id)}
              <tr>
                <td class="center">{i + 1}</td>
                <td
                  >{new Date(report.meeting_date).toLocaleDateString(
                    "fr-FR",
                  )}</td
                >
                <td>{formatDuration(report.duration)}</td>
                <td>
                  <Badge
                    variant="info"
                    label={MEETING_TYPE_LABELS[report.meeting_type] ??
                      report.meeting_type}
                  />
                </td>
                <td class="topics">{report.topics || "-"}</td>
                <td>
                  <select
                    class="status-select status-{report.status}"
                    value={report.status}
                    disabled={updatingStatusId === report.id}
                    onchange={(e) =>
                      updateStatus(
                        report.id,
                        (e.target as HTMLSelectElement).value,
                      )}
                  >
                    <option value="a_faire">À faire</option>
                    <option value="en_cours">En cours</option>
                    <option value="termine">Terminé</option>
                  </select>
                </td>
                <td class="observation">{report.observation ?? "-"}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
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

  .form-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-lg);
    margin-bottom: var(--spacing-lg);

    h3 {
      font-size: var(--text-base);
      font-weight: 600;
      font-family: var(--font-sans);
      color: var(--color-text);
      margin: 0 0 var(--spacing-md) 0;
    }

    form {
      display: flex;
      flex-direction: column;
      gap: var(--spacing-md);
    }
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-md);
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm);
    padding-top: var(--spacing-sm);
  }

  .error-banner {
    padding: 0.75rem 1rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-danger) 20%, transparent);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-danger);
    margin-bottom: var(--spacing-md);
  }

  .empty-reports {
    text-align: center;
    padding: 3rem 1rem;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
  }

  .table-wrapper {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
  }

  th {
    text-align: left;
    padding: 0.75rem 1rem;
    font-weight: 600;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
    border-bottom: 2px solid var(--color-border);
    white-space: nowrap;
  }

  td {
    padding: 0.75rem 1rem;
    color: var(--color-text);
    border-bottom: 1px solid var(--color-border);
    vertical-align: top;
  }

  tr:last-child td {
    border-bottom: none;
  }

  tr:hover td {
    background: var(--color-background-100);
  }

  .center {
    text-align: center;
    font-weight: 600;
    color: var(--color-text-muted);
  }

  .topics {
    max-width: 260px;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .observation {
    max-width: 200px;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    white-space: pre-wrap;
    word-break: break-word;
  }

  .status-select {
    appearance: none;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    padding: 0.25rem 0.6rem;
    font-size: var(--text-xs);
    font-family: var(--font-sans);
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.15s;

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }

  .status-a_faire {
    background: color-mix(in srgb, var(--color-text-muted) 12%, transparent);
    color: var(--color-text-muted);
    border-color: color-mix(in srgb, var(--color-text-muted) 25%, transparent);
  }

  .status-en_cours {
    background: color-mix(in srgb, var(--color-warning) 12%, transparent);
    color: var(--color-warning);
    border-color: color-mix(in srgb, var(--color-warning) 30%, transparent);
  }

  .status-termine {
    background: color-mix(in srgb, var(--color-success) 12%, transparent);
    color: var(--color-success);
    border-color: color-mix(in srgb, var(--color-success) 30%, transparent);
  }
</style>
