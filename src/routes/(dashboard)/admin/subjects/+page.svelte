<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { Check, Clock, X, Sparkles, Star } from "lucide-svelte";
  import { admin } from "$lib/api";
  import { showToast } from "$lib/utils/toast";
  import type {
    PfeSubject,
    Profile,
    Domain,
    ValidatorRecommendation,
  } from "$lib/types";
  import {
    SUBJECT_STATUS_LABELS,
    SUBJECT_STATUS_VARIANTS,
    GROUP_TYPE_LABELS,
    ROLE_LABELS,
    formatDate,
  } from "$lib/constants/labels";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();
  const subjects: PfeSubject[] = $derived(data.subjects);
  const teachers: Profile[] = $derived(data.teachers);

  let statusFilter = $state("");

  const filteredSubjects = $derived(() => {
    if (!statusFilter) return subjects;

    return subjects.filter((s) => s.status === statusFilter);
  });

  let showAssignModal = $state(false);
  let assignSubjectId = $state(0);
  let assignError = $state("");
  let validator1Id = $state(0);
  let validator2Id = $state(0);
  let loadingRecommendations = $state(false);

  let recommendations: ValidatorRecommendation[] = $state([]);
  let subjectDomains: Domain[] = $state([]);

  async function openAssignModal(subjectId: number) {
    assignSubjectId = subjectId;
    assignError = "";
    validator1Id = 0;
    validator2Id = 0;
    recommendations = [];
    subjectDomains = [];
    showAssignModal = true;
    await loadRecommendations(subjectId);
  }

  async function loadRecommendations(subjectId: number) {
    loadingRecommendations = true;
    try {
      const result = await admin.recommendJury(subjectId);
      recommendations = result.recommended ?? [];
      subjectDomains = result.subject_domains ?? [];
    } catch {
      recommendations = [];
    } finally {
      loadingRecommendations = false;
    }
  }

  function selectRecommended(rec: ValidatorRecommendation, slot: 1 | 2) {
    const teacherId = rec.teacher.id;

    if (slot === 1) {
      validator1Id = teacherId;
    } else {
      validator2Id = teacherId;
    }
  }

  async function assignValidatorsAction() {
    assignError = "";
    if (!validator1Id || !validator2Id) {
      assignError = "Veuillez sélectionner les deux validateurs";
      return;
    }
    if (validator1Id === validator2Id) {
      assignError = "Les deux validateurs doivent être différents";
      return;
    }
    try {
      await admin.subjectAction(assignSubjectId, "assign-validators", {
        validator1_id: validator1Id,
        validator2_id: validator2Id,
      });
      showAssignModal = false;
      showToast.success("Validateurs assignés avec succès");
      await invalidateAll();
    } catch (err) {
      assignError = err instanceof Error ? err.message : "Erreur réseau";
    }
  }

  async function unblockSubject(id: number) {
    try {
      await admin.subjectAction(id, "unblock");
      showToast.success("Sujet débloqué avec succès");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }
</script>

<Page
  title="Sujets"
  subtitle="Gérer tous les sujets PFE, assigner des validateurs, débloquer les sujets en attente"
>
  <div class="header">
    <div class="filters">
      <Button
        variant={!statusFilter ? "secondary" : "ghost"}
        onclick={() => (statusFilter = "")}
      >
        Tous
      </Button>
      <Button
        variant={statusFilter === "en_attente" ? "secondary" : "ghost"}
        onclick={() => (statusFilter = "en_attente")}
      >
        <Clock size={16} />
        En attente
      </Button>
      <Button
        variant={statusFilter === "valide" ? "secondary" : "ghost"}
        onclick={() => (statusFilter = "valide")}
      >
        <Check size={16} />
        Validés
      </Button>
      <Button
        variant={statusFilter === "refuse" ? "secondary" : "ghost"}
        onclick={() => (statusFilter = "refuse")}
      >
        <X size={16} />
        Refusés
      </Button>
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th>Titre</th>
        <th>Groupe</th>
        <th>Proposé par</th>
        <th>Statut</th>
        <th>Date</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each filteredSubjects()
        .sort((a, b) => b.created_at.localeCompare(a.created_at))

        .sort((a, b) => {
          if (a.status === "en_attente" && b.status !== "en_attente") return -1;
          if (a.status !== "en_attente" && b.status === "en_attente") return 1;
          return 0;
        }) as subject}
        <tr>
          <td>
            <a href="/admin/subjects/{subject.id}">{subject.title}</a>
          </td>
          <td>{GROUP_TYPE_LABELS[subject.group_type] ?? subject.group_type}</td>
          <td class="proposer-cell">
            <span class="proposer-name">
              {subject.proposer?.full_name ?? "-"}
            </span>
            {#if subject.proposer_role === "company" && subject.company?.company_name}
              <span class="proposer-company"
                >{subject.company.company_name}</span
              >
            {:else}
              <span class="proposer-role"
                >{ROLE_LABELS[subject.proposer_role] ??
                  subject.proposer_role}</span
              >
            {/if}
          </td>
          <td class="badges">
            <div class="wrapper">
              <Badge
                variant={SUBJECT_STATUS_VARIANTS[subject.status] ?? "info"}
                label={SUBJECT_STATUS_LABELS[subject.status] ?? subject.status}
              />

              {#if subject.status === "en_attente" && subject.validator1_id && subject.validator2_id}
                <Badge variant="info" label="Validateurs assignés" />
              {/if}

              {#if subject.co_supervisor_id}
                <Badge variant="success" label="Co-encadrant assigné" />
              {/if}
            </div>
          </td>
          <td>{formatDate(subject.created_at)}</td>
          <td>
            <div class="actions-cell">
              {#if subject.status === "en_attente" || subject.status === "accepte_sous_reserve"}
                <Button
                  size="sm"
                  disabled={subject.status !== "accepte_sous_reserve" &&
                    !!(subject.validator1_id || subject.validator2_id)}
                  onclick={() => openAssignModal(subject.id)}
                >
                  Assigner validateurs
                </Button>
              {/if}

              {#if subject.status === "refuse"}
                <Button size="sm" onclick={() => unblockSubject(subject.id)}>
                  Débloquer
                </Button>
              {/if}
            </div>
          </td>
        </tr>
      {:else}
        <tr>
          <td colspan="6" style="text-align: center; padding: 2rem;">
            Aucun sujet trouvé.
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</Page>

<!-- Assign Validators Modal with Recommendations -->
<Modal
  open={showAssignModal}
  title="Assigner des validateurs"
  onClose={() => (showAssignModal = false)}
  width="640px"
>
  {#if assignError}
    <div class="error-banner">{assignError}</div>
  {/if}

  <!-- Recommendation section -->
  {#if subjectDomains.length > 0}
    <div class="domain-tags">
      <span class="domain-label">Domaines du sujet :</span>
      {#each subjectDomains as domain}
        <span class="domain-tag">{domain.name}</span>
      {/each}
    </div>
  {/if}

  {#if loadingRecommendations}
    <div class="rec-loading">Chargement des recommandations...</div>
  {:else if recommendations.length > 0}
    <div class="rec-section">
      <div class="rec-header">
        <Sparkles size={16} />
        <span>Recommandations basées sur les domaines</span>
      </div>
      <div class="rec-list">
        {#each recommendations.slice(0, 6) as rec}
          <div class="rec-card" class:has-match={rec.score > 0}>
            <div class="rec-info">
              <span class="rec-name"
                >{rec.teacher.profile?.full_name ?? "Enseignant"}</span
              >
              {#if rec.score > 0}
                <div class="rec-domains">
                  {#each rec.matching_domains as d}
                    <span class="match-tag">
                      <Star size={10} />
                      {d.name}
                    </span>
                  {/each}
                </div>
              {:else}
                <span class="rec-no-match">Aucun domaine en commun</span>
              {/if}
            </div>
            <div class="rec-score">
              <Badge
                variant={rec.score > 0 ? "success" : "info"}
                label="{rec.score} match{rec.score !== 1 ? 'es' : ''}"
              />
            </div>
            <div class="rec-actions">
              <button
                class="rec-btn"
                class:selected={validator1Id === rec.teacher.id}
                onclick={() => selectRecommended(rec, 1)}
              >
                V1
              </button>
              <button
                class="rec-btn"
                class:selected={validator2Id === rec.teacher.id}
                onclick={() => selectRecommended(rec, 2)}
              >
                V2
              </button>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <div class="modal-form">
    <FormField label="Validateur 1" required>
      <select required class="input" bind:value={validator1Id}>
        <option value={0}>Sélectionner un enseignant</option>
        {#each teachers.filter((t) => t.teacher) as t}
          <option value={t.teacher!.id}>{t.full_name}</option>
        {/each}
      </select>
    </FormField>

    <FormField label="Validateur 2" required>
      <select required class="input" bind:value={validator2Id}>
        <option value={0}>Sélectionner un enseignant</option>
        {#each teachers.filter((t) => t.teacher) as t}
          <option value={t.teacher!.id}>{t.full_name}</option>
        {/each}
      </select>
    </FormField>
  </div>

  <div class="form-actions">
    <Button onclick={() => (showAssignModal = false)}>Annuler</Button>
    <Button variant="primary" onclick={assignValidatorsAction}>Assigner</Button>
  </div>
</Modal>

<style>
  .header {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: end;
    margin-bottom: var(--spacing-lg);
  }

  .filters {
    display: flex;
    gap: var(--spacing-xs);
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
  }

  thead th {
    text-align: left;
    padding: var(--spacing-sm) var(--spacing-md);
    font-weight: 600;
    color: var(--color-text-muted);
    border-bottom: 2px solid var(--color-border);
    white-space: nowrap;
  }

  tbody tr {
    border-bottom: 1px solid var(--color-border);
  }

  tbody tr:hover {
    background: var(--color-background-100);
  }

  tbody td {
    padding: var(--spacing-sm) var(--spacing-md);
    color: var(--color-text);
  }

  tbody td a {
    color: var(--color-accent);
    text-decoration: none;
  }

  tbody td a:hover {
    text-decoration: underline;
  }

  .proposer-cell {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .proposer-name {
    font-weight: 600;
    color: var(--color-text);
  }

  .proposer-company {
    font-size: var(--text-xs);
    color: var(--color-accent);
    font-weight: 500;
  }

  .proposer-role {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .badges {
    .wrapper {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
  }

  .actions-cell {
    display: flex;
    gap: var(--spacing-xs);
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

  /* Recommendation styles */
  .domain-tags {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: var(--spacing-md);
    padding: var(--spacing-sm);
    background: var(--color-background);
    border-radius: 8px;
  }

  .domain-label {
    font-size: var(--text-xs);
    font-weight: 600;
    color: var(--color-text-muted);
  }

  .domain-tag {
    font-size: var(--text-xs);
    padding: 0.2rem 0.5rem;
    border-radius: 999px;
    background: color-mix(in srgb, var(--color-accent) 15%, transparent);
    color: var(--color-accent);
    font-weight: 500;
  }

  .rec-loading {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    padding: var(--spacing-md);
    text-align: center;
  }

  .rec-section {
    margin-bottom: var(--spacing-lg);
  }

  .rec-header {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
    margin-bottom: var(--spacing-sm);
  }

  .rec-list {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .rec-card {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background: var(--color-surface);
    transition: border-color var(--transition-fast);
  }

  .rec-card.has-match {
    border-color: color-mix(in srgb, var(--color-success) 40%, transparent);
    background: color-mix(in srgb, var(--color-success) 3%, transparent);
  }

  .rec-info {
    flex: 1;
    min-width: 0;
  }

  .rec-name {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
  }

  .rec-domains {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
    margin-top: 0.2rem;
  }

  .match-tag {
    display: inline-flex;
    align-items: center;
    gap: 0.2rem;
    font-size: 0.7rem;
    padding: 0.1rem 0.4rem;
    border-radius: 999px;
    background: color-mix(in srgb, var(--color-success) 15%, transparent);
    color: var(--color-success);
    font-weight: 500;
  }

  .rec-no-match {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .rec-score {
    flex-shrink: 0;
  }

  .rec-actions {
    display: flex;
    gap: 0.25rem;
    flex-shrink: 0;
  }

  .rec-btn {
    padding: 0.25rem 0.5rem;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    background: var(--color-background);
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    font-weight: 600;
    font-family: var(--font-sans);
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .rec-btn:hover {
    border-color: var(--color-accent);
    color: var(--color-accent);
  }

  .rec-btn.selected {
    background: var(--color-accent);
    border-color: var(--color-accent);
    color: #fff;
  }
</style>
