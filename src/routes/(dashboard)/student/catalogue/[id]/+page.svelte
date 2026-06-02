<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { goto } from "$app/navigation";
    import { Users, Building2, GraduationCap, Heart } from "lucide-svelte";
    import { student } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";
    import { GROUP_TYPE_LABELS } from "$lib/constants/labels";

    let { data } = $props();

    const {
        subject,
        alreadyWished,
        alreadyAssigned,
        subjectTaken,
        wishesCount,
    } = $derived(data);

    let addError = $state("");

    const groupTypeLabel = $derived(
        subject
            ? (GROUP_TYPE_LABELS[subject.group_type] ?? subject.group_type)
            : "",
    );

    async function addWishAction() {
        addError = "";
        if (!subject) return;
        try {
            await student.createWish({ subject_id: subject.id });
            await invalidateAll();
        } catch (err: unknown) {
            addError = err instanceof Error ? err.message : "Erreur reseau";
        }
    }
</script>

{#if !subject}
    <Page
        title="Sujet introuvable"
        subtitle="Ce sujet n'existe pas ou a ete supprime."
    >
        <p class="error-message">Impossible de charger le sujet.</p>
    </Page>
{:else}
    <Page title={subject.title} subtitle="Detail complet du sujet PFE">
        <div class="detail-layout">
            <div class="detail-main">
                <section>
                    <h2>Description</h2>
                    <p class="description">{subject.description}</p>
                </section>

                <section>
                    <h2>Informations</h2>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="info-label">Domaines</span>
                            <span class="info-value"
                                >{subject.domains
                                    ?.map((d) => d.name)
                                    .join(", ") || "-"}</span
                            >
                        </div>
                        <div class="info-item">
                            <span class="info-label">Type de groupe</span>
                            <Badge variant="info" label={groupTypeLabel} />
                        </div>
                        <div class="info-item">
                            <span class="info-label">Encadreur</span>
                            <span class="info-value">
                                {#if subject.proposer_role === "company"}
                                    <Building2 size={14} />
                                {:else}
                                    <Users size={14} />
                                {/if}
                                {subject.proposer?.full_name ??
                                    subject.proposer_role}
                            </span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">Type de groupe</span>
                            <span class="info-value">{groupTypeLabel}</span>
                        </div>
                    </div>
                </section>
            </div>

            <aside class="detail-sidebar">
                <div class="sidebar-card">
                    <h3>Actions</h3>

                    {#if alreadyAssigned}
                        <p class="info-message">
                            Vous avez déjà un PFE affecté. La soumission de vœux
                            est désactivée.
                        </p>
                        <Button variant="ghost" href="/student/my-pfe">
                            Voir mon PFE
                        </Button>
                    {:else if subjectTaken}
                        <p class="info-message is-taken">
                            Ce sujet est déjà pris.
                        </p>
                    {:else if alreadyWished}
                        <p class="info-message is-success">
                            <Heart size={16} />
                            Ce sujet est déjà dans vos vœux.
                        </p>
                        <Button
                            variant="ghost"
                            onclick={() => goto("/student/voeux")}
                        >
                            Voir mes vœux
                        </Button>
                    {:else}
                        <Button variant="primary" onclick={addWishAction}>
                            Ajouter à mes vœux
                        </Button>
                    {/if}

                    {#if addError}
                        <p class="error-message">{addError}</p>
                    {/if}
                </div>
            </aside>
        </div>
    </Page>
{/if}

<style>
    .detail-layout {
        display: grid;
        grid-template-columns: 1fr 300px;
        gap: var(--spacing-lg);

        @media screen and (max-width: 768px) {
            & {
                grid-template-columns: 1fr;
            }
        }
    }

    .detail-main {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-lg);
    }

    section {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
    }

    h2 {
        font-size: var(--text-lg);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text);
        margin: 0 0 var(--spacing-md);
    }

    .description {
        font-size: var(--text-sm);
        color: var(--color-text);
        font-family: var(--font-sans);
        line-height: 1.7;
        margin: 0;
        white-space: pre-wrap;
    }

    .info-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: var(--spacing-md);

        @media screen and (max-width: 480px) {
            & {
                grid-template-columns: 1fr;
            }
        }
    }

    .info-item {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .info-label {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .info-value {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
        font-family: var(--font-sans);
    }

    .sidebar-card {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);

        h3 {
            font-size: var(--text-base);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
        }
    }

    .info-message {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        margin: 0;
        padding: 0.75rem;
        border-radius: 8px;
        background: var(--color-background);
    }

    .info-message.is-success {
        color: var(--color-success);
        background: color-mix(
            in srgb,
            var(--color-success) 10%,
            var(--color-surface)
        );
    }

    .info-message.is-taken {
        color: var(--color-danger);
        background: color-mix(
            in srgb,
            var(--color-danger) 10%,
            var(--color-surface)
        );
    }

    .error-message {
        font-size: var(--text-sm);
        color: var(--color-error);
        font-family: var(--font-sans);
        margin: 0;
    }
</style>
