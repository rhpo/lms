<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { teacher } from "$lib/api";
  import {
    Calendar,
    MapPin,
    Users,
    BookOpen,
    User,
    Shield,
    CheckCircle,
    Clock,
    AlertCircle,
  } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import type { GradeContext } from "./+page.ts";

  let { data } = $props();
  const { duties } = $derived(data);

  const STATUS_LABELS: Record<string, string> = {
    scheduled: "Planifiée",
    done: "Passée",
    postponed: "Reportée",
  };
  const STATUS_VARIANTS: Record<
    string,
    "info" | "warning" | "success" | "danger"
  > = {
    scheduled: "info",
    done: "success",
    postponed: "warning",
  };

  function getStudentNames(defense: any): string {
    const a = defense.assignment;
    if (!a) return "-";
    const names: string[] = [];
    if (a.student?.profile?.full_name) names.push(a.student.profile.full_name);
    if (a.student2?.profile?.full_name)
      names.push(a.student2.profile.full_name);
    if (a.student3?.profile?.full_name)
      names.push(a.student3.profile.full_name);
    return names.length > 0 ? names.join(", ") : "-";
  }

  function getRoleVariant(role: string): "info" | "warning" {
    return role === "Président" ? "warning" : "info";
  }

  type ArchiveDecision =
    | "archivable"
    | "minor_corrections"
    | "major_corrections"
    | "";

  type MemberGradeForm = {
    c1: number;
    c2: number;
    c3: number;
    c4: number;
    archiveDecision: ArchiveDecision;
    loading: boolean;
    error: string;
    success: boolean;

    presidentChoice: "member" | "new" | "";
  };

  let gradeForms = $state<Record<number, MemberGradeForm>>({});

  const defaultForm = (): MemberGradeForm => ({
    c1: 0,
    c2: 0,
    c3: 0,
    c4: 0,
    archiveDecision: "",
    loading: false,
    error: "",
    success: false,
    presidentChoice: "",
  });

  const ARCHIVE_OPTIONS: { value: ArchiveDecision; label: string }[] = [
    { value: "archivable", label: "Le mémoire peut être archivé" },
    {
      value: "minor_corrections",
      label: "Peut être archivé après des corrections mineures",
    },
    {
      value: "major_corrections",
      label: "Ne peut être archivé, nécessite des corrections majeures",
    },
  ];

  const CRITERIA_LABELS = [
    "Rédaction du mémoire",
    "Présentation",
    "Réponses aux questions",
    "Réalisation et qualité des résultats obtenus",
  ];

  $effect(() => {
    for (const { defense, gradeCtx } of duties) {
      if (!gradeForms[defense.id]) {
        const f = defaultForm();

        if (gradeCtx?.my_grade) {
          const g = gradeCtx.my_grade;
          f.c1 = g.criterion1 ?? 0;
          f.c2 = g.criterion2 ?? 0;
          f.c3 = g.criterion3 ?? 0;
          f.c4 = g.criterion4 ?? 0;
          f.archiveDecision = (g.archive_decision ?? "") as ArchiveDecision;
        }

        gradeForms[defense.id] = f;
      }
    }
  });

  function formFor(defenseId: number): MemberGradeForm {
    return gradeForms[defenseId] ?? defaultForm();
  }

  function totalFor(defenseId: number): number {
    const f = formFor(defenseId);
    return +(f.c1 + f.c2 + f.c3 + f.c4).toFixed(2);
  }

  function updateCriterion(
    defenseId: number,
    key: "c1" | "c2" | "c3" | "c4",
    value: number,
  ) {
    const f = gradeForms[defenseId];
    if (f) f[key] = Math.min(4, Math.max(0, value));
  }

  async function handleSubmitGrade(defenseId: number) {
    const f = gradeForms[defenseId];
    if (!f) return;
    f.error = "";
    f.success = false;

    if ([f.c1, f.c2, f.c3, f.c4].some((v) => v < 0 || v > 4 || isNaN(v))) {
      f.error = "Chaque critère doit être entre 0 et 4.";
      return;
    }
    if (f.c1 + f.c2 + f.c3 + f.c4 === 0) {
      f.error = "La note totale ne peut pas être 0.";
      return;
    }

    f.loading = true;
    try {
      await teacher.submitGrade(defenseId, {
        criterion1: f.c1,
        criterion2: f.c2,
        criterion3: f.c3,
        criterion4: f.c4,
        archive_decision: "",
      });
      f.success = true;
      await invalidateAll();
    } catch (err: unknown) {
      f.error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      f.loading = false;
    }
  }

  async function handleSubmitFinalGrade(
    defenseId: number,
    gradeCtx: GradeContext,
  ) {
    const f = gradeForms[defenseId];
    if (!f) return;
    f.error = "";
    f.success = false;

    if (!f.presidentChoice) {
      f.error = "Veuillez choisir quelle évaluation utiliser.";
      return;
    }
    if (!f.archiveDecision) {
      f.error = "Veuillez sélectionner une décision d'archivage.";
      return;
    }

    if (f.presidentChoice === "new") {
      if ([f.c1, f.c2, f.c3, f.c4].some((v) => v < 0 || v > 4 || isNaN(v))) {
        f.error = "Chaque critère doit être entre 0 et 4.";
        return;
      }
      if (f.c1 + f.c2 + f.c3 + f.c4 === 0) {
        f.error = "La note totale ne peut pas être 0.";
        return;
      }
    }

    f.loading = true;
    try {
      const body: any = {
        choice: f.presidentChoice,
        archive_decision: f.archiveDecision,
      };
      if (f.presidentChoice === "new") {
        body.criterion1 = f.c1;
        body.criterion2 = f.c2;
        body.criterion3 = f.c3;
        body.criterion4 = f.c4;
      }
      await teacher.submitFinalGrade(defenseId, body);
      f.success = true;
      await invalidateAll();
    } catch (err: unknown) {
      f.error = err instanceof Error ? err.message : "Erreur inconnue";
    } finally {
      f.loading = false;
    }
  }

  function onPresidentChoiceChange(
    defenseId: number,
    choice: "member" | "new",
    gradeCtx: GradeContext,
  ) {
    const f = gradeForms[defenseId];
    if (!f) return;
    f.presidentChoice = choice;
    if (choice === "new" && gradeCtx.member_grade) {
      f.c1 = gradeCtx.member_grade.criterion1 ?? 0;
      f.c2 = gradeCtx.member_grade.criterion2 ?? 0;
      f.c3 = gradeCtx.member_grade.criterion3 ?? 0;
      f.c4 = gradeCtx.member_grade.criterion4 ?? 0;
    }
  }
