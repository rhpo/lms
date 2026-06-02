<script lang="ts">
    import {
        Calendar,
        Clock,
        MapPin,
        Users,
        Award,
        Archive,
    } from "lucide-svelte";
    import type { ArchiveDecision } from "$lib/types";

    const ARCHIVE_LABELS: Record<ArchiveDecision, string> = {
        archivable: "Le mémoire peut être archivé",
        minor_corrections: "Peut être archivé après des corrections mineures",
        major_corrections:
            "Ne peut être archivé - corrections majeures requises",
    };

    import Badge from "$lib/components/ui/Badge.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { defense, grades, supervisorNote } = $derived(data);

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
    {#if !defense || !defense.scheduled_at}
        <div class="empty-state">
            <Calendar size={48} />
            <h2>Soutenance non planifiée</h2>
            <p>
                Votre soutenance n'a pas encore été planifiée. Vous serez
                notifié dès qu'une date sera fixée par l'administration.
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
                        label={defenseResultLabel[defense.result] ??
                            defense.result}
                    />
                {/if}
            </div>

            <div class="defense-datetime">
                <div class="datetime-item">
                    <Calendar size={18} />
                    <span>
                        {new Date(defense.scheduled_at!).toLocaleDateString(
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
                        {new Date(defense.scheduled_at!).toLocaleTimeString(
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
                        {defense.jury?.president?.profile?.full_name ??
                            "Non assigne"}
                    </span>
                </div>
                <div class="jury-card">
                    <span class="jury-role">Membre</span>
                    <span class="jury-name">
                        {defense.jury?.member?.profile?.full_name ??
                            "Non assigne"}
                    </span>
                </div>
            </div>
        </section>

        <!-- Archivage section: shown only when at least one jury member submitted -->
        {#if grades.length > 0 || defense.final_grade !== null}
            <section>
                <h3><Archive size={16} /> Archivage du document</h3>
                {#if grades.length > 0}
                    <div class="archive-decisions">
                        {#each grades as g}
                            {#if g.archive_decision}
                                <div class="archive-row">
                                    <span class="archive-member"
                                        >{g.jury_member?.profile?.full_name ??
                                            "Membre"}</span
                                    >
                                    <span class="archive-label-val"
                                        >{ARCHIVE_LABELS[
                                            g.archive_decision
                                        ]}</span
                                    >
                                </div>
                            {/if}
                        {/each}
                    </div>
                {:else}
                    <p class="pending-note">Décision d'archivage en attente.</p>
                {/if}
            </section>
        {/if}

        <section>
            <h3>
                <Award size={16} />
                Note finale
            </h3>
            {#if defense.final_grade !== null && defense.final_grade !== undefined}
                <div class="grade-section">
                    <div class="final-grade">
                        <span class="grade-value"
                            >{typeof defense.final_grade === "number"
                                ? defense.final_grade.toFixed(2)
                                : defense.final_grade}/20</span
                        >
                        <span class="grade-sub"
                            >{typeof defense.final_grade === "number" &&
                            defense.final_grade >= 10
                                ? "✓ Admis"
                                : "✗ Non admis"}</span
                        >
                    </div>
                    <div class="grade-breakdown">
                        <div class="criterion">
                            <span class="criterion-label"
                                >Jury (critères C1–C4)</span
                            >
                            <span class="criterion-value">
                                {supervisorNote !== null
                                    ? (
                                          (defense.final_grade ?? 0) -
                                          supervisorNote
                                      ).toFixed(2)
                                    : "-"} / 16
                            </span>
                        </div>
                        <div class="criterion supervisor-row">
                            <span class="criterion-label"
                                >Note encadrant (C5)</span
                            >
                            <span class="criterion-value">
                                {#if supervisorNote !== null}
                                    {supervisorNote} / 4
                                {:else}
                                    <em class="not-submitted">Non disponible</em
                                    >
                                {/if}
                            </span>
                        </div>
                        <div class="criterion total-row">
                            <span class="criterion-label"
                                ><strong>Total</strong></span
                            >
                            <span class="criterion-value">
                                <strong
                                    >{typeof defense.final_grade === "number"
                                        ? defense.final_grade.toFixed(2)
                                        : defense.final_grade} / 20</strong
                                >
                            </span>
                        </div>
                    </div>
                </div>
            {:else if supervisorNote !== null}
                <p class="pending-note">
                    Note encadrant reçue ({supervisorNote}/4). En attente des
                    notes du jury.
                </p>
            {:else}
                <p class="pending-note">
                    Les notes n'ont pas encore été soumises par le jury.
                </p>
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
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 1.5rem;
        gap: 0.25rem;
    }

    .grade-value {
        font-size: 2.5rem;
        font-weight: 700;
        font-family: var(--font-sans);
        color: var(--color-accent);
    }

    .grade-sub {
        font-size: var(--text-sm);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text-muted);
    }

    .total-row {
        background: color-mix(
            in srgb,
            var(--color-accent) 8%,
            var(--color-surface)
        );
    }

    .archive-decisions {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .archive-row {
        display: flex;
        flex-direction: column;
        gap: 0.1rem;
        padding: 0.6rem 0.9rem;
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-background-100);
    }

    .archive-member {
        font-size: var(--text-xs);
        font-weight: 600;
        color: var(--color-text-muted);
        text-transform: uppercase;
        letter-spacing: 0.04em;
    }

    .archive-label-val {
        font-size: var(--text-sm);
        color: var(--color-text);
        font-weight: 500;
    }

    .grade-breakdown {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        border: 1px solid var(--color-border);
        border-radius: 8px;
        overflow: hidden;
    }

    .criterion {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.6rem 1rem;
        border-bottom: 1px solid var(--color-border);

        &:last-child {
            border-bottom: none;
        }
    }

    .criterion.supervisor-row {
        background: color-mix(
            in srgb,
            var(--color-accent) 5%,
            var(--color-surface)
        );
    }

    .criterion-label {
        font-size: var(--text-sm);
        color: var(--color-text);
        font-family: var(--font-sans);
    }

    .criterion-value {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
        font-family: var(--font-sans);
        white-space: nowrap;
    }

    .not-submitted {
        font-weight: 400;
        color: var(--color-text-muted);
        font-style: italic;
    }

    .pending-note {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        font-style: italic;
        margin: 0;
        padding: 1rem;
        text-align: center;
    }
</style>
