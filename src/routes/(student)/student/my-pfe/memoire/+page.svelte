<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { Upload, FileText, AlertTriangle } from "lucide-svelte";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { pfe, deadlinePassed, memoireDeadline } = $derived(data);

    let uploading = $state(false);
    let uploadError = $state("");
    let uploadSuccess = $state(false);

    async function handleFileChange(e: Event) {
        const input = e.target as HTMLInputElement;
        const file = input.files?.[0];
        if (!file) return;

        if (file.type !== "application/pdf") {
            uploadError = "Seuls les fichiers PDF sont acceptes.";
            return;
        }

        if (file.size > 10 * 1024 * 1024) {
            uploadError = "Le fichier ne doit pas depasser 10 Mo.";
            return;
        }

        uploading = true;
        uploadError = "";
        uploadSuccess = false;

        try {
            const formData = new FormData();
            formData.append("file", file);

            const res = await fetch(`/api/pfe/upload-memoire`, {
                method: "POST",
                body: formData,
            });

            if (res.ok) {
                uploadSuccess = true;
                await invalidateAll();
            } else {
                const err = await res.json();
                uploadError =
                    err.message ?? "Erreur lors de l'envoi du memoire.";
            }
        } catch {
            uploadError = "Erreur reseau. Veuillez reessayer.";
        } finally {
            uploading = false;
            input.value = "";
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

<Page title="Memoire" subtitle="Soumettez votre memoire de fin d'etudes.">
    {#if !pfe}
        <div class="empty-state">
            <FileText size={48} />
            <h2>Aucun PFE en cours</h2>
            <p>
                Vous n'avez pas de PFE actif pour lequel soumettre un memoire.
            </p>
            <Button variant="primary" href="/student/my-pfe">
                Voir mon PFE
            </Button>
        </div>
    {:else}
        <div class="card">
            <div class="card-header">
                <h2>Soumettre le memoire</h2>
                <Badge
                    variant={pfe.memoire_url ? "success" : "warning"}
                    label={pfe.memoire_url ? "Soumis" : "Non soumis"}
                />
            </div>

            {#if uploadSuccess}
                <div class="success-banner">
                    Votre memoire a ete soumis avec succes. L'encadrant et
                    l'administration ont ete notifies.
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
                        Date depassee{deadlineFormatted
                            ? ` (${deadlineFormatted})`
                            : ""} : vous ne pouvez plus soumettre votre Memoire.
                    </span>
                </div>
                <Button
                    variant="primary"
                    disabled={true}
                    title="Date depassee, vous ne pouvez plus soumettre votre Memoire."
                >
                    <Upload size={14} />
                    Soumettre le memoire
                </Button>
            {:else}
                {#if pfe.memoire_url}
                    <div class="current-file">
                        <FileText size={16} />
                        <span>Memoire deja soumis.</span>
                        <a
                            href={pfe.memoire_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            class="link"
                        >
                            Voir le fichier
                        </a>
                    </div>
                {/if}

                <div class="upload-zone">
                    <label class="upload-label" for="memoire-upload">
                        <Upload size={24} />
                        <span class="upload-text">
                            {uploading
                                ? "Envoi en cours..."
                                : "Cliquez pour selectionner un fichier PDF"}
                        </span>
                        <span class="upload-hint"
                            >PDF uniquement, max 10 Mo</span
                        >
                    </label>
                    <input
                        id="memoire-upload"
                        type="file"
                        accept=".pdf"
                        disabled={uploading}
                        onchange={handleFileChange}
                        class="file-input"
                    />
                </div>

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

    .upload-zone {
        border: 2px dashed var(--color-border);
        border-radius: 12px;
        padding: 2rem;
        text-align: center;
        cursor: pointer;
        transition:
            border-color 0.2s,
            background 0.2s;

        &:hover {
            border-color: var(--color-accent);
            background: color-mix(in srgb, var(--color-accent) 5%, transparent);
        }
    }

    .upload-label {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
        cursor: pointer;
        color: var(--color-text-muted);
    }

    .upload-text {
        font-size: var(--text-sm);
        font-family: var(--font-sans);
    }

    .upload-hint {
        font-size: var(--text-xs);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
    }

    .file-input {
        display: none;
    }

    .deadline-info {
        font-size: var(--text-xs);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin-top: var(--spacing-sm);
        text-align: center;
    }
</style>
