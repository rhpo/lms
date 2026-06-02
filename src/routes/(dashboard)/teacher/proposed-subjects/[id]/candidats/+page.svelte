<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { ArrowLeft, Check, Lock } from "lucide-svelte";
  import { teacher } from "$lib/api";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import { GROUP_TYPE_LABELS } from "$lib/constants/labels";

  let { data } = $props();

  const { subject, wishes } = $derived(data);

  const GROUP_LIMITS: Record<string, number> = {
    monome: 1,
    binome: 2,
    trinome: 3,
  };

  const maxStudents = $derived(subject ? GROUP_LIMITS[subject.group_type] ?? 1 : 1);


  const alreadyAssigned = $derived(wishes.some((w) => w.status === "accepte"));


  let selectedIds = $state<number[]>([]);
  $effect(() => {
    selectedIds = wishes.filter((w) => w.status === "accepte").map((w) => w.student_id);
  });

  let error = $state("");
  let loading = $state(false);

  function toggleStudent(id: number) {

    if (wishes.find((w) => w.student_id === id)?.status === "accepte") return;
    if (selectedIds.includes(id)) {
      selectedIds = selectedIds.filter((s) => s !== id);
    } else {
      if (selectedIds.length >= maxStudents) return;
      selectedIds = [...selectedIds, id];
    }
  }

  async function handleAssign() {
    if (selectedIds.length === 0) {
      error = "Selectionnez au moins un etudiant.";
      return;
    }
    if (!subject) return;

    error = "";
    loading = true;

    try {
      await teacher.acceptCandidat(subject.id, { student_ids: selectedIds });

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
  title="Candidats - {subject?.title ?? 'Sujet'}"
  subtitle="Selectionnez les etudiants a affecter a ce sujet."
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

  {#if !subject}
    <div class="empty"><p>Sujet introuvable.</p></div>
  {:else}
    {#if error}
      <div class="error">{error}</div>
    {/if}

    <div class="info">
      <span class="info-label">Type de groupe :</span>
      <Badge
        variant="info"
        label={GROUP_TYPE_LABELS[subject.group_type] ?? subject.group_type}
      />
      <span class="info-label">Places :</span>
      <span class="info-value"
        >{selectedIds.length} / {maxStudents} selectionne(s)</span
      >
    </div>

    {#if wishes.length === 0}
      <div class="empty">
        <p>Aucun etudiant n'a mis ce sujet dans ses voeux pour le moment.</p>
      </div>
    {:else}
      <div class="student-list">
        {#each wishes as wish (wish.id)}
          {@const studentName = wish.student?.profile?.full_name ?? "Étudiant inconnu"}
          {@const specialityName = wish.student?.speciality?.name ?? null}
          {@const level = wish.student?.level ?? null}
          {@const detail = [specialityName, level].filter(Boolean).join(" · ") || "Spécialité non renseignée"}
          {@const isAccepted = wish.status === "accepte"}
          {@const isRefused = wish.status === "refuse"}
          {@const isChecked = selectedIds.includes(wish.student_id)}
          {@const isLocked = isAccepted || isRefused}
          <button
            class="student-card"
            class:selected={isChecked && !isRefused}
            class:locked={isLocked}
            class:refused={isRefused}
            onclick={() => toggleStudent(wish.student_id)}
            disabled={isLocked || (!isChecked && selectedIds.length >= maxStudents)}
          >
            <div class="checkbox" class:accepted={isAccepted}>
              {#if isAccepted}
                <Check size={16} />
              {:else if isChecked}
                <Check size={16} />
              {/if}
            </div>
            <div class="student-info">
              <span class="student-name">{studentName}</span>
              <span class="student-detail">{detail}</span>
            </div>
            {#if isAccepted}
              <span class="status-tag tag-accepted">
                <Lock size={12} />
                Affecté
              </span>
            {:else if isRefused}
              <span class="status-tag tag-refused">Non retenu</span>
            {/if}
          </button>
        {/each}
      </div>

      <div class="actions">
        <Button
          variant="ghost"
          onclick={() => goto("/teacher/proposed-subjects")}
        >
          Annuler
        </Button>
        <span title={alreadyAssigned ? "Soumission déjà faite" : undefined}>
          <Button
            variant="primary"
            onclick={handleAssign}
            disabled={alreadyAssigned || selectedIds.length === 0 || loading}
          >
            {loading ? "Affectation..." : "Confirmer l'affectation"}
          </Button>
        </span>
      </div>
    {/if}
  {/if}
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

  .info {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: var(--spacing-md);
    font-size: var(--text-sm);
  }

  .info-label {
    font-weight: 600;
    color: var(--color-text-muted);
  }

  .info-value {
    color: var(--color-text);
  }

  .empty {
    text-align: center;
    padding: 3rem 1rem;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .student-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: var(--spacing-lg);
  }

  .student-card {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background: var(--color-surface);
    cursor: pointer;
    text-align: left;
    width: 100%;
    font-family: var(--font-sans);
    transition:
      border-color var(--transition-fast),
      background var(--transition-fast);

    &:hover:not(:disabled) {
      border-color: var(--color-accent);
    }

    &.selected {
      border-color: var(--color-accent);
      background: color-mix(in srgb, var(--color-accent) 8%, transparent);
    }

    &.locked {
      cursor: default;
    }

    &.refused {
      opacity: 0.55;
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }

  .checkbox {
    width: 22px;
    height: 22px;
    border: 2px solid var(--color-border);
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: var(--color-accent);
    transition:
      border-color var(--transition-fast),
      background var(--transition-fast);

    &.accepted {
      border-color: var(--color-success, #22c55e);
      background: var(--color-success, #22c55e);
      color: white;
    }
  }

  .selected .checkbox {
    border-color: var(--color-accent);
    background: var(--color-accent);
    color: white;
  }

  .status-tag {
    margin-left: auto;
    display: flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.75rem;
    font-weight: 600;
    font-family: var(--font-sans);
    padding: 0.2rem 0.55rem;
    border-radius: 999px;
    white-space: nowrap;
  }

  .tag-accepted {
    background: color-mix(in srgb, var(--color-success, #22c55e) 12%, transparent);
    color: var(--color-success, #22c55e);
  }

  .tag-refused {
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
  }

  .student-info {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .student-name {
    font-weight: 600;
    font-size: var(--text-sm);
    color: var(--color-text);
  }

  .student-detail {
    font-size: 0.8rem;
    color: var(--color-text-muted);
  }

  .actions {
    display: flex;
    gap: var(--spacing-sm);
    justify-content: flex-end;
  }
</style>
