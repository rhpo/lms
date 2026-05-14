<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { ArrowLeft, CheckCircle, XCircle, Star } from "lucide-svelte";
  import { admin } from "$lib/api";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { subject, teachers } = $derived(data);

  // === Toast ====
  let toastMessage = $state("");
  let toastType = $state<"success" | "error">("success");
  let toastVisible = $state(false);
  let toastTimer: ReturnType<typeof setTimeout> | null = null;

  function showToast(message: string, type: "success" | "error" = "success") {
    toastMessage = message;
    toastType = type;
    toastVisible = true;
    if (toastTimer) clearTimeout(toastTimer);
    toastTimer = setTimeout(() => {
      toastVisible = false;
    }, 3000);
  }

  // === Assign validators =======
  let showAssignModal = $state(false);
  let assignError = $state("");
  let validator1Id = $state("");
  let validator2Id = $state("");

  function openAssignModal() {
    assignError = "";
    validator1Id = "";
    validator2Id = "";
    showAssignModal = true;
  }

  async function assignValidatorsAction() {
    assignError = "";
    try {
      await admin.subjectAction(subject.id, "assign-validators", {
        validator1_id: validator1Id,
        validator2_id: validator2Id,
      });
      showAssignModal = false;
      showToast("Validateurs assignes avec succes");
      await invalidateAll();
    } catch (err) {
      assignError = err instanceof Error ? err.message : "Erreur reseau";
    }
  }

  async function unblockSubject() {
    try {
      await admin.subjectAction(subject.id, "unblock");
      showToast("Sujet debloque avec succes");
      await invalidateAll();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Erreur reseau", "error");
    }
  }

  function statusVariant(
    status: string,
  ): "success" | "warning" | "danger" | "info" | "neutral" {
    switch (status) {
      case "valide":
        return "success";
      case "en_attente":
        return "warning";
      case "refuse":
        return "danger";
      case "expire":
        return "neutral";
      default:
        return "info";
    }
  }

  function reviewVariant(
    decision: string | null,
  ): "success" | "warning" | "danger" | "info" {
    switch (decision) {
      case "accepte":
        return "success";
      case "accepte_sous_reserve":
        return "warning";
      case "refuse":
        return "danger";
      default:
        return "info";
    }
  }

  // === Recommendation helpers ================
  function matchStars(score: number): string {
    return "★".repeat(score) + "☆".repeat(3 - score);
  }

  function matchLabel(score: number): string {
    if (score === 3) return "Specialite identique";
    if (score === 1) return "Specialite proche";
    return "Aucune correspondance";
  }
</script>

