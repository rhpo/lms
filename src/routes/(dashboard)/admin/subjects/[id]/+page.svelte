<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { ArrowLeft, UserCheck } from "lucide-svelte";
  import { admin } from "$lib/api";
  import { showToast } from "$lib/utils/toast";
  import {
    SUBJECT_STATUS_LABELS,
    SUBJECT_STATUS_VARIANTS,
    GROUP_TYPE_LABELS,
    ROLE_LABELS,
    REVIEW_DECISION_LABELS,
    REVIEW_DECISION_VARIANTS,
    formatDate,
  } from "$lib/constants/labels";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";

  let { data } = $props();

  const subject = $derived(data.subject);
  const teachers = $derived(data.teachers ?? []);

  let showAssignModal = $state(false);
  let validator1Id: number | null = $state(null);
  let validator2Id: number | null = $state(null);
  let assigning = $state(false);

  async function unblockSubject() {
    try {
      await admin.subjectAction(subject!.id, "unblock");
      showToast.success("Sujet débloqué avec succès");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function assignValidators() {
    if (!validator1Id) {
      showToast.error("Validateur 1 requis");
      return;
    }
    assigning = true;
    try {
      await admin.subjectAction(subject!.id, "assign-validators", {
        validator1_id: validator1Id,
        validator2_id: validator2Id ?? 0,
      });
      showToast.success("Validateurs assignés - sujet repassé en attente");
      showAssignModal = false;
      validator1Id = null;
      validator2Id = null;
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    } finally {
      assigning = false;
    }
  }
</script>

{#if subject}
  <Page title={subject.title} subtitle="Détail du sujet PFE">
    {#snippet actions()}
      <a href="/admin/subjects">
        <Button variant="ghost" Icon={ArrowLeft}>Retour</Button>
      </a>
      {#if subject.status === "en_attente" || subject.status === "accepte_sous_reserve"}
        <Button
          variant="primary"
          Icon={UserCheck}
          onclick={() => (showAssignModal = true)}
        >
          Assigner des validateurs
        </Button>
      {/if}
      {#if subject.status === "refuse"}
        <Button variant="secondary" onclick={unblockSubject}>Débloquer</Button>
      {/if}
    {/snippet}

    <div class="detail-grid">
      <section class="card">
        <h2>Informations Générales</h2>
        <dl>
          <dt>Titre</dt>
          <dd>{subject.title}</dd>

          <dt>Description</dt>
          <dd>{subject.description}</dd>

          <dt>Type de groupe</dt>
          <dd>{GROUP_TYPE_LABELS[subject.group_type] ?? subject.group_type}</dd>

          <dt>Statut</dt>
          <dd>
            <Badge
              variant={SUBJECT_STATUS_VARIANTS[subject.status] ?? "info"}
              label={SUBJECT_STATUS_LABELS[subject.status] ?? subject.status}
            />
          </dd>

          <dt>Proposé par</dt>
          <dd>{ROLE_LABELS[subject.proposer_role] ?? subject.proposer_role}</dd>

          <dt>Créé le</dt>
          <dd>{formatDate(subject.created_at)}</dd>

          <dt>Modifié le</dt>
          <dd>{formatDate(subject.updated_at)}</dd>
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
                  variant={REVIEW_DECISION_VARIANTS[
                    subject.validator1_decision
                  ] ?? "info"}
                  label={REVIEW_DECISION_LABELS[subject.validator1_decision] ??
                    "-"}
                />
                {#if subject.validator1_comment}
                  <p class="comment">{subject.validator1_comment}</p>
                {/if}
              {:else}
                <p class="pending">En attente de décision</p>
              {/if}
            </div>

            <div class="validator">
              <h3>Validateur 2</h3>
              {#if subject.validator2_decision}
                <Badge
                  variant={REVIEW_DECISION_VARIANTS[
                    subject.validator2_decision
                  ] ?? "info"}
                  label={REVIEW_DECISION_LABELS[subject.validator2_decision] ??
                    "-"}
                />
                {#if subject.validator2_comment}
                  <p class="comment">{subject.validator2_comment}</p>
                {/if}
              {:else}
                <p class="pending">En attente de décision</p>
              {/if}
            </div>
          </div>
        {:else}
          <p class="pending">Aucun validateur assigné pour le moment.</p>
        {/if}
      </section>
    </div>
  </Page>

  <!-- Assign Validators Modal -->
  {#if showAssignModal}
    <Modal
      open={true}
      title="Assigner des validateurs"
      onClose={() => (showAssignModal = false)}
    >
      <div class="modal-form">
        {#if subject.status === "accepte_sous_reserve"}
          <p class="info-note">
            Ce sujet a été modifié par l'entreprise suite à une réserve. Les
            assigner repassera le sujet en <strong>En attente</strong> et effacera
            les décisions précédentes.
          </p>
        {/if}
        <div class="field">
          <label for="v1">Validateur 1 <span class="required">*</span></label>
          <select id="v1" class="input" bind:value={validator1Id}>
            <option value={null}>Sélectionner un enseignant</option>
            {#each teachers as t}
              {@const isProposer =
                subject.proposer_role === "teacher" &&
                t.id === subject.proposer_id}
              <option value={t.teacher?.id ?? t.id} disabled={isProposer}>
                {t.full_name}{isProposer ? " - (Proposeur du sujet)" : ""}
              </option>
            {/each}
          </select>
        </div>
        <div class="field">
          <label for="v2"
            >Validateur 2 <span class="optional">(optionnel)</span></label
          >
          <select id="v2" class="input" bind:value={validator2Id}>
            <option value={null}>Sélectionner un enseignant</option>
            {#each teachers as t}
              {@const isProposer =
                subject.proposer_role === "teacher" &&
                t.id === subject.proposer_id}
              <option value={t.teacher?.id ?? t.id} disabled={isProposer}>
                {t.full_name}{isProposer ? " - (Proposeur du sujet)" : ""}
              </option>
            {/each}
          </select>
        </div>
        <div class="modal-actions">
          <Button
            variant="ghost"
            onclick={() => (showAssignModal = false)}
            disabled={assigning}
          >
            Annuler
          </Button>
          <Button
            onclick={assignValidators}
            disabled={assigning || !validator1Id}
          >
            {assigning ? "Assignation..." : "Assigner"}
          </Button>
        </div>
      </div>
    </Modal>
  {/if}
{:else}
  <Page title="Sujet introuvable" subtitle="">
    <p>Ce sujet n'existe pas ou a été supprimé.</p>
    <a href="/admin/subjects">
      <Button variant="ghost" Icon={ArrowLeft}>Retour</Button>
    </a>
  </Page>
{/if}

<style>
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

  .modal-form {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .info-note {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    background: color-mix(
      in srgb,
      var(--color-warning) 8%,
      var(--color-surface)
    );
    border: 1px solid color-mix(in srgb, var(--color-warning) 25%, transparent);
    border-radius: 8px;
    padding: var(--spacing-sm) var(--spacing-md);
    margin: 0;
    line-height: 1.5;
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

  option:disabled {
    color: var(--color-text-muted);
    font-style: italic;
  }

  .required {
    color: var(--color-danger);
  }
  .optional {
    font-weight: 400;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-sm);
  }

  @media screen and (max-width: 768px) {
    .detail-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
