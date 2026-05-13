<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { Check, CheckCircle, Clock, X, XCircle } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import Select from "$lib/components/ui/Select.svelte";

  const specialities = [
    { id: "spec-1", name: "ISIL", year_type: "ingenieur" },
    { id: "spec-2", name: "CHIM", year_type: "master" },
    { id: "spec-3", name: "ELEC", year_type: "ingenieur" },
  ];

  const niveaux = [
    { value: "licence", label: "Licence" },
    { value: "master", label: "Master" },
    { value: "ingenieur", label: "Ingenieur" },
  ];

  const subjects: Array<{
    id: string;
    title: string;
    specialty: string;
    group_type_label: string;
    proposer_role_label: string;
    status: string;
    status_variant: string;
    status_label: string;
    created_at_formatted: string;
  }> = [
    {
      id: "subj-1",
      title: "Application mobile de suivi des stages",
      specialty: "ISIL",
      group_type_label: "Binome",
      proposer_role_label: "Enseignant",
      status: "en_attente",
      status_variant: "warning",
      status_label: "En attente",
      created_at_formatted: "2025-03-10",
    },
    {
      id: "subj-2",
      title: "Systeme de recommandation pour e-commerce",
      specialty: "ISIL",
      group_type_label: "Trinome",
      proposer_role_label: "Entreprise",
      status: "valide",
      status_variant: "success",
      status_label: "Valide",
      created_at_formatted: "2025-02-20",
    },
    {
      id: "subj-3",
      title: "Analyse des donnees de production",
      specialty: "ELEC",
      group_type_label: "Monome",
      proposer_role_label: "Enseignant",
      status: "refuse",
      status_variant: "danger",
      status_label: "Refuse",
      created_at_formatted: "2025-01-15",
    },
  ];

  const teachers = [
    { id: "tch-1", profile_name: "Dr. Amine Benali" },
    { id: "tch-2", profile_name: "Dr. Sarah Khelif" },
    { id: "tch-3", profile_name: "Pr. Mohamed Lamine" },
  ];

  let currentNiveau = $state($page.url.searchParams.get("niveau") ?? "");
  let currentSpecialty = $state($page.url.searchParams.get("specialty") ?? "");
  let statusFilter = $state($page.url.searchParams.get("status") ?? "");

  const filteredSpecialities = $derived(
    currentNiveau
      ? specialities.filter((s) => s.year_type === currentNiveau)
      : specialities,
  );

  function onNiveauChange(val: string) {
    currentNiveau = val;
    currentSpecialty = "";
    goto(`?niveau=${val || ""}&specialty=${""}&status=${statusFilter}`);
  }

  function onSpecialtyChange(val: string) {
    currentSpecialty = val;
    goto(
      `?niveau=${currentNiveau || ""}&specialty=${val || ""}&status=${statusFilter}`,
    );
  }

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
  let assignSubjectId = $state("");
  let assignError = $state("");
  let validator1Id = $state("");
  let validator2Id = $state("");

  function openAssignModal(subjectId: string) {
    assignSubjectId = subjectId;
    assignError = "";
    validator1Id = "";
    validator2Id = "";
    showAssignModal = true;
  }

  async function assignValidatorsAction() {
    assignError = "";
    try {
      const res = await fetch("/api/subjects/assign-validators", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          subject_id: assignSubjectId,
          validator1_id: validator1Id,
          validator2_id: validator2Id,
        }),
      });
      if (res.ok) {
        showAssignModal = false;
        assignSubjectId = "";
        showToast("Validateurs assignes avec succes");
        await invalidateAll();
      } else {
        const err = await res.json();
        assignError = err.message ?? "Erreur lors de l'assignation";
      }
    } catch {
      assignError = "Erreur reseau";
    }
  }

  async function unblockSubject(id: string) {
    try {
      const res = await fetch("/api/subjects/unblock", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id }),
      });
      if (!res.ok) {
        const err = await res.json();
        showToast(err.message ?? "Erreur lors du deblocage", "error");
        return;
      }
      showToast("Sujet debloque avec succes");
      await invalidateAll();
    } catch {
      showToast("Erreur reseau", "error");
    }
  }
