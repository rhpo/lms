<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { Plus, Calendar, Clock, Tag, FileText } from "lucide-svelte";

    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";

    let { data } = $props();
    const { pfe, progressReports } = $derived(data);

    let showForm = $state(false);
    let submitting = $state(false);
    let submitError = $state("");

    let form = $state({
        meeting_date: "",
        student_notes: "",
    });

    function resetForm() {
        form = { meeting_date: "", student_notes: "" };
        submitError = "";
    }

    async function handleSubmit(e: Event) {
        e.preventDefault();
        if (!form.meeting_date) {
            submitError = "La date est obligatoire.";
            return;
        }
        submitting = true;
        submitError = "";

        try {
            if (!pfe) {
                submitError =
                    "Aucun PFE trouve. Impossible d'ajouter un suivi.";
                return;
            }
            const res = await fetch("/api/progress-reports/create", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    assignmentId: pfe.id,
                    meeting_date: form.meeting_date,
                    student_notes: form.student_notes || "",
                }),
            });
            if (res.ok) {
                resetForm();
                showForm = false;
                await invalidateAll();
            } else {
                const err = await res.json();
                submitError = err.message ?? "Erreur lors de l'ajout du suivi.";
            }
        } catch {
            submitError = "Erreur reseau. Veuillez reessayer.";
        } finally {
            submitting = false;
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
</script>

<Page title="Suivi / Meetings" subtitle="Journal de suivi de votre PFE.">
    {#if !pfe}
        <div class="empty-state">
            <FileText size={48} />
            <h2>Aucun PFE en cours</h2>
            <p>Vous n'avez pas de PFE actif.</p>
            <Button variant="primary" href="/student/my-pfe">
                Voir mon PFE
            </Button>
        </div>
    {:else}
        <div class="toolbar">
            <Button variant="primary" onclick={() => (showForm = !showForm)}>
                <Plus size={14} />
                {showForm ? "Annuler" : "Ajouter un suivi"}
            </Button>
        </div>

        {#if showForm}
            <div class="form-card">
                <h3>Nouvelle entree de suivi</h3>
                {#if submitError}
                    <div class="error-banner">{submitError}</div>
                {/if}
                <form onsubmit={handleSubmit}>
                    <label>
                        <span>Date et heure</span>
                        <input
                            type="datetime-local"
                            bind:value={form.meeting_date}
                            required
                        />
                    </label>
                    <label>
                        <span>Notes (optionnel)</span>
                        <textarea
                            bind:value={form.student_notes}
                            placeholder="Resume de la reunion, points abordes..."
                            rows={4}
                        ></textarea>
                    </label>
                    <div class="form-actions">
                        <Button
                            variant="secondary"
                            type="button"
                            onclick={() => {
                                showForm = false;
                                resetForm();
                            }}
                        >
                            Annuler
                        </Button>
                        <Button
                            variant="primary"
                            type="submit"
                            disabled={submitting}
                        >
                            {submitting ? "Envoi..." : "Enregistrer"}
                        </Button>
                    </div>
                </form>
            </div>
        {/if}

        {#if progressReports.length === 0}
            <div class="empty-reports">
                <Calendar size={40} />
                <p>Aucune entree de suivi pour le moment.</p>
            </div>
        {:else}
            <div class="reports-list">
                {#each progressReports as report (report.id)}
                    <div class="report-card">
                        <div class="report-header">
                            <div class="report-meta">
                                <Calendar size={14} />
                                <span>{formatDate(report.meeting_date)}</span>
                            </div>
                            <Badge
                                variant={report.signed_by_teacher
                                    ? "success"
                                    : "warning"}
                                label={report.signed_by_teacher
                                    ? "Signe"
                                    : "En attente"}
                            />
                        </div>
                        {#if report.student_notes}
                            <p class="report-notes">
                                <Tag size={14} />
                                {report.student_notes}
                            </p>
                        {/if}
                        {#if report.teacher_feedback || report.supervisor_notes}
                            <div class="teacher-feedback">
                                <strong>Retour de l'encadrant :</strong>
                                <p>
                                    {report.teacher_feedback ||
                                        report.supervisor_notes}
                                </p>
                            </div>
                        {/if}
                    </div>
                {/each}
            </div>
        {/if}
    {/if}
</Page>

<style>
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 1rem;
        padding: 3rem;
        text-align: center;
        color: var(--color-text-muted);

        h2 {
            font-size: var(--text-lg);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
        }

        p {
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            margin: 0;
            max-width: 400px;
        }
    }

    .toolbar {
        display: flex;
        justify-content: flex-end;
        margin-bottom: var(--spacing-md);
    }

    .form-card {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        margin-bottom: var(--spacing-lg);

        h3 {
            font-size: var(--text-base);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0 0 var(--spacing-md) 0;
        }

        form {
            display: flex;
            flex-direction: column;
            gap: var(--spacing-md);
        }

        label {
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            color: var(--color-text);

            span {
                font-weight: 500;
            }

            input,
            textarea {
                padding: 0.625rem 0.75rem;
                border: 1px solid var(--color-border);
                border-radius: 8px;
                font-size: var(--text-sm);
                font-family: var(--font-sans);
                background: var(--color-background);
                color: var(--color-text);
            }

            textarea {
                resize: vertical;
            }
        }
    }

    .form-actions {
        display: flex;
        justify-content: flex-end;
        gap: var(--spacing-sm);
    }

    .error-banner {
        padding: 0.75rem 1rem;
        background: color-mix(
            in srgb,
            var(--color-error) 10%,
            var(--color-surface)
        );
        border: 1px solid
            color-mix(in srgb, var(--color-error) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-error);
        margin-bottom: var(--spacing-md);
    }

    .empty-reports {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.75rem;
        padding: 3rem;
        color: var(--color-text-muted);
        font-size: var(--text-sm);
        font-family: var(--font-sans);
    }

    .reports-list {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);
    }

    .report-card {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-md);
    }

    .report-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: var(--spacing-sm);
    }

    .report-meta {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
    }

    .report-notes {
        display: flex;
        align-items: flex-start;
        gap: 0.5rem;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text);
        margin: var(--spacing-sm) 0;
        padding: 0.75rem;
        background: var(--color-background);
        border-radius: 8px;
    }

    .teacher-feedback {
        margin-top: var(--spacing-sm);
        padding: 0.75rem;
        background: color-mix(
            in srgb,
            var(--color-accent) 5%,
            var(--color-background)
        );
        border: 1px solid
            color-mix(in srgb, var(--color-accent) 15%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text);

        strong {
            display: block;
            margin-bottom: 0.25rem;
            font-weight: 600;
            color: var(--color-accent);
        }

        p {
            margin: 0;
        }
    }
</style>