</script>

<Page
  title="Mes soutenances"
  subtitle="Soutenances où vous êtes membre du jury."
>
  {#if duties.length === 0}
    <div class="empty">
      <p>Aucune soutenance de jury pour le moment.</p>
    </div>
  {:else}
    <div class="list">
      {#each duties as { defense, gradeCtx }}
        {@const role =
          gradeCtx?.my_role === "president" ? "Président" : "Examinateur"}
        <div class="card">
          <!-- Header -->
          <div class="card-header">
            <div class="card-title-row">
              <h3 class="subject-title">
                {defense.assignment?.subject?.title ??
                  `Soutenance #${defense.id}`}
              </h3>
              <div class="header-badges">
                <Badge variant={getRoleVariant(role)} label={role} />
                <Badge
                  variant={STATUS_VARIANTS[defense.status] ?? "info"}
                  label={STATUS_LABELS[defense.status] ?? defense.status}
                />
              </div>
            </div>
            {#if defense.assignment?.pfe_code}
              <span class="pfe-code">{defense.assignment.pfe_code}</span>
            {/if}
          </div>

          <!-- Body: info grid -->
          <div class="card-body">
            <div class="info-grid">
              <div class="info-item">
                <Calendar size={14} />
                <div>
                  <span class="info-label">Date &amp; heure</span>
                  <span class="info-value">
                    {defense.scheduled_at
                      ? new Date(defense.scheduled_at).toLocaleDateString(
                          "fr-FR",
                          {
                            weekday: "long",
                            day: "numeric",
                            month: "long",
                            year: "numeric",
                            hour: "2-digit",
                            minute: "2-digit",
                          },
                        )
                      : "Date non définie"}
                  </span>
                </div>
              </div>
              <div class="info-item">
                <MapPin size={14} />
                <div>
                  <span class="info-label">Salle</span>
                  <span class="info-value">{defense.room ?? "Non définie"}</span
                  >
                </div>
              </div>
              <div class="info-item">
                <Users size={14} />
                <div>
                  <span class="info-label">Étudiant(s)</span>
                  <span class="info-value">{getStudentNames(defense)}</span>
                </div>
              </div>
              <div class="info-item">
                <User size={14} />
                <div>
                  <span class="info-label">Encadrant</span>
                  <span class="info-value">
                    {defense.assignment?.supervisor?.profile?.full_name ?? "-"}
                  </span>
                </div>
              </div>
            </div>

            {#if defense.jury}
              <div class="jury-section">
                <span class="info-label"
                  ><Shield size={12} /> Composition du jury</span
                >
                <div class="jury-members">
                  <span class="jury-member">
                    <strong>Président :</strong>
                    {defense.jury.president?.profile?.full_name ?? "-"}
                  </span>
                  <span class="jury-member">
                    <strong>Examinateur :</strong>
                    {defense.jury.member?.profile?.full_name ?? "-"}
                  </span>
                </div>
              </div>
            {/if}
          </div>

          <!-- Grade section -->
          {#if defense.status === "scheduled" || defense.status === "done"}
            <div class="grade-section">
              <h4><BookOpen size={14} /> Notation</h4>

              {#if gradeCtx === null}
                <p class="loading-note">Chargement du contexte de notation…</p>
              {:else if gradeCtx.final_grade_set}
                <!-- Final grade already set - show read-only summary -->
                <div class="final-done-banner">
                  <CheckCircle size={16} />
                  Note finale soumise. La note a été enregistrée pour cet étudiant.
                </div>
              {:else}
                <!-- ─── ENCADRANT EVALUATION STATUS ────────────────────────── -->
                <div class="eval-status-row">
                  {#if gradeCtx.supervisor_submitted}
                    <span class="eval-chip ok">
                      <CheckCircle size={12} /> Évaluation encadrant reçue ({gradeCtx
                        .supervisor_eval?.criterion5}/4)
                    </span>
                  {:else}
                    <span class="eval-chip pending">
                      <Clock size={12} /> En attente de l'évaluation de l'encadrant
                    </span>
                  {/if}

                  {#if gradeCtx.my_role === "president"}
                    {#if gradeCtx.member_submitted}
                      <span class="eval-chip ok">
                        <CheckCircle size={12} /> Évaluation examinateur reçue
                      </span>
                    {:else}
                      <span class="eval-chip pending">
                        <Clock size={12} /> En attente de l'évaluation de l'examinateur
                      </span>
                    {/if}
                  {/if}
                </div>

                {#if formFor(defense.id).success}
                  <div class="success-banner">Note soumise avec succès !</div>
                {/if}
                {#if formFor(defense.id).error}
                  <div class="error-banner">{formFor(defense.id).error}</div>
                {/if}

                <!-- ─── MEMBER VIEW ─────────────────────────────────────────── -->
                {#if gradeCtx.my_role === "member"}
                  <div class="criteria-grid">
                    {#each CRITERIA_LABELS as label, i}
                      {@const key = `c${i + 1}` as "c1" | "c2" | "c3" | "c4"}
                      <div class="criterion">
                        <label for="c-{defense.id}-{i}">
                          {label}<span class="criterion-max">/4</span>
                        </label>
                        <input
                          id="c-{defense.id}-{i}"
                          type="number"
                          min="0"
                          max="4"
                          step="0.5"
                          value={formFor(defense.id)[key]}
                          oninput={(e) =>
                            updateCriterion(
                              defense.id,
                              key,
                              Number(e.currentTarget.value),
                            )}
                          class="input criterion-input"
                        />
                      </div>
                    {/each}
                  </div>

                  <div class="grade-footer">
                    <div class="total">
                      Total : <strong>{totalFor(defense.id)} / 16</strong>
                    </div>
                    <Button
                      variant="primary"
                      onclick={() => handleSubmitGrade(defense.id)}
                      disabled={formFor(defense.id).loading}
                    >
                      {formFor(defense.id).loading
                        ? "Envoi…"
                        : gradeCtx.my_grade
                          ? "Mettre à jour mon évaluation"
                          : "Soumettre mon évaluation"}
                    </Button>
                  </div>

                  <!-- ─── PRESIDENT VIEW ───────────────────────────────────────── -->
                {:else if gradeCtx.my_role === "president"}
                  {#if !gradeCtx.member_submitted || !gradeCtx.supervisor_submitted}
                    <div class="blocked-notice">
                      <AlertCircle size={16} />
                      <span>
                        Vous pourrez soumettre la note finale une fois que
                        {#if !gradeCtx.member_submitted && !gradeCtx.supervisor_submitted}
                          l'examinateur et l'encadrant auront soumis leurs
                          évaluations.
                        {:else if !gradeCtx.member_submitted}
                          l'examinateur aura soumis son évaluation.
                        {:else}
                          l'encadrant aura soumis son évaluation.
                        {/if}
                      </span>
                    </div>
                  {:else}
                    <!-- Show member's evaluation for reference -->
                    {#if gradeCtx.member_grade}
                      <div class="ref-block">
                        <p class="ref-label">
                          Évaluation de l'examinateur (référence)
                        </p>
                        <div class="ref-criteria">
                          {#each CRITERIA_LABELS as label, i}
                            {@const key = `criterion${i + 1}` as
                              | "criterion1"
                              | "criterion2"
                              | "criterion3"
                              | "criterion4"}
                            <div class="ref-row">
                              <span class="ref-crit-label">{label}</span>
                              <span class="ref-crit-val"
                                >{gradeCtx.member_grade![key] ?? 0} / 4</span
                              >
                            </div>
                          {/each}
                          <div class="ref-row total-row">
                            <span class="ref-crit-label"
                              >Total jury (examinateur)</span
                            >
                            <span class="ref-crit-val">
                              {(
                                (gradeCtx.member_grade.criterion1 ?? 0) +
                                (gradeCtx.member_grade.criterion2 ?? 0) +
                                (gradeCtx.member_grade.criterion3 ?? 0) +
                                (gradeCtx.member_grade.criterion4 ?? 0)
                              ).toFixed(2)} / 16
                            </span>
                          </div>
                        </div>
                      </div>
                    {/if}

                    <!-- Choice selector -->
                    <div class="choice-section">
                      <p class="archive-label">
                        Quelle évaluation utiliser pour la note finale ?
                      </p>
                      <div class="choice-options">
                        <label
                          class="choice-option"
                          class:selected={formFor(defense.id)
                            .presidentChoice === "member"}
                        >
                          <input
                            type="radio"
                            name="pres-choice-{defense.id}"
                            value="member"
                            checked={formFor(defense.id).presidentChoice ===
                              "member"}
                            oninput={() =>
                              onPresidentChoiceChange(
                                defense.id,
                                "member",
                                gradeCtx,
                              )}
                          />
                          <div class="choice-info">
                            <strong
                              >Utiliser l'évaluation de l'examinateur</strong
                            >
                            <span
                              >Les critères C1-C4 de l'examinateur seront
                              retenus.</span
                            >
                          </div>
                        </label>
                        <label
                          class="choice-option"
                          class:selected={formFor(defense.id)
                            .presidentChoice === "new"}
                        >
                          <input
                            type="radio"
                            name="pres-choice-{defense.id}"
                            value="new"
                            checked={formFor(defense.id).presidentChoice ===
                              "new"}
                            oninput={() =>
                              onPresidentChoiceChange(
                                defense.id,
                                "new",
                                gradeCtx,
                              )}
                          />
                          <div class="choice-info">
                            <strong>Saisir une nouvelle évaluation</strong>
                            <span
                              >Vous définissez les critères C1-C4 vous-même.</span
                            >
                          </div>
                        </label>
                      </div>
                    </div>

                    <!-- If president chose "new", show input fields -->
                    {#if formFor(defense.id).presidentChoice === "new"}
                      <div class="criteria-grid">
                        {#each CRITERIA_LABELS as label, i}
                          {@const key = `c${i + 1}` as
                            | "c1"
                            | "c2"
                            | "c3"
                            | "c4"}
                          <div class="criterion">
                            <label for="pres-c-{defense.id}-{i}">
                              {label}<span class="criterion-max">/4</span>
                            </label>
                            <input
                              id="pres-c-{defense.id}-{i}"
                              type="number"
                              min="0"
                              max="4"
                              step="0.5"
                              value={formFor(defense.id)[key]}
                              oninput={(e) =>
                                updateCriterion(
                                  defense.id,
                                  key,
                                  Number(e.currentTarget.value),
                                )}
                              class="input criterion-input"
                            />
                          </div>
                        {/each}
                      </div>
                      <div
                        class="total"
                        style="margin-bottom: var(--spacing-sm)"
                      >
                        Total jury : <strong>{totalFor(defense.id)} / 16</strong
                        >
                      </div>
                    {/if}

                    <!-- Archivage decision for president -->
                    {#if formFor(defense.id).presidentChoice}
                      <div class="archive-section">
                        <p class="archive-label">Décision d'archivage</p>
                        {#each ARCHIVE_OPTIONS as opt}
                          <label class="archive-option">
                            <input
                              type="radio"
                              name="pres-archive-{defense.id}"
                              value={opt.value}
                              checked={formFor(defense.id).archiveDecision ===
                                opt.value}
                              oninput={() => {
                                if (!gradeForms[defense.id])
                                  gradeForms[defense.id] = defaultForm();
                                gradeForms[defense.id].archiveDecision =
                                  opt.value as any;
                              }}
                            />
                            <span>{opt.label}</span>
                          </label>
                        {/each}
                      </div>

                      <!-- Final grade preview -->
                      {#if formFor(defense.id).presidentChoice && gradeCtx.supervisor_eval?.criterion5 !== undefined}
                        {@const juryTotal =
                          formFor(defense.id).presidentChoice === "member"
                            ? (gradeCtx.member_grade?.criterion1 ?? 0) +
                              (gradeCtx.member_grade?.criterion2 ?? 0) +
                              (gradeCtx.member_grade?.criterion3 ?? 0) +
                              (gradeCtx.member_grade?.criterion4 ?? 0)
                            : totalFor(defense.id)}
                        {@const supNote =
                          gradeCtx.supervisor_eval?.criterion5 ?? 0}
                        <div class="grade-preview">
                          <span class="grade-preview-label"
                            >Aperçu note finale :</span
                          >
                          <span class="grade-preview-val">
                            {juryTotal.toFixed(2)} (jury) + {supNote} (encadrant)
                            =
                            <strong
                              >{(juryTotal + supNote).toFixed(2)} / 20</strong
                            >
                          </span>
                        </div>
                      {/if}

                      <div class="grade-footer">
                        <div></div>
                        <Button
                          variant="primary"
                          onclick={() =>
                            handleSubmitFinalGrade(defense.id, gradeCtx)}
                          disabled={formFor(defense.id).loading ||
                            !formFor(defense.id).presidentChoice}
                        >
                          {formFor(defense.id).loading
                            ? "Envoi…"
                            : "Soumettre la note finale"}
                        </Button>
                      </div>
                    {/if}
                  {/if}
                {/if}
              {/if}
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
    gap: var(--spacing-lg);
  }

  .card {
    border: 1px solid var(--color-border);
    border-radius: 12px;
    background: var(--color-surface);
    overflow: hidden;
  }

  .card-header {
    padding: var(--spacing-lg);
    border-bottom: 1px solid var(--color-border);
  }

  .card-title-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }

  .subject-title {
    font-family: var(--font-sans);
    font-size: var(--text-lg);
    font-weight: 700;
    color: var(--color-text);
    margin: 0;
    flex: 1;
    min-width: 0;
  }

  .header-badges {
    display: flex;
    gap: 0.35rem;
    flex-shrink: 0;
  }

  .pfe-code {
    display: inline-block;
    margin-top: 0.35rem;
    font-family: var(--font-mono, monospace);
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .card-body {
    padding: var(--spacing-lg);
  }

  .info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-md);
  }

  .info-item {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
    color: var(--color-text-muted);
    & > div {
      display: flex;
      flex-direction: column;
      gap: 0.1rem;
    }
  }

  .info-label {
    font-size: var(--text-xs);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
    font-family: var(--font-sans);
    display: flex;
    align-items: center;
    gap: 0.3rem;
  }

  .info-value {
    font-size: var(--text-sm);
    color: var(--color-text);
    font-family: var(--font-sans);
  }

  .jury-section {
    margin-top: var(--spacing-md);
    padding-top: var(--spacing-md);
    border-top: 1px solid var(--color-border);
  }

  .jury-members {
    display: flex;
    gap: var(--spacing-lg);
    margin-top: 0.35rem;
  }

  .jury-member {
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
    & strong {
      font-weight: 600;
      color: var(--color-text-muted);
    }
  }

  /* ── Grade section ─────────────────────────── */
  .grade-section {
    padding: var(--spacing-lg);
    border-top: 1px solid var(--color-border);
    background: var(--color-background);

    & h4 {
      font-size: var(--text-sm);
      font-weight: 700;
      font-family: var(--font-sans);
      color: var(--color-text);
      margin: 0 0 var(--spacing-md);
      display: flex;
      align-items: center;
      gap: 0.35rem;
    }
  }

  .loading-note {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    font-style: italic;
    margin: 0;
  }

  .final-done-banner {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.6rem 0.9rem;
    background: color-mix(
      in srgb,
      var(--color-success) 10%,
      var(--color-surface)
    );
    border: 1px solid color-mix(in srgb, var(--color-success) 25%, transparent);
    border-radius: 8px;
    color: var(--color-success);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    font-weight: 500;
  }

  /* ── Status chips ───────────────────────────── */
  .eval-status-row {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: var(--spacing-md);
  }

  .eval-chip {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    font-size: var(--text-xs);
    font-weight: 600;
    padding: 0.25rem 0.6rem;
    border-radius: 99px;
    font-family: var(--font-sans);

    &.ok {
      background: color-mix(in srgb, var(--color-success) 12%, transparent);
      color: var(--color-success);
      border: 1px solid
        color-mix(in srgb, var(--color-success) 25%, transparent);
    }

    &.pending {
      background: color-mix(in srgb, var(--color-warning) 12%, transparent);
      color: var(--color-warning);
      border: 1px solid
        color-mix(in srgb, var(--color-warning) 25%, transparent);
    }
  }

  /* ── Banners ─────────────────────────────────── */
  .success-banner {
    padding: 0.5rem 0.75rem;
    background: color-mix(
      in srgb,
      var(--color-success) 10%,
      var(--color-surface)
    );
    border: 1px solid color-mix(in srgb, var(--color-success) 20%, transparent);
    border-radius: 6px;
    font-size: var(--text-sm);
    color: var(--color-success);
    margin-bottom: var(--spacing-md);
  }

  .error-banner {
    padding: 0.5rem 0.75rem;
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-danger) 20%, transparent);
    border-radius: 6px;
    color: var(--color-danger);
    font-size: var(--text-sm);
    margin-bottom: var(--spacing-md);
  }

  .blocked-notice {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    background: color-mix(
      in srgb,
      var(--color-warning) 8%,
      var(--color-surface)
    );
    border: 1px solid color-mix(in srgb, var(--color-warning) 20%, transparent);
    border-radius: 8px;
    color: var(--color-warning);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
  }

  /* ── Reference block (president view) ──────── */
  .ref-block {
    margin-bottom: var(--spacing-md);
    background: var(--color-background-100);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    overflow: hidden;
  }

  .ref-label {
    font-size: var(--text-xs);
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
    padding: 0.55rem 0.9rem;
    border-bottom: 1px solid var(--color-border);
    margin: 0;
    background: var(--color-surface);
  }

  .ref-criteria {
    display: flex;
    flex-direction: column;
  }

  .ref-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0.9rem;
    border-bottom: 1px solid var(--color-border);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    &:last-child {
      border-bottom: none;
    }
    &.total-row {
      background: color-mix(
        in srgb,
        var(--color-accent) 5%,
        var(--color-surface)
      );
    }
  }

  .ref-crit-label {
    color: var(--color-text);
  }
  .ref-crit-val {
    font-weight: 600;
    color: var(--color-text);
  }

  /* ── Choice section (president picks which eval) ─ */
  .choice-section {
    margin-bottom: var(--spacing-md);
  }

  .choice-options {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }

  .choice-option {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 0.75rem 0.9rem;
    border: 1.5px solid var(--color-border);
    border-radius: 8px;
    cursor: pointer;
    transition:
      border-color 0.15s,
      background 0.15s;
    background: var(--color-surface);

    &:hover {
      border-color: var(--color-accent);
    }

    &.selected {
      border-color: var(--color-accent);
      background: color-mix(
        in srgb,
        var(--color-accent) 6%,
        var(--color-surface)
      );
    }

    input[type="radio"] {
      margin-top: 2px;
      accent-color: var(--color-accent);
      flex-shrink: 0;
    }
  }

  .choice-info {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;

    & strong {
      font-size: var(--text-sm);
      color: var(--color-text);
    }
    & span {
      font-size: var(--text-xs);
      color: var(--color-text-muted);
    }
  }

  /* ── Grade preview ─────────────────────────── */
  .grade-preview {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.6rem 0.9rem;
    background: color-mix(
      in srgb,
      var(--color-accent) 8%,
      var(--color-surface)
    );
    border: 1px solid color-mix(in srgb, var(--color-accent) 20%, transparent);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    margin-bottom: var(--spacing-md);
    flex-wrap: wrap;
    gap: 0.35rem;
  }

  .grade-preview-label {
    color: var(--color-text-muted);
  }
  .grade-preview-val {
    color: var(--color-text);
    & strong {
      color: var(--color-accent);
    }
  }

  /* ── Criteria grid ─────────────────────────── */
  .criteria-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-md);
  }

  .criterion {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
    & label {
      font-size: var(--text-sm);
      font-weight: 600;
      color: var(--color-text);
    }
  }

  .criterion-max {
    font-weight: 400;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
  }
  .criterion-input {
    max-width: 80px;
  }

  /* ── Archive section ───────────────────────── */
  .archive-section {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    background: var(--color-background-100);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    margin-bottom: var(--spacing-md);
  }

  .archive-label {
    font-size: var(--text-xs);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
    margin: 0 0 0.25rem;
  }

  .archive-option {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: var(--text-sm);
    color: var(--color-text);
    cursor: pointer;
    input[type="radio"] {
      accent-color: var(--color-accent);
    }
  }

  /* ── Footer ────────────────────────────────── */
  .grade-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: var(--spacing-sm);
  }

  .total {
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
  }

  /* ── Responsive ────────────────────────────── */
  @media screen and (max-width: 640px) {
    .info-grid {
      grid-template-columns: 1fr;
    }
    .criteria-grid {
      grid-template-columns: 1fr;
    }
    .jury-members {
      flex-direction: column;
      gap: 0.35rem;
    }
    .card-title-row {
      flex-direction: column;
    }
  }
</style>
