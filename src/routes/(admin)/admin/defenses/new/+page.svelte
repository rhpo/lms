<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { goto } from "$app/navigation";

  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { pfeOptions, preSelectedId } = $derived(data);

  // Capture initial value for form field (intentionally not reactive)
  let assignmentId = $state("");

  $effect(() => {
    assignmentId = preSelectedId ?? "";
  });
  let scheduledAt = $state("");
  let room = $state("");
  let presidentId = $state("");
  let memberId = $state("");
  let error = $state("");
  let loading = $state(false);
  let recommendedTeachers = $state<{ id: string; full_name: string }[]>([]);

  async function loadRecommendedTeachers(id: string) {
    if (!id) {
      recommendedTeachers = [];
      presidentId = "";
      memberId = "";
      return;
    }

    const res = await fetch(
      `/api/defenses/recommended-teachers?assignment_id=${id}`,
    );
    if (res.ok) {
      const body = await res.json();
      recommendedTeachers = body.teachers;
      presidentId = "";
      memberId = "";
    }
  }

  // Load teachers if pre-selected
  $effect(() => {
    if (preSelectedId) {
      loadRecommendedTeachers(preSelectedId);
    }
  });

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    try {
      const res = await fetch("/api/defenses/schedule", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          assignment_id: assignmentId,
          scheduled_at: scheduledAt,
          room,
          president_id: presidentId,
          member_id: memberId,
        }),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.message || "Erreur lors de la planification");
      }

      await invalidateAll();
      goto("/admin/defenses");
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      loading = false;
    }
  }

  const availableMembers = $derived(
    recommendedTeachers.filter((t) => t.id !== presidentId),
  );
</script>

<Page
  title="Planifier une soutenance"
  subtitle="Selectionnez le PFE, la date, la salle et le jury"
>
  <form onsubmit={handleSubmit}>
    {#if error}
      <div class="error">{error}</div>
    {/if}

    <div class="field">
      <label for="assignment">PFE concerne</label>
      <select
        id="assignment"
        bind:value={assignmentId}
        onchange={() => loadRecommendedTeachers(assignmentId)}
        required
      >
        <option value="">Selectionnez un PFE</option>
        {#each pfeOptions as opt}
          <option value={opt.id}>{opt.label}</option>
        {/each}
      </select>
    </div>

    <div class="field">
      <label for="scheduled_at">Date et heure</label>
      <input
        id="scheduled_at"
        type="datetime-local"
        bind:value={scheduledAt}
        required
      />
    </div>

    <div class="field">
      <label for="room">Salle</label>
      <input
        id="room"
        type="text"
        bind:value={room}
        required
        placeholder="Ex. Amphi A"
      />
    </div>

    <fieldset class="jury">
      <legend>Jury</legend>

      <div class="field">
        <label for="president">President du jury</label>
        <select id="president" bind:value={presidentId} required>
          <option value="">Selectionnez le president</option>
          {#each recommendedTeachers as teacher}
            <option value={teacher.id}>{teacher.full_name}</option>
          {/each}
        </select>
      </div>

      <div class="field">
        <label for="member">Membre du jury</label>
        <select id="member" bind:value={memberId} required>
          <option value="">Selectionnez le membre</option>
          {#each availableMembers as teacher}
            <option value={teacher.id}>{teacher.full_name}</option>
          {/each}
        </select>
      </div>
    </fieldset>

    <div class="actions">
      <Button variant="primary" type="submit" disabled={loading}>
        {loading ? "Planification..." : "Confirmer la planification"}
      </Button>
      <Button variant="ghost" onclick={() => goto("/admin/defenses")}
        >Annuler</Button
      >
    </div>
  </form>
</Page>

<style>
  form {
    max-width: 600px;
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .error {
    padding: 0.75rem 1rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    color: var(--color-danger);
    border-radius: 8px;
    font-size: var(--text-sm);
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  label {
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
  }

  select,
  input {
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-surface);
    color: var(--color-text);
  }

  .jury {
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-md);
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
  }

  legend {
    font-weight: 600;
    font-size: var(--text-sm);
    color: var(--color-text);
    padding: 0 0.5rem;
  }

  .actions {
    display: flex;
    gap: var(--spacing-sm);
    align-items: center;
  }
</style>
