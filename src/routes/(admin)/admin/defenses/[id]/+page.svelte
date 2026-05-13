<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { goto } from "$app/navigation";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const {
    defense,
    grades,
    supervisorEval,
    jury,
    presidentName,
    memberName,
    currentUserId,
  } = $derived(data);

  const isPresident = $derived(jury?.president_id === currentUserId);
  const isJuryMember = $derived(
    jury?.president_id === currentUserId || jury?.member_id === currentUserId,
  );

  // Step management
  // step: 'grade_form' | 'waiting_for_member' | 'resolve'
  let step = $state<"grade_form" | "waiting_for_member" | "resolve">(
    "grade_form",
  );
  let error = $state("");
  let loading = $state(false);

  // Grade form state
  let c1 = $state(0);
  let c2 = $state(0);
  let c3 = $state(0);
  let c4 = $state(0);
  const total = $derived(c1 + c2 + c3 + c4);

  // Resolve state
  let resolveChoice = $state<"own" | "member" | "new">("own");
  let nc1 = $state(0);
  let nc2 = $state(0);
  let nc3 = $state(0);
  let nc4 = $state(0);
  const newTotal = $derived(nc1 + nc2 + nc3 + nc4);

  // Determine initial step based on existing grades
  $effect(() => {
    if (grades.length === 0) {
      step = "grade_form";
    } else if (grades.length === 1) {
      // One grade submitted
      const myGrade = grades.find((g) => g.jury_member_id === currentUserId);
      if (myGrade) {
        // I already submitted, waiting for the other member
        step = "waiting_for_member";
      } else {
        // I haven't submitted yet
        step = "grade_form";
      }
    } else {
      // Both grades submitted
      if (isPresident) {
        step = "resolve";
      } else {
        step = "waiting_for_member";
      }
    }
  });

  const STATUS_LABELS: Record<string, string> = {
    scheduled: "Planifiee",
    done: "Passee",
    postponed: "Reportee",
  };

  function statusVariant(
    status: string,
  ): "success" | "warning" | "danger" | "info" | "neutral" {
    switch (status) {
      case "scheduled":
        return "warning";
      case "done":
        return "success";
      case "postponed":
        return "danger";
      default:
        return "neutral";
    }
  }

  function formatDate(iso: string): string {
    return new Date(iso).toLocaleDateString("fr-FR", {
      day: "numeric",
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  async function handleSubmitGrade(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    try {
      const res = await fetch(`/api/defenses/${defense.id}/submit-grade`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          criterion1: c1,
          criterion2: c2,
          criterion3: c3,
          criterion4: c4,
        }),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.message || "Erreur lors de la soumission");
      }

      await invalidateAll();
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      loading = false;
    }
  }

  async function handleContinue() {
    error = "";
    if (grades.length < 2) {
      error =
        "Impossible de continuer tant que le membre du jury n'a pas soumis son evaluation";
      return;
    }
    step = "resolve";
  }

  async function handleResolve(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    try {
      const payload: { choice: string; newGrade?: Record<string, number> } = {
        choice: resolveChoice,
      };
      if (resolveChoice === "new") {
        payload.newGrade = {
          criterion1: nc1,
          criterion2: nc2,
          criterion3: nc3,
          criterion4: nc4,
        };
      }

      const res = await fetch(`/api/defenses/${defense.id}/resolve-grade`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.message || "Erreur lors de la resolution");
      }

      await invalidateAll();
    } catch (err: unknown) {
      error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      loading = false;
    }
  }

  const myGrade = $derived(
    grades.find((g) => g.jury_member_id === currentUserId),
  );
  const otherGrade = $derived(
    grades.find((g) => g.jury_member_id !== currentUserId),
  );
</script>

<Page
  title="Detail de la soutenance"
  subtitle={formatDate(defense.scheduled_at)}
