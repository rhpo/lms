<script lang="ts">
  import { goto } from "$app/navigation";
  import { ArrowLeft } from "lucide-svelte";
  import { admin } from "$lib/api";
  import { showToast } from "$lib/utils/toast";

  import Button from "$lib/components/ui/Button.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const assignment = $derived(data.assignment);
  const teachers = $derived(data.teachers || []);

  let presidentId = $state<number | "">("");
  let memberId = $state<number | "">("");
  let scheduledAt = $state("");
  let room = $state("");

  const now = new Date();
  now.setMinutes(now.getMinutes() + 1, 0, 0);
  const minDatetime = now.toISOString().slice(0, 16);

  let disabled = $derived(
    !presidentId || !memberId || !scheduledAt || !room || !assignment,
  );
  let submitting = $state(false);
  let errorMsg = $state("");

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (disabled || !assignment) return;
    submitting = true;
    errorMsg = "";

    if (new Date(scheduledAt) <= new Date()) {
      errorMsg = "La date de soutenance doit être dans le futur";
      submitting = false;
      return;
    }

    if (room.trim().length < 1) {
      errorMsg = "La salle est obligatoire";
      submitting = false;
      return;
    }

    if (presidentId === memberId) {
      errorMsg =
        "Le président et l'examinateur doivent être des personnes différentes";
      submitting = false;
      return;
    }

    try {
      const scheduledAtRFC3339 =
        scheduledAt.length === 16 ? `${scheduledAt}:00Z` : scheduledAt;
      await admin.createDefense({
        assignment_id: Number(assignment.id),
        president_id: Number(presidentId),
        member_id: Number(memberId),
        scheduled_at: scheduledAtRFC3339,
        room,
      });
      showToast.success("Soutenance programmée avec succès");
      goto("/admin/pfe");
    } catch (err) {
      errorMsg = err instanceof Error ? err.message : "Erreur réseau";
      showToast.error(errorMsg);
    } finally {
      submitting = false;
    }
  }
</script>

{#if assignment}
  <Page
    title="Planifier une soutenance"
    subtitle="Programmer une nouvelle soutenance pour ce PFE"
  >
    {#snippet actions()}
      <a href="/admin/pfe">
        <Button variant="ghost" Icon={ArrowLeft}>Retour aux PFE</Button>
      </a>
    {/snippet}

    <div class="form-container">
      <div class="card">
        {#if errorMsg}
          <div class="error-banner">{errorMsg}</div>
        {/if}

        <div class="assignment-info">
          <div class="info-group">
            <span class="label">PFE:</span>
            <span class="value"
              >[{assignment.pfe_code || "PFE"}] {assignment.subject?.title ??
                "Sujet inconnu"}</span
            >
          </div>
          <div class="info-group">
            <span class="label">Étudiant:</span>
            <span class="value"
              >{assignment.student?.profile?.full_name ?? "Inconnu"}</span
            >
          </div>
        </div>

        <form onsubmit={handleSubmit} class="form-grid">
          <div class="separator"></div>

          <div class="grid-2">
            <FormField label="Date et heure" required>
              <input
                type="datetime-local"
                bind:value={scheduledAt}
                min={minDatetime}
                class="input"
                required
              />
            </FormField>

            <FormField label="Salle" required>
              <input
                type="text"
                bind:value={room}
                class="input"
                placeholder="Ex: Salle 12, Amphi A"
                required
              />
            </FormField>
          </div>

          <div class="separator"></div>

          <div class="grid-2">
            <FormField label="Président du Jury" required>
              <select bind:value={presidentId} class="input" required>
                <option value="">Sélectionner l'enseignant</option>
                {#each teachers as t}
                  {@const teacherEntityId = t.teacher?.id ?? null}
                  {@const isSupervisor =
                    t.id === assignment.supervisor?.profile_id}
                  {@const isAlreadyMember =
                    teacherEntityId !== null && teacherEntityId === memberId}
                  {#if teacherEntityId}
                    <option
                      value={teacherEntityId}
                      disabled={isSupervisor || isAlreadyMember}
                    >
                      {t.full_name}{isSupervisor
                        ? " - (Encadrant de ce PFE)"
                        : isAlreadyMember
                          ? " - (Déjà Membre)"
                          : ""}
                    </option>
                  {/if}
                {/each}
              </select>
            </FormField>

            <FormField label="Membre du Jury (Examinateur)" required>
              <select bind:value={memberId} class="input" required>
                <option value="">Sélectionner l'enseignant</option>
                {#each teachers as t}
                  {@const teacherEntityId = t.teacher?.id ?? null}
                  {@const isSupervisor =
                    t.id === assignment.supervisor?.profile_id}
                  {@const isAlreadyPresident =
                    teacherEntityId !== null && teacherEntityId === presidentId}
                  {#if teacherEntityId}
                    <option
                      value={teacherEntityId}
                      disabled={isSupervisor || isAlreadyPresident}
                    >
                      {t.full_name}{isSupervisor
                        ? " - (Encadrant de ce PFE)"
                        : isAlreadyPresident
                          ? " - (Déjà Président)"
                          : ""}
                    </option>
                  {/if}
                {/each}
              </select>
            </FormField>
          </div>

          <div class="form-actions">
            <a href="/admin/pfe">
              <Button variant="ghost" type="button">Annuler</Button>
            </a>
            <Button
              variant="primary"
              type="submit"
              disabled={disabled || submitting}
            >
              {submitting ? "Création en cours..." : "Programmer la soutenance"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  </Page>
{:else}
  <Page title="PFE introuvable" subtitle="">
    <p>Cette affectation PFE n'existe pas ou a été supprimée.</p>
    <a href="/admin/pfe">
      <Button variant="ghost" Icon={ArrowLeft}>Retour</Button>
    </a>
  </Page>
{/if}

<style>
  .form-container {
    max-width: 900px;
    margin: 0 auto;
  }

  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-xl);
    box-shadow:
      0 4px 6px -1px rgb(0 0 0 / 0.1),
      0 2px 4px -2px rgb(0 0 0 / 0.1);
  }

  .assignment-info {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
    background: var(--color-background-100);
    padding: var(--spacing-md);
    border-radius: 8px;
    border: 1px solid var(--color-border);
    margin-bottom: var(--spacing-lg);
  }

  .info-group {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    font-size: var(--text-sm);
  }

  .label {
    font-weight: 600;
    color: var(--color-text-muted);
    min-width: 80px;
  }

  .value {
    color: var(--color-text);
    font-weight: 500;
  }

  .form-grid {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .grid-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-xl);
  }

  .separator {
    height: 1px;
    background: var(--color-border);
    margin: var(--spacing-md) 0;
    opacity: 0.5;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-md);
    margin-top: var(--spacing-xl);
    padding-top: var(--spacing-lg);
    border-top: 1px solid var(--color-border);
  }

  option:disabled {
    color: var(--color-text-muted);
    font-style: italic;
  }

  .error-banner {
    padding: var(--spacing-md);
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-danger) 20%, transparent);
    border-radius: 8px;
    color: var(--color-danger);
    font-size: var(--text-sm);
    margin-bottom: var(--spacing-lg);
    display: flex;
    align-items: center;
  }

  @media screen and (max-width: 768px) {
    .grid-2 {
      grid-template-columns: 1fr;
      gap: var(--spacing-lg);
    }

    .card {
      padding: var(--spacing-lg);
    }
  }
</style>
