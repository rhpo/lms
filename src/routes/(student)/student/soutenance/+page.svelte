<script lang="ts">
    import { Calendar, Clock, MapPin, Users, Award } from "lucide-svelte";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { defense, grades } = $derived(data);

    const defenseStatusLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            scheduled: "Planifiee",
            done: "Passee",
            postponed: "Reportee",
        };
        return labels;
    });

    const defenseStatusVariant = $derived.by(() => {
        const variants: Record<string, "warning" | "success" | "info"> = {
            scheduled: "warning",
            done: "success",
            postponed: "info",
        };
        return variants;
    });

    const defenseResultLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            admitted: "Admis",
            corrections_required: "Corrections requises",
            not_admitted: "Non admis",
        };
        return labels;
    });

    const defenseResultVariant = $derived.by(() => {
        const variants: Record<string, "success" | "warning" | "danger"> = {
            admitted: "success",
            corrections_required: "warning",
            not_admitted: "danger",
        };
        return variants;
    });
</script>

<Page
    title="Ma soutenance"
    subtitle="Informations sur votre soutenance de PFE."
>
    {#if !defense}
        <div class="empty-state">
            <Calendar size={48} />
            <h2>Aucune soutenance planifiee</h2>
            <p>
                Votre soutenance n'a pas encore ete planifiee. Vous serez
                notifie des qu'une date sera fixee.
            </p>
        </div>
    {:else}
        <div class="defense-header">
            <div class="defense-status">
                <Badge
                    variant={defenseStatusVariant[defense.status] ?? "info"}
                    label={defenseStatusLabel[defense.status] ?? defense.status}
                />
                {#if defense.result}
                    <Badge
                        variant={defenseResultVariant[defense.result] ?? "info"}
                        label={defenseResultLabel[defense.result] ?? defense.result}
                    />
                {/if}
            </div>

            <div class="defense-datetime">
                <div class="datetime-item">
                    <Calendar size={18} />
                    <span>
                        {new Date(defense.scheduled_at).toLocaleDateString(
                            "fr-FR",
                            {
                                weekday: "long",
                                year: "numeric",
                                month: "long",
                                day: "numeric",
                            },
                        )}
                    </span>
                </div>
                <div class="datetime-item">
                    <Clock size={18} />
                    <span>
                        {new Date(defense.scheduled_at).toLocaleTimeString(
                            "fr-FR",
                            {
                                hour: "2-digit",
                                minute: "2-digit",
                            },
                        )}
                    </span>
                </div>
                <div class="datetime-item">
                    <MapPin size={18} />
                    <span>{defense.room}</span>
                </div>
            </div>
        </div>

        <section>
            <h3>
                <Users size={16} />
                Composition du jury
            </h3>
            <div class="jury-grid">
                <div class="jury-card">
                    <span class="jury-role">President</span>
                    <span class="jury-name">
                        {defense.president_name ?? "Non assigne"}
                    </span>
                </div>
                <div class="jury-card">
                    <span class="jury-role">Membre</span>
                    <span class="jury-name">
                        {defense.member_name ?? "Non assigne"}
                    </span>
                </div>
            </div>
        </section>

        {#if defense.final_grade !== null}
            <section>
                <h3>
                    <Award size={16} />
                    Note finale
                </h3>
                <div class="grade-section">
                    <div class="final-grade">
                        <span class="grade-value">
                            {defense.final_grade}/20
                        </span>
                    </div>

                    {#if grades.length > 0}
                        <div class="grade-details">
                            <h4>Detail des notes</h4>
                            {#each grades as grade}
                                <div class="grade-card">
                                    <span class="grade-jury">
                                        {grade.jury_member_name ?? "Membre"}
                                    </span>
                                    <div class="grade-criteria">
                                        <div class="criterion">
                                            <span class="criterion-label">
                                                Criteres
                                            </span>
                                            <span class="criterion-value">
                                                {grade.criterion1} + {grade.criterion2} + {grade.criterion3} + {grade.criterion4}
                                            </span>
                                        </div>
                                        <div class="criterion total">
                                            <span class="criterion-label">
                                                Total
                                            </span>
                                            <span class="criterion-value">
                                                {grade.total}/20
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            </section>
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

    .defense-header {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        margin-bottom: var(--spacing-lg);
    }

    .defense-status {
        display: flex;
        gap: var(--spacing-sm);
        margin-bottom: var(--spacing-md);
    }

    .defense-datetime {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .datetime-item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: var(--text-base);
        color: var(--color-text);
        font-family: var(--font-sans);

        span {
            font-weight: 500;
        }
    }

    section {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        margin-bottom: var(--spacing-lg);
    }

    h3 {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: var(--text-lg);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text);
        margin: 0 0 var(--spacing-md);
    }

    .jury-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: var(--spacing-md);

        @media screen and (max-width: 480px) {
            & {
                grid-template-columns: 1fr;
            }
        }
    }

    .jury-card {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
        padding: var(--spacing-md);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-background);
    }

    .jury-role {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .jury-name {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
        font-family: var(--font-sans);
    }

    .grade-section {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-lg);
    }

    .final-grade {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 1.5rem;
    }

    .grade-value {
        font-size: 2.5rem;
        font-weight: 700;
        font-family: var(--font-sans);
        color: var(--color-accent);
    }

    .grade-details {
        h4 {
            font-size: var(--text-sm);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text-muted);
            margin: 0 0 var(--spacing-sm);
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }
    }

    .grade-card {
        padding: var(--spacing-md);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-background);
        margin-bottom: var(--spacing-sm);
    }

    .grade-jury {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
        font-family: var(--font-sans);
        display: block;
        margin-bottom: 0.5rem;
    }

    .grade-criteria {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
    }

    .criterion {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .criterion-label {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    .criterion-value {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
        font-family: var(--font-sans);
    }

    .criterion.total .criterion-value {
        color: var(--color-accent);
    }
</style>
