<script lang="ts">
  import { goto } from "$app/navigation";
  import { invalidateAll } from "$app/navigation";
  import { ArrowLeft } from "lucide-svelte";
  import { teacher } from "$lib/api";
  import { atomic } from "$lib/stores/atomic.svelte";

  const domains = $derived(atomic.domains);

  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  const specialities = $derived(atomic.specialities);

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
  let selectedDomainIds = $state<number[]>([]);
  let error = $state("");
  let loading = $state(false);

  function toggleDomain(id: number) {
    if (selectedDomainIds.includes(id)) {
      selectedDomainIds = selectedDomainIds.filter((d) => d !== id);
    } else {
      selectedDomainIds = [...selectedDomainIds, id];
    }
  }

  const filteredSpecialities = $derived(
    selectedNiveau
      ? specialities.filter((s) => s.year_type === selectedNiveau)
      : [],
  );

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    if (title.trim().length < 5) {
      error = "Le titre doit contenir au moins 5 caractères.";
      loading = false;
      return;
    }
    if (title.trim().length > 200) {
      error = "Le titre ne doit pas dépasser 200 caractères.";
      loading = false;
      return;
    }
    if (description.trim().length < 20) {
      error = "La description doit contenir au moins 20 caractères.";
      loading = false;
      return;
    }
    if (!selectedSpecialtyId) {
      error = "Veuillez selectionner une specialite.";
      loading = false;
      return;
    }

    try {
      await teacher.createProposedSubject({
        title,
        description,
        group_type: groupType,
        domain_ids:
          selectedDomainIds.length > 0 ? selectedDomainIds : undefined,
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
          minlength={5}
          maxlength={200}
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
            <option value={spec.id}>{spec.code} - {spec.name}</option>
          {/each}
        </select>
      </div>

      <div class="field full-width">
        <span class="field-heading">
          Domaines de spécialité
          {#if selectedDomainIds.length > 0}
            <span class="domain-count"
              >({selectedDomainIds.length} sélectionné{selectedDomainIds.length >
              1
                ? "s"
                : ""})</span
            >
          {/if}
        </span>
        <p class="domain-hint">
          Sélectionnez les domaines liés à votre sujet - ils servent à
          identifier les meilleurs validateurs.
        </p>
        <div class="domains-grid">
          {#each domains as dom}
            <label class="checkbox-label">
              <input
                type="checkbox"
                checked={selectedDomainIds.includes(dom.id)}
                onchange={() => toggleDomain(dom.id)}
              />
              <span class="checkbox-text">{dom.name}</span>
            </label>
          {/each}
        </div>
      </div>
    </div>

    <div class="field full-width">
      <label for="description">Description</label>
      <textarea
        id="description"
        bind:value={description}
        required
        minlength={20}
        rows={6}
        placeholder="Decrivez le sujet en detail : objectifs, technologies, livrables attendus..."
      ></textarea>
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

  .domain-hint {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    font-family: var(--font-sans);
    margin: 0 0 0.5rem;
  }

  .domain-count {
    font-weight: 400;
    color: var(--color-accent);
    margin-left: 0.35rem;
  }

  .domains-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    background: color-mix(
      in srgb,
      var(--color-surface) 50%,
      var(--color-background)
    );
    border: 1px solid var(--color-border);
    border-radius: 8px;
  }

  .field-heading {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
    display: block;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
    font-weight: 400;
    width: fit-content;

    [type="checkbox"] {
      width: fit-content;
    }
  }

  .checkbox-label input[type="checkbox"] {
    accent-color: var(--color-accent);
    cursor: pointer;
  }

  .checkbox-text {
    font-size: var(--text-sm);
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
