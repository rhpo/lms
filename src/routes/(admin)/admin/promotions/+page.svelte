<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { Plus } from "lucide-svelte";

  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  type Speciality = {
    id: string;
    name: string;
    code: string;
    year_type: "licence" | "master" | "ingenieur";
  };

  const specialities: Speciality[] = [
    {
      id: "spec-1",
      name: "Ingenierie des Systemes d'Information",
      code: "ISIL",
      year_type: "ingenieur",
    },
    {
      id: "spec-2",
      name: "Chimie des Materiaux",
      code: "CHIM",
      year_type: "master",
    },
    {
      id: "spec-3",
      name: "Electrotechnique",
      code: "ELEC",
      year_type: "ingenieur",
    },
  ];

  const promotions = [
    {
      id: "promo-1",
      label: "Promotion ISIL 2024-2025",
      specialty: "ISIL",
      academic_year_label: "2024/2025",
      created_at_formatted: "2024-09-15",
    },
    {
      id: "promo-2",
      label: "Promotion CHIM 2024-2025",
      specialty: "CHIM",
      academic_year_label: "2024/2025",
      created_at_formatted: "2024-09-15",
    },
  ];

  const academicYears = [
    { id: "ay-1", label: "2024/2025" },
    { id: "ay-2", label: "2025/2026" },
  ];

  const niveaux = [
    { value: "licence", label: "Licence" },
    { value: "master", label: "Master" },
    { value: "ingenieur", label: "Ingenieur" },
  ] as const;

  let showCreateModal = $state(false);
  let createError = $state("");

  let label = $state("");
  let selectedNiveau = $state<string>("");
  let selectedSpecialtyId = $state<string>("");
  let academicYearId = $state("");

  const filteredSpecialities = $derived(
    selectedNiveau
      ? specialities.filter((s) => s.year_type === selectedNiveau)
      : [],
  );

  async function createPromotion() {
    createError = "";
    try {
      const res = await fetch("/api/promotions/create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          label,
          specialty: selectedSpecialtyId,
          academic_year_id: academicYearId,
        }),
      });
      if (res.ok) {
        showCreateModal = false;
        label = "";
        selectedSpecialtyId = "";
        selectedNiveau = "";
        academicYearId = "";
        await invalidateAll();
      } else {
        const err = await res.json();
        createError = err.message ?? "Erreur lors de la creation";
      }
    } catch {
      createError = "Erreur reseau";
    }
  }
</script>

<Page
  title="Promotions"
  subtitle="Gérer les promotions par specialite et annee universitaire"
>
  {#snippet actions()}
    <Button
      variant="primary"
      Icon={Plus}
      onclick={() => (showCreateModal = true)}
    >
      Nouvelle promotion
    </Button>
  {/snippet}

  <table>
    <thead>
      <tr>
        <th>Libelle</th>
        <th>Specialite</th>
        <th>Annee universitaire</th>
        <th>Date de creation</th>
      </tr>
    </thead>
    <tbody>
      {#each promotions as promo}
        <tr>
          <td>{promo.label}</td>
          <td>{promo.specialty}</td>
          <td>{promo.academic_year_label}</td>
          <td>{promo.created_at_formatted}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</Page>

<Modal
  open={showCreateModal}
  title="Nouvelle promotion"
  onClose={() => (showCreateModal = false)}
>
  {#if createError}
    <div class="error-banner">{createError}</div>
  {/if}

  <div class="modal-form">
    <FormField label="Libelle" required>
      <input
        type="text"
        bind:value={label}
        placeholder="ex. Promotion ISIL 2024-2025"
        required
        class="input"
      />
    </FormField>

    <FormField label="Niveau" required>
      <select
        bind:value={selectedNiveau}
        required
        class="input"
        onchange={() => (selectedSpecialtyId = "")}
      >
        <option value="">Selectionner un niveau</option>
        {#each niveaux as n}
          <option value={n.value}>{n.label}</option>
        {/each}
      </select>
    </FormField>

    <FormField label="Specialite" required>
      <select
        bind:value={selectedSpecialtyId}
        required
        disabled={!selectedNiveau}
        class="input"
      >
        <option value="">
          {selectedNiveau
            ? "Selectionner une specialite"
            : "Choisissez d'abord un niveau"}
        </option>
        {#each filteredSpecialities as spec}
          <option value={spec.id}>{spec.code} — {spec.name}</option>
        {/each}
      </select>
    </FormField>

    <FormField label="Annee universitaire" required>
      <select bind:value={academicYearId} required class="input">
        <option value="">Selectionner une annee</option>
        {#each academicYears as year}
          <option value={year.id}>{year.label}</option>
        {/each}
      </select>
    </FormField>
  </div>

  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showCreateModal = false)}
      >Annuler</Button
    >
    <Button variant="primary" onclick={createPromotion}>Creer</Button>
  </div>
</Modal>

<style>
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
