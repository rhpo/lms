<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { ArrowLeft } from "lucide-svelte";
  import { teacher } from "$lib/api";
  import { authStore } from "$lib/stores/auth";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { pfe, progressReports, supervisorEval } = $derived(data);
  const myProfileId = $derived(authStore.profile?.id);
  const isCoSupervisor = $derived(
    pfe != null &&
      pfe.supervisor?.profile_id !== myProfileId &&
      pfe.co_supervisor?.profile_id === myProfileId,
  );

  const PFE_STATUS_LABELS: Record<string, string> = {
    en_cours: "En cours",
    soutenance_planifiee: "Soutenance planifiée",
    valide: "Validé",
    refuse: "Refusé",
  };

  const PFE_STATUS_VARIANTS: Record<
    string,
    "info" | "warning" | "success" | "danger"
  > = {
    en_cours: "info",
    soutenance_planifiee: "warning",
    valide: "success",
    refuse: "danger",
  };

  const DURATION_OPTIONS = [
    { value: 30, label: "30 min" },
    { value: 45, label: "45 min" },
    { value: 60, label: "1h" },
    { value: 90, label: "1h30" },
    { value: 120, label: "2h" },
  ];

  const MEETING_TYPE_LABELS: Record<string, string> = {
    presentiel: "Présentiel",
    visio: "Visioconférence",
  };

  const MEETING_STATUS_LABELS: Record<string, string> = {
    a_faire: "À faire",
    en_cours: "En cours",
    termine: "Terminé",
  };

  function formatDuration(mins: number): string {
    if (mins < 60) return `${mins} min`;
    const h = Math.floor(mins / 60);
    const m = mins % 60;
    return m ? `${h}h${String(m).padStart(2, "0")}` : `${h}h`;
  }

  let newMeeting = $state({
    date: new Date().toISOString().split("T")[0],
    duration: 60,
    type: "presentiel" as "presentiel" | "visio",
    status: "a_faire" as "a_faire" | "en_cours" | "termine",
    topics: "",
    observation: "",
  });
  let meetingError = $state("");
  let meetingLoading = $state(false);

  function resetMeetingForm() {
    newMeeting = {
      date: new Date().toISOString().split("T")[0],
      duration: 60,
      type: "presentiel",
      status: "a_faire",
      topics: "",
      observation: "",
    };
  }

  async function handleAddMeeting() {
    if (!newMeeting.topics.trim()) {
      meetingError = "Les points discutés sont requis.";
      return;
    }
    meetingError = "";
    meetingLoading = true;
    try {
      await teacher.addMeeting(pfe!.id, {
        meeting_date: newMeeting.date,
        duration: newMeeting.duration,
        meeting_type: newMeeting.type,
        topics: newMeeting.topics.trim(),
        status: newMeeting.status,
        observation: newMeeting.observation.trim() || undefined,
      });
      resetMeetingForm();
      await invalidateAll();
    } catch (err: unknown) {
      meetingError = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      meetingLoading = false;
    }
  }

  let criterion5 = $state(0);
  let evalComment = $state("");
  let evalError = $state("");
  let evalSuccess = $state("");
  let evalLoading = $state(false);

  async function handleSubmitEvaluation() {
    evalError = "";
    evalSuccess = "";
    if (criterion5 < 0 || criterion5 > 4) {
      evalError = "La note doit être comprise entre 0 et 4.";
      return;
    }
    evalLoading = true;
    try {
      await teacher.submitEvaluation(pfe!.id, { criterion5 });
      evalSuccess = "Évaluation soumise avec succès.";
      await invalidateAll();
    } catch (err: unknown) {
      evalError = err instanceof Error ? err.message : "Erreur réseau.";
    } finally {
      evalLoading = false;
    }
  }
</script>

