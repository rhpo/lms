<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { ArrowLeft, Check } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { subject, students } = $derived(data);

  const GROUP_LIMITS: Record<string, number> = {
    monome: 1,
    binome: 2,
    trinome: 3,
  };

  const maxStudents = $derived(GROUP_LIMITS[subject.group_type] ?? 1);

  let selectedIds = $state<string[]>([]);
  let error = $state("");
  let loading = $state(false);

  function toggleStudent(id: string) {
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

    error = "";
    loading = true;

    try {
      const res = await fetch(`/api/subjects/${subject.id}/assign`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ student_ids: selectedIds }),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.error?.message || "Erreur lors de l'affectation");
      }

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
  title="Candidats - {subject.title}"
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

  {#if error}
    <div class="error">{error}</div>
  {/if}

  <div class="info">
    <span class="info-label">Type de groupe :</span>
    <Badge
      variant="info"
      label={subject.group_type === "monome"
        ? "Monome"
        : subject.group_type === "binome"
          ? "Binome"
          : "Trinome"}
    />
    <span class="info-label">Places :</span>
    <span class="info-value"
      >{selectedIds.length} / {maxStudents} selectionne(s)</span
    >
  </div>

  {#if students.length === 0}
    <div class="empty">
      <p>Aucun etudiant n'a mis ce sujet dans ses voeux pour le moment.</p>
    </div>
  {:else}
    <div class="student-list">
      {#each students as student (student.id)}
        <button
          class="student-card"
          class:selected={selectedIds.includes(student.id)}
          onclick={() => toggleStudent(student.id)}
          disabled={!selectedIds.includes(student.id) &&
            selectedIds.length >= maxStudents}
        >
          <div class="checkbox">
            {#if selectedIds.includes(student.id)}
              <Check size={16} />
            {/if}
          </div>
          <div class="student-info">
            <span class="student-name"
              >{student.first_name} {student.last_name}</span
            >
            <span class="student-detail"
              >{student.specialty || "Specialite non renseignee"}</span
            >
          </div>
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
      <Button
        variant="primary"
        onclick={handleAssign}
        disabled={selectedIds.length === 0 || loading}
      >
        {loading ? "Affectation..." : "Confirmer l'affectation"}
      </Button>
    </div>
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
  }

  .selected .checkbox {
    border-color: var(--color-accent);
    background: var(--color-accent);
    color: white;
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
