<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { untrack } from "svelte";
  import { teacher as teacherApi } from "$lib/api";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import DateInput from "$lib/components/ui/DateInput.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();
  const { teacher } = $derived(data);


  let status = $state(
    untrack(() => teacher?.availability_status ?? "disponible"),
  );
  let unavailableUntil = $state(
    untrack(() => teacher?.unavailable_until ?? null),
  );

  const STATUS_LABELS: Record<string, string> = {
    disponible: "Disponible",
    indisponible: "Indisponible",
    indisponible_jusqu_au: "Indisponible jusqu'a une date",
  };

  const STATUS_VARIANTS: Record<
    string,
    "info" | "warning" | "success" | "danger" | "neutral"
  > = {
    disponible: "success",
    indisponible: "danger",
    indisponible_jusqu_au: "warning",
  };

  let saving = $state(false);
  let error = $state("");

  async function handleSave() {
    saving = true;
    error = "";

    try {
      await teacherApi.updateAvailability({
        availability_status: status,
        unavailable_until:
          status === "indisponible_jusqu_au" ? unavailableUntil : undefined,
      });

      await invalidateAll();
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      saving = false;
    }
  }
</script>

<Page title="Disponibilité" subtitle="Gérer votre statut de disponibilité.">
  <div class="current">
    <span class="label">Statut actuel :</span>
    <Badge
      variant={STATUS_VARIANTS[teacher.availability_status] || "info"}
      label={STATUS_LABELS[teacher.availability_status] ||
        teacher.availability_status}
    />
  </div>

  <div class="form">
    {#if error}
      <div class="error">{error}</div>
    {/if}

    <div class="field">
      <label>
        <input
          type="radio"
          name="status"
          value="disponible"
          bind:group={status}
        />
        Disponible
      </label>
    </div>

    <div class="field">
      <label>
        <input
          type="radio"
          name="status"
          value="indisponible"
          bind:group={status}
        />
        Indisponible
      </label>
    </div>

    <div class="field">
      <label>
        <input
          type="radio"
          name="status"
          value="indisponible_jusqu_au"
          bind:group={status}
        />
        Indisponible jusqu'au
      </label>

      {#if status === "indisponible_jusqu_au"}
        <br />
        <DateInput
          bind:value={unavailableUntil}
          min={new Date().toISOString()}
        />
      {/if}
    </div>

    <div class="actions">
      <Button variant="primary" onclick={handleSave} disabled={saving}>
        {saving ? "Enregistrement..." : "Enregistrer"}
      </Button>
    </div>
  </div>
</Page>

<style>
  .current {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    margin-bottom: var(--spacing-lg);
    padding: var(--spacing-md);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
  }

  .label {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
  }

  .form {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: var(--spacing-lg);
  }

  .error {
    padding: 0.5rem 0.75rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    color: var(--color-danger);
    border-radius: 6px;
    font-size: var(--text-sm);
    margin-bottom: var(--spacing-md);
  }

  .field {
    margin-bottom: var(--spacing-md);

    & label {
      display: flex;
      align-items: center;
      gap: var(--spacing-sm);
      font-size: var(--text-sm);
      color: var(--color-text);
      cursor: pointer;
    }

    & input[type="radio"] {
      accent-color: var(--color-accent);
    }
  }

  .actions {
    margin-top: var(--spacing-lg);
    display: flex;
    justify-content: flex-end;
  }
</style>