<Page title={subject.title} subtitle="Detail du sujet PFE">
  {#snippet actions()}
    <a href="/admin/subjects">
      <Button variant="ghost" Icon={ArrowLeft}>Retour</Button>
    </a>
    {#if subject.status === "en_attente"}
      <Button variant="primary" onclick={openAssignModal}
        >Assigner validateurs</Button
      >
    {/if}
    {#if subject.status === "refuse"}
      <Button variant="secondary" onclick={unblockSubject}>Debloquer</Button>
    {/if}
  {/snippet}

  <!-- Toast notification -->
  {#if toastVisible}
    <div
      class="toast"
      class:toast-success={toastType === "success"}
      class:toast-error={toastType === "error"}
    >
      {#if toastType === "success"}
        <CheckCircle size={18} />
      {:else}
        <XCircle size={18} />
      {/if}
      <span>{toastMessage}</span>
    </div>
  {/if}

  <div class="detail-grid">
    <section class="card">
      <h2>Informations Générales</h2>

      <dl>
        <dt>Titre</dt>
        <dd>{subject.title}</dd>

        <dt>Description</dt>
        <dd>{subject.description}</dd>

        <dt>Specialite</dt>
        <dd>{subject.specialty_name}</dd>

        <dt>Type de groupe</dt>
        <dd>{subject.group_type_label}</dd>

        <dt>Statut</dt>
        <dd>
          <Badge
            variant={statusVariant(subject.status)}
            label={subject.status_label}
          />
        </dd>

        <dt>Propose par</dt>
        <dd>{subject.proposer_role_label}</dd>

        <dt>Cree le</dt>
        <dd>{subject.created_at_formatted}</dd>

        <dt>Modifie le</dt>
        <dd>{subject.updated_at_formatted}</dd>
      </dl>
    </section>

    <section class="card">
      <h2>Validation</h2>

      {#if subject.validator1_id || subject.validator2_id}
        <div class="validators">
          <div class="validator">
            <h3>Validateur 1</h3>
            {#if subject.validator1_decision}
              <Badge
                variant={reviewVariant(subject.validator1_decision)}
                label={subject.validator1_decision_label ?? "—"}
              />
              {#if subject.validator1_comment}
                <p class="comment">{subject.validator1_comment}</p>
              {/if}
            {:else}
              <p class="pending">En attente de decision</p>
            {/if}
          </div>

          <div class="validator">
            <h3>Validateur 2</h3>
            {#if subject.validator2_decision}
              <Badge
                variant={reviewVariant(subject.validator2_decision)}
                label={subject.validator2_decision_label ?? "—"}
              />
              {#if subject.validator2_comment}
                <p class="comment">{subject.validator2_comment}</p>
              {/if}
            {:else}
              <p class="pending">En attente de decision</p>
            {/if}
          </div>
        </div>
      {:else}
        <p class="pending">Aucun validateur assigne pour le moment.</p>
      {/if}
    </section>
  </div>
</Page>

<Modal
  open={showAssignModal}
  title="Assigner des validateurs"
  onClose={() => (showAssignModal = false)}
>
  {#if assignError}
    <div class="error-banner">{assignError}</div>
  {/if}

  <div class="recommendation-note">
    <Star size={16} />
    <span
      >Enseignants tries par pertinence — <strong
        >{subject.specialty_name}</strong
      ></span
    >
  </div>

  <div class="modal-form">
    <FormField label="Validateur 1" required>
      <select bind:value={validator1Id} required class="input">
        <option value="">Selectionner un enseignant</option>
        {#each teachers as teacher}
          <option value={teacher.id}>
            {teacher.profile_name} — {matchStars(teacher.match_score)}
            {matchLabel(teacher.match_score)}
          </option>
        {/each}
      </select>
    </FormField>

    <FormField label="Validateur 2" required>
      <select bind:value={validator2Id} required class="input">
        <option value="">Selectionner un enseignant</option>
        {#each teachers as teacher}
          <option value={teacher.id}>
            {teacher.profile_name} — {matchStars(teacher.match_score)}
            {matchLabel(teacher.match_score)}
          </option>
        {/each}
      </select>
    </FormField>
  </div>

  <div class="teacher-list">
    <h3>Recommandations</h3>
    {#each teachers as teacher}
      <div class="teacher-row" class:recommended={teacher.match_score > 0}>
        <div class="teacher-info">
          <span class="teacher-name">{teacher.profile_name}</span>
          <span class="teacher-specialties"
            >{teacher.specialty_names.join(", ") || "Aucune specialite"}</span
          >
        </div>
        <div class="teacher-match">
          {#if teacher.match_score > 0}
            <span class="stars">{matchStars(teacher.match_score)}</span>
            <span
              class="match-badge"
              class:exact={teacher.match_score === 3}
              class:partial={teacher.match_score === 1}
            >
              {matchLabel(teacher.match_score)}
            </span>
          {:else}
            <span class="no-match">Non recommandé</span>
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showAssignModal = false)}
      >Annuler</Button
    >
    <Button variant="primary" onclick={assignValidatorsAction}>Assigner</Button>
  </div>
</Modal>

<style>
  /* == Toast ==== */
  .toast {
    position: fixed;
    top: 1rem;
    right: 1rem;
    z-index: 9999;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    border-radius: 10px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    font-weight: 500;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    animation: slideIn 0.2s ease-out;
  }

  .toast-success {
    background: color-mix(in srgb, var(--color-success) 15%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-success) 30%, transparent);
    color: var(--color-success);
  }

  .toast-error {
    background: color-mix(in srgb, var(--color-danger) 15%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-danger) 30%, transparent);
    color: var(--color-danger);
  }

  @keyframes slideIn {
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  .detail-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-lg);
  }

  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-lg);
  }

  h2 {
    font-family: var(--font-sans);
    font-size: var(--text-lg);
    font-weight: 700;
    color: var(--color-text);
    margin: 0 0 var(--spacing-md);
  }

  dl {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: var(--spacing-sm) var(--spacing-md);
    margin: 0;
  }

  dt {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
    white-space: nowrap;
  }

  dd {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    color: var(--color-text);
    margin: 0;
  }

  .validators {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .validator {
    padding: var(--spacing-md);
    background: var(--color-background);
    border-radius: 8px;
    border: 1px solid var(--color-border);
  }

  .validator h3 {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
    margin: 0 0 var(--spacing-sm);
  }

  .comment {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: var(--spacing-sm) 0 0;
    font-style: italic;
  }

  .pending {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: 0;
    font-style: italic;
  }

  .error-banner {
    padding: var(--spacing-sm) var(--spacing-md);
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-danger) 20%, transparent);
    border-radius: 8px;
    color: var(--color-danger);
    font-size: var(--text-sm);
    margin-bottom: var(--spacing-md);
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-lg);
  }

  .modal-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  /* == Recommendation note ===== */
  .recommendation-note {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    margin-bottom: var(--spacing-md);
    background: color-mix(in srgb, var(--color-accent) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-accent) 20%, transparent);
    border-radius: 8px;
    font-size: var(--text-sm);
    color: var(--color-accent);
  }

  /* == Teacher recommendation list =========== */
  .teacher-list {
    margin-top: var(--spacing-lg);
    border-top: 1px solid var(--color-border);
    padding-top: var(--spacing-md);
  }

  .teacher-list h3 {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
    margin: 0 0 var(--spacing-sm);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .teacher-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0.75rem;
    border-radius: 8px;
    transition: background 0.15s;
  }

  .teacher-row.recommended {
    background: color-mix(in srgb, var(--color-accent) 5%, transparent);
  }

  .teacher-row + .teacher-row {
    border-top: 1px solid var(--color-border);
  }

  .teacher-info {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .teacher-name {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
  }

  .teacher-specialties {
    font-family: var(--font-sans);
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .teacher-match {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }

  .stars {
    font-size: var(--text-sm);
    color: #f59e0b;
    letter-spacing: 0.05em;
  }

  .match-badge {
    font-size: var(--text-xs);
    font-weight: 500;
    padding: 0.15rem 0.5rem;
    border-radius: 999px;
    white-space: nowrap;
  }

  .match-badge.exact {
    background: color-mix(in srgb, var(--color-success) 15%, transparent);
    color: var(--color-success);
  }

  .match-badge.partial {
    background: color-mix(in srgb, var(--color-warning) 15%, transparent);
    color: var(--color-warning);
  }

  .no-match {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    font-style: italic;
  }

  @media screen and (max-width: 768px) {
    .detail-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