</script>

<Page
  title="Sujets"
  subtitle="Gérer tous les sujets PFE, assigner des validateurs, debloquer les sujets en attente"
>
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

  <div class="header">
    <div class="specialities">
      <Select
        value={currentNiveau}
        options={niveaux}
        placeholder="Tous les niveaux"
        label="Niveau"
        onChange={onNiveauChange}
      />
      <Select
        value={currentSpecialty}
        options={(filteredSpecialities ?? []).map((s) => ({
          value: s.name,
          label: s.name,
        }))}
        placeholder="Toutes les specialites"
        label="Specialite"
        onChange={onSpecialtyChange}
      />
    </div>

    <div class="filters">
      <a href="?" class:is-active={!statusFilter}>Tous</a>

      <a
        href="?status=en_attente"
        class:is-active={statusFilter === "en_attente"}
      >
        <Clock size="16px" />
        En attente</a
      >
      <a href="?status=valide" class:is-active={statusFilter === "valide"}>
        <Check size="16px" />
        Valides</a
      >
      <a href="?status=refuse" class:is-active={statusFilter === "refuse"}>
        <X size="16px" />
        Refuses</a
      >
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th>Titre</th>
        <th>Specialite</th>
        <th>Groupe</th>
        <th>Propose par</th>
        <th>Statut</th>
        <th>Date</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each subjects as subject}
        <tr>
          <td>
            <a href="/admin/subjects/{subject.id}">{subject.title}</a>
          </td>
          <td>{subject.specialty}</td>
          <td>{subject.group_type_label}</td>
          <td>{subject.proposer_role_label}</td>
          <td>
            <Badge
              variant={subject.status_variant}
              label={subject.status_label}
            />
          </td>
          <td>{subject.created_at_formatted}</td>
          <td>
            <div class="actions-cell">
              {#if subject.status === "en_attente"}
                <Button
                  variant="ghost"
                  size="sm"
                  onclick={() => openAssignModal(subject.id)}
                >
                  Assigner validateurs
                </Button>
              {/if}
              {#if subject.status === "refuse"}
                <Button
                  variant="ghost"
                  size="sm"
                  onclick={() => unblockSubject(subject.id)}
                >
                  Debloquer
                </Button>
              {/if}
            </div>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</Page>

<Modal
  open={showAssignModal}
  title="Assigner des validateurs"
  onClose={() => (showAssignModal = false)}
>
  {#if assignError}
    <div class="error-banner">{assignError}</div>
  {/if}

  <div class="modal-form">
    <FormField label="Validateur 1" required>
      <select bind:value={validator1Id} required class="input">
        <option value="">Selectionner un enseignant</option>
        {#each teachers as teacher}
          <option value={teacher.id}>{teacher.profile_name}</option>
        {/each}
      </select>
    </FormField>

    <FormField label="Validateur 2" required>
      <select bind:value={validator2Id} required class="input">
        <option value="">Selectionner un enseignant</option>
        {#each teachers as teacher}
          <option value={teacher.id}>{teacher.profile_name}</option>
        {/each}
      </select>
    </FormField>
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

  .header {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: end;

    margin-bottom: var(--spacing-lg);
  }

  .specialities {
    display: flex;
    gap: var(--spacing-md);
    flex-wrap: wrap;
  }

  .filters {
    display: flex;
    gap: var(--spacing-xs);

    a {
      padding: var(--spacing-xs) var(--spacing-md);
      border-radius: 999px;
      font-size: var(--text-sm);
      font-family: var(--font-sans);
      color: var(--color-text-muted);
      text-decoration: none;
      border: 1px solid var(--color-border);
      transition:
        background var(--transition-fast),
        color var(--transition-fast);

      display: flex;
      align-items: center;
      gap: 0.5rem;

      &:hover {
        background: var(--color-background-100);
      }

      &:active {
        background: var(--color-accent);
        color: #fff;
        border-color: var(--color-accent);
      }
    }
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
</style>
