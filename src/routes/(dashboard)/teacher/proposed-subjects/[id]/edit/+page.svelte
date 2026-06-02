<script lang="ts">
  import { goto } from "$app/navigation";
  import { ArrowLeft, AlertTriangle } from "lucide-svelte";
  import { teacher } from "$lib/api";
  import { atomic } from "$lib/stores/atomic.svelte";

  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();
  const { subject } = $derived(data);

  const domains = $derived(atomic.domains);
  const specialities = $derived(atomic.specialities);

  const niveaux = [
    { value: "licence", label: "Licence" },
    { value: "master", label: "Master" },
    { value: "ingenieur", label: "Ingénieur" },
  ] as const;

  let title = $state(subject?.title ?? "");
  let description = $state(subject?.description ?? "");
  let groupType = $state<"monome" | "binome" | "trinome">(
    (subject?.group_type as "monome" | "binome" | "trinome") ?? "monome",
  );
  let selectedDomainIds = $state<number[]>(
    subject?.domains?.map((d: { id: number }) => d.id) ?? [],
  );
  let selectedNiveau = $state<string>("");
  let selectedSpecialtyId = $state<string>("");

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

  const validator1Comment = $derived(subject?.validator1_comment ?? null);
  const validator2Comment = $derived(subject?.validator2_comment ?? null);
  const hasComments = $derived(
    Boolean(validator1Comment) || Boolean(validator2Comment),
  );

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (!subject) return;
    error = "";

    if (title.trim().length < 5) {
      error = "Le titre doit contenir au moins 5 caractères.";
      return;
    }
    if (title.trim().length > 200) {
      error = "Le titre ne doit pas dépasser 200 caractères.";
      return;
    }
    if (description.trim().length < 20) {
      error = "La description doit contenir au moins 20 caractères.";
      return;
    }

    loading = true;
    try {
      await teacher.resubmitProposedSubject(subject.id, {
        title,
        description,
        group_type: groupType,
        domain_ids:
          selectedDomainIds.length > 0 ? selectedDomainIds : undefined,
      });
      goto("/teacher/proposed-subjects");
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      loading = false;
    }
  }
</script>

{#if !subject}
  <Page
    title="Sujet introuvable"
    subtitle="Ce sujet n'existe pas ou vous n'y avez pas accès."
  >
    <Button
      variant="ghost"
      Icon={ArrowLeft}
      onclick={() => goto("/teacher/proposed-subjects")}
    >
      Retour
    </Button>
  </Page>
{:else}
  <Page
    title="Modifier et resoumettre"
    subtitle="Corrigez votre sujet et resoumettez-le pour validation."
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

    {#if hasComments}
      <div class="comments-card">
        <div class="comments-header">
          <AlertTriangle size={16} />
          <span>Remarques des validateurs</span>
        </div>
        <div class="comments-body">
          {#if validator1Comment}
            <div class="comment-item">
              <span class="comment-label">Validateur 1</span>
              <p class="comment-text">{validator1Comment}</p>
            </div>
          {/if}
          {#if validator2Comment}
            <div class="comment-item">
              <span class="comment-label">Validateur 2</span>
              <p class="comment-text">{validator2Comment}</p>
            </div>
          {/if}
        </div>
      </div>
    {/if}

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
            placeholder="Ex. Conception d'un système de recommandation"
          />
        </div>

        <div class="field">
          <label for="group_type">Type de groupe</label>
          <select id="group_type" bind:value={groupType} required>
            <option value="monome">Monôme (1 étudiant)</option>
            <option value="binome">Binôme (2 étudiants)</option>
            <option value="trinome">Trinôme (3 étudiants)</option>
          </select>
        </div>

        <div class="field">
          <label for="niveau">Niveau</label>
          <select
            id="niveau"
            bind:value={selectedNiveau}
            onchange={() => (selectedSpecialtyId = "")}
          >
            <option value="">Sélectionner un niveau</option>
            {#each niveaux as n}
              <option value={n.value}>{n.label}</option>
            {/each}
          </select>
        </div>

        <div class="field">
          <label for="specialty">Spécialité</label>
          <select
            id="specialty"
            bind:value={selectedSpecialtyId}
            disabled={!selectedNiveau}
          >
            <option value="">
              {selectedNiveau
                ? "Sélectionner une spécialité"
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
              <span class="domain-count">
                ({selectedDomainIds.length} sélectionné{selectedDomainIds.length >
                1
                  ? "s"
                  : ""})
              </span>
            {/if}
          </span>
          <p class="domain-hint">
            Sélectionnez les domaines liés à votre sujet.
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
          placeholder="Décrivez le sujet en détail : objectifs, technologies, livrables attendus..."
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
          {loading ? "Resoumission..." : "Resoumettre pour validation"}
        </Button>
      </div>
    </form>
  </Page>
{/if}

<style>
  .comments-card {
    background: color-mix(
      in srgb,
      var(--color-warning) 8%,
      var(--color-surface)
    );
    border: 1px solid color-mix(in srgb, var(--color-warning) 30%, transparent);
    border-radius: 12px;
    padding: var(--spacing-md);
    margin-bottom: var(--spacing-lg);
  }

  .comments-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: var(--text-sm);
    font-weight: 600;
    font-family: var(--font-sans);
    color: var(--color-warning);
    margin-bottom: var(--spacing-sm);
  }

  .comments-body {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .comment-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .comment-label {
    font-size: var(--text-xs);
    font-weight: 600;
    font-family: var(--font-sans);
    color: var(--color-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .comment-text {
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
    margin: 0;
    line-height: 1.5;
  }

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

  .field-heading {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
    display: block;
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
    margin-top: var(--spacing-lg);
  }
</style>
