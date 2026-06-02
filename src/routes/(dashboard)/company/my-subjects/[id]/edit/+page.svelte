<script lang="ts">
  import { goto } from "$app/navigation";
  import { ArrowLeft, AlertTriangle } from "lucide-svelte";
  import { company } from "$lib/api";
  import { showToast } from "$lib/utils/toast";

  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();
  const subject = $derived(data.subject);

  let title = $state(subject?.title ?? "");
  let description = $state(subject?.description ?? "");
  let groupType = $state<"monome" | "binome" | "trinome">(
    (subject?.group_type as "monome" | "binome" | "trinome") ?? "monome",
  );

  let loading = $state(false);
  let error = $state("");


  const validator1Comment = $derived(subject?.validator1_comment ?? null);
  const validator2Comment = $derived(subject?.validator2_comment ?? null);
  const hasComments = $derived(Boolean(validator1Comment) || Boolean(validator2Comment));

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
      await company.updateSubject(subject.id, {
        title,
        description,
        group_type: groupType,
      });
      showToast.success("Sujet mis à jour avec succès");
      goto("/company/my-subjects");
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      loading = false;
    }
  }
</script>

{#if !subject}
  <Page title="Sujet introuvable" subtitle="">
    <Button variant="ghost" Icon={ArrowLeft} onclick={() => goto("/company/my-subjects")}>
      Retour
    </Button>
  </Page>
{:else}
  <Page
    title="Modifier le sujet"
    subtitle="Corrigez votre sujet suite aux remarques des validateurs."
  >
    {#snippet actions()}
      <Button variant="ghost" Icon={ArrowLeft} onclick={() => goto("/company/my-subjects")}>
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
      <div class="error-banner">{error}</div>
    {/if}

    <form onsubmit={handleSubmit}>
      <div class="form-grid">
        <div class="field">
          <label for="title">Titre du sujet</label>
          <input id="title" type="text" bind:value={title} required minlength={5} maxlength={200} />
        </div>

        <div class="field">
          <label for="group_type">Type de groupe</label>
          <select id="group_type" bind:value={groupType} required>
            <option value="monome">Monôme (1 étudiant)</option>
            <option value="binome">Binôme (2 étudiants)</option>
            <option value="trinome">Trinôme (3 étudiants)</option>
          </select>
        </div>
      </div>

      <div class="field">
        <label for="description">Description</label>
        <textarea id="description" bind:value={description} required rows={6} minlength={20}></textarea>
      </div>

      <div class="actions">
        <Button variant="ghost" onclick={() => goto("/company/my-subjects")} disabled={loading}>
          Annuler
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? "Enregistrement..." : "Enregistrer les modifications"}
        </Button>
      </div>
    </form>
  </Page>
{/if}

<style>
  .comments-card {
    background: color-mix(in srgb, var(--color-warning) 8%, var(--color-surface));
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

  .error-banner {
    padding: 0.75rem 1rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    color: var(--color-danger);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    margin-bottom: var(--spacing-md);
  }

  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-md);

    @media screen and (max-width: 600px) {
      & { grid-template-columns: 1fr; }
    }
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    margin-bottom: var(--spacing-md);
  }

  label {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
  }

  input, select, textarea {
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-surface);
    color: var(--color-text);
    width: 100%;
    box-sizing: border-box;
    transition: border-color 0.15s;

    &:focus {
      outline: none;
      border-color: var(--color-accent);
      box-shadow: 0 0 0 2px color-mix(in srgb, var(--color-accent) 20%, transparent);
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
