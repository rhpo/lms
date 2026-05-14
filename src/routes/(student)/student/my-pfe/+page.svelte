<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import {
        FileText,
        Calendar,
        User,
        Building2,
        MessageSquare,
    } from "lucide-svelte";
    import { student } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";
    import FormField from "$lib/components/ui/FormField.svelte";
    import DateInput from "$lib/components/ui/DateInput.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { pfe, progressReports } = $derived(data);

    let showClaimModal = $state(false);
    let claimMessage = $state("");
    let claimError = $state("");
    let claimSuccess = $state(false);

    let newMeetingDate = $state("");
    let newMeetingNotes = $state("");
    let meetingError = $state("");

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
            await student.addMyMeeting({
                meeting_date: newMeetingDate,
                student_notes: newMeetingNotes,
            });
            newMeetingDate = "";
            newMeetingNotes = "";
            await invalidateAll();
        } catch (err: unknown) {
            meetingError = err instanceof Error ? err.message : "Erreur reseau";
        }
    }

    async function submitClaimAction() {
        claimError =
            "Erreur: Fonctionnalité non disponible. Contactez l'administration.";
        claimSuccess = false;
        if (!claimMessage) {
            claimError = "Veuillez saisir un message.";
            return;
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
                        Encadreur : {pfe.supervisor_name ?? "Non assigne"}
                    </span>
                </div>
                {#if pfe.co_supervisor_id}
                    <div class="meta-item">
                        <Building2 size={14} />
                        <span>
                            Co-promoteur : {pfe.co_supervisor_name ??
                                "Entreprise"}
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
                            <p class="meeting-notes">{report.student_notes}</p>
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

        <div class="actions-section">
            <Button variant="ghost" onclick={() => (showClaimModal = true)}>
                <MessageSquare size={14} />
                Soumettre une reclamation
            </Button>
        </div>
    {/if}
</Page>

<Modal
    open={showClaimModal}
    title="Soumettre une reclamation"
    onClose={() => {
        showClaimModal = false;
        claimError = "";
        claimSuccess = false;
    }}
>
    {#if claimSuccess}
        <div class="success-banner">
            Votre reclamation a ete envoyee a l'administration.
        </div>
    {:else}
        {#if claimError}
            <div class="error-banner">{claimError}</div>
        {/if}
        <FormField label="Message" required>
            <textarea
                bind:value={claimMessage}
                rows="5"
                placeholder="Expliquez votre reclamation..."
                required
                class="input"
            ></textarea>
        </FormField>
        <div class="form-actions">
            <Button
                variant="ghost"
                onclick={() => {
                    showClaimModal = false;
                    claimError = "";
                }}
            >
                Annuler
            </Button>
            <Button variant="primary" onclick={submitClaimAction}>
                Envoyer
            </Button>
        </div>
    {/if}
</Modal>

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

    .actions-section {
        display: flex;
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
    }

    .form-actions {
        display: flex;
        justify-content: flex-end;
        gap: var(--spacing-sm);
        margin-top: var(--spacing-lg);
    }
</style>
