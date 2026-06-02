<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { goto } from "$app/navigation";
    import { Heart, X, FileText, Lock } from "lucide-svelte";
    import { student } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { wishes, maxWishes, hasPfe } = $derived(data);

    let removeError = $state("");

    const wishStatusLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            en_attente: "En attente",
            accepte: "Accepté",
            refuse: "Refusé",
        };
        return labels;
    });

    const wishStatusVariant = $derived.by(() => {
        const variants: Record<
            string,
            "warning" | "success" | "danger" | "info"
        > = {
            en_attente: "warning",
            accepte: "success",
            refuse: "danger",
        };
        return variants;
    });

    const acceptedWish = $derived(
        wishes.find((w: { status: string }) => w.status === "accepte"),
    );

    async function removeWishAction(wishId: number) {
        removeError = "";
        try {
            await student.deleteWish(wishId);
            await invalidateAll();
        } catch (err: unknown) {
            removeError = err instanceof Error ? err.message : "Erreur reseau";
        }
    }
</script>

<Page title="Mes voeux" subtitle="Consultez et Gérer vos voeux de sujets PFE.">
    {#if hasPfe}
        <div class="locked-banner">
            <Lock size={16} />
            <span>
                Vous avez déjà un PFE affecté. Vos vœux sont en lecture seule.
                <a href="/student/my-pfe">Voir mon PFE</a>
            </span>
        </div>
    {:else}
        <div class="quota-bar">
            <Heart size={16} />
            <span>
                {wishes.length} voeu{wishes.length !== 1 ? "x" : ""} sur {maxWishes}
                maximum
            </span>
        </div>
    {/if}

    {#if acceptedWish}
        <div class="accepted-banner">
            <FileText size={16} />
            <span>
                Vous avez été accept sur un sujet.
                <a href="/student/my-pfe">Voir mon PFE</a>
            </span>
        </div>
    {/if}

    {#if removeError}
        <div class="error-banner">{removeError}</div>
    {/if}

    {#if wishes.length === 0}
        <p class="empty">
            Vous n'avez soumis aucun voeu pour le moment.
            <a href="/student/catalogue">Consulter le catalogue</a>
        </p>
    {:else}
        <div class="wish-list">
            {#each wishes as wish}
                {@const subjectTitle =
                    wish.subject?.title ?? wish.subject_title ?? null}
                {@const proposerName =
                    wish.subject?.proposer?.full_name ?? null}
                {@const domains = wish.subject?.domains ?? []}
                <div class="wish-card">
                    <div class="wish-info">
                        <h3>
                            <a href="/student/catalogue/{wish.subject_id}">
                                {subjectTitle ?? "Chargement..."}
                            </a>
                        </h3>
                        {#if proposerName || domains.length > 0}
                            <div class="wish-subject-meta">
                                {#if proposerName}
                                    <span class="subject-proposer"
                                        >Proposé par {proposerName}</span
                                    >
                                {/if}
                                {#if domains.length > 0}
                                    <span class="subject-domains"
                                        >{domains
                                            .map(
                                                (d: { name: string }) => d.name,
                                            )
                                            .join(", ")}</span
                                    >
                                {/if}
                            </div>
                        {/if}
                        <div class="wish-meta">
                            <Badge
                                variant={wishStatusVariant[wish.status] ??
                                    "info"}
                                label={wishStatusLabel[wish.status] ??
                                    wish.status}
                            />
                            <span class="wish-date">
                                Ajouté le{" "}
                                {new Date(wish.created_at).toLocaleDateString(
                                    "fr-FR",
                                )}
                            </span>
                        </div>
                    </div>

                    {#if wish.status === "en_attente" && !hasPfe}
                        <Button
                            variant="ghost"
                            size="sm"
                            onclick={() => removeWishAction(wish.id)}
                        >
                            <X size={14} />
                            Retirer
                        </Button>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</Page>

<style>
    .locked-banner {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: color-mix(in srgb, var(--color-accent) 8%, var(--color-surface));
        border: 1px solid color-mix(in srgb, var(--color-accent) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-accent);
        margin-bottom: var(--spacing-lg);

        a {
            color: var(--color-accent);
            font-weight: 600;
            text-decoration: underline;
        }
    }

    .quota-bar {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin-bottom: var(--spacing-lg);
    }

    .accepted-banner {
        display: flex;
        align-items: center;
        gap: 0.5rem;
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

        a {
            color: var(--color-success);
            font-weight: 600;
            text-decoration: underline;
        }
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

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        font-style: italic;
        font-size: var(--text-sm);
        padding: 2rem;
        font-family: var(--font-sans);

        a {
            color: var(--color-accent);
            text-decoration: underline;
        }
    }

    .wish-list {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-sm);
    }

    .wish-card {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: var(--spacing-md);
        padding: var(--spacing-md);
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        transition: box-shadow var(--transition-normal);

        &:hover {
            box-shadow: var(--shadow-sm);
        }
    }

    .wish-info {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
        flex: 1;
        min-width: 0;

        h3 {
            font-size: var(--text-base);
            font-weight: 600;
            font-family: var(--font-sans);
            margin: 0;

            a {
                color: var(--color-text);
                text-decoration: none;

                &:hover {
                    color: var(--color-accent);
                    text-decoration: underline;
                }
            }
        }
    }

    .wish-subject-meta {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        flex-wrap: wrap;
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    .subject-proposer::after {
        content: "·";
        margin-left: 0.5rem;
    }

    .subject-domains:last-child ~ .subject-proposer::after {
        content: "";
    }

    .subject-domains {
        font-style: italic;
    }

    .wish-meta {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        flex-wrap: wrap;
    }

    .wish-date {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }
</style>
