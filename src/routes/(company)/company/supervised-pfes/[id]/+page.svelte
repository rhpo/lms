<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { Calendar, User, Building2, ArrowLeft } from "lucide-svelte";
    import { company } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import FormField from "$lib/components/ui/FormField.svelte";
    import DateInput from "$lib/components/ui/DateInput.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { pfe, progressReports, supervisorEval } = $derived(data);

    let newMeetingDate = $state("");
    let newMeetingNotes = $state("");
    let meetingError = $state("");

    // Evaluation form state
    let criterion5 = $state(0);
    let evalComment = $state("");
    let evalError = $state("");
    let evalSuccess = $state("");
    let evalLoading = $state(false);

    async function submitEvaluation() {
        evalError = "";
        evalSuccess = "";
        if (criterion5 < 0 || criterion5 > 4) {
            evalError = "La note doit etre comprise entre 0 et 4.";
            return;
        }
        evalLoading = true;
        try {
            await company.submitEvaluation(pfe.id, {
                criterion5,
                comment: evalComment,
            });
            evalSuccess = "Evaluation soumise avec succes.";
            await invalidateAll();
        } catch (err: unknown) {
            evalError = err instanceof Error ? err.message : "Erreur reseau.";
        } finally {
            evalLoading = false;
        }
    }

    const pfeStatusLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            en_cours: "En cours",
            soutenance_planifiee: "Soutenance planifiee",
            valide: "Valide",
            refuse: "Refuse",
        };
        return labels;
    });

    const pfeStatusVariant = $derived.by(() => {
        const variants: Record<
            string,
            "info" | "warning" | "success" | "danger"
        > = {
            en_cours: "info",
            soutenance_planifiee: "warning",
            valide: "success",
            refuse: "danger",
        };
        return variants;
    });

    async function addMeetingAction() {
        meetingError = "";
        if (!newMeetingDate || !newMeetingNotes) {
            meetingError = "Veuillez remplir tous les champs.";
            return;
        }
        try {
            await company.addMeeting(pfe.id, {
                meeting_date: newMeetingDate,
                notes: newMeetingNotes,
            });
            newMeetingDate = "";
            newMeetingNotes = "";
            await invalidateAll();
        } catch (err: unknown) {
            meetingError = err instanceof Error ? err.message : "Erreur reseau";
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
            <h2>{pfe.subject?.title ?? "Sujet"}</h2>
            <Badge
                variant={pfeStatusVariant[pfe.status] ?? "info"}
                label={pfeStatusLabel[pfe.status] ?? pfe.status}
            />
        </div>

        <div class="pfe-meta">
            <div class="meta-item">
                <User size={14} />
                <span>
                    Etudiant(s) : {pfe.student_name ?? "Non assigne"}
                    {#if pfe.student2_id}
                        , {pfe.student2_name ?? "Etudiant 2"}
                    {/if}
                    {#if pfe.student3_id}
                        , {pfe.student3_name ?? "Etudiant 3"}
                    {/if}
                </span>
            </div>
            {#if pfe.co_supervisor_id}
                <div class="meta-item">
                    <Building2 size={14} />
                    <span>
                        Co-promoteur : {pfe.co_supervisor_name ?? "Enseignant"}
                    </span>
                </div>
            {/if}
        </div>
    </div>

    <section>
        <div class="section-header">
            <h3>Journal de suivi</h3>
        </div>

        {#if progressReports.length === 0}
            <p class="empty-section">
                Aucune entrée de reunion pour le moment.
            </p>
        {:else}
            <div class="meeting-list">
                {#each progressReports as report}
                    <div class="meeting-card">
                        <div class="meeting-header">
                            <Calendar size={14} />
                            <span class="meeting-date">
                                {new Date(
                                    report.meeting_date,
                                ).toLocaleDateString("fr-FR")}
                            </span>
                        </div>
                        <p class="meeting-notes">
                            {report.supervisor_notes ?? report.student_notes}
                        </p>
                    </div>
                {/each}
            </div>
        {/if}

        <div class="add-meeting">
            <h4>Ajouter une entrée de reunion</h4>
            {#if meetingError}
                <div class="error-banner">{meetingError}</div>
            {/if}
            <FormField label="Date de la reunion" required>
                <DateInput bind:value={newMeetingDate} required />
            </FormField>
            <FormField label="Notes" required>
                <textarea
                    bind:value={newMeetingNotes}
                    rows="3"
                    placeholder="Description de l'avancement..."
                    required
                    class="input"
                ></textarea>
            </FormField>
            <Button variant="primary" onclick={addMeetingAction}>
                Ajouter
            </Button>
        </div>
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
                    Note attribuee : <strong>{supervisorEval.total}/4</strong>
                </p>
                {#if supervisorEval.comment}
                    <p>Commentaire : {supervisorEval.comment}</p>
                {/if}
            </div>
        {:else}
            <div class="eval-form">
                <p class="eval-desc">
                    Attribuez une note sur 4 points pour l'encadrement de ce
                    PFE.
                </p>
                <FormField label="Note (/4)" required>
                    <input
                        type="number"
                        min="0"
                        max="4"
                        step="0.5"
                        bind:value={criterion5}
                        class="input"
                        required
                    />
                </FormField>
                <FormField label="Commentaire (optionnel)">
                    <textarea
                        bind:value={evalComment}
                        rows="3"
                        placeholder="Votre avis sur le travail..."
                        class="input"
                    ></textarea>
                </FormField>
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
        padding: 1rem;
    }

    .meeting-list {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-sm);
        margin-bottom: var(--spacing-lg);
    }

    .meeting-card {
        padding: var(--spacing-md);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-background);
    }

    .meeting-header {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        margin-bottom: 0.5rem;
    }

    .meeting-date {
        font-size: var(--text-xs);
        font-weight: 600;
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .meeting-notes {
        font-size: var(--text-sm);
        color: var(--color-text);
        font-family: var(--font-sans);
        margin: 0;
        line-height: 1.5;
        white-space: pre-wrap;
    }

    .add-meeting {
        border-top: 1px solid var(--color-border);
        padding-top: var(--spacing-lg);
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);

        h4 {
            font-size: var(--text-base);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
        }

        textarea {
            width: 100%;
            padding: 0.5rem 0.75rem;
            border: 1px solid var(--color-border);
            border-radius: 8px;
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            background: var(--color-background);
            color: var(--color-text);
            resize: vertical;
        }
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
    }

    .eval-desc {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        margin: 0;
    }
</style>
