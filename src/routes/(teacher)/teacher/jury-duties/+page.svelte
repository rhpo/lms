<script lang="ts">
  import { invalidateAll } from "$app/navigation";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { duties } = $derived(data);

  const STATUS_LABELS: Record<string, string> = {
    scheduled: "Planifiee",
    done: "Passee",
    postponed: "Reportee",
  };

  const STATUS_VARIANTS: Record<
    string,
    "info" | "warning" | "success" | "danger" | "neutral"
  > = {
    scheduled: "info",
    done: "success",
    postponed: "warning",
  };

  // Grade form state per defense
  let gradeForms: Record<
    string,
    {
      c1: number;
      c2: number;
      c3: number;
      c4: number;
      loading: boolean;
      error: string;
    }
  > = $state({});

  function getForm(defenseId: string) {
    if (!gradeForms[defenseId]) {
      gradeForms[defenseId] = {
        c1: 0,
        c2: 0,
        c3: 0,
        c4: 0,
        loading: false,
        error: "",
      };
    }
    return gradeForms[defenseId];
  }

  function getTotal(defenseId: string) {
    const f = getForm(defenseId);
    return f.c1 + f.c2 + f.c3 + f.c4;
  }

  function updateCriterion(
    defenseId: string,
    key: "c1" | "c2" | "c3" | "c4",
    value: number,
  ) {
    const f = getForm(defenseId);
    f[key] = value;
  }

  async function handleSubmitGrade(defenseId: string) {
    const f = getForm(defenseId);
    f.error = "";
    f.loading = true;

    try {
      const res = await fetch(`/api/defenses/${defenseId}/submit-grade`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          criterion1: f.c1,
          criterion2: f.c2,
          criterion3: f.c3,
          criterion4: f.c4,
        }),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.error?.message || "Erreur lors de l'envoi");
      }

      await invalidateAll();
    } catch (err: unknown) {
      f.error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      f.loading = false;
    }
  }
</script>

<Page title="Mes jurys" subtitle="Soutenances ou vous etes membre du jury.">
  {#if duties.length === 0}
    <div class="empty">
      <p>Aucune soutenance de jury pour le moment.</p>
    </div>
  {:else}
    <div class="list">
      {#each duties as { defense, grades }}
        <div class="card">
          <div class="card-header">
            <div class="card-info">
              <span class="card-title"
                >Soutenance #{defense.id.slice(0, 8)}</span
              >
              <Badge
                variant={STATUS_VARIANTS[defense.status] || "info"}
                label={STATUS_LABELS[defense.status] || defense.status}
              />
            </div>
            <span class="card-date">
              {new Date(defense.scheduled_at).toLocaleDateString("fr-FR", {
                day: "numeric",
                month: "long",
                year: "numeric",
                hour: "2-digit",
                minute: "2-digit",
              })}
            </span>
          </div>

          <div class="card-body">
            <div class="detail">
              <span class="label">Salle</span>
              <span>{defense.room}</span>
            </div>
            <div class="detail">
              <span class="label">Role</span>
              <span>Membre du jury</span>
            </div>
          </div>

          {#if defense.status === "done" && grades.length === 0}
            <div class="grade-section">
              <h4>Saisie de la note</h4>
              {#if getForm(defense.id).error}
                <div class="error">{getForm(defense.id).error}</div>
              {/if}
              <div class="criteria">
                <label>
                  Criter 1 (sur 5)
                  <input
                    type="number"
                    min="0"
                    max="5"
                    value={getForm(defense.id).c1}
                    oninput={(e) =>
                      updateCriterion(
                        defense.id,
                        "c1",
                        Number(e.currentTarget.value),
                      )}
                  />
                </label>
                <label>
                  Criter 2 (sur 5)
                  <input
                    type="number"
                    min="0"
                    max="5"
                    value={getForm(defense.id).c2}
                    oninput={(e) =>
                      updateCriterion(
                        defense.id,
                        "c2",
                        Number(e.currentTarget.value),
                      )}
                  />
                </label>
                <label>
                  Criter 3 (sur 5)
                  <input
                    type="number"
                    min="0"
                    max="5"
                    value={getForm(defense.id).c3}
                    oninput={(e) =>
                      updateCriterion(
                        defense.id,
                        "c3",
                        Number(e.currentTarget.value),
                      )}
                  />
                </label>
                <label>
                  Criter 4 (sur 5)
                  <input
                    type="number"
                    min="0"
                    max="5"
                    value={getForm(defense.id).c4}
                    oninput={(e) =>
                      updateCriterion(
                        defense.id,
                        "c4",
                        Number(e.currentTarget.value),
                      )}
                  />
                </label>
              </div>
              <div class="total">
                Total : <strong>{getTotal(defense.id)} / 20</strong>
              </div>
              <Button
                variant="primary"
                onclick={() => handleSubmitGrade(defense.id)}
                disabled={getForm(defense.id).loading}
              >
                {getForm(defense.id).loading ? "Envoi..." : "Soumettre la note"}
              </Button>
            </div>
          {:else if grades.length > 0}
            <div class="grade-info">
              <Badge variant="success" label="Note saisie" />
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</Page>

<style>
  .empty {
    text-align: center;
    padding: 3rem 1rem;
    color: var(--color-text-muted);

    & p {
      font-size: var(--text-sm);
      margin: 0;
    }
  }

  .list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .card {
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background: var(--color-surface);
    overflow: hidden;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--spacing-md);
    border-bottom: 1px solid var(--color-border);
  }

  .card-info {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
  }

  .card-title {
    font-weight: 600;
    font-size: var(--text-sm);
    color: var(--color-text);
  }

  .card-date {
    font-size: 0.8rem;
    color: var(--color-text-muted);
  }

  .card-body {
    padding: var(--spacing-md);
    display: flex;
    gap: var(--spacing-lg);
  }

  .detail {
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

  .grade-section {
    padding: var(--spacing-md);
    border-top: 1px solid var(--color-border);
  }

  .grade-section h4 {
    font-size: var(--text-sm);
    font-weight: 600;
    margin: 0 0 var(--spacing-sm);
    color: var(--color-text);
  }

  .error {
    padding: 0.5rem 0.75rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    color: var(--color-danger);
    border-radius: 6px;
    font-size: var(--text-sm);
    margin-bottom: var(--spacing-sm);
  }

  .criteria {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: var(--spacing-sm);
    margin-bottom: var(--spacing-sm);
  }

  .criteria label {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-size: var(--text-sm);
    color: var(--color-text);
  }

  .criteria input {
    padding: 0.4rem 0.6rem;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    color: var(--color-text);
    background: var(--color-background);
    width: 100%;
    max-width: 80px;

    &:focus {
      outline: none;
      border-color: var(--color-accent);
    }
  }

  .total {
    font-size: var(--text-sm);
    color: var(--color-text);
    margin-bottom: var(--spacing-sm);
  }

  .grade-info {
    padding: var(--spacing-md);
    border-top: 1px solid var(--color-border);
  }
</style>
