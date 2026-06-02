<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import {
        FileText,
        User,
        Building2,
        Upload,
        CheckCircle,
        Award,
    } from "lucide-svelte";
    import { student } from "$lib/api";
    import { upload } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import FormField from "$lib/components/ui/FormField.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { pfe, progressReports, defense, supervisorNote } = $derived(data);

    const hasFinalGrade = $derived(
        !!(defense?.final_grade !== null && defense?.final_grade !== undefined),
    );

    let showMeetingForm = $state(false);
    let meetingError = $state("");
    let meetingSubmitting = $state(false);

    let newMeeting = $state({
        meeting_date: "",
        duration: 60,
        meeting_type: "presentiel" as "presentiel" | "visio",
        topics: "",
        status: "a_faire" as "a_faire" | "en_cours" | "termine",
        observation: "",
    });

    const DURATION_OPTIONS = [
        { value: 30, label: "30 min" },
        { value: 60, label: "1 h" },
        { value: 90, label: "1 h 30" },
        { value: 120, label: "2 h" },
        { value: 150, label: "2 h 30" },
        { value: 180, label: "3 h" },
    ];

    const MEETING_TYPE_LABELS: Record<string, string> = {
        presentiel: "Présentiel",
        visio: "Visio",
    };

    const STATUS_LABELS: Record<string, string> = {
        a_faire: "À faire",
        en_cours: "En cours",
        termine: "Terminé",
    };

    let memoireFile = $state<File | null>(null);
    let memoireUploading = $state(false);
    let memoireError = $state("");
    let memoireSuccess = $state(false);

    const hasMemoire = $derived(!!pfe?.memoire_url);

    async function submitMemoireAction() {
        if (!memoireFile) {
            memoireError = "Veuillez sélectionner un fichier PDF.";
            return;
        }
        if (!memoireFile.name.toLowerCase().endsWith(".pdf")) {
            memoireError = "Seuls les fichiers PDF sont acceptés.";
            return;
        }
        if (memoireFile.size > 50 * 1024 * 1024) {
            memoireError = "Le fichier ne doit pas dépasser 50 Mo.";
            return;
        }
        memoireError = "";
        memoireUploading = true;
        try {
            const formData = new FormData();
            formData.append("file", memoireFile);
            const { url } = await upload.memoire(formData);

            await student.submitMemoire({ memoire_url: url });

            memoireSuccess = true;
            memoireFile = null;
            await invalidateAll();
        } catch (err: unknown) {
            memoireError =
                err instanceof Error
                    ? err.message
                    : "Erreur lors du dépôt du mémoire.";
        } finally {
            memoireUploading = false;
        }
    }

    let updatingStatusId = $state<number | null>(null);

    async function updateMeetingStatus(reportId: number, newStatus: string) {
        updatingStatusId = reportId;
        try {
            await student.updateMyMeeting(reportId, { status: newStatus });
            await invalidateAll();
        } catch {
        } finally {
            updatingStatusId = null;
        }
    }

    function formatDuration(minutes: number): string {
        if (minutes < 60) return `${minutes} min`;
        const h = Math.floor(minutes / 60);
        const m = minutes % 60;
        return m > 0 ? `${h} h ${m}` : `${h} h`;
    }

    function resetMeetingForm() {
        newMeeting = {
            meeting_date: "",
            duration: 60,
            meeting_type: "presentiel",
            topics: "",
            status: "a_faire",
            observation: "",
        };
        meetingError = "";
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

    async function addMeetingAction(e: Event) {
        e.preventDefault();
        meetingError = "";
        if (!newMeeting.meeting_date) {
            meetingError = "La date est obligatoire.";
            return;
        }
        meetingSubmitting = true;
        try {
            await student.addMyMeeting({
                meeting_date: newMeeting.meeting_date,
                duration: newMeeting.duration,
                meeting_type: newMeeting.meeting_type,
                topics: newMeeting.topics,
                status: newMeeting.status,
                observation: newMeeting.observation || undefined,
            });
            resetMeetingForm();
            showMeetingForm = false;
            await invalidateAll();
        } catch (err: unknown) {
            meetingError =
                err instanceof Error ? err.message : "Erreur réseau.";
        } finally {
            meetingSubmitting = false;
        }
    }
</script>

<Page
    title="Mon PFE"
    subtitle="Suivez l'avancement de votre projet de fin d'etudes."
>
    {#if !pfe}
        <div class="empty-state">
            <FileText size={48} />
            <h2>Vous n'avez pas encore de PFE</h2>
            <p>
                Consultez le catalogue pour soumettre vos voeux et attendre
                qu'un encadreur vous accepte.
            </p>
            <Button variant="primary" href="/student/catalogue">
                Consulter le catalogue
            </Button>
        </div>
    {:else}
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
                        Encadreur : {pfe.supervisor?.profile?.full_name ??
                            "Non assigne"}
                    </span>
                </div>
                {#if pfe.co_supervisor_id}
                    <div class="meta-item">
                        <Building2 size={14} />
                        <span>
                            Co-promoteur : {pfe.co_supervisor?.profile
                                ?.full_name ?? "Entreprise"}
                        </span>
                    </div>
                {/if}
            </div>
        </div>

        <section class="memoire-section">
            <div class="section-header">
                <h3>Mémoire</h3>
            </div>

            <!-- Final grade banner -->
            {#if hasFinalGrade && defense}
                <div class="final-grade-banner">
                    <div class="fg-left">
                        <Award size={20} />
                        <div>
                            <p class="fg-title">Note finale disponible</p>
                            <p class="fg-sub">
                                Votre soutenance a été évaluée. Note :
                                <strong
                                    >{typeof defense.final_grade === "number"
                                        ? defense.final_grade.toFixed(2)
                                        : defense.final_grade} / 20</strong
                                >
                                - {typeof defense.final_grade === "number" &&
                                defense.final_grade >= 10
                                    ? "✓ Admis"
                                    : "✗ Non admis"}
                            </p>
                        </div>
                    </div>
                    <Button
                        variant="primary"
                        href="/student/soutenance"
                        size="sm"
                    >
                        Voir le détail
                    </Button>
                </div>
            {/if}

            {#if hasMemoire}
                <div class="memoire-done">
                    <CheckCircle size={20} />
                    <div>
                        <p class="memoire-done-title">Mémoire déposé</p>
                        <p class="memoire-done-hint">
                            Le mémoire a été soumis avec succès. Un seul dépôt
                            par groupe est nécessaire.
                        </p>
                    </div>
                </div>
            {:else}
                <p class="memoire-hint">
                    Un seul membre du {pfe.subject?.group_type === "monome"
                        ? "monôme"
                        : pfe.subject?.group_type === "binome"
                          ? "binôme"
                          : "trinôme"} doit déposer le mémoire (PDF uniquement).
                </p>

                {#if memoireSuccess}
                    <div class="success-banner">
                        Mémoire déposé avec succès.
                    </div>
                {/if}

                {#if memoireError}
                    <div class="error-banner">{memoireError}</div>
                {/if}

                <div class="memoire-upload">
                    <label class="file-input-label" for="memoire-input">
                        <Upload size={16} />
                        {memoireFile
                            ? memoireFile.name
                            : "Choisir un fichier PDF"}
                    </label>
                    <input
                        id="memoire-input"
                        type="file"
                        accept=".pdf,application/pdf"
                        class="file-input-hidden"
                        onchange={(e) => {
                            const input = e.target as HTMLInputElement;
                            memoireFile = input.files?.[0] ?? null;
                            memoireError = "";
                            memoireSuccess = false;
                        }}
                    />
                    <Button
                        variant="primary"
                        onclick={submitMemoireAction}
                        disabled={!memoireFile || memoireUploading}
                    >
                        {memoireUploading
                            ? "Envoi en cours..."
                            : "Déposer le mémoire"}
                    </Button>
                </div>
            {/if}
        </section>

        <section>
            <div class="section-header">
                <h3>Journal de suivi</h3>
                <Button
                    variant="primary"
                    size="sm"
                    onclick={() => (showMeetingForm = !showMeetingForm)}
                >
                    {showMeetingForm ? "Annuler" : "Ajouter un suivi"}
                </Button>
            </div>

            {#if showMeetingForm}
                <div class="meeting-form-card">
                    <h4>Nouvelle entrée de suivi</h4>
                    {#if meetingError}
                        <div class="error-banner">{meetingError}</div>
                    {/if}
                    <form onsubmit={addMeetingAction}>
                        <div class="meeting-form-grid">
                            <FormField label="Date" required>
                                <input
                                    type="date"
                                    bind:value={newMeeting.meeting_date}
                                    required
                                    class="input"
                                />
                            </FormField>

                            <FormField label="Durée" required>
                                <select
                                    bind:value={newMeeting.duration}
                                    class="input"
                                >
                                    {#each DURATION_OPTIONS as opt}
                                        <option value={opt.value}
                                            >{opt.label}</option
                                        >
                                    {/each}
                                </select>
                            </FormField>

                            <FormField label="Type de réunion" required>
                                <select
                                    bind:value={newMeeting.meeting_type}
                                    class="input"
                                >
                                    <option value="presentiel"
                                        >Présentiel</option
                                    >
                                    <option value="visio">Visio</option>
                                </select>
                            </FormField>

                            <FormField label="État" required>
                                <select
                                    bind:value={newMeeting.status}
                                    class="input"
                                >
                                    <option value="a_faire">À faire</option>
                                    <option value="en_cours">En cours</option>
                                    <option value="termine">Terminé</option>
                                </select>
                            </FormField>
                        </div>

                        <FormField label="Points discutés">
                            <textarea
                                bind:value={newMeeting.topics}
                                placeholder="Résumé des points abordés lors de la réunion..."
                                rows={3}
                                class="input"
                            ></textarea>
                        </FormField>

                        <FormField label="Observation">
                            <textarea
                                bind:value={newMeeting.observation}
                                placeholder="Remarques supplémentaires (optionnel)..."
                                rows={2}
                                class="input"
                            ></textarea>
                        </FormField>

                        <div class="meeting-form-actions">
                            <Button
                                variant="ghost"
                                type="button"
                                onclick={() => {
                                    showMeetingForm = false;
                                    resetMeetingForm();
                                }}
                            >
                                Annuler
                            </Button>
                            <Button
                                variant="primary"
                                type="submit"
                                disabled={meetingSubmitting}
                            >
                                {meetingSubmitting
                                    ? "Enregistrement..."
                                    : "Enregistrer"}
                            </Button>
                        </div>
                    </form>
                </div>
            {/if}

            {#if progressReports.length === 0}
                <p class="empty-section">
                    Aucune entrée de suivi pour le moment.
                </p>
            {:else}
                <div class="meeting-table-wrapper">
                    <table class="meeting-table">
                        <thead>
                            <tr>
                                <th>N° réunion</th>
                                <th>Date</th>
                                <th>Durée</th>
                                <th>Type de réunion</th>
                                <th>Points discutés</th>
                                <th>État</th>
                                <th>Observation</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each [...progressReports] as report, i (report.id)}
                                <tr>
                                    <td class="meeting-num">{i + 1}</td>
                                    <td
                                        >{new Date(
                                            report.meeting_date,
                                        ).toLocaleDateString("fr-FR")}</td
                                    >
                                    <td>{formatDuration(report.duration)}</td>
                                    <td>
                                        <Badge
                                            variant="info"
                                            label={MEETING_TYPE_LABELS[
                                                report.meeting_type
                                            ] ?? report.meeting_type}
                                        />
                                    </td>
                                    <td class="meeting-topics"
                                        >{report.topics || "-"}</td
                                    >
                                    <td>
                                        <select
                                            class="meeting-status-select meeting-status-{report.status}"
                                            value={report.status}
                                            disabled={updatingStatusId ===
                                                report.id}
                                            onchange={(e) =>
                                                updateMeetingStatus(
                                                    report.id,
                                                    (
                                                        e.target as HTMLSelectElement
                                                    ).value,
                                                )}
                                        >
                                            <option value="a_faire"
                                                >À faire</option
                                            >
                                            <option value="en_cours"
                                                >En cours</option
                                            >
                                            <option value="termine"
                                                >Terminé</option
                                            >
                                        </select>
                                    </td>
                                    <td class="meeting-observation"
                                        >{report.observation ?? "-"}</td
                                    >
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            {/if}
        </section>
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

    .meeting-form-card {
        background: var(--color-background);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        padding: var(--spacing-md);
        margin-bottom: var(--spacing-md);

        h4 {
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
    }

    .meeting-form-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: var(--spacing-md);
    }

    .meeting-form-actions {
        display: flex;
        justify-content: flex-end;
        gap: var(--spacing-sm);
        padding-top: var(--spacing-sm);
    }

    .meeting-table-wrapper {
        overflow-x: auto;
    }

    .meeting-table {
        width: 100%;
        border-collapse: collapse;
        font-family: var(--font-sans);
        font-size: var(--text-sm);

        th {
            text-align: left;
            padding: 0.75rem 1rem;
            font-weight: 600;
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.04em;
            color: var(--color-text-muted);
            border-bottom: 2px solid var(--color-border);
            white-space: nowrap;
        }

        td {
            padding: 0.75rem 1rem;
            color: var(--color-text);
            border-bottom: 1px solid var(--color-border);
            vertical-align: top;
        }

        tr:last-child td {
            border-bottom: none;
        }

        tr:hover td {
            background: var(--color-surface);
        }
    }

    .meeting-num {
        text-align: center;
        font-weight: 600;
        color: var(--color-text-muted);
    }

    .meeting-topics {
        max-width: 240px;
        white-space: pre-wrap;
        word-break: break-word;
    }

    .meeting-observation {
        max-width: 180px;
        color: var(--color-text-muted);
        font-size: var(--text-xs);
        white-space: pre-wrap;
        word-break: break-word;
    }

    .meeting-status-select {
        appearance: none;
        border: 1px solid var(--color-border);
        border-radius: 6px;
        padding: 0.25rem 0.6rem;
        font-size: var(--text-xs);
        font-family: var(--font-sans);
        font-weight: 500;
        cursor: pointer;
        transition: opacity 0.15s;

        &:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
    }

    .meeting-status-a_faire {
        background: color-mix(
            in srgb,
            var(--color-text-muted) 12%,
            transparent
        );
        color: var(--color-text-muted);
        border-color: color-mix(
            in srgb,
            var(--color-text-muted) 25%,
            transparent
        );
    }

    .meeting-status-en_cours {
        background: color-mix(in srgb, var(--color-warning) 12%, transparent);
        color: var(--color-warning);
        border-color: color-mix(in srgb, var(--color-warning) 30%, transparent);
    }

    .meeting-status-termine {
        background: color-mix(in srgb, var(--color-success) 12%, transparent);
        color: var(--color-success);
        border-color: color-mix(in srgb, var(--color-success) 30%, transparent);
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
    }

    .success-banner {
        padding: 0.75rem 1rem;
        background: color-mix(
            in srgb,
            var(--color-success) 10%,
            var(--color-surface)
        );
        border: 1px solid
            color-mix(in srgb, var(--color-success) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-success);
        margin-bottom: var(--spacing-sm);
    }

    .memoire-hint {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        margin: 0 0 var(--spacing-md);
    }

    .memoire-upload {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
    }

    .file-input-hidden {
        display: none;
    }

    .file-input-label {
        display: inline-flex;
        align-items: center;
        gap: 0.4rem;
        padding: 0.5rem 1rem;
        border: 1px dashed var(--color-border);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        cursor: pointer;
        transition: all var(--transition-fast);
        flex: 1;
        min-width: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;

        &:hover {
            border-color: var(--color-accent);
            color: var(--color-accent);
        }
    }

    .memoire-done {
        display: flex;
        align-items: flex-start;
        gap: var(--spacing-sm);
        padding: var(--spacing-md);
        background: color-mix(
            in srgb,
            var(--color-success) 8%,
            var(--color-surface)
        );
        border: 1px solid
            color-mix(in srgb, var(--color-success) 20%, transparent);
        border-radius: 8px;
        color: var(--color-success);
    }

    .memoire-done-title {
        font-size: var(--text-sm);
        font-weight: 600;
        font-family: var(--font-sans);
        margin: 0;
    }

    .memoire-done-hint {
        font-size: var(--text-xs);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin: 0.25rem 0 0;
    }
</style>
