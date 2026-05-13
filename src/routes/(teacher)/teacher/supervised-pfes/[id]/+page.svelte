<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { ArrowLeft, Plus } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import DateInput from "$lib/components/ui/DateInput.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { pfe, progressReports, supervisorEval } = $derived(data);

  const STATUS_LABELS: Record<string, string> = {
    en_cours: "En cours",
    soutenance_planifiee: "Soutenance planifiee",
    valide: "Valide",
    refuse: "Refuse",
  };

  const STATUS_VARIANTS: Record<
    string,
    "info" | "warning" | "success" | "danger"
  > = {
    en_cours: "info",
    soutenance_planifiee: "warning",
    valide: "success",
    refuse: "danger",
  };

  // Meeting form state
  let meetingDate = $state(new Date().toISOString().split("T")[0]);
  let meetingNotes = $state("");
  let meetingProgress = $state("");
  let meetingError = $state("");
  let meetingLoading = $state(false);

  // Evaluation form state
  let criterion5 = $state(0);
  let evalComment = $state("");
  let evalError = $state("");
  let evalSuccess = $state("");
  let evalLoading = $state(false);

  async function handleSubmitEvaluation() {
    evalError = "";
    evalSuccess = "";
    if (criterion5 < 0 || criterion5 > 4) {
      evalError = "La note doit etre comprise entre 0 et 4.";
      return;
    }
    evalLoading = true;
    try {
      const res = await fetch("/api/evaluations/supervisor/submit", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          assignmentId: pfe.id,
          criterion5,
        }),
      });
      if (res.ok) {
        evalSuccess = "Evaluation soumise avec succes.";
        await invalidateAll();
      } else {
        const body = await res.json().catch(() => ({}));
        evalError =
          body.message ?? "Erreur lors de la soumission de l evaluation.";
      }
    } catch {
      evalError = "Erreur reseau.";
    } finally {
      evalLoading = false;
    }
  }

  async function handleAddMeeting() {
    if (!meetingNotes.trim()) {
      meetingError = "Veuillez saisir des notes.";
      return;
    }

    meetingError = "";
    meetingLoading = true;

    try {
      const res = await fetch(`/api/assignments/${pfe.id}/log-meeting`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          meeting_date: meetingDate,
          notes: meetingNotes,
          progress: meetingProgress,
        }),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.error?.message || "Erreur lors de l'ajout");
      }

      meetingNotes = "";
      meetingProgress = "";
      await invalidateAll();
    } catch (err: unknown) {
      meetingError = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      meetingLoading = false;
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
    <h2>PFE #{pfe.id.slice(0, 8)}</h2>

    <div class="meta">
      <div class="meta-item">
        <span class="label">Sujet</span>
        <span>{pfe.subject_id}</span>
      </div>
      <div class="meta-item">
        <span class="label">Etudiant(s)</span>
        <span>
          {pfe.student_id}
          {#if pfe.student2_id}, {pfe.student2_id}{/if}
          {#if pfe.student3_id}, {pfe.student3_id}{/if}
        </span>
      </div>
      <div class="meta-item">
        <span class="label">Statut</span>
        <Badge
          variant={STATUS_VARIANTS[pfe.status] || "info"}
          label={STATUS_LABELS[pfe.status] || pfe.status}
        />
      </div>
      {#if pfe.co_supervisor_id}
        <div class="meta-item">
          <span class="label">Co-encadrant</span>
          <span>{pfe.co_supervisor_id}</span>
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
      <div class="form-row">
        <label>
          Date
          <DateInput bind:value={meetingDate} />
        </label>
      </div>
      <div class="form-row">
        <label>
          Notes
          <textarea
            bind:value={meetingNotes}
            rows={3}
            placeholder="Notes de la reunion..."
          ></textarea>
        </label>
      </div>
      <div class="form-row">
        <label>
          Avancement
          <textarea
            bind:value={meetingProgress}
            rows={2}
            placeholder="Avancement du projet..."
          ></textarea>
        </label>
      </div>
      <Button
        variant="primary"
        onclick={handleAddMeeting}
        disabled={meetingLoading}
      >
        {meetingLoading ? "Ajout..." : "Ajouter"}
      </Button>
    </div>

    {#if progressReports.length === 0}
      <div class="empty">
        <p>Aucune entrée de journal pour le moment.</p>
      </div>
    {:else}
      <div class="entries">
        {#each progressReports as entry}
          <div class="entry">
            <div class="entry-header">
              <span class="entry-date"
                >{new Date(entry.meeting_date).toLocaleDateString(
                  "fr-FR",
                )}</span
              >
            </div>
            <p class="entry-notes">{entry.student_notes}</p>
            {#if entry.teacher_feedback}
              <div class="entry-feedback">
                <strong>Feedback :</strong>
                {entry.teacher_feedback}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </section>

  <section class="evaluation">
    <h3>Evaluation de l'encadrant (critere 5)</h3>

    {#if evalError}
      <div class="error">{evalError}</div>
    {/if}
    {#if evalSuccess}
      <div class="success">{evalSuccess}</div>
    {/if}

    {#if supervisorEval}
      <div class="eval-submitted">
        <p>Vous avez deja soumis votre evaluation.</p>
        <p>
          Note attribuee : <strong>{supervisorEval.total}/4</strong>
        </p>
        {#if supervisorEval.comment}
          <p>Commentaire : {supervisorEval.comment}</p>
        {/if}
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
  </section>
</Page>

<style>
  .error {
    padding: 0.75rem 1rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    color: var(--color-danger);
    border-radius: 8px;
    font-size: var(--text-sm);
    margin-bottom: var(--spacing-md);
  }

  .pfe-info {
    margin-bottom: var(--spacing-lg);
  }

  .pfe-info h2 {
    font-size: var(--text-lg);
    font-weight: 700;
    margin: 0 0 var(--spacing-md);
    color: var(--color-text);
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

  .journal {
    border-top: 1px solid var(--color-border);
    padding-top: var(--spacing-lg);
  }

  .journal h3 {
    font-size: var(--text-md);
    font-weight: 600;
    margin: 0 0 var(--spacing-md);
    color: var(--color-text);
  }

  .add-meeting {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: var(--spacing-md);
    margin-bottom: var(--spacing-lg);
  }

  .add-meeting h4 {
    font-size: var(--text-sm);
    font-weight: 600;
    margin: 0 0 var(--spacing-sm);
    color: var(--color-text);
  }

  .form-row {
    margin-bottom: var(--spacing-sm);
  }

  .form-row label {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-text);
  }

  .form-row textarea {
    padding: 0.5rem 0.65rem;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    color: var(--color-text);
    background: var(--color-background);
    transition: border-color var(--transition-fast);

    &:focus {
      outline: none;
      border-color: var(--color-accent);
    }
  }

  .empty {
    text-align: center;
    padding: 2rem 1rem;
    color: var(--color-text-muted);

    & p {
      font-size: var(--text-sm);
      margin: 0;
    }
  }

  .entries {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .entry {
    padding: var(--spacing-sm) var(--spacing-md);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background: var(--color-surface);
  }

  .entry-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.35rem;
  }

  .entry-date {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--color-text-muted);
  }

  .entry-notes {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--color-text);
    line-height: 1.5;
  }

  .entry-feedback {
    margin-top: 0.5rem;
    padding-top: 0.5rem;
    border-top: 1px solid var(--color-border);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
</style>
