<script lang="ts">
  import { goto } from "$app/navigation";
  import { invalidateAll } from "$app/navigation";
  import { ArrowLeft } from "lucide-svelte";
  import { teacher } from "$lib/api";

  import Button from "$lib/components/ui/Button.svelte";
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

  const niveaux = [
    { value: "licence", label: "Licence" },
    { value: "master", label: "Master" },
    { value: "ingenieur", label: "Ingenieur" },
  ] as const;

  let selectedNiveau = $state<string>("");
  let selectedSpecialtyId = $state<string>("");
  let title = $state("");
  let description = $state("");
  let groupType = $state<"monome" | "binome" | "trinome">("monome");
  let error = $state("");
  let loading = $state(false);

  const filteredSpecialities = $derived(
    selectedNiveau
      ? specialities.filter((s) => s.year_type === selectedNiveau)
      : [],
  );

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    if (!selectedSpecialtyId) {
      error = "Veuillez selectionner une specialite.";
      loading = false;
      return;
    }

    try {
      await teacher.createProposedSubject({
        title,
        description,
        specialty: selectedSpecialtyId,
        level: selectedNiveau,
        group_type: groupType,
      });

      await invalidateAll();
      goto("/teacher/proposed-subjects");
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      loading = false;
    }
  }
</script>

<Page
  title="Proposer un nouveau sujet"
  subtitle="Remplissez le formulaire pour soumettre un sujet PFE."
>
  {#snippet actions()}
    <Button
      variant="ghost"
      Icon={ArrowLeft}
      onclick={() => goto("/teacher/proposed-subjects")}
    >
      Retour
    </Button>
  {/snippet}

  {#if error}
    <div class="error">{error}</div>
  {/if}

  <form onsubmit={handleSubmit}>
    <div class="form-grid">
      <div class="field">
        <label for="title">Titre du sujet</label>
        <input
          id="title"
          type="text"
          bind:value={title}
          required
          placeholder="Ex. Conception d'un systeme de recommandation"
        />
      </div>

      <div class="field">
        <label for="group_type">Type de groupe</label>
        <select id="group_type" bind:value={groupType} required>
          <option value="monome">Monome (1 etudiant)</option>
          <option value="binome">Binome (2 etudiants)</option>
          <option value="trinome">Trinome (3 etudiants)</option>
        </select>
      </div>

      <div class="field">
        <label for="niveau">Niveau</label>
        <select
          id="niveau"
          bind:value={selectedNiveau}
          required
          onchange={() => (selectedSpecialtyId = "")}
        >
          <option value="">Selectionner un niveau</option>
          {#each niveaux as n}
            <option value={n.value}>{n.label}</option>
          {/each}
        </select>
      </div>

      <div class="field">
        <label for="specialty">Specialite</label>
        <select
          id="specialty"
          bind:value={selectedSpecialtyId}
          required
          disabled={!selectedNiveau}
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
      </div>

      <div class="field full-width">
        <label for="description">Description</label>
        <textarea
          id="description"
          bind:value={description}
          required
          rows={6}
          placeholder="Decrivez le sujet en detail : objectifs, technologies, livrables attendus..."
        ></textarea>
      </div>
    </div>

    <div class="actions">
      <Button
        variant="ghost"
        onclick={() => goto("/teacher/proposed-subjects")}
      >
        Annuler
      </Button>
      <Button variant="primary" type="submit" disabled={loading}>
        {loading ? "Soumission..." : "Soumettre le sujet"}
      </Button>
    </div>
  </form>
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

  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-lg);

    @media screen and (max-width: 600px) {
      & {
        grid-template-columns: 1fr;
      }
    }
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .full-width {
    grid-column: 1 / -1;
  }

  label {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
  }

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
    min-height: 120px;
  }

  .actions {
    display: flex;
    gap: var(--spacing-sm);
    justify-content: flex-end;
  }
</style>
