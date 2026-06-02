<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { User, Building2, ArrowLeft } from "lucide-svelte";
    import { company } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { pfe, progressReports, supervisorEval } = $derived(data);

    const PFE_STATUS_LABELS: Record<string, string> = {
        en_cours: "En cours",
        soutenance_planifiee: "Soutenance planifiée",
        valide: "Validé",
        refuse: "Refusé",
    };

    const PFE_STATUS_VARIANTS: Record<
        string,
        "info" | "warning" | "success" | "danger"
    > = {
        en_cours: "info",
        soutenance_planifiee: "warning",
        valide: "success",
        refuse: "danger",
    };

    const DURATION_OPTIONS = [
        { value: 30, label: "30 min" },
        { value: 45, label: "45 min" },
        { value: 60, label: "1h" },
        { value: 90, label: "1h30" },
        { value: 120, label: "2h" },
    ];

    const MEETING_TYPE_LABELS: Record<string, string> = {
        presentiel: "Présentiel",
        visio: "Visioconférence",
    };

    const MEETING_STATUS_LABELS: Record<string, string> = {
        a_faire: "À faire",
        en_cours: "En cours",
        termine: "Terminé",
    };

    function formatDuration(mins: number): string {
        if (mins < 60) return `${mins} min`;
        const h = Math.floor(mins / 60);
        const m = mins % 60;
        return m ? `${h}h${String(m).padStart(2, "0")}` : `${h}h`;
    }

    let newMeeting = $state({
        date: new Date().toISOString().split("T")[0],
        duration: 60,
        type: "presentiel" as "presentiel" | "visio",
        status: "a_faire" as "a_faire" | "en_cours" | "termine",
        topics: "",
        observation: "",
    });
    let meetingError = $state("");
    let meetingLoading = $state(false);

    function resetMeetingForm() {
        newMeeting = {
            date: new Date().toISOString().split("T")[0],
            duration: 60,
            type: "presentiel",
            status: "a_faire",
            topics: "",
            observation: "",
        };
    }

    async function addMeetingAction() {
        if (!newMeeting.topics.trim()) {
            meetingError = "Les points discutés sont requis.";
            return;
        }
        meetingError = "";
        meetingLoading = true;
        try {
            await company.addMeeting(pfe!.id, {
                meeting_date: newMeeting.date,
                duration: newMeeting.duration,
                meeting_type: newMeeting.type,
                topics: newMeeting.topics.trim(),
                status: newMeeting.status,
                observation: newMeeting.observation.trim() || undefined,
            });
            resetMeetingForm();
            await invalidateAll();
        } catch (err: unknown) {
            meetingError = err instanceof Error ? err.message : "Erreur réseau";
        } finally {
            meetingLoading = false;
        }
    }

    let criterion5 = $state(0);
    let evalComment = $state("");
    let evalError = $state("");
    let evalSuccess = $state("");
    let evalLoading = $state(false);

    async function submitEvaluation() {
        evalError = "";
        evalSuccess = "";
        if (criterion5 < 0 || criterion5 > 4) {
            evalError = "La note doit être comprise entre 0 et 4.";
            return;
        }
        evalLoading = true;
        try {
            await company.submitEvaluation(pfe!.id, { criterion5 });
            evalSuccess = "Évaluation soumise avec succès.";
            await invalidateAll();
        } catch (err: unknown) {
            evalError = err instanceof Error ? err.message : "Erreur réseau.";
        } finally {
            evalLoading = false;
        }
    }
</script>

