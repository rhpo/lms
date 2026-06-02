<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { Plus } from "lucide-svelte";
  import { admin, ref } from "$lib/api";
  import { formatDate } from "$lib/constants/labels";

  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();
  const { promotions } = $derived(data);

  let showCreateModal = $state(false);
  let createError = $state("");
  let label = $state("");
  let academicYearId = $state(0);
  let academicYears = $state<{ id: number; label: string }[]>([]);


  $effect(() => {
    if (showCreateModal && academicYears.length === 0) {
      admin.listAcademicYears().then((years) => {
        academicYears = years ?? [];
        if (years?.length) academicYearId = years[0].id;
      }).catch(() => {});
    }
  });

  async function createPromotion() {
    createError = "";
    if (!label.trim() || !academicYearId) {
      createError = "Veuillez remplir tous les champs.";
      return;
    }
    try {
      await admin.createPromotion({ label, academic_year_id: academicYearId });
      showCreateModal = false;
      label = "";
      academicYearId = 0;
      await invalidateAll();
    } catch (err) {
      createError = err instanceof Error ? err.message : "Erreur reseau";
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
        <th>Annee universitaire</th>
        <th>Date de creation</th>
      </tr>
    </thead>
    <tbody>
      {#each promotions as promo}
        <tr>
          <td>{promo.label}</td>
          <td>{promo.academic_year_id}</td>
          <td>{formatDate(promo.created_at)}</td>
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

    <FormField label="Annee universitaire" required>
      <select required class="input" value={academicYearId} onchange={(e) => academicYearId = Number(e.currentTarget.value)}>
        <option value={0}>Selectionner une annee</option>
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
