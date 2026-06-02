<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { Link, FileText, AlertTriangle } from "lucide-svelte";
    import { student } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { pfe, deadlinePassed, memoireDeadline } = $derived(data);

    let submitting = $state(false);
    let uploadError = $state("");
    let uploadSuccess = $state(false);
    let memoireUrl = $state("");

    async function handleSubmit(e: Event) {
        e.preventDefault();
        const url = memoireUrl.trim();
        if (!url) {
            uploadError = "Veuillez saisir l'URL de votre mémoire.";
            return;
        }

        try {
            new URL(url);
        } catch {
            uploadError = "L'URL saisie n'est pas valide.";
            return;
        }

        submitting = true;
        uploadError = "";
        uploadSuccess = false;

        try {
            await student.submitMemoire({ memoire_url: url });
            uploadSuccess = true;
            memoireUrl = "";
            await invalidateAll();
        } catch (err: unknown) {
            uploadError =
                err instanceof Error
                    ? err.message
                    : "Erreur réseau. Veuillez réessayer.";
        } finally {
            submitting = false;
        }
    }

    const deadlineFormatted = $derived.by(() => {
        if (!memoireDeadline) return null;
        return new Date(memoireDeadline).toLocaleDateString("fr-FR", {
            day: "numeric",
            month: "long",
            year: "numeric",
        });
    });
</script>

<Page title="Mémoire" subtitle="Soumettez l'URL de votre mémoire de fin d'études.">
    {#if !pfe}
        <div class="empty-state">
            <FileText size={48} />
            <h2>Aucun PFE en cours</h2>
            <p>
                Vous n'avez pas de PFE actif pour lequel soumettre un mémoire.
            </p>
            <Button variant="primary" href="/student/my-pfe">
                Voir mon PFE
            </Button>
        </div>
    {:else}
        <div class="card">
            <div class="card-header">
                <h2>Soumettre le mémoire</h2>
                <Badge
                    variant={pfe.memoire_url ? "success" : "warning"}
                    label={pfe.memoire_url ? "Soumis" : "Non soumis"}
                />
            </div>

            {#if uploadSuccess}
                <div class="success-banner">
                    Votre mémoire a été soumis avec succès. L'encadrant et
                    l'administration ont été notifiés.
                </div>
            {/if}

            {#if uploadError}
                <div class="error-banner">
                    <AlertTriangle size={14} />
                    {uploadError}
                </div>
            {/if}

            {#if deadlinePassed}
                <div class="deadline-banner">
                    <AlertTriangle size={16} />
                    <span>
                        Date dépassée{deadlineFormatted
                            ? ` (${deadlineFormatted})`
                            : ""} : vous ne pouvez plus soumettre votre mémoire.
                    </span>
                </div>
            {:else}
                {#if pfe.memoire_url}
                    <div class="current-file">
                        <FileText size={16} />
                        <span>Mémoire déjà soumis.</span>
                        <a
                            href={pfe.memoire_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            class="link"
                        >
                            Voir le document
                        </a>
                    </div>
                {/if}

                <form onsubmit={handleSubmit} class="url-form">
                    <label for="memoire-url" class="url-label">
                        <Link size={16} />
                        <span>URL du mémoire (lien public vers le PDF)</span>
                    </label>
                    <div class="url-input-row">
                        <input
                            id="memoire-url"
                            type="url"
                            bind:value={memoireUrl}
                            placeholder="https://drive.google.com/..."
                            class="input"
                            disabled={submitting}
                            required
                        />
                        <Button
                            variant="primary"
                            type="submit"
                            disabled={submitting || !memoireUrl.trim()}
                        >
                            {submitting ? "Envoi…" : "Soumettre"}
                        </Button>
                    </div>
                    <p class="url-hint">
                        Hébergez votre mémoire sur Google Drive, OneDrive ou tout autre service et collez le lien public ici.
                    </p>
                </form>

                {#if deadlineFormatted}
                    <p class="deadline-info">
                        Date limite de remise : {deadlineFormatted}
                    </p>
                {/if}
            {/if}
        </div>
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

    .card {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
    }

    .card-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: var(--spacing-lg);

        h2 {
            font-size: var(--text-lg);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
        }
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
        margin-bottom: var(--spacing-md);
    }

    .error-banner {
        display: flex;
        align-items: center;
        gap: 0.5rem;
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

    .deadline-banner {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: color-mix(
            in srgb,
            var(--color-warning) 10%,
            var(--color-surface)
        );
        border: 1px solid
            color-mix(in srgb, var(--color-warning) 30%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-warning);
        margin-bottom: var(--spacing-md);
    }

    .current-file {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: var(--color-background);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text);
        margin-bottom: var(--spacing-md);

        a {
            color: var(--color-accent);
            text-decoration: underline;
        }
    }

    .url-form {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        margin-top: var(--spacing-md);
    }

    .url-label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        font-weight: 500;
        color: var(--color-text);
    }

    .url-input-row {
        display: flex;
        gap: 0.75rem;
        align-items: center;
    }

    .input {
        flex: 1;
        padding: 0.5rem 0.75rem;
        border: 1px solid var(--color-border);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        background: var(--color-background);
        color: var(--color-text);
        outline: none;

        &:focus {
            border-color: var(--color-accent);
            box-shadow: 0 0 0 2px color-mix(in srgb, var(--color-accent) 20%, transparent);
        }
    }

    .url-hint {
        font-size: var(--text-xs);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin: 0;
    }

    .deadline-info {
        font-size: var(--text-xs);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin-top: var(--spacing-sm);
        text-align: center;
    }
</style>