>
  {#if error}
    <div class="error">{error}</div>
  {/if}

  <div class="grid">
    <div class="card">
      <h3>Informations</h3>
      <dl>
        <dt>Date</dt>
        <dd>{formatDate(defense.scheduled_at)}</dd>

        <dt>Salle</dt>
        <dd>{defense.room}</dd>

        <dt>Statut</dt>
        <dd>
          <Badge
            variant={statusVariant(defense.status)}
            label={STATUS_LABELS[defense.status] ?? defense.status}
          />
        </dd>

        <dt>President</dt>
        <dd>{presidentName}</dd>

        <dt>Membre</dt>
        <dd>{memberName}</dd>

        {#if defense.final_grade !== null}
          <dt>Note finale</dt>
          <dd><strong>{defense.final_grade}/20</strong></dd>

          <dt>Decision</dt>
          <dd>
            {#if defense.tiebreak_choice === "president"}
              Note du president retenue
            {:else if defense.tiebreak_choice === "member"}
              Note du membre retenue
            {:else if defense.tiebreak_choice === "average"}
              Nouvelle evaluation du president
            {:else}
              —
            {/if}
          </dd>
          {#if supervisorEval}
            <dt>Evaluation encadrant</dt>
            <dd><strong>{supervisorEval.technical_quality}/4</strong></dd>
          {/if}
        {/if}
      </dl>
    </div>

    {#if isJuryMember && defense.status !== "done"}
      <div class="card">
        <h3>Evaluation</h3>

        {#if step === "grade_form"}
          <form onsubmit={handleSubmitGrade}>
            <p class="hint">
              Saisissez votre evaluation (4 criteres sur 5 chacun)
            </p>

            <div class="criteria-grid">
              <label>
                Critere 1
                <input
                  type="number"
                  min="0"
                  max="5"
                  step="0.5"
                  bind:value={c1}
                  required
                />
              </label>
              <label>
                Critere 2
                <input
                  type="number"
                  min="0"
                  max="5"
                  step="0.5"
                  bind:value={c2}
                  required
                />
              </label>
              <label>
                Critere 3
                <input
                  type="number"
                  min="0"
                  max="5"
                  step="0.5"
                  bind:value={c3}
                  required
                />
              </label>
              <label>
                Critere 4
                <input
                  type="number"
                  min="0"
                  max="5"
                  step="0.5"
                  bind:value={c4}
                  required
                />
              </label>
            </div>

            <p class="total">Total : <strong>{total}/20</strong></p>

            <Button variant="primary" type="submit" disabled={loading}>
              {loading ? "Soumission..." : "Soumettre mon evaluation"}
            </Button>
          </form>
        {:else if step === "waiting_for_member"}
          <p class="hint">
            {#if myGrade}
              Votre evaluation a ete soumise ({myGrade.total}/20). En attente de
              l'evaluation du {isPresident ? "membre" : "president"} du jury.
            {:else}
              En attente de votre evaluation.
            {/if}
          </p>

          {#if isPresident && myGrade}
            <div class="actions">
              <Button variant="primary" onclick={handleContinue}>
                Continuer
              </Button>
            </div>
          {/if}
        {:else if step === "resolve"}
          <form onsubmit={handleResolve}>
            <p class="hint">
              Les deux evaluations ont ete soumises. En tant que president,
              choisissez la note finale.
            </p>

            <div class="grades-comparison">
              <div class="grade-card">
                <h4>Votre evaluation</h4>
                {#if myGrade}
                  <p>
                    C1: {myGrade.criterion1} | C2: {myGrade.criterion2} | C3: {myGrade.criterion3}
                    | C4: {myGrade.criterion4}
                  </p>
                  <p class="total"><strong>{myGrade.total}/20</strong></p>
                {/if}
              </div>
              <div class="grade-card">
                <h4>Evaluation du membre</h4>
                {#if otherGrade}
                  <p>
                    C1: {otherGrade.criterion1} | C2: {otherGrade.criterion2} | C3:
                    {otherGrade.criterion3} | C4: {otherGrade.criterion4}
                  </p>
                  <p class="total"><strong>{otherGrade.total}/20</strong></p>
                {/if}
              </div>
            </div>

            <div class="resolve-options">
              <label>
                <input
                  type="radio"
                  name="resolve"
                  value="own"
                  bind:group={resolveChoice}
                />
                Ma note jury ({myGrade?.total}/20) + encadrant ({supervisorEval?.technical_quality ??
                  "?"}/4)
              </label>
              <label>
                <input
                  type="radio"
                  name="resolve"
                  value="member"
                  bind:group={resolveChoice}
                />
                Note du membre ({otherGrade?.total}/20) + encadrant ({supervisorEval?.technical_quality ??
                  "?"}/4)
              </label>
              <label>
                <input
                  type="radio"
                  name="resolve"
                  value="new"
                  bind:group={resolveChoice}
                />
                Nouvelle evaluation
              </label>
            </div>

            {#if resolveChoice === "new"}
              <div class="criteria-grid">
                <label>
                  Critere 1
                  <input
                    type="number"
                    min="0"
                    max="5"
                    step="0.5"
                    bind:value={nc1}
                    required
                  />
                </label>
                <label>
                  Critere 2
                  <input
                    type="number"
                    min="0"
                    max="5"
                    step="0.5"
                    bind:value={nc2}
                    required
                  />
                </label>
                <label>
                  Critere 3
                  <input
                    type="number"
                    min="0"
                    max="5"
                    step="0.5"
                    bind:value={nc3}
                    required
                  />
                </label>
                <label>
                  Critere 4
                  <input
                    type="number"
                    min="0"
                    max="5"
                    step="0.5"
                    bind:value={nc4}
                    required
                  />
                </label>
              </div>
              <p class="total">
                Sous-total jury : <strong>{newTotal}/20</strong>
              </p>
              {#if supervisorEval}
                <p class="total">
                  + Critere 5 encadrant : <strong
                    >{supervisorEval.technical_quality}/4</strong
                  >
                  = Note finale :
                  <strong
                    >{newTotal +
                      (supervisorEval.technical_quality ?? 0)}/20</strong
                  >
                </p>
              {/if}
            {:else if supervisorEval}
              <p class="supervisor-note">
                Critere 5 (encadrant) : <strong
                  >{supervisorEval.technical_quality}/4</strong
                >
                — Note finale estimee :
                <strong
                  >{(resolveChoice === "own"
                    ? (myGrade?.total ?? 0)
                    : (otherGrade?.total ?? 0)) +
                    (supervisorEval.technical_quality ?? 0)}/20</strong
                >
              </p>
            {/if}

            <div class="actions">
              <Button variant="primary" type="submit" disabled={loading}>
                {loading ? "Validation..." : "Valider la note finale"}
              </Button>
            </div>
          </form>
        {/if}
      </div>
    {/if}

    {#if grades.length > 0}
      <div class="card">
        <h3>Notes du jury</h3>
        <table>
          <thead>
            <tr>
              <th>Membre</th>
              <th>C1 (5)</th>
              <th>C2 (5)</th>
              <th>C3 (5)</th>
              <th>C4 (5)</th>
              <th>Total (20)</th>
            </tr>
          </thead>
          <tbody>
            {#each grades as grade}
              <tr>
                <td>
                  {grade.jury_member_id === jury?.president_id
                    ? presidentName
                    : grade.jury_member_id === jury?.member_id
                      ? memberName
                      : grade.jury_member_id}
                </td>
                <td>{grade.criterion1}</td>
                <td>{grade.criterion2}</td>
                <td>{grade.criterion3}</td>
                <td>{grade.criterion4}</td>
                <td><strong>{grade.total}</strong></td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>

  <div class="actions">
    <Button variant="ghost" onclick={() => goto("/admin/defenses")}
      >Retour a la liste</Button
    >
  </div>
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

  .grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-lg);

    @media screen and (max-width: 768px) {
      & {
        grid-template-columns: 1fr;
      }
    }
  }

  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-lg);
  }

  h3 {
    margin: 0 0 var(--spacing-md);
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
  }

  h4 {
    margin: 0 0 0.5rem;
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--color-text);
  }

  dl {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 0.5rem 1rem;
    margin: 0;
  }

  dt {
    font-weight: 600;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  dd {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--color-text);
  }

  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin-bottom: var(--spacing-md);
  }

  .criteria-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-sm);
    margin-bottom: var(--spacing-md);
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
  }

  input[type="number"] {
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-surface);
    color: var(--color-text);
    width: 100%;
    box-sizing: border-box;
  }

  .total {
    font-size: 1rem;
    color: var(--color-text);
    margin-bottom: var(--spacing-md);
  }

  .grades-comparison {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-md);
  }

  .grade-card {
    background: var(--color-background);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: var(--spacing-md);
  }

  .grade-card p {
    font-size: var(--text-sm);
    color: var(--color-text);
    margin: 0.25rem 0;
  }

  .resolve-options {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: var(--spacing-md);
  }

  .resolve-options label {
    flex-direction: row;
    align-items: center;
    gap: 0.5rem;
    font-weight: 400;
    cursor: pointer;
  }

  .resolve-options input[type="radio"] {
    margin: 0;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-sans);
    font-size: 0.85rem;
  }

  th {
    text-align: left;
    padding: 0.5rem 0.75rem;
    font-weight: 600;
    color: var(--color-text);
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    border-bottom: 1px solid var(--color-border);
  }

  td {
    padding: 0.5rem 0.75rem;
    color: var(--color-text);
    border-bottom: 1px solid var(--color-border);
  }

  tr:last-child td {
    border-bottom: none;
  }

  .actions {
    margin-top: var(--spacing-lg);
    display: flex;
    gap: var(--spacing-sm);
    align-items: center;
  }
</style>
