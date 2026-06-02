<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { ArrowLeft } from "lucide-svelte";
  import { teacher } from "$lib/api";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { subject } = $derived(data);

  let decision = $state<"valide" | "accepte_sous_reserve" | "refuse" | "">("");
  let comment = $state("");
  let error = $state("");
  let loading = $state(false);
  let alreadyDecided = $state(false);

  const GROUP_LABELS: Record<string, string> = {
    monome: "Monome",
    binome: "Binome",
    trinome: "Trinome",
  };

  const PROPOSER_LABELS: Record<string, string> = {
    teacher: "Enseignant",
    company: "Entreprise",
    student: "Etudiant",
  };

  async function handleSubmit() {
    if (!decision) {
      error = "Veuillez selectionner une decision.";
      return;
    }

    if (
      (decision === "refuse" || decision === "accepte_sous_reserve") &&
      !comment.trim()
    ) {
      error = "Un commentaire est obligatoire pour cette décision.";
      return;
    }

    if (
      (decision === "refuse" || decision === "accepte_sous_reserve") &&
      comment.trim().length < 10
    ) {
      error = "Le commentaire doit contenir au moins 10 caractères.";
      return;
    }

    error = "";
    loading = true;

    try {
      await teacher.validateSubject(subject!.id, {
        decision: decision as import("$lib/types").ReviewDecision,
        comment,
      });

      await invalidateAll();
      goto("/teacher/subjects-to-validate");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Erreur inconnue";
      if (msg.includes("déjà soumis")) {
        alreadyDecided = true;
      } else {
        error = msg;
      }
    } finally {
      loading = false;
    }
  }
</script>

<Page
  title="Examiner le sujet"
  subtitle="Detail du sujet et formulaire de decision."
>
  {#snippet actions()}
    <Button
      variant="ghost"
      Icon={ArrowLeft}
      onclick={() => goto("/teacher/subjects-to-validate")}
    >
      Retour
    </Button>
  {/snippet}

  {#if error}
    <div class="error">{error}</div>
  {/if}

  <section class="subject-detail">
    <h2>{subject!.title}</h2>

    <div class="meta">
      <div class="meta-item">
        <span class="label">Auteur</span>
        <span
          >{PROPOSER_LABELS[subject!.proposer_role] ??
            subject!.proposer_role}</span
        >
      </div>
      <div class="meta-item">
        <span class="label">Domaines</span>
        <span>{subject!.domains?.map((d) => d.name).join(", ") || "-"}</span>
      </div>
      <div class="meta-item">
        <span class="label">Type de groupe</span>
        <Badge
          variant="info"
          label={GROUP_LABELS[subject!.group_type] ?? subject!.group_type}
        />
      </div>
      <div class="meta-item">
        <span class="label">Statut</span>
        <Badge
          variant={subject!.status === "en_attente" ? "warning" : "info"}
          label={subject!.status === "en_attente"
            ? "En attente"
            : subject!.status}
        />
      </div>
    </div>

    <div class="description">
      <h3>Description</h3>
      <p>{subject!.description}</p>
    </div>
  </section>

  <section class="decision-form">
    <h3>Votre decision</h3>

    {#if alreadyDecided}
      <div class="already-decided">
        <p>
          Vous avez déjà soumis votre décision pour ce sujet. Elle est en cours
          de traitement.
        </p>
        <Button
          variant="ghost"
          onclick={() => goto("/teacher/subjects-to-validate")}
        >
          Retour à la liste
        </Button>
      </div>
    {:else}
      <div class="radio-group">
        <label class="radio">
          <input
            type="radio"
            name="decision"
            value="valide"
            bind:group={decision}
          />
          <span>Accepter</span>
        </label>
        <label class="radio">
          <input
            type="radio"
            name="decision"
            value="accepte_sous_reserve"
            bind:group={decision}
          />
          <span>Accepter sous reserve</span>
        </label>
        <label class="radio">
          <input
            type="radio"
            name="decision"
            value="refuse"
            bind:group={decision}
          />
          <span>Refuser</span>
        </label>
      </div>

      {#if decision === "refuse" || decision === "accepte_sous_reserve"}
        <div class="comment-field">
          <label for="comment">Commentaire (obligatoire)</label>
          <textarea
            id="comment"
            bind:value={comment}
            placeholder="Expliquez les raisons de votre decision..."
            rows={4}
          ></textarea>
        </div>
      {/if}

      <div class="actions">
        <Button
          variant="ghost"
          onclick={() => goto("/teacher/subjects-to-validate")}
        >
          Annuler
        </Button>
        <Button
          variant="primary"
          onclick={handleSubmit}
          disabled={!decision || loading}
        >
          {loading ? "Envoi..." : "Soumettre la decision"}
        </Button>
      </div>
    {/if}
  </section>
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

  .subject-detail {
    margin-bottom: var(--spacing-lg);
  }

  .subject-detail h2 {
    font-size: var(--text-lg);
    font-weight: 700;
    margin: 0 0 var(--spacing-md);
    color: var(--color-text);
  }

  .meta {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-md);
  }

  .meta-item {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .label {
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
  }

  .description {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: var(--spacing-md);
  }

  .description h3 {
    font-size: var(--text-sm);
    font-weight: 600;
    margin: 0 0 var(--spacing-sm);
    color: var(--color-text);
  }

  .description p {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    line-height: 1.6;
    white-space: pre-wrap;
  }

  .decision-form {
    border-top: 1px solid var(--color-border);
    padding-top: var(--spacing-lg);
  }

  .decision-form h3 {
    font-size: var(--text-md);
    font-weight: 600;
    margin: 0 0 var(--spacing-md);
    color: var(--color-text);
  }

  .radio-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: var(--spacing-md);
  }

  .radio {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    font-size: var(--text-sm);
    color: var(--color-text);
  }

  .radio input[type="radio"] {
    accent-color: var(--color-accent);
  }

  .comment-field {
    margin-bottom: var(--spacing-md);
  }

  .comment-field label {
    display: block;
    font-size: var(--text-sm);
    font-weight: 600;
    margin-bottom: 0.35rem;
    color: var(--color-text);
  }

  .comment-field textarea {
    width: 100%;
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    color: var(--color-text);
    background: var(--color-surface);
    resize: vertical;
    transition: border-color var(--transition-fast);

    &:focus {
      outline: none;
      border-color: var(--color-accent);
    }
  }

  .actions {
    display: flex;
    gap: var(--spacing-sm);
    justify-content: flex-end;
  }

  .already-decided {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
    padding: var(--spacing-md);
    background: color-mix(in srgb, var(--color-accent) 8%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-accent) 25%, transparent);
    border-radius: 8px;

    & p {
      margin: 0;
      font-size: var(--text-sm);
      color: var(--color-text);
    }
  }
</style>