<Page title="Detail de l'encadrement" subtitle="Suivez l'avancement du PFE.">
    <div class="back-link">
        <a href="/company/supervised-pfes">
            <ArrowLeft size={14} />
            Retour a mes encadrements
        </a>
    </div>

    <div class="pfe-header">
        <div class="pfe-title-section">
            <h2>{pfe!.subject?.title ?? "Sujet"}</h2>
            <Badge
                variant={PFE_STATUS_VARIANTS[pfe!.status] ?? "info"}
                label={PFE_STATUS_LABELS[pfe!.status] ?? pfe!.status}
            />
        </div>

        <div class="pfe-meta">
            <div class="meta-item">
                <User size={14} />
                <span>
                    Etudiant(s) : {pfe!.student?.profile?.full_name ??
                        "Non assigne"}
                    {#if pfe!.student2_id}
                        , {pfe!.student2?.profile?.full_name ?? "Etudiant 2"}
                    {/if}
                    {#if pfe!.student3_id}
                        , {pfe!.student3?.profile?.full_name ?? "Etudiant 3"}
                    {/if}
                </span>
            </div>
            {#if pfe!.co_supervisor_id}
                <div class="meta-item">
                    <Building2 size={14} />
                    <span>
                        Co-promoteur : {pfe!.co_supervisor?.profile
                            ?.full_name ?? "Enseignant"}
                    </span>
                </div>
            {/if}
        </div>
    </div>

    <section>
        <div class="section-header">
            <h3>Journal de suivi</h3>
        </div>

        {#if meetingError}
            <div class="error-banner">{meetingError}</div>
        {/if}

        <div class="add-meeting">
            <h4>Ajouter une entrée</h4>
            <div class="form-grid">
                <div class="field">
                    <label for="m-date">Date</label>
                    <input
                        id="m-date"
                        type="date"
                        class="input"
                        bind:value={newMeeting.date}
                    />
                </div>

                <div class="field">
                    <label for="m-duration">Durée</label>
                    <select
                        id="m-duration"
                        class="input"
                        bind:value={newMeeting.duration}
                    >
                        {#each DURATION_OPTIONS as opt}
                            <option value={opt.value}>{opt.label}</option>
                        {/each}
                    </select>
                </div>

                <div class="field">
                    <label for="m-type">Type de réunion</label>
                    <select
                        id="m-type"
                        class="input"
                        bind:value={newMeeting.type}
                    >
                        <option value="presentiel">Présentiel</option>
                        <option value="visio">Visioconférence</option>
                    </select>
                </div>

                <div class="field">
                    <label for="m-status">État</label>
                    <select
                        id="m-status"
                        class="input"
                        bind:value={newMeeting.status}
                    >
                        <option value="a_faire">À faire</option>
                        <option value="en_cours">En cours</option>
                        <option value="termine">Terminé</option>
                    </select>
                </div>

                <div class="field full-width">
                    <label for="m-topics">Points discutés</label>
                    <textarea
                        id="m-topics"
                        class="input"
                        bind:value={newMeeting.topics}
                        rows={3}
                        placeholder="Sujets abordés lors de cette réunion..."
                    ></textarea>
                </div>

                <div class="field full-width">
                    <label for="m-obs">Observation (optionnel)</label>
                    <textarea
                        id="m-obs"
                        class="input"
                        bind:value={newMeeting.observation}
                        rows={2}
                        placeholder="Remarques, décisions, prochaines étapes..."
                    ></textarea>
                </div>
            </div>

            <Button
                variant="primary"
                onclick={addMeetingAction}
                disabled={meetingLoading}
            >
                {meetingLoading ? "Ajout..." : "Ajouter l'entrée"}
            </Button>
        </div>

        {#if progressReports.length === 0}
            <p class="empty-section">
                Aucune entrée de journal pour le moment.
            </p>
        {:else}
            <div class="table-wrapper">
                <table>
                    <thead>
                        <tr>
                            <th>Date</th>
                            <th>Durée</th>
                            <th>Type</th>
                            <th>État</th>
                            <th>Points discutés</th>
                            <th>Observation</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each progressReports as report}
                            <tr>
                                <td class="cell-date">
                                    {new Date(
                                        report.meeting_date,
                                    ).toLocaleDateString("fr-FR")}
                                </td>
                                <td class="cell-meta"
                                    >{formatDuration(report.duration)}</td
                                >
                                <td class="cell-meta"
                                    >{MEETING_TYPE_LABELS[
                                        report.meeting_type
                                    ] ?? report.meeting_type}</td
                                >
                                <td>
                                    <span
                                        class="status-badge status-{report.status}"
                                    >
                                        {MEETING_STATUS_LABELS[report.status] ??
                                            report.status}
                                    </span>
                                </td>
                                <td class="cell-topics">{report.topics}</td>
                                <td class="cell-obs"
                                    >{report.observation ?? "-"}</td
                                >
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </section>

    <section class="evaluation-section">
        <div class="section-header">
            <h3>Evaluation de l'encadrant (critere 5)</h3>
        </div>

        {#if evalError}
            <div class="error-banner">{evalError}</div>
        {/if}
        {#if evalSuccess}
            <div class="success-banner">{evalSuccess}</div>
        {/if}

        {#if supervisorEval}
            <div class="eval-submitted">
                <p>Vous avez deja soumis votre evaluation.</p>
                <p>
                    Note attribuee : <strong
                        >{supervisorEval.criterion5}/4</strong
                    >
                </p>
                <!-- no comment field in SupervisorEvaluation type -->
            </div>
        {:else}
            <div class="eval-form">
                <p class="eval-desc">
                    Attribuez une note sur 4 points pour l'encadrement de ce
                    PFE.
                </p>
                <label>
                    Note (/4)
                    <input
                        type="number"
                        min="0"
                        max="4"
                        step="0.5"
                        bind:value={criterion5}
                        class="input"
                        required
                    />
                </label>
                <label>
                    Commentaire (optionnel)
                    <textarea
                        bind:value={evalComment}
                        rows={3}
                        placeholder="Votre avis sur le travail..."
                        class="input"
                    ></textarea>
                </label>
                <Button
                    variant="primary"
                    onclick={submitEvaluation}
                    disabled={evalLoading}
                >
                    {evalLoading ? "Soumission..." : "Soumettre l'evaluation"}
                </Button>
            </div>
        {/if}
    </section>
</Page>

<style>
    .back-link {
        margin-bottom: var(--spacing-md);

        a {
            display: inline-flex;
            align-items: center;
            gap: 0.35rem;
            font-size: var(--text-sm);
            color: var(--color-accent);
            text-decoration: none;
            font-family: var(--font-sans);

            &:hover {
                text-decoration: underline;
            }
        }
    }

    .pfe-header {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        margin-bottom: var(--spacing-lg);
    }

    .pfe-title-section {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: var(--spacing-sm);
        margin-bottom: var(--spacing-md);

        h2 {
            font-size: var(--text-xl);
            font-weight: 700;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
            flex: 1;
        }
    }

    .pfe-meta {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .meta-item {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    section {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        margin-bottom: var(--spacing-lg);
    }

    .section-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: var(--spacing-md);

        h3 {
            font-size: var(--text-lg);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
        }
    }

    .empty-section {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-style: italic;
        font-family: var(--font-sans);
        text-align: center;
        padding: 1.5rem 1rem;
    }

    /* ── Add meeting form ─────────────────────────── */
    .add-meeting {
        margin-bottom: var(--spacing-lg);

        h4 {
            font-size: var(--text-base);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0 0 var(--spacing-md);
        }
    }

    .form-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: var(--spacing-sm) var(--spacing-md);
        margin-bottom: var(--spacing-md);

        @media (max-width: 600px) {
            grid-template-columns: 1fr;
        }
    }

    .field {
        display: flex;
        flex-direction: column;
        gap: 0.3rem;

        label {
            font-family: var(--font-sans);
            font-size: var(--text-sm);
            font-weight: 600;
            color: var(--color-text-muted);
        }
    }

    .full-width {
        grid-column: 1 / -1;
    }

    /* Shared input style - matches teacher page */
    input.input,
    textarea.input,
    select.input {
        width: 100%;
        padding: 0.6rem 0.75rem;
        border: 1px solid var(--color-border);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        background: var(--color-surface);
        color: var(--color-text);
        box-sizing: border-box;
        transition: border-color var(--transition-fast);
        resize: vertical;

        &:focus {
            outline: none;
            border-color: var(--color-accent);
            box-shadow: 0 0 0 2px
                color-mix(in srgb, var(--color-accent) 20%, transparent);
        }

        &:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
    }

    /* ── Table ────────────────────────────────────── */
    .table-wrapper {
        overflow-x: auto;
        margin-top: var(--spacing-sm);
    }

    table {
        width: 100%;
        border-collapse: collapse;
        font-family: var(--font-sans);
        font-size: var(--text-sm);
    }

    th {
        text-align: left;
        padding: 0.6rem 0.75rem;
        font-size: 0.72rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.04em;
        color: var(--color-text-muted);
        border-bottom: 1px solid var(--color-border);
        white-space: nowrap;
    }

    td {
        padding: 0.7rem 0.75rem;
        border-bottom: 1px solid var(--color-border);
        color: var(--color-text);
        vertical-align: top;
    }

    tr:last-child td {
        border-bottom: none;
    }

    .cell-date {
        white-space: nowrap;
        font-weight: 600;
        color: var(--color-text-muted);
    }
    .cell-meta {
        white-space: nowrap;
        color: var(--color-text-muted);
    }
    .cell-topics {
        max-width: 260px;
        line-height: 1.5;
    }
    .cell-obs {
        max-width: 200px;
        line-height: 1.5;
        color: var(--color-text-muted);
        font-style: italic;
    }

    .status-badge {
        display: inline-block;
        padding: 0.2rem 0.55rem;
        border-radius: 999px;
        font-size: 0.72rem;
        font-weight: 600;
        white-space: nowrap;
    }

    .status-a_faire {
        background: color-mix(
            in srgb,
            var(--color-text-muted) 12%,
            transparent
        );
        color: var(--color-text-muted);
    }

    .status-en_cours {
        background: color-mix(
            in srgb,
            var(--color-warning, #f59e0b) 15%,
            transparent
        );
        color: var(--color-warning, #f59e0b);
    }

    .status-termine {
        background: color-mix(
            in srgb,
            var(--color-success, #22c55e) 15%,
            transparent
        );
        color: var(--color-success, #22c55e);
    }

    .error-banner {
        padding: 0.75rem 1rem;
        background: color-mix(in srgb, var(--color-danger) 10%, transparent);
        border: 1px solid
            color-mix(in srgb, var(--color-danger) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-danger);
    }

    .success-banner {
        padding: 0.75rem 1rem;
        background: color-mix(in srgb, var(--color-success) 10%, transparent);
        border: 1px solid
            color-mix(in srgb, var(--color-success) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-success);
        margin-bottom: var(--spacing-md);
    }

    .evaluation-section {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
    }

    .eval-submitted {
        text-align: center;
        padding: 1.5rem;
        color: var(--color-text-muted);
        font-family: var(--font-sans);

        p {
            margin: 0 0 0.35rem;
            font-size: var(--text-sm);
        }

        strong {
            color: var(--color-text);
            font-size: var(--text-lg);
        }
    }

    .eval-form {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);

        label {
            display: flex;
            flex-direction: column;
            gap: 0.35rem;
            font-family: var(--font-sans);
            font-size: var(--text-sm);
            font-weight: 600;
            color: var(--color-text-muted);
        }
    }

    .eval-desc {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        margin: 0;
    }
</style>