<Page title="Detail du PFE" subtitle="Informations et journal de suivi.">
  {#snippet actions()}
    <Button
      variant="ghost"
      Icon={ArrowLeft}
      onclick={() => goto("/teacher/supervised-pfes")}
    >
      Retour
    </Button>
  {/snippet}

  <section class="pfe-info">
    <h2>{pfe!.subject?.title ?? pfe!.subject_title ?? "PFE sans titre"}</h2>

    <div class="meta">
      <div class="meta-item">
        <span class="label">Code PFE</span>
        <span>{pfe!.pfe_code}</span>
      </div>
      <div class="meta-item">
        <span class="label">Étudiant(s)</span>
        <span>
          {pfe!.student?.profile?.full_name ?? "Étudiant inconnu"}
          {#if pfe!.student2_id}, {pfe!.student2?.profile?.full_name ??
              "Étudiant inconnu"}{/if}
          {#if pfe!.student3_id}, {pfe!.student3?.profile?.full_name ??
              "Étudiant inconnu"}{/if}
        </span>
      </div>
      <div class="meta-item">
        <span class="label">Statut</span>
        <Badge
          variant={PFE_STATUS_VARIANTS[pfe!.status] || "info"}
          label={PFE_STATUS_LABELS[pfe!.status] || pfe!.status}
        />
      </div>
      <div class="meta-item">
        <span class="label">Votre rôle</span>
        {#if isCoSupervisor}
          <Badge variant="warning" label="Co-encadrant" />
        {:else}
          <Badge variant="info" label="Encadrant principal" />
        {/if}
      </div>
      {#if pfe!.co_supervisor_id}
        <div class="meta-item">
          <span class="label">Co-encadrant</span>
          <span
            >{pfe!.co_supervisor?.profile?.full_name ??
              pfe!.co_supervisor_id}</span
          >
        </div>
      {/if}
    </div>
  </section>

  <section class="journal">
    <h3>Journal de suivi</h3>

    {#if meetingError}
      <div class="error">{meetingError}</div>
    {/if}

    <div class="add-meeting">
      <h4>Ajouter une entrée</h4>
      <div class="form-grid">
        <div class="field">
          <label for="m-date">Date</label>
          <input id="m-date" type="date" bind:value={newMeeting.date} />
        </div>

        <div class="field">
          <label for="m-duration">Durée</label>
          <select id="m-duration" bind:value={newMeeting.duration}>
            {#each DURATION_OPTIONS as opt}
              <option value={opt.value}>{opt.label}</option>
            {/each}
          </select>
        </div>

        <div class="field">
          <label for="m-type">Type de réunion</label>
          <select id="m-type" bind:value={newMeeting.type}>
            <option value="presentiel">Présentiel</option>
            <option value="visio">Visioconférence</option>
          </select>
        </div>

        <div class="field">
          <label for="m-status">État</label>
          <select id="m-status" bind:value={newMeeting.status}>
            <option value="a_faire">À faire</option>
            <option value="en_cours">En cours</option>
            <option value="termine">Terminé</option>
          </select>
        </div>

        <div class="field full-width">
          <label for="m-topics">Points discutés</label>
          <textarea
            id="m-topics"
            bind:value={newMeeting.topics}
            rows={3}
            placeholder="Sujets abordés lors de cette réunion..."
          ></textarea>
        </div>

        <div class="field full-width">
          <label for="m-obs">Observation (optionnel)</label>
          <textarea
            id="m-obs"
            bind:value={newMeeting.observation}
            rows={2}
            placeholder="Remarques, décisions, prochaines étapes..."
          ></textarea>
        </div>
      </div>

      <Button
        variant="primary"
        onclick={handleAddMeeting}
        disabled={meetingLoading}
      >
        {meetingLoading ? "Ajout..." : "Ajouter l'entrée"}
      </Button>
    </div>

    {#if progressReports.length === 0}
      <div class="empty">
        <p>Aucune entrée de journal pour le moment.</p>
      </div>
    {:else}
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Date</th>
              <th>Durée</th>
              <th>Type</th>
              <th>État</th>
              <th>Points discutés</th>
              <th>Observation</th>
            </tr>
          </thead>
          <tbody>
            {#each progressReports as entry}
              <tr>
                <td class="cell-date">
                  {new Date(entry.meeting_date).toLocaleDateString("fr-FR")}
                </td>
                <td class="cell-meta">{formatDuration(entry.duration)}</td>
                <td class="cell-meta"
                  >{MEETING_TYPE_LABELS[entry.meeting_type] ??
                    entry.meeting_type}</td
                >
                <td>
                  <span class="status-badge status-{entry.status}">
                    {MEETING_STATUS_LABELS[entry.status] ?? entry.status}
                  </span>
                </td>
                <td class="cell-topics">{entry.topics}</td>
                <td class="cell-obs">{entry.observation ?? "-"}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </section>

  <section class="evaluation">
    <h3>Evaluation de l'encadrant (critere 5)</h3>

    {#if isCoSupervisor}
      <p class="co-supervisor-note">
        En tant que co-encadrant, vous ne pouvez pas soumettre l'évaluation.
        Cette action est réservée à l'encadrant principal.
      </p>
    {:else}
      {#if evalError}
        <div class="error">{evalError}</div>
      {/if}
      {#if evalSuccess}
        <div class="success-msg">{evalSuccess}</div>
      {/if}

      {#if supervisorEval}
        <div class="eval-submitted">
          <p>Vous avez deja soumis votre evaluation.</p>
          <p>Note attribuee : <strong>{supervisorEval.criterion5}/4</strong></p>
        </div>
      {:else}
        <div class="eval-form">
          <p class="eval-desc">
            Attribuez une note sur 4 points pour l'encadrement de ce PFE.
          </p>
          <div class="form-row">
            <label>
              Note (/4)
              <input
                type="number"
                min="0"
                max="4"
                step="0.5"
                bind:value={criterion5}
                required
              />
            </label>
          </div>
          <div class="form-row">
            <label>
              Commentaire (optionnel)
              <textarea
                bind:value={evalComment}
                rows={3}
                placeholder="Votre avis sur le travail..."
              ></textarea>
            </label>
          </div>
          <Button
            variant="primary"
            onclick={handleSubmitEvaluation}
            disabled={evalLoading}
          >
            {evalLoading ? "Soumission..." : "Soumettre l'evaluation"}
          </Button>
        </div>
      {/if}
    {/if}
  </section>
</Page>

<style>
  /* ── Shared input style ─────────────────────────────── */
  input,
  select,
  textarea {
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-surface);
    color: var(--color-text);
    width: 100%;
    box-sizing: border-box;
    transition: border-color var(--transition-fast);

    &:focus {
      outline: none;
      border-color: var(--color-accent);
      box-shadow: 0 0 0 2px
        color-mix(in srgb, var(--color-accent) 20%, transparent);
    }
  }

  textarea {
    resize: vertical;
    min-height: 80px;
  }

  /* ── Alerts ─────────────────────────────────────────── */
  .error {
    padding: 0.75rem 1rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    color: var(--color-danger);
    border-radius: 8px;
    font-size: var(--text-sm);
    margin-bottom: var(--spacing-md);
  }

  .success-msg {
    padding: 0.75rem 1rem;
    background: color-mix(
      in srgb,
      var(--color-success, #22c55e) 10%,
      transparent
    );
    color: var(--color-success, #22c55e);
    border: 1px solid
      color-mix(in srgb, var(--color-success, #22c55e) 25%, transparent);
    border-radius: 8px;
    font-size: var(--text-sm);
    margin-bottom: var(--spacing-md);
  }

  /* ── PFE info ───────────────────────────────────────── */
  .pfe-info {
    margin-bottom: var(--spacing-lg);

    h2 {
      font-size: var(--text-lg);
      font-weight: 700;
      margin: 0 0 var(--spacing-md);
      color: var(--color-text);
    }
  }

  .meta {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-md);
  }

  .meta-item {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .label {
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
  }

  /* ── Journal section ────────────────────────────────── */
  .journal {
    border-top: 1px solid var(--color-border);
    padding-top: var(--spacing-lg);

    h3 {
      font-size: var(--text-md);
      font-weight: 600;
      margin: 0 0 var(--spacing-md);
      color: var(--color-text);
    }
  }

  /* ── Add meeting form ───────────────────────────────── */
  .add-meeting {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: var(--spacing-md);
    margin-bottom: var(--spacing-lg);

    h4 {
      font-size: var(--text-sm);
      font-weight: 600;
      margin: 0 0 var(--spacing-md);
      color: var(--color-text);
    }
  }

  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-sm) var(--spacing-md);
    margin-bottom: var(--spacing-md);

    @media (max-width: 600px) {
      grid-template-columns: 1fr;
    }
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;

    label {
      font-family: var(--font-sans);
      font-size: var(--text-sm);
      font-weight: 600;
      color: var(--color-text-muted);
    }
  }

  .full-width {
    grid-column: 1 / -1;
  }

  /* ── Empty ──────────────────────────────────────────── */
  .empty {
    text-align: center;
    padding: 2rem 1rem;
    color: var(--color-text-muted);

    p {
      font-size: var(--text-sm);
      margin: 0;
    }
  }

  /* ── Table ──────────────────────────────────────────── */
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
    padding: 0.6rem 0.75rem;
    font-size: 0.72rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
    border-bottom: 1px solid var(--color-border);
    white-space: nowrap;
  }

  td {
    padding: 0.7rem 0.75rem;
    border-bottom: 1px solid var(--color-border);
    color: var(--color-text);
    vertical-align: top;
  }

  tr:last-child td {
    border-bottom: none;
  }

  .cell-date {
    white-space: nowrap;
    font-weight: 600;
    color: var(--color-text-muted);
  }
  .cell-meta {
    white-space: nowrap;
    color: var(--color-text-muted);
  }
  .cell-topics {
    max-width: 260px;
    line-height: 1.5;
  }
  .cell-obs {
    max-width: 200px;
    line-height: 1.5;
    color: var(--color-text-muted);
    font-style: italic;
  }

  .status-badge {
    display: inline-block;
    padding: 0.2rem 0.55rem;
    border-radius: 999px;
    font-size: 0.72rem;
    font-weight: 600;
    white-space: nowrap;
  }

  .status-a_faire {
    background: color-mix(in srgb, var(--color-text-muted) 12%, transparent);
    color: var(--color-text-muted);
  }

  .status-en_cours {
    background: color-mix(
      in srgb,
      var(--color-warning, #f59e0b) 15%,
      transparent
    );
    color: var(--color-warning, #f59e0b);
  }

  .status-termine {
    background: color-mix(
      in srgb,
      var(--color-success, #22c55e) 15%,
      transparent
    );
    color: var(--color-success, #22c55e);
  }

  /* ── Evaluation section ─────────────────────────────── */
  .evaluation {
    border-top: 1px solid var(--color-border);
    padding-top: var(--spacing-lg);
    margin-top: var(--spacing-lg);

    h3 {
      font-size: var(--text-md);
      font-weight: 600;
      margin: 0 0 var(--spacing-md);
      color: var(--color-text);
    }
  }

  .co-supervisor-note {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    font-style: italic;
    padding: var(--spacing-md);
    background: color-mix(
      in srgb,
      var(--color-warning) 8%,
      var(--color-surface)
    );
    border: 1px solid color-mix(in srgb, var(--color-warning) 25%, transparent);
    border-radius: 8px;
    margin: 0;
  }

  .eval-submitted {
    padding: var(--spacing-md);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-family: var(--font-sans);

    p {
      margin: 0 0 0.25rem;
      font-size: var(--text-sm);
      color: var(--color-text-muted);
    }

    strong {
      color: var(--color-text);
      font-size: var(--text-lg);
    }
  }

  .eval-form {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);

    label {
      display: flex;
      flex-direction: column;
      gap: 0.3rem;
      font-family: var(--font-sans);
      font-size: var(--text-sm);
      font-weight: 600;
      color: var(--color-text-muted);
    }
  }

  .eval-desc {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    font-family: var(--font-sans);
    margin: 0;
  }
</style>
